package main

import (
	"gitlab.xinghuolive.com/birds-backend/phoenix/middleware"
	"os"
)

// 只需要把需要修改的sql粘贴到下面变量即可
var sqls = []string{
	``,
}

func main() {
	infra := middleware.NewInfra("migrations")
	for _, sql := range sqls {
		r, err := infra.DB.Exec(sql)
		if err != nil {
			infra.Logger.WithError(err).Warn("sql error")
			os.Exit(1)
		} else {
            infra.Logger.WithError(err).Info("affected:", r.RowsAffected())
        }
	}
	os.Exit(0)
}
