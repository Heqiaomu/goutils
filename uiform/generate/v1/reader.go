package v1

import (
	"github.com/Heqiaomu/goutil/uiform/conform"
	v0 "github.com/Heqiaomu/goutil/uiform/generate/v0"
	"github.com/Heqiaomu/goutil/uiform/yaml"
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
	// 构造uiData
	// 缓存：输入项主键2输入项
	fieldId2Field := make(map[string]*yaml.Field, 16)
	for _, field := range fields.Fields {
		fieldId2Field[field.ID] = field
	}
	// 缓存：输入域主键2输入域
	inputId2Input := make(map[string]*yaml.Input, 16)
	for _, input := range inputs.Inputs {
		inputId2Input[input.ID] = input
	}
	uiData := &conform.UIData{
		Fields:         fields.Fields,
		Inputs:         inputs.Inputs,
		FieldReactions: reactions.FieldReactions,
		Id2UIData: &conform.Id2UIData{
			FieldKey2Field:     fieldId2Field,
			InputKey2Input:     inputId2Input,
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
			err = AddDepend(input, uiData, dependsDir)
			if err != nil {
				return nil, err
			}
		}
	}
	return uiData, nil
}
