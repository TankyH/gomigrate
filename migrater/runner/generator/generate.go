package generator

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	migrationError "migrations/migrate/runner/errors"
	"migrations/migrate/runner/indexstruct"
	"migrations/migrate/runner/utils"
)

/*
this module is to generate the file of migration script
1. load old migration index nubmer
2. create new migration files, by template
3. (optional) check whether the old migrationg has conlict with
*/

type Generator interface {
	Run()
}

type MigrationGenerator struct {
	DestDir      string // should be migrations/${version}/xxxx/main.go
	Template     string // static file
	Version      string
	Modules      []string
	ExistModules []string
	GenModules   []string
	makeFlag     string
	IndexJson    indexstruct.IndexConfig
	RunIndex     map[string]int
	Index        int //
	WorkingPath  string
	VersionPath  string
	ConfigPath   string
}

const (
	IndexName     = "index.json"
	LogFolderName = "log"
)

const (
	TemplatePathPrefix = "migrate/template"
)

func Init(modules []string, version string, migrationPath string, makeFlag string) *MigrationGenerator {
	g := new(MigrationGenerator)
	g.Version = version
	// run test case need to handle inject path, for production not need to inject
	g.WorkingPath = migrationPath
	g.GenModules = modules
	g.makeFlag = makeFlag
	g.RunIndex = make(map[string]int)
	g.initPath()
	return g
}

func (g *MigrationGenerator) initPath() {
	switch {
	case g.makeFlag == "table":
		g.Template = filepath.Join(g.WorkingPath, TemplatePathPrefix, "main-table.tpl")
	case g.makeFlag == "sql":
		g.Template = filepath.Join(g.WorkingPath, TemplatePathPrefix, "main-sql.tpl")
	default:
		g.Template = filepath.Join(g.WorkingPath, TemplatePathPrefix, "main.tpl")
	}

	g.VersionPath = filepath.Join(g.WorkingPath, "version", g.Version)
	g.ConfigPath = filepath.Join(g.VersionPath, IndexName)
}

func (g *MigrationGenerator) Check() error {
	if _, err := os.Stat(g.VersionPath); err != nil && !os.IsNotExist(err) {
		log.Infof("version path : %v, %v", g.VersionPath, err)
		return migrationError.VersionExist
	}
	err := os.MkdirAll(g.VersionPath, os.ModePerm)
	if err != nil {
		log.Errorf("mkdir error:", err)
		return err
	}
	return nil
}

func (g *MigrationGenerator) Load() error {
	_, err := os.Stat(filepath.Join(g.VersionPath, IndexName))
	var IndexNotExist bool
	if os.IsNotExist(err) {
		IndexNotExist = true
	} else {
		var configErr error
		g.IndexJson, configErr = utils.ReadIndexConfig(filepath.Join(g.VersionPath, IndexName))
		if configErr != nil {
			log.Errorf("config error: %v", configErr)
			return configErr
		}
		for _, item := range g.IndexJson {
			g.RunIndex[item.ModuleName] = item.Seq
		}
	}
	// 获取对应版本的文件目录
	files, err := filepath.Glob(g.VersionPath + "/*")
	if err != nil {
		log.Fatalf("get folder from run path:", err)
	}
	if len(files) != 0 && IndexNotExist {
		log.Errorf("index not found: %v, %v", files, IndexNotExist)
		return migrationError.IndexNotExistError
	}
	mapDict := make(map[string]bool)
	// 对比目录中实际模块和index.json中的模块
	for _, item := range files {
		fileName := strings.Split(item, g.VersionPath)
		filename := strings.Trim(fileName[len(fileName)-1], "\\")
		filename = strings.Trim(filename, "/")
		if filename == IndexName || filename == LogFolderName {
			continue
		}
		ListName := strings.Split(filename, "_")
		moduleName := ListName[len(ListName)-1]
		if _, ok := g.RunIndex[moduleName]; !ok {
			log.Infof("module %s is self-created, warning:", moduleName)
			continue
		}
		g.ExistModules = append(g.ExistModules, moduleName)
		mapDict[moduleName] = true
	}
	//  对比目录中实际模块和index.json中的模块
	for key, _ := range g.RunIndex {
		if _, ok := mapDict[key]; !ok {
			log.Info("module %s has in index.json, but exist in real module:", key)
			return migrationError.ModuleNotExistError
		}
		g.Index += 1
	}
	for _, name := range g.GenModules {
		if _, ok := mapDict[name]; ok {
			return migrationError.ModuleAlreayCreateError
		}
	}
	return nil
}

func (g *MigrationGenerator) Generate() error {
	tplFile, err := os.Open(g.Template)
	data, err := ioutil.ReadAll(tplFile)
	if err != nil {
		log.Fatalf("read template data error: %v, %v", g.Template, err)
	}
	_ = tplFile.Close()
	for _, module := range g.GenModules {
		g.Index += 1
		ts := time.Now().Unix()
		stringIndex := strconv.Itoa(g.Index)
		stringTs := strconv.FormatInt(ts, 10)
		folderName := stringIndex + "_" + stringTs + "_" + module
		destFolder := filepath.Join(g.VersionPath, folderName)
		err := os.Mkdir(destFolder, 0755)
		if err != nil {
			log.Fatalf("create folder error, module name is %s", module)
			return errors.New("")
		}
		desFile, err := os.Create(filepath.Join(destFolder, "main.go"))
		if err != nil {
			return err
		}
		_, err = desFile.Write(data)
		_ = desFile.Close()
		if g.IndexJson == nil {
			g.IndexJson = make(indexstruct.IndexConfig, 0)
		}
		indexItem := indexstruct.Index{ModuleName: module, FullName: folderName, Seq: g.Index}
		g.IndexJson = append(g.IndexJson, indexItem)
	}
	// and change index.json
	newIndex, err := json.MarshalIndent(g.IndexJson, "", "\t")
	err = ioutil.WriteFile(filepath.Join(g.VersionPath, IndexName), newIndex, os.ModePerm)
	if err != nil {
		log.Info("update new index.json error:", err)
	}
	return nil
}

func (g *MigrationGenerator) Run() {
	err := g.Check()
	if err != migrationError.VersionExist && err != nil {
		return
	}
	err = g.Load()
	if err != nil {
		log.Fatal("generator loading error:", err)
		return
	}
	err = g.Generate()
	if err != nil {
		log.Fatal("generating error:", err)
		return
	}
}
