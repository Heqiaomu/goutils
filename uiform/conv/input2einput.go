package conv

import (
	"github.com/Heqiaomu/protocol/ui"
)

func Input2EInput(input *ui.Input) map[string]interface{} {
	eInput := dealInput(input)
	rootInputID := input.Id
	rootEInput := make(map[string]interface{})
	rootEInput[rootInputID] = eInput
	return rootEInput
}

func dealInput(input *ui.Input) interface{} {
	// rootEInput.Set(rootInputID,?) ? 是个什么？可能是一个[]EInput，可能是 {"":"","":""}，我们现在不知道，要往下遍历
	// 处理 Input节点 的 子Input
	if len(input.SubInputs) > 0 {
		inputs := make(map[string]interface{})
		for _, subInput := range input.SubInputs {
			subInputID := subInput.Id
			subInput := dealInput(subInput)
			inputs[subInputID] = subInput
		}
		return inputs
	}
	if len(input.Fields) > 0 {
		inputFields := make(map[string]interface{})
		for _, inputField := range input.Fields {
			value, reactInputs := dealInputField(inputField)
			inputFieldID := inputField.Id
			if value != nil {
				inputFields[inputFieldID] = value
			}
			if reactInputs != nil {
				inputFields[reactID(inputFieldID)] = reactInputs
			}
		}
		return inputFields
	}
	return nil
}

func dealInputField(inputField *ui.InputField) (interface{}, interface{}) {
	inputType := inputField.InputType
	var value interface{}
	// 批量
	if inputType == "inputs" || inputType == "numberInputs" || inputType == "files" {
		value = inputField.Value
	} else {
		if len(inputField.Value) > 0 {
			value = inputField.Value[0]
		} else {
			value = ""
		}
	}

	// 如果本身只是一个按钮，单纯为了触发联动，那么直接返回联动对象
	if len(inputField.Buttons) > 0 {
		inputFields := make(map[string]interface{})
		for _, button := range inputField.Buttons {
			value, reactInputs := dealInputField(button)

			subInputFieldID := button.Id
			if value != nil {
				inputFields[subInputFieldID] = value
			}
			if reactInputs != nil {
				inputFields[reactID(subInputFieldID)] = reactInputs
			}
		}
		return value, inputFields
	}

	// 如果有联动，那么联动就是一个Input数组
	if len(inputField.InputReacts) > 0 {
		reactInputs := make(map[string]interface{})
		for _, react := range inputField.InputReacts {
			// 如果没有指定目标ID，说明就是放在这里面
			if react.TargetInputId == "" {
				for _, reactInput := range react.Inputs {
					reactInputID := reactInput.Id
					input := dealInput(reactInput)
					if input != nil {
						reactInputs[reactInputID] = input
					}
				}
			}
		}
		return value, reactInputs
	}

	return value, nil
}

func reactID(ID string) string {
	return ID + "Reacts"
}
