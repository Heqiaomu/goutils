package generate

import (
	"fmt"
	"github.com/Heqiaomu/goutil/uiform/conform"
	v0 "github.com/Heqiaomu/goutil/uiform/generate/v0"
	"github.com/Heqiaomu/goutil/uiform/generate/v1"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	yamlv2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

func ReadInfo(driverDir string) (*yaml.Info, error) {
	infoFilePath := filepath.Join(driverDir, "info.yaml")
	return ReadInfoFromInfoYaml(infoFilePath)
}

func ReadInfoFromInfoYaml(infoFilePath string) (*yaml.Info, error) {
	isExist, err := v0.PathExists(infoFilePath)
	if err != nil {
		return nil, fmt.Errorf("finding info.yaml=%s, err=%s", infoFilePath, err.Error())
	}
	if !isExist {
		return nil, fmt.Errorf("cannot find info.yaml=%s ", infoFilePath)
	}
	infoBytes, err := ioutil.ReadFile(infoFilePath)
	if err != nil {
		return nil, fmt.Errorf("read info.yaml=%s, err=%s", infoFilePath, err.Error())
	}
	var info yaml.Info
	err = yamlv2.Unmarshal(infoBytes, &info)
	if err != nil {
		return nil, fmt.Errorf("unmarshal info.yaml=%s, err=%s", infoFilePath, err.Error())
	}
	err = checkInfoYaml(infoFilePath, &info)
	if err != nil {
		return nil, err
	}
	err = ModifyInfo(infoFilePath, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// ReadUiformAction 将 uiformActionDir 目录下的fields.yaml、inputs.yaml、reactions.yaml读取成 UIData
func ReadUiformAction(uiformVersion string, uiformActionDir string, dependsDir string) (*conform.UIData, error) {
	switch uiformVersion {
	case "v1":
		return v1.ReadUiformAction(uiformActionDir, dependsDir)
	case "v2":
	}
	return nil, fmt.Errorf("unsupport uiform verison")
}
