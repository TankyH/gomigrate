package errors

import (
	"errors"
)

var (
	IndexNotExistError      = errors.New("index.json not exist")
	IndexNotMatchError      = errors.New("index.json not match exist folder")
	VersionExist            = errors.New("version has been created, skip create folder")
	ManualCreateError       = errors.New("module is manual create, not in index.json")
	ModuleNotExistError     = errors.New("module not exist, but in index.json")
	ModuleAlreayCreateError = errors.New("module is already create, no need to make migrate, change the script inside")
	VersionNotExist         = errors.New("version migration is not created yet")
	AlreadyRunError         = errors.New("Migration Already Run")
	GetMigrationError       = errors.New("Get Migration History Error")
)
