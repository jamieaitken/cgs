package mysql

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	maxLifetime = time.Second * 60
	maxOpenCons = 0
	maxIdleCons = 5
)

type MySQL struct {
	maxLifetime time.Duration
	maxOpenCons int
	maxIdleCons int
	client      *sql.DB
	addr        string
}

type Option func(*MySQL)

func New(addr string, opts ...Option) (*MySQL, error) {
	m := &MySQL{
		maxLifetime: maxLifetime,
		maxIdleCons: maxIdleCons,
		maxOpenCons: maxOpenCons,
		addr:        addr,
	}

	m.add(opts...)

	client, err := sql.Open("mysql", m.addr)
	if err != nil {
		return nil, err
	}

	client.SetConnMaxLifetime(m.maxLifetime)
	client.SetMaxOpenConns(m.maxOpenCons)
	client.SetMaxIdleConns(m.maxIdleCons)

	m.client = client

	return m, nil
}

func (m *MySQL) add(opts ...Option) {
	for _, opt := range opts {
		opt(m)
	}
}

func (m *MySQL) Client() *sql.DB {
	return m.client
}

func (m *MySQL) Addr() string {
	return m.addr
}

func (m *MySQL) MaxLifetime() time.Duration {
	return m.maxLifetime
}

func (m *MySQL) MaxIdleCons() int {
	return m.maxIdleCons
}

func (m *MySQL) MaxOpenCons() int {
	return m.maxOpenCons
}
