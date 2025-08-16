package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"

	"oauth2-server/internal/config"
	"oauth2-server/internal/model"
)

var (
	dumpvar bool
	portvar int
)

func init() {
	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
}

func main() {
	flag.Parse()

	// 加载配置
	var c config.Config
	conf.MustLoad("etc/oauth2-api.yaml", &c)

	if dumpvar {
		log.Println("Dumping requests")
	}

	// 创建OAuth2管理器
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// 使用Redis存储token
	// redisStore := redis.MustNewRedis(c.Redis)
	// 暂时使用内存存储，后续可以改为Redis存储
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// 生成JWT访问令牌
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(c.Auth.AccessSecret), jwt.SigningMethodHS512))
	// manager.MapAccessGenerate(generates.NewAccessGenerate())

	// 创建数据库连接
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	clientModel := model.NewClientModel(conn)

	// 创建客户端存储
	clientStore := store.NewClientStore()

	// 从数据库加载客户端
	clients, err := clientModel.FindAll(context.Background())

	if err == nil {
		for _, client := range clients {
			log.Println(client.ID, " ", client.Secret, " ", client.RedirectURL)
			clientStore.Set(client.ID, &models.Client{
				ID:     client.ID,
				Secret: client.Secret,
				Domain: client.RedirectURL,
			})
		}
	}

	manager.MapClientStorage(clientStore)

	// 创建OAuth2服务器
	srv := server.NewServer(server.NewConfig(), manager)

	// 设置密码授权处理器
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		// 这里应该实现真实的用户认证逻辑
		if username == "test" && password == "test" {
			userID = "test_user"
		} else {
			err = errors.New("invalid username or password")
		}
		return
	})

	// 设置用户授权处理器
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	// 设置内部错误处理器
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	// 设置响应错误处理器
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	// 创建HTTP服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册路由
	registerRoutes(server, srv, clientModel, c)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func registerRoutes(server *rest.Server, srv *server.Server, clientModel model.ClientModel, c config.Config) {
	// 客户端注册接口
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/client/register",
		Handler: clientRegisterHandler(clientModel),
	})

	// 登录页面
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/login",
		Handler: loginHandler,
	})

	// 登录页面
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/login",
		Handler: loginHandler,
	})

	// 授权页面
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/auth",
		Handler: authHandler,
	})

	// OAuth2授权端点
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/oauth/authorize",
		Handler: authorizeHandler(srv),
	})

	// OAuth2授权端点
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/oauth/authorize",
		Handler: authorizeHandler(srv),
	})

	// OAuth2令牌端点
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/oauth/token",
		Handler: tokenHandler(srv),
	})

	// 用户信息端点
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/oauth/userinfo",
		Handler: userInfoHandler(srv),
	})
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r)
	}
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	log.Println("userAuthorizeHandler uid: ", uid)
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "login", r)
	}
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if r.Form.Get("username") == "test" && r.Form.Get("password") == "test" {
			store.Set("LoggedInUserID", r.Form.Get("username"))
			store.Save()

			w.Header().Set("Location", "/auth")
			w.WriteHeader(http.StatusFound)
			return
		} else {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "auth", r)
	}
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func clientRegisterHandler(clientModel model.ClientModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name        string `json:"name"`
			RedirectURL string `json:"redirect_url"`
			GrantType   string `json:"grant_type"`
			Scope       string `json:"scope"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client := &model.Client{
			Name:        req.Name,
			RedirectURL: req.RedirectURL,
			GrantType:   req.GrantType,
			Scope:       req.Scope,
		}

		_, err := clientModel.Insert(r.Context(), client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"client_id":     client.ID,
			"client_secret": client.Secret,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func authorizeHandler(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			dumpRequest(os.Stdout, "authorize", r)
		}

		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		r.Form = form

		store.Delete("ReturnUri")
		store.Save()

		err = srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func tokenHandler(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "token", r)
		}

		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func userInfoHandler(srv *server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpvar {
			_ = dumpRequest(os.Stdout, "userinfo", r)
		}

		token, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := map[string]interface{}{
			"userid":   token.GetUserID(),
			"username": "test_user",
			"phone":    "13800138000",
		}

		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	}
}
