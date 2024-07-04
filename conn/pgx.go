package conn

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)
type ctxKey string

const (
	_CTX_KEY_PGX = "ctxPgxKey"
	CtxPgxKey    = ctxKey(_CTX_KEY_PGX)
)

type PgxConn interface {
	Ping() error
	Close() error
	GetPgxConn() *pgx.Conn
}
func SetPgxConnToGin(c *gin.Context, clt PgxConn) *gin.Context {
	c.Set(_CTX_KEY_PGX, clt)
	return c
}
func GetPgxConnFromGin(c *gin.Context) PgxConn {
	cltInter ,_:= c.Get(string(CtxPgxKey))
	if dbclt, ok := cltInter.(PgxConn); ok {
		return dbclt
	}
	return nil
}

func SetPgxConnToReq(req *http.Request, clt PgxConn) *http.Request {
	return req.WithContext(SetPgxConnToCtx(req.Context(), clt))
}
func GetPgxConnFromReq(req *http.Request) PgxConn {
	return GetPgxConnFromCtx(req.Context())
}
func SetPgxConnToCtx(ctx context.Context, clt PgxConn) context.Context {
	return context.WithValue(ctx, CtxPgxKey, clt)
}
func GetPgxConnFromCtx(ctx context.Context) PgxConn {
	cltInter := ctx.Value(CtxPgxKey)
	if dbclt, ok := cltInter.(PgxConn); ok {
		return dbclt
	}
	return nil
}
type pgxClientImpl struct {
	ctx             context.Context
	cancel          context.CancelFunc
	pgx         	*pgx.Conn
}

func (p *pgxClientImpl) GetPgxConn() *pgx.Conn {
	return p.pgx
}

func (p *pgxClientImpl) Ping() error {
	return p.pgx.Ping(p.ctx)
}

func (p *pgxClientImpl) Close() error{
	err :=p.pgx.Close(p.ctx)
	if err != nil {
		return err
	}
	p.cancel()
	return nil
}