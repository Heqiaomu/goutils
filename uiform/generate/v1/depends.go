package v1

import (
	"bytes"
	"encoding/gob"
	"github.com/Heqiaomu/goutil/uiform/conform"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	"path/filepath"
	"strings"
)

func AddDepend(targetInput *yaml.Input, targetUIData *conform.UIData, dir string) error {
	resourcesId := targetInput.Depend

	dependDir := filepath.Join(dir, resourcesId)
	// 加载依赖模板：已有资源、创建资源
	rawUIData, err := ReadUiformAction(dependDir, dir)
	if err != nil {
		return err
	}
	// 深拷贝
	uiData, err := deepCopy(rawUIData)
	if err != nil {
		return err
	}

	// 前缀
	targetInputId := targetInput.ID
	prefix := strings.TrimSuffix(targetInputId, "/Input")
	// 依赖资源的InputId
	dependResourcesInputId := "depend/" + resourcesId + "/Input"
	dependResourcesInputId = prefix + "/" + dependResourcesInputId
	// 为加载依赖模板中所有Input、InputFiled增加前缀
	AddPrefix(uiData, prefix)
	// 新建一个Input
	dependInputId := prefix + "/depend/Input"
	dependInput := &yaml.Input{
		ID:        dependInputId,
		ShowType:  "list",
		SubInputs: []string{dependResourcesInputId, targetInputId},
	}
	// 需要寻找原来引用targetInput的位置，把targetInput替换成dependInput
	ReplaceInput(targetUIData, targetInputId, dependInputId)
	//if targetInput.Fields != nil {
	//	// 把其中的fields拿出来放到一个新的input中
	//	targetFieldsToInputId := dependInputId + "." + targetInputId
	//	// 将targetInput的subInputs改造成 资源模版+原资源的形式
	//	targetInput.SubInputs = []string{dependInputId, targetFieldsToInputId}
	//	// 如果targetUIData中没有记录该targetFieldsToInput，则记录
	//	if targetUIData.Id2UIData.InputKey2Input[targetFieldsToInputId] == nil {
	//		targetFieldsToInput := &Input{
	//			ID:       targetFieldsToInputId,
	//			ShowType: targetInput.ShowType,
	//			Fields:   targetInput.Fields,
	//		}
	//		targetUIData.Inputs = append(targetUIData.Inputs, targetFieldsToInput)
	//		targetUIData.Id2UIData.InputKey2Input[targetFieldsToInputId] = targetFieldsToInput
	//	}
	//} else {
	//	// 改造targetInput，在其中的subInputs中插入依赖Input，然后
	//	targetInput.SubInputs = append([]string{dependInputId}, targetInput.SubInputs...)
	//}
	if targetUIData.Id2UIData.InputKey2Input[dependInputId] == nil {
		// 依赖中的UI数据添加到目标UI数据中
		conform.Combine(uiData, targetUIData)
		// 将dependInput放到目标Input的前面，这样做的目的是如果目标Input是主函数，dependInput就成为主函数
		for index, targetInputInSplice := range targetUIData.Inputs {
			if targetInputInSplice.ID == targetInput.ID {
				rear := append([]*yaml.Input{}, targetUIData.Inputs[index:]...)
				targetUIData.Inputs = append(append(targetUIData.Inputs[:index], dependInput), rear...)
				break
			}
		}
		targetUIData.Id2UIData.InputKey2Input[dependInputId] = dependInput
	}
	// 清除这次depend
	targetInput.Depend = ""
	return nil
}

func AddPrefix(uiData *conform.UIData, prefix string) {
	for _, input := range uiData.Inputs {
		input.ID = prefix + "/" + input.ID
		uiData.Id2UIData.InputKey2Input[input.ID] = input

		for index := range input.SubInputs {
			input.SubInputs[index] = prefix + "/" + input.SubInputs[index]
		}
		for index := range input.Fields {
			input.Fields[index] = prefix + "/" + input.Fields[index]
		}
	}
	for _, field := range uiData.Fields {
		field.ID = prefix + "/" + field.ID
		for index := range field.Buttons {
			field.Buttons[index] = prefix + "/" + field.Buttons[index]
		}
		uiData.Id2UIData.FieldKey2Field[field.ID] = field
	}
	for _, fieldReaction := range uiData.FieldReactions {
		fieldReaction.Field = prefix + "/" + fieldReaction.Field
		uiData.Id2UIData.FieldKey2Reactions[fieldReaction.Field] = fieldReaction.Reactions

		for _, reaction := range fieldReaction.Reactions {
			if reaction.TargetInputId != "" {
				reaction.TargetInputId = prefix + "/" + reaction.TargetInputId
			}
			if reaction.InputId != "" {
				reaction.InputId = prefix + "/" + reaction.InputId
			}
		}
	}
}
func ReplaceInput(uiData *conform.UIData, oldInputId string, newInput string) {
	for _, input := range uiData.Inputs {
		for index := range input.SubInputs {
			if input.SubInputs[index] == oldInputId {
				input.SubInputs[index] = newInput
			}
		}
	}
	for _, fieldReaction := range uiData.FieldReactions {
		for _, reaction := range fieldReaction.Reactions {
			if reaction.TargetInputId == oldInputId {
				reaction.TargetInputId = newInput
			}
			if reaction.InputId == oldInputId {
				reaction.InputId = newInput
			}
		}
	}
}

func deepCopy(src *conform.UIData) (*conform.UIData, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	var dst conform.UIData
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&dst); err != nil {
		return nil, err
	}
	return &dst, nil
}
