package db

import (
	"gitlab.xinghuolive.com/birds-backend/phoenix/models/pgmodels"
	"time"
)

type Schema struct {
	pgmodels.Model

	tableName  struct{}  `pg:"common.migration,alias:migration,discard_unknown_columns"`
	ModuleName string    `pg:"module_name"` // 模块名称
	FullName   string    `pg:"full_name"`   // 运行模块的全程（make_migration)生成
	ExecTime   time.Time `pg:"exec_time"`   // 运行时间
	Version    string    `pg:"version"`     // 迁移版本
	Force      bool      `pg:"force"`       // 是否为强制运行
	Success    bool      `pg:"success"`     // 是否运行成功（不报错即为运行成功)
	IsFake     bool      `pg:"is_fake"`     // 是否是人工标记为执行成功
}

func (s Schema) CreateTable(db pg.DB) error {
	_, err := db.Exec(`
create table if not exists common.migration
(
	id bigserial primary key,
	created_at timestamp with time zone default timezone('utc'::text, now()),
	updated_at timestamp with time zone default timezone('utc'::text, now()),
	deleted_at timestamp with time zone,
	is_delete boolean,
	module_name varchar(255),
	full_name  varchar(255),
	exec_time timestamp,
	version   varchar(255),
	force    boolean,
	success  boolean,
	is_fake  boolean
)
;
`)
	return err
}

func (s Schema) CreateIndex(db pg.DB) error {
	var commandList = []string{
		ddl.Index("migration", "created_at"),
	}
	return pgu.ExecCommand(db, &s, commandList)
}

func (s Schema) CreateTrigger(db pg.DB) error {
	var commandList = []string{
		ddl.CreateUpdatedAtTriggerIfNotExist("common", "migration"),
	}
	return pgu.ExecCommand(db, &s, commandList)
}

func (s Schema) CreateConstraint(db pg.DB) error {
	return nil
}
