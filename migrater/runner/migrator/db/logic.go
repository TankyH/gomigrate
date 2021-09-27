package db

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
