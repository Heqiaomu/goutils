package v1

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
		if input.Title != nil && input.Title.Text == "" {
			return fmt.Errorf("inputs[%d].title.text is empty", index)
		}
		if input.ShowType != "" {
			if !v0.InputShowType[input.ShowType] {
				return fmt.Errorf("inputs[%d].type is illegal. actual type is %s while %s is excepted ", index, input.ShowType, v0.KeywordMapToString(v0.InputShowType))
			}
		}
		if input.Depend != "" {
			if !v0.InputDependType[input.Depend] {
				return fmt.Errorf("inputs[%d].depend is illegal. actual depend is %s while %s is excepted ", index, input.Depend, v0.KeywordMapToString(v0.InputDependType))
			}
		}
		if input.SubInputs != nil {
			for i, subInputId := range input.SubInputs {
				if id2Data.InputKey2Input[subInputId] == nil {
					return fmt.Errorf("inputs[%d](id=%s).sub[%d](id=%s) is not defined in fields.yaml", index, input.ID, i, subInputId)
				}
			}
		} else if input.Fields != nil {
			for i, fieldId := range input.Fields {
				if id2Data.FieldKey2Field[fieldId] == nil {
					return fmt.Errorf("inputs[%d](id=%s).fields[%d](id=%s) is not defined in fields.yaml", index, input.ID, i, fieldId)
				}
			}
		}
	}
	fields := data.Fields
	for index, field := range fields {
		if field.Title != nil {
			if field.Title.Text == "" {
				return fmt.Errorf("fields[%d].title.text is empty", index)
			}
		}
		if field.InputType != "" {
			if !v0.InputFieldInputType[field.InputType] {
				return fmt.Errorf("fields[%d].type is illegal. actual type is %s while %s is excepted ", index, field.InputType, v0.KeywordMapToString(v0.InputFieldInputType))
			}
		}
		if field.Buttons != nil {
			for index, fieldId := range field.Buttons {
				if id2Data.FieldKey2Field[fieldId] == nil {
					return fmt.Errorf("field.buttons[%d](id=%s) is not defined in fields.yaml", index, fieldId)
				}
			}
		}
	}
	fieldReactions := data.FieldReactions
	for i, fieldReaction := range fieldReactions {
		for j, reaction := range fieldReaction.Reactions {
			if reaction.ReactType != "" {
				if !v0.InputReactReactType[reaction.ReactType] {
					return fmt.Errorf("fieldReactions[%d].reactions[%d].type is illegal. actual type is %s while %s is excepted ", i, j, reaction.ReactType, v0.KeywordMapToString(v0.InputReactReactType))
				}
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
