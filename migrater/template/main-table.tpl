package main

import (
	"gitlab.xinghuolive.com/birds-backend/phoenix/middleware"
	"gitlab.xinghuolive.com/birds-backend/phoenix/models/createtables"
	"os"
)

func main() {
	infra := middleware.NewInfra("migration")
	createtables.Execute(infra.DB, []createtables.DataDefinition{
		//publishversion.Schema{},
	}...)
	os.Exit(0)
}