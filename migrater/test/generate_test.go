package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"migrations/migrate/runner/errors"
	"migrations/migrate/runner/generator"
	"migrations/migrate/runner/indexstruct"
)

type TestCase struct {
	Name        string
	ExpectError error
	SetUp       func(string)
	CleanUp     func(string)
	Version     string
}

func TestCheckFunction(t *testing.T) {
	tests := []TestCase{
		{Name: "version folder doesn't exist",
			ExpectError: nil,
			Version:     "v9.2.0",
			CleanUp: func(version string) {
				_ = os.Remove(version)
			},
		},
		{Name: "version folder  exist",
			ExpectError: errors.VersionExist,
			Version:     "v9.1.0",
		},
	}
	for _, testcase := range tests {
		generateor := generator.Init([]string{"oooo"}, testcase.Version, "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration", false, false)
		err := generateor.Check()
		if testcase.CleanUp != nil {
			testcase.CleanUp(testcase.Version)
		}
		if err != testcase.ExpectError {
			t.Fatalf("test %s failed, expected %v, real %v", testcase.Name, testcase.ExpectError, err)
		}
	}
}

func TestLoadFunction(t *testing.T) {
	tests := []TestCase{
		{
			Name:        "index.json not exist, folder is empty",
			ExpectError: nil,
			Version:     "v19.9.9",
			SetUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				_ = os.MkdirAll(path, os.ModePerm)
			},
			CleanUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				os.RemoveAll(path)
			},
		},
		{
			Name:        "index.json not exist, but folder is not empty",
			ExpectError: errors.IndexNotExistError,
			Version:     "v19.9.9",
			SetUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				_ = os.MkdirAll(path, os.ModePerm)
				_ = os.Mkdir(path+"\\kkkk", os.ModePerm)
			},
			CleanUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				_ = os.RemoveAll(path)
			},
		},
		{
			Name:        "index exist, but folder not exist",
			ExpectError: errors.ModuleNotExistError,
			Version:     "19.10.1",
			SetUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				_ = os.MkdirAll(path, os.ModePerm)
				content := map[string]int{
					"kkk2222": 1,
				}
				jsonContent, _ := json.Marshal(&content)
				err := ioutil.WriteFile(path+"\\index.json", jsonContent, 0755)
				if err != nil {
					fmt.Printf("Unable to write file: %v", err)
				}
			},
			CleanUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				os.RemoveAll(path)
			},
		},
		{
			Name:        "index exist, but folder not in index",
			ExpectError: errors.ManualCreateError,
			Version:     "v19.10.1",
			SetUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				_ = os.MkdirAll(path, os.ModePerm)
				_ = os.Mkdir(path+"\\kkk", os.ModePerm)
				content := map[string]int{
					"kkk2222": 1,
				}
				jsonContent, _ := json.Marshal(&content)
				err := ioutil.WriteFile(path+"\\index.json", jsonContent, 0755)
				if err != nil {
					fmt.Printf("Unable to write file: %v", err)
				}
			},
			CleanUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				os.RemoveAll(path)
			},
		},
		{
			Name:        "index exist, and match",
			ExpectError: nil,
			Version:     "v19.10.2",
			SetUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				_ = os.MkdirAll(path, os.ModePerm)
				_ = os.Mkdir(path+"\\kkk", os.ModePerm)
				content := map[string]int{
					"kkk": 1,
				}
				jsonContent, _ := json.Marshal(&content)
				err := ioutil.WriteFile(path+"\\index.json", jsonContent, 0755)
				if err != nil {
					fmt.Printf("Unable to write file: %v", err)
				}
			},
			CleanUp: func(version string) {
				path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
				os.RemoveAll(path)
			},
		},
	}
	for _, testcase := range tests {
		gen := generator.Init([]string{testcase.Name}, testcase.Version, "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration")
		gen.VersionPath = gen.WorkingPath + "\\version\\" + testcase.Version
		testcase.SetUp(testcase.Version)
		err := gen.Load()
		testcase.CleanUp(testcase.Version)
		if err != testcase.ExpectError {
			t.Fatalf("")
		}
	}
}

type GenerateTestCase struct {
	TestCase
	Param []string
}

func TestGenerateFunction(t *testing.T) {
	tests := []GenerateTestCase{
		{
			TestCase: TestCase{Name: "Create new module, for newVersion",
				ExpectError: nil,
				Version:     "v19.9.9",
				SetUp: func(version string) {
					path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
					_ = os.MkdirAll(path, os.ModePerm)
				},
				CleanUp: func(version string) {
					path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
					os.RemoveAll(path)
				}},
			Param: []string{"module1"},
		},
		{
			TestCase: TestCase{Name: "Create new module, for ExistVersion",
				ExpectError: nil,
				Version:     "v19.9.9",
				SetUp: func(version string) {
					path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
					_ = os.MkdirAll(path, os.ModePerm)
				},
				CleanUp: func(version string) {
					path := "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration\\version\\" + version
					os.RemoveAll(path)
				}},
			Param: []string{},
		},
	}
	for _, testcase := range tests {
		gen := generator.Init([]string{"kkkk"}, testcase.Version, "D:\\backend\\src\\gitlab.xinghuolive.com\\birds-backend\\migrations\\migration")
		gen.VersionPath = gen.WorkingPath + "\\version\\" + testcase.Version
		gen.GenModules = testcase.Param
		gen.IndexJson = make(indexstruct.IndexConfig, 0)
		testcase.SetUp(testcase.Version)
		err := gen.Generate()
		testcase.CleanUp(testcase.Version)
		if err != testcase.ExpectError {
			t.Fatalf("")
		}
	}
}
