package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ClientModel interface {
	Insert(ctx context.Context, data *Client) (sql.Result, error)
	FindOne(ctx context.Context, id string) (*Client, error)
	FindAll(ctx context.Context) ([]*Client, error)
	Update(ctx context.Context, data *Client) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*Client, error)
}

type defaultClientModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewClientModel(conn sqlx.SqlConn) ClientModel {
	return &defaultClientModel{
		conn:  conn,
		table: "`client`",
	}
}

func (m *defaultClientModel) Insert(ctx context.Context, data *Client) (sql.Result, error) {
	// 生成客户端ID和密钥
	if data.ID == "" {
		data.ID = "client_" + uuid.New().String()[:8]
	}
	if data.Secret == "" {
		data.Secret = uuid.New().String()
	}

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	query := `insert into ` + m.table + ` (` + clientRowsExpectAutoSet + `) values (?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.ExecCtx(ctx, query, data.ID, data.Secret, data.Name, data.RedirectURL, data.GrantType, data.Scope, data.CreatedAt, data.UpdatedAt)
}

func (m *defaultClientModel) FindOne(ctx context.Context, id string) (*Client, error) {
	query := `select ` + clientRows + ` from ` + m.table + ` where id = ? limit 1`
	var resp Client
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultClientModel) FindAll(ctx context.Context) ([]*Client, error) {
	query := `select ` + clientRows + ` from ` + m.table
	var resp []*Client
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *defaultClientModel) Update(ctx context.Context, data *Client) error {
	data.UpdatedAt = time.Now()
	query := `update ` + m.table + ` set ` + clientRowsWithPlaceHolder + ` where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, data.Secret, data.Name, data.RedirectURL, data.GrantType, data.Scope, data.UpdatedAt, data.ID)
	return err
}

func (m *defaultClientModel) Delete(ctx context.Context, id string) error {
	query := `delete from ` + m.table + ` where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultClientModel) FindByID(ctx context.Context, id string) (*Client, error) {
	return m.FindOne(ctx, id)
}

var (
	clientRows                = "id, secret, name, redirect_url, grant_type, scope, created_at, updated_at"
	clientRowsExpectAutoSet   = "id, secret, name, redirect_url, grant_type, scope, created_at, updated_at"
	clientRowsWithPlaceHolder = "secret = ?, name = ?, redirect_url = ?, grant_type = ?, scope = ?, updated_at = ?"
)

var ErrNotFound = sql.ErrNoRows
