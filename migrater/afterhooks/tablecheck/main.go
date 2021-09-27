package main

import (
	"os"

	"gitlab.xinghuolive.com/birds-backend/phoenix/middleware"
)

const (
	sql = `with cte1 as (
    SELECT a.attname AS field
    FROM pg_class c,
         pg_attribute a
             LEFT OUTER JOIN pg_description b ON a.attrelid = b.objoid AND a.attnum = b.objsubid,
         pg_type t,
         pg_namespace n
    WHERE c.relname = 'user_student'
      and a.attnum > 0
      and n.nspname = 'oto'
      and a.attrelid = c.oid
      and a.atttypid = t.oid
      and c.relnamespace = n.oid
    ORDER BY a.attnum),
     cte2 as (
         SELECT a.attname AS field
         FROM pg_class c,
              pg_attribute a
                  LEFT OUTER JOIN pg_description b ON a.attrelid = b.objoid AND a.attnum = b.objsubid,
              pg_type t,
              pg_namespace n
         WHERE c.relname = 'temp_student'
           and a.attnum > 0
           and n.nspname = 'oto'
           and a.attrelid = c.oid
           and a.atttypid = t.oid
           and c.relnamespace = n.oid
         ORDER BY a.attnum
     )
select cte3.field
from cte1 as cte3
     left join cte2 as cte4 on cte3.field = cte4.field
where cte4.field is null;`
)

func main() {
	infra := middleware.NewInfra("migrations")
	var fields []string
	_, err := infra.DB.Query(&fields, sql)
	if err != nil {
		infra.Logger.WithError(err).Warn("sql error")
		os.Exit(1)
	}

	if len(fields) > 0 {
		infra.Logger.Error("temp_student loss fileds: ", fields)
		os.Exit(1)
	}

	os.Exit(0)
}
