package v2

import (
	"fmt"
	"github.com/Heqiaomu/goutil/uiform/conform"
	v0 "github.com/Heqiaomu/goutil/uiform/generate/v0"
)

// Check 检查yaml结构是否填写正确
func Check(data *conform.UIData) error {
	id2Data := data.Id2UIData
	inputs := data.Inputs
	for index, input := range inputs {
		// 输入域的标题可以为空，但是如果不为空则标题必须有名字
		if input.Title != nil && input.Title.Text == "" {
			return fmt.Errorf("inputs[%d].title.text is empty", index)
		}
		// 输入域的显示类型是枚举的
		if input.ShowType != "" {
			if !v0.InputShowType[input.ShowType] {
				return fmt.Errorf("inputs[%d].type is illegal. actual type is %s while %s is excepted ", index, input.ShowType, v0.KeywordMapToString(v0.InputShowType))
			}
		}
		// 输入域的依赖是枚举的
		if input.Depend != "" {
			if !v0.InputDependType[input.Depend] {
				return fmt.Errorf("inputs[%d].depend is illegal. actual depend is %s while %s is excepted ", index, input.Depend, v0.KeywordMapToString(v0.InputDependType))
			}
		}
		if input.SubInputs != nil {
			// 输入域的子输入域，其Key值必须是已经存在的输入域
			for i, subInputKey := range input.SubInputs {
				if id2Data.InputKey2Input[subInputKey] == nil {
					return fmt.Errorf("inputs[%d](key=%s).sub[%d](key=%s) is not defined in fields.yaml", index, input.Key, i, subInputKey)
				}
			}
		} else if input.Fields != nil {
			// 输入域的输入项，其Key值必须是已经存在的输入项
			for i, fieldKey := range input.Fields {
				if id2Data.FieldKey2Field[fieldKey] == nil {
					return fmt.Errorf("inputs[%d](key=%s).fields[%d](key=%s) is not defined in fields.yaml", index, input.Key, i, fieldKey)
				}
			}
		}
	}
	fields := data.Fields
	for index, field := range fields {
		// 输入项的标题可以为空，如果不为空标题名是必填
		if field.Title != nil {
			if field.Title.Text == "" {
				return fmt.Errorf("fields[%d].title.text is empty", index)
			}
		}
		// 输入项的类型是枚举的
		if field.InputType != "" {
			if !v0.InputFieldInputType[field.InputType] {
				return fmt.Errorf("fields[%d].type is illegal. actual type is %s while %s is excepted ", index, field.InputType, v0.KeywordMapToString(v0.InputFieldInputType))
			}
		}
		// 输入项的按钮，比如当前是选择已有主机的选择框，这个选择框旁边有一个按钮，用来创建主机
		if field.Buttons != nil {
			for index, fieldId := range field.Buttons {
				if id2Data.FieldKey2Field[fieldId] == nil {
					return fmt.Errorf("field.buttons[%d](key=%s) is not defined in fields.yaml", index, fieldId)
				}
			}
		}
	}
	fieldReactions := data.FieldReactions
	for i, fieldReaction := range fieldReactions {
		for j, reaction := range fieldReaction.Reactions {
			if reaction.ReactType != "" {
				// 联动类型是枚举的
				if !v0.InputReactReactType[reaction.ReactType] {
					return fmt.Errorf("fieldReactions[%d].reactions[%d].type is illegal. actual type is %s while %s is excepted ", i, j, reaction.ReactType, v0.KeywordMapToString(v0.InputReactReactType))
				}
				// 联动的 触发显示的输入域 url触发 目标输入域 不能都为空
				if reaction.InputId == "" && reaction.UrlReact == nil && reaction.TargetInputId == "" {
					return fmt.Errorf("fieldReactions[%d].reactions[%d].input and target is empty and it's urlReact is also nil", i, j)
				}
				if reaction.UrlReact != nil {
					if reaction.UrlReact.Body != "" {
						if id2Data.InputKey2Input[reaction.UrlReact.Body] == nil {
							return fmt.Errorf("fieldReactions[%d].reactions[%d].urlReact.body is not defined in inputs.yaml", i, j)
						}
					}
				}
			}
		}
	}
	return nil
}
