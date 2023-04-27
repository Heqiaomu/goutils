package v2

import (
	"fmt"
	"github.com/Heqiaomu/goutil/uiform/conform"
	v0 "github.com/Heqiaomu/goutil/uiform/generate/v0"
	v1 "github.com/Heqiaomu/goutil/uiform/generate/v1"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	"github.com/google/uuid"
	"strconv"
)

// ReadUiformAction 将 uiformActionDir 目录下的fields.yaml、inputs.yaml、reactions.yaml读取成 UIData
func ReadUiformAction(uiformActionDir string, dependsDir string) (*conform.UIData, error) {
	// 读取fields.yaml文件
	fields, err := v0.ReadUiformActionFields(uiformActionDir)
	if err != nil {
		return nil, err
	}
	// 读取inputs.yaml文件
	inputs, err := v0.ReadUiformActionInputs(uiformActionDir)
	if err != nil {
		return nil, err
	}
	// 读取reactions.yaml文件
	reactions, fieldId2Reactions, err := v0.ReadUiformActionReactions(uiformActionDir)
	if err != nil {
		return nil, err
	}

	// v2 - 处理 输入域 和 输入项表单的ID
	for _, field := range fields.Fields {
		field.ID = strconv.Itoa(int(uuid.New().ID()))
	}
	for _, input := range inputs.Inputs {
		input.ID = strconv.Itoa(int(uuid.New().ID()))
	}
	// 缓存：输入项主键2输入项
	fieldKey2Field := make(map[string]*yaml.Field, 16)
	for _, field := range fields.Fields {
		if fieldKey2Field[field.Key] != nil {
			return nil, fmt.Errorf("field key existed, key=%s", field.Key)
		}
		fieldKey2Field[field.Key] = field
	}
	// 缓存：输入域主键2输入域
	inputKey2Input := make(map[string]*yaml.Input, 16)
	for _, input := range inputs.Inputs {
		if inputKey2Input[input.Key] != nil {
			return nil, fmt.Errorf("input key existed, key=%s", input.Key)
		}
		input.UIFormVersion = "v2"
		inputKey2Input[input.Key] = input
	}
	uiData := &conform.UIData{
		Fields:         fields.Fields,
		Inputs:         inputs.Inputs,
		FieldReactions: reactions.FieldReactions,
		Id2UIData: &conform.Id2UIData{
			FieldKey2Field:     fieldKey2Field,
			InputKey2Input:     inputKey2Input,
			FieldKey2Reactions: fieldId2Reactions,
		},
	}
	// 校验
	err = Check(uiData)
	if err != nil {
		return nil, err
	}
	// 处理depend
	for _, input := range uiData.Inputs {
		// 如果具有依赖，则交由依赖转化为一个新的input
		if input.Depend != "" {
			err = v1.AddDepend(input, uiData, dependsDir)
			if err != nil {
				return nil, err
			}
		}
	}
	return uiData, nil
}
