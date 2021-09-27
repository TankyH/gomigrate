package migrator

import (
	"bytes"
	"fmt"
	"io"
	"migrations/migrate/runner/hook"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.xinghuolive.com/birds-backend/phoenix/middleware"
	"migrations/migrate/runner/errors"
	"migrations/migrate/runner/indexstruct"
	"migrations/migrate/runner/migrator/db"
	"migrations/migrate/runner/utils"
)

const (
	EnvVar          = "ENV"
	CheckOwlStudent = "owl"
)

const (
	LogLevelDebug = 1
)

var (
	Local = "local"
	Pre   = "pre"
	Prod  = "prod"
	Dev   = "dev"
	Test  = "test"
	// After Hook Path, 以后可以做成配置，目前只需要
	afterHookPrefix  = "migrate/afterhooks"
	afterHookFolders = []string{"tablecheck"}
)

type InputModule struct {
	ModuleName string // 模块名称
	FullName   string // 生成的名称 : 序号_时间戳_模块名(1_1577890090_student)
}

type Migrator struct {
	/*
	 如果强制运行，模块在此RunModules map对应的值为true, 否则为false,
	 如果Module -- all 加 --force， 则全部模块都为true
	*/
	Version string
	// Flags
	AllFlag       bool   // 全部执行标志
	Force         bool   // 强制运行标志
	FakeFlag      bool   // 标记为已运行标志
	Env           string // 工具执行中的环境变量
	OriginEnv     string // 原始运行的环境变量
	ResetEnv      bool   // 是否需要回复运行的环境变量
	AfterHookFlag bool   // 只运行after hook 检查标志

	// Paths
	VersionPath string // 对应的版本路径
	WorkingPath string // 工作路径
	// config
	Infra        *middleware.Infra
	RunRecrod    []db.Schema   // 运行每个脚本后的记录
	InputModules []InputModule // 需要运行的模块
	RunModules   map[string]bool
	RunIndex     indexstruct.IndexConfig // 读取index.json 存放到对象内存中
	LogLevel     int
	Check        string
	// AfterHookPath, 是以migration目录为根目录的相对路径，用于运行监测migration执行完成后时候有同步的工作，暂时写死在init的函数中
	AfterHookPath []string
}

func (m *Migrator) Debug() bool {
	return m.LogLevel == 1
}

func (m *Migrator) setAfterHookOnly() {
	m.AfterHookFlag = true
}

func Init(all bool, modules []string, version string, force bool, workDir string, fake bool, env, check string) (*Migrator, error) {
	m := new(Migrator)
	m.Version = version
	m.Infra = middleware.NewInfra("migrations")
	m.AllFlag = all
	m.WorkingPath = workDir
	m.Force = force
	m.FakeFlag = fake
	m.VersionPath = filepath.Join(m.WorkingPath, "version", m.Version)
	m.RunModules = make(map[string]bool)
	m.Env = env
	m.Check = check
	if env != Local || env != "" {
		m.OriginEnv = os.Getenv(EnvVar)
		m.ResetEnv = true
		_ = os.Setenv(EnvVar, m.Env)
	}
	m.initAfterHookPath(afterHookFolders)

	// 读取运行配置文件
	var configErr error
	m.RunIndex, configErr = utils.ReadIndexConfig(filepath.Join(m.VersionPath, "index.json"))
	if configErr != nil {
		log.Fatal("get run index error, please check version index.json:", configErr)
		return m, errors.IndexNotExistError
	}
	// 初始化需要运行的模块
	if m.AllFlag {
		// 全量执行， 则读取index.json 获取执行顺序
		m.InputModules = make([]InputModule, len(m.RunIndex))
		for i, item := range m.RunIndex {
			m.InputModules[i] = InputModule{ModuleName: item.ModuleName, FullName: item.FullName}
			m.RunModules[item.ModuleName] = m.Force
		}
	} else {
		// 如果不是全量运行，需要判断是否在index中，不存在，则直接退出
		runIndexMap := make(map[string]int)
		fullNameIndex := make(map[string]string) // module name as key , full path as value
		for _, item := range m.RunIndex {
			runIndexMap[item.ModuleName] = item.Seq
			fullNameIndex[item.ModuleName] = item.FullName
		}
		for _, item := range modules {
			if _, ok := runIndexMap[item]; !ok {
				return m, errors.ModuleNotExistError
			}
			m.InputModules = append(m.InputModules, InputModule{ModuleName: item, FullName: fullNameIndex[item]})
			m.RunModules[item] = m.Force
		}
	}

	log.Debugf("m property: %#v\n", m)
	return m, nil
}

func (m *Migrator) initAfterHookPath(folders []string) {
	for _, folderName := range folders {
		fileFullPath := filepath.Join(m.WorkingPath, afterHookPrefix, folderName)
		m.AfterHookPath = append(m.AfterHookPath, fileFullPath)
	}
}

func (m *Migrator) resetEnv() error {
	err := os.Setenv(EnvVar, m.OriginEnv)
	if err != nil {
		log.Println("reset env path environment failed")
		return err
	}
	return nil
}

func (m *Migrator) Validate() error {
	/*
		if module has run success and not force run again, will be return error
	*/
	migrationMapper := db.New(m.Infra.DB)
	finishedMigrations := make([]db.Schema, 0)
	err := m.Infra.DB.Model(&db.Schema{}).
		Where(migrationMapper.Tag(&migrationMapper.Version).Eq(), m.Version).
		Where(migrationMapper.Tag(&migrationMapper.Success).Is(), true).
		Order(migrationMapper.Tag(&migrationMapper.CreatedAt).Asc()).
		Select(&finishedMigrations)
	if err != nil {
		return errors.GetMigrationError
	}
	finishedMigrationsMap := make(map[string]int64)
	for _, item := range finishedMigrations {
		finishedMigrationsMap[item.ModuleName] = item.CreatedAt.Unix()
	}
	skipModules := make([]string, 0)
	for moduleName, isForce := range m.RunModules {
		if _, ok := finishedMigrationsMap[moduleName]; ok {
			if isForce {
				continue
			}
			skipModules = append(skipModules, moduleName)
		}
	}
	for _, module := range skipModules {
		delete(m.RunModules, module)
	}
	_ = os.Mkdir(filepath.Join(m.VersionPath, "log"), os.ModePerm)
	return nil
}

type CaptureWriter struct {
	buf bytes.Buffer
	w   io.Writer
}

func NewCaptureWriter(w io.Writer) *CaptureWriter {
	return &CaptureWriter{w: w}
}

func (w *CaptureWriter) Write(d []byte) (int, error) {
	w.buf.Write(d)
	return w.w.Write(d)
}

func (w *CaptureWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func (m *Migrator) run(fileName string, moduleName string, startTime string) error {
	path := filepath.Join(m.VersionPath, fileName, "main.go")
	outLogPath, _ := os.Create(filepath.Join(m.VersionPath, "log", fmt.Sprintf("out-%s-%s.log", moduleName, startTime)))
	errLogPath, _ := os.Create(filepath.Join(m.VersionPath, "log", fmt.Sprintf("err-%s-%s.log", moduleName, startTime)))
	cmd := exec.Command("go", "run", path)
	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("cmd StdoutPipe fail with error:", err)
		return err
	}
	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		log.Error("cmd StderrPipe fail with error:", err)
		return err
	}
	stdout := NewCaptureWriter(os.Stdout)
	stderr := NewCaptureWriter(os.Stderr)
	if err := cmd.Start(); err != nil {
		log.Error("cmd Start fail with error:", err)
		return err
	}
	log.Infof("================= Start [version:%v, module:%v] ================", m.Version, moduleName)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			if p := recover(); p != nil {
				stackInfo := debug.Stack()
				log.Error(string(stackInfo), p)
			}
		}()
		_, errStdout := io.Copy(stdout, stdoutIn)
		if errStdout != nil {
			log.Error("failed to capture stdout")
		}
		wg.Done()
	}()

	_, errStderr := io.Copy(stderr, stderrIn)
	if errStderr != nil {
		log.Error("failed to capture stderr")
	}
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Error("cmd.Run fail with error:", err)
	}
	_, _ = stdout.buf.WriteTo(outLogPath)
	_, _ = stderr.buf.WriteTo(errLogPath)

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		log.Errorf("================= End [exit code: %v]================", exitCode)
	} else {
		log.Info("================= End ================")
	}

	return err
}

func (m *Migrator) Migrate() error {
	for _, module := range m.InputModules {
		if _, ok := m.RunModules[module.ModuleName]; !ok {
			log.Printf("===module %s is already run, and this is not force run, this run will skip====", module.ModuleName)
			continue
		}
		startAt := time.Now()
		formattedStartAt := strconv.FormatInt(startAt.Unix(), 10)
		dbRecord := db.Schema{}
		dbRecord.ExecTime = startAt
		dbRecord.ModuleName = module.ModuleName
		dbRecord.FullName = module.FullName
		dbRecord.Version = m.Version
		dbRecord.Force = m.Force
		err := m.run(module.FullName, module.ModuleName, formattedStartAt)
		if err != nil {
			dbRecord.Success = false
		} else {
			dbRecord.Success = true
		}
		_, insertErr := m.Infra.DB.Model(&dbRecord).Insert()
		if insertErr != nil {
			log.Fatalf("insert record to db failed")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) Fake() {
	for _, module := range m.InputModules {
		startAt := time.Now()
		dbRecord := db.Schema{}
		dbRecord.ExecTime = startAt
		dbRecord.ModuleName = module.ModuleName
		dbRecord.FullName = module.FullName
		dbRecord.Version = m.Version
		dbRecord.Force = m.Force
		dbRecord.Success = true
		dbRecord.IsFake = true
		m.RunRecrod = append(m.RunRecrod, dbRecord)
	}
	_, err := m.Infra.DB.Model(&m.RunRecrod).Insert()
	if err != nil {
		log.Fatalf("insert record to db failed")
	}
	log.Print("insert run record success, fake success all ")
}

/*
这个函数是用于运行通用的每次版本升级后都需要运行的脚本。目前使用变量存在了Migration的源代码中，如果以后变多以后
需要单独读写配置来进行检测
*/
func (m *Migrator) afterHookRun() hook.HookErr {
	for _, item := range m.AfterHookPath {
		workingPath := make([]string, 0)
		_, suffix := filepath.Split(item)
		log.Debugf("hook item: %v", item)
		_ = filepath.Walk(item, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && suffix != info.Name() {
				workingPath = append(workingPath, path)
			}
			return nil
		})
		log.Debugf("hook working path: %v", workingPath)
		for _, path := range workingPath {
			err := m.runhook(path)
			if err.HasError() {
				m.Infra.Logger.Logger.Error("run after hook failed! err is", err)
				return err
			}
		}
	}
	return hook.HookErr{}
}

func (m *Migrator) runhook(scriptPath string) hook.HookErr {
	modulePath, folderName := filepath.Split(scriptPath)
	outLogPath, err := os.Create(filepath.Join(m.VersionPath, "log", fmt.Sprintf("out-after-hook-%s.log", folderName)))
	if err != nil {
		log.Error("create log error:", err)
	}
	errLogPath, _ := os.Create(filepath.Join(m.VersionPath, "log", fmt.Sprintf("err-afterhook-%s.log", folderName)))

	path := filepath.Join(scriptPath, "main.go")
	cmd := exec.Command("go", "run", path)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	_ = cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	log.Infof("=============== hook end [exit code is %v] ==================", exitCode)
	_, _ = out.WriteTo(outLogPath)
	_, _ = stderr.WriteTo(errLogPath)
	if exitCode != 0 {
		return hook.HookErr{Code: exitCode, ModuleName: modulePath, FolderName: folderName}
	}
	return hook.HookErr{}
}

func (m *Migrator) Run() {
	err := m.Validate()
	if err != nil {
		log.Fatalf("check Meet error, error is %v", err)
		return
	}
	if !m.FakeFlag {
		execErr := m.Migrate()
		if m.ResetEnv {
			err = m.resetEnv()
		}
		if execErr != nil {
			log.Printf(execErr.Error())
			return
		}
	} else {
		m.Fake()
	}
	if m.Check == CheckOwlStudent {
		afterHookErr := m.afterHookRun()
		if afterHookErr.HasError() {
			m.Infra.Logger.Errorf("after check error, module is %s, folder is %s, err is %v",
				afterHookErr.ModuleName,
				afterHookErr.FolderName,
				afterHookErr,
			)
			os.Exit(1)
		}
	}
}
