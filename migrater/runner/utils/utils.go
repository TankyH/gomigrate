package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"migrations/migrate/runner/indexstruct"
)

func ReadIndexConfig(filePath string) (indexstruct.IndexConfig, error) {
	/*
		filePath is /..../index.json
	*/
	log.Debug("version path:", filePath)
	result := make(indexstruct.IndexConfig, 0)
	indexFile, err := os.Open(filePath)
	if err != nil {
		log.Errorf("read file path error: %v, %v", filePath, err)
		return result, err
	}
	data, err := ioutil.ReadAll(indexFile)
	if err != nil {
		log.Errorf("read file error: %v, %v", filePath, err)
		return result, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Error("json unmarshal error:", data, err)
		return result, err
	}
	return result, nil
}
