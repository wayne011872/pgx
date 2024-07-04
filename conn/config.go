package conn

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type PgxDI interface {
	NewPgxConn(ctx context.Context) (PgxConn, error)
	SetAuth(user, pwd string)
	GetUri() string
}

type PgxConf struct {
	Uri       string `yaml:"uri"`
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	DefaultDB string `yaml:"defaul"`
	SslMode   string `yaml:"sslmode"`

	authUri  string	
}

func (pc *PgxConf) SetAuth(user, pwd string) {
	pc.authUri = strings.Replace(pc.Uri, "{User}", user, 1)
	pc.authUri = strings.Replace(pc.authUri, "{Pwd}", pwd, 1)
}

func (pc *PgxConf) GetUri() string {
	if pc.Uri == "" {
		panic("Postgresql uri not set")
	}
	if pc.DefaultDB == "" {
		panic("Postgresql default db not set")
	}
	return pc.Uri + "/" + pc.DefaultDB + "?sslmode="+pc.SslMode
}
func (pc *PgxConf) GetAuthUri() string {
	if pc.Uri == "" {
		panic("Postgresql uri not set")
	}
	if pc.DefaultDB == "" {
		panic("Postgresql default db not set")
	}
	pc.SetAuth(pc.User, pc.Pass)
	return pc.authUri+ "/" + pc.DefaultDB + "?sslmode="+pc.SslMode
}

func (pc *PgxConf) GetDb()string{
	return pc.DefaultDB
}

func (pc *PgxConf) NewPgxConn(ctx context.Context) (PgxConn, error) {
	var uri string
	if pc.User == "" || pc.Pass == "" {
		uri = pc.GetUri()
	}else{
		uri = pc.GetAuthUri()
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	pgx, err := pgx.Connect(ctx, uri)
	if err != nil {
		cancel()
		return nil, err
	}
	err = pgx.Ping(ctx)
	if err != nil {
		cancel()
		return nil, err
	}

	return &pgxClientImpl{
		ctx:     ctx,
		cancel:  cancel,
		pgx: pgx,
	}, nil

}