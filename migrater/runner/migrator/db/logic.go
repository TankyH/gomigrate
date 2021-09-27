package db

import (
	"gitlab.xinghuolive.com/birds-backend/chameleon/store/postgres/pg"
	"gitlab.xinghuolive.com/birds-backend/phoenix/pgu"
	"gitlab.xinghuolive.com/birds-backend/phoenix/pgu/meta"
)

var (
	S Schema
	m meta.Meta
)

func init() {
	m.Init(&S)
}

type Mapper struct {
	*pgu.Mapper
	*Schema
}

func New(db pg.DB) *Mapper {
	m := &Mapper{
		pgu.New(db, &m),
		&S,
	}
	return m
}
