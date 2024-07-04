package conn

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	PgxOpsQueued = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pgxPool_connection",
		Help: "Number of pgx connection.",
	})
)
type PgxOptsDI interface {
	NewDbConnWithOpts(ctx context.Context, db string) (PgxConn, error)
	SetAuth(user, pwd string)
	GetUri() string
	GetDb() string
}

func (pc *PgxConf) NewDbConnWithOpts(ctx context.Context, db string) (PgxConn, error) {
	result, err := pc.NewPgxConn(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxOptsImpl{
		PgxConn: result,
	}, nil
}

type pgxOptsImpl struct {
	PgxConn
}

func (p *pgxOptsImpl) Close() error{
	err := p.PgxConn.Close()
	if err != nil {
		return err
	}
	return nil
}