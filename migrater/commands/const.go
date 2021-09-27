package commands

const (
	FlagVersion  = "version"
	FlagModule   = "module"
	FlagMake     = "make"
	FlagFake     = "fake"
	FlagForce    = "force"
	FlagEnv      = "env"
	FlagCheck    = "check"
	FlagLogLevel = "log"
)

// Run env. migrate 运行时所需环境变量
var (
	VipPathEnv = "vip_path"     // phoneix 的运行配置路径的环境变量名（因为连接数据库需要用到）
	MigrateEnv = "migrate_path" // migration 目录
)
