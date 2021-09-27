package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"migrations/migrate/commands"
	"os"
)

/*
env环境配置, redis, pg配置,
和晓服务, 教学, 题库等secret, url配置等等,
都要按部署环境做好配置
*/
func main() {
	app := &cli.App{}
	app.Flags = commands.NewFlags()
	app.Action = commands.NewAction()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
