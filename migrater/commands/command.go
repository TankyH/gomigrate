package commands

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    FlagLogLevel,
			Aliases: []string{"l"},
			Value:   "info",
			Usage:   "set log level",
		},
		&cli.StringFlag{
			Name:     FlagVersion,
			Aliases:  []string{"v"},
			Value:    "v1.0.0",
			Usage:    "version of migration",
			Required: true,
		},
		&cli.StringFlag{
			Name:    FlagModule,
			Aliases: []string{"m"},
			Usage:   "module of migration",
		},
		&cli.StringFlag{
			Name:  FlagMake,
			Value: "",
			Usage: "make migrate, option: sql, table",
		},
		&cli.BoolFlag{
			Name:  FlagFake,
			Usage: "fake migrate: only add a migrate record, but not running migration",
		},
		&cli.BoolFlag{
			Name:  FlagForce,
			Usage: "force migrate: force migrate, even migration has run",
		},
		&cli.StringFlag{
			Name:  FlagEnv,
			Usage: "run env setting",
		},
		&cli.StringFlag{
			Name:  FlagCheck,
			Usage: "check something",
		},
	}
	return flags
}

func NewAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		mk := c.String(FlagMake)
		v := c.String(FlagVersion)
		mdl := c.String(FlagModule)
		env := c.String(FlagEnv)
		fake := c.Bool(FlagFake)
		force := c.Bool(FlagForce)
		logLevel := c.String(FlagLogLevel)
		check := c.String(FlagCheck)

		SetLogLevel(logLevel)
		log.Debugf(`version: %v, module: %v, fake: %v, force: %v, make: %v, log: %v`,
			v, mdl, fake, force, mk, logLevel)

		if len(mk) != 0 {
			log.Info("make migration:")
			p := MakeMigrateParam{
				version: v,
				module:  mdl,
				make:    mk,
			}
			makeMigration(&p)
		} else {
			log.Info("migrate:")
			p := MigrateParam{
				version: v,
				module:  mdl,
				env:     env,
				check:   check,
				fake:    fake,
				force:   force,
			}
			migrate(&p)
		}

		return nil
	}
}

func SetLogLevel(lvl string) {
	m := map[string]log.Level{
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
	}
	if level, ok := m[lvl]; ok {
		log.SetLevel(level)
	}
}

func workPath() string {
	workingPath := os.Getenv(MigrateEnv)
	if len(workingPath) != 0 {
		return workingPath
	}
	workingPath, _ = os.Getwd()
	return workingPath
}

type MakeMigrateParam struct {
	version string
	module  string
	make    string
}

type MigrateParam struct {
	version string
	module  string
	env     string
	check   string
	fake    bool
	force   bool
}

func migrate(p *MigrateParam) {
	log.Debug("vip_path:", os.Getenv(VipPathEnv))
	log.Debug("migrate_path:", os.Getenv(MigrateEnv))
	vipPath := os.Getenv(VipPathEnv)
	if len(vipPath) == 0 {
		log.Fatalf("please set environment `vip_path` to /path/to/project/vip_path")
	}
	migratePath := os.Getenv(MigrateEnv)
	if len(migratePath) == 0 {
		log.Fatalf("please set environment `migrate_path` to /path/to/migration/version")
	}
	var allFlag bool
	if len(p.module) == 0 {
		allFlag = true
	}
	var err error
	workingPath := workPath()
	Runner, err := migrator.Init(allFlag, []string{p.module}, p.version, p.force, workingPath, p.fake, p.env, p.check)
	if err != nil {
		log.Fatalf("init migrator failed:", err)
	}
	Runner.Run()
}

func makeMigration(p *MakeMigrateParam) {
	if len(p.module) == 0 {
		log.Fatalf("module can't be zero")
	}
	workingPath := workPath()
	Runner := generator.Init([]string{p.module}, p.version, workingPath, p.make)
	Runner.Run()
}
