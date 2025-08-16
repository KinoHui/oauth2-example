package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AuthorizationModel interface {
	Insert(ctx context.Context, data *Authorization) (sql.Result, error)
	FindOne(ctx context.Context, id int64) (*Authorization, error)
	FindByCode(ctx context.Context, code string) (*Authorization, error)
	Update(ctx context.Context, data *Authorization) error
	Delete(ctx context.Context, id int64) error
}

type defaultAuthorizationModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewAuthorizationModel(conn sqlx.SqlConn) AuthorizationModel {
	return &defaultAuthorizationModel{
		conn:  conn,
		table: "`authorization`",
	}
}

func (m *defaultAuthorizationModel) Insert(ctx context.Context, data *Authorization) (sql.Result, error) {
	// 生成授权码
	if data.Code == "" {
		data.Code = uuid.New().String()
	}

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	query := `insert into ` + m.table + ` (` + authorizationRowsExpectAutoSet + `) values (?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.ExecCtx(ctx, query, data.ClientID, data.UserID, data.Scope, data.Code, data.Status, data.CreatedAt, data.UpdatedAt)
}

func (m *defaultAuthorizationModel) FindOne(ctx context.Context, id int64) (*Authorization, error) {
	query := `select ` + authorizationRows + ` from ` + m.table + ` where id = ? limit 1`
	var resp Authorization
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

func (m *defaultAuthorizationModel) FindByCode(ctx context.Context, code string) (*Authorization, error) {
	query := `select ` + authorizationRows + ` from ` + m.table + ` where code = ? limit 1`
	var resp Authorization
	err := m.conn.QueryRowCtx(ctx, &resp, query, code)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAuthorizationModel) Update(ctx context.Context, data *Authorization) error {
	data.UpdatedAt = time.Now()
	query := `update ` + m.table + ` set ` + authorizationRowsWithPlaceHolder + ` where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, data.ClientID, data.UserID, data.Scope, data.Code, data.Status, data.UpdatedAt, data.ID)
	return err
}

func (m *defaultAuthorizationModel) Delete(ctx context.Context, id int64) error {
	query := `delete from ` + m.table + ` where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

var (
	authorizationRows                = "id, client_id, user_id, scope, code, status, created_at, updated_at"
	authorizationRowsExpectAutoSet   = "client_id, user_id, scope, code, status, created_at, updated_at"
	authorizationRowsWithPlaceHolder = "client_id = ?, user_id = ?, scope = ?, code = ?, status = ?, updated_at = ?"
)
