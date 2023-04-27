package v0

import (
	"bytes"
	"fmt"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	yamlv2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadUiformActionFields(uiformActionDir string) (*yaml.Fields, error) {
	// fields 文件提取
	var fields yaml.Fields
	fieldsFilePath := filepath.Join(uiformActionDir, "fields.yaml")
	isExist, err := PathExists(fieldsFilePath)
	if err != nil {
		return nil, fmt.Errorf("finding fields.yaml=%s err, err=%s", fieldsFilePath, err.Error())
	}
	if !isExist {
		return nil, fmt.Errorf("cannot find fields.yaml=%s ", fieldsFilePath)
	}
	fieldsBytes, err := ioutil.ReadFile(fieldsFilePath)
	if err != nil {
		return nil, fmt.Errorf("read fields.yaml=%s err, err=%s", fieldsFilePath, err.Error())
	}

	err = yamlv2.Unmarshal(fieldsBytes, &fields)
	if err != nil {
		return nil, fmt.Errorf("unmarshal fields.yaml=%s err, err=%s", fieldsFilePath, err.Error())
	}
	return &fields, nil
}

func ReadUiformActionInputs(uiformActionDir string) (*yaml.Inputs, error) {
	// inputs 文件提取
	var inputs yaml.Inputs
	inputsFilePath := filepath.Join(uiformActionDir, "inputs.yaml")
	isExist, err := PathExists(inputsFilePath)
	if err != nil {
		return nil, fmt.Errorf("finding inputs.yaml=%s err, err=%s", inputsFilePath, err.Error())
	}
	if !isExist {
		return nil, fmt.Errorf("cannot find inputs.yaml=%s ", inputsFilePath)
	}
	inputsBytes, err := ioutil.ReadFile(inputsFilePath)
	if err != nil {
		return nil, fmt.Errorf("read inputs.yaml=%s err, err=%s", inputsFilePath, err.Error())
	}
	err = yamlv2.Unmarshal(inputsBytes, &inputs)
	if err != nil {
		return nil, fmt.Errorf("unmarshal inputs.yaml=%s err, err=%s", inputsFilePath, err.Error())
	}
	return &inputs, nil
}

func ReadUiformActionReactions(uiformActionDir string) (*yaml.Reactions, map[string][]*yaml.Reaction, error) {
	// reactions 提取
	reactions := yaml.Reactions{FieldReactions: []*yaml.FieldReactions{}}
	// 缓存：输入项主键2联动
	fieldId2Reactions := make(map[string][]*yaml.Reaction, 16)
	reactionsFilePath := filepath.Join(uiformActionDir, "reactions.yaml")
	isExist, err := PathExists(reactionsFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("finding inputs.yaml=%s err, err=%s", reactionsFilePath, err.Error())
	}
	if isExist {
		reactionsBytes, err := ioutil.ReadFile(reactionsFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("read reactions.yaml=%s err, err=%s", reactionsFilePath, err.Error())
		}
		err = yamlv2.Unmarshal(reactionsBytes, &reactions)
		if err != nil {
			return nil, nil, fmt.Errorf("unmarshal reactions.yaml=%s err, err=%s", reactionsFilePath, err.Error())
		}
		for _, filedReaction := range reactions.FieldReactions {
			fieldId2Reactions[filedReaction.Field] = filedReaction.Reactions
		}
	}
	return &reactions, fieldId2Reactions, nil
}

// PathExists 判断文件或文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// KeywordMapToString 将type映射转化为数组列表
func KeywordMapToString(m map[string]bool) string {
	var buf bytes.Buffer
	for key := range m {
		buf.WriteString(key)
		buf.WriteString(", ")
	}
	return buf.String()
}
