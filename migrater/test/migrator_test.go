package test

import (
	"gitlab.xinghuolive.com/birds-backend/phoenix/middleware"
	"migrations/migrate/runner/migrator"
	"migrations/migrate/runner/migrator/db"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMigrate(t *testing.T) {
	m := migrator.Migrator{}
	version := "v9.2.0"
	m.Version = version
	m.VersionPath = "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
	_ = os.Mkdir(m.VersionPath, os.ModePerm)
	m.InputModules = make([]migrator.InputModule, 0)
	m.InputModules = append(m.InputModules, migrator.InputModule{ModuleName: "vvvv", FullName: "1_171727218_vvvv"})
	err := m.Migrate()
	if err != nil {
		t.Fatalf("do migrate error")
	}
}

func TestInsertRecord(t *testing.T) {
	m := migrator.Migrator{}
	m.Infra = middleware.NewInfra("")
	m.RunRecrod = make([]db.Schema, 0)
	newRecord := db.Schema{}
	newRecord.Force = false
	newRecord.ModuleName = "testInsert"
	newRecord.FullName = "1" + strconv.FormatInt(time.Now().Unix(), 10) + "testInsert"
	newRecord.ExecTime = time.Now()
	newRecord.Success = true
	newRecord.Version = "v66666"
	m.RunRecrod = append(m.RunRecrod, newRecord)
	m.InsertRecords()
}
