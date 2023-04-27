package v1

import (
	"github.com/Heqiaomu/goutil/uiform/conform"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	"github.com/Heqiaomu/protocol/ui"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"strings"
)

//const (
//	ELRegex = `\$\{([a-zA-Z/-]+)\}`
//)

// Input2InputProto 将配置结构Input转化为渲染结构Input
func Input2InputProto(input *yaml.Input, id2Data *conform.Id2UIData) *ui.Input {
	var inputProto ui.Input
	// 将驱动的名字作为ID的前缀
	inputProto.Id = id2Data.DriverID + "/" + input.ID
	if input.Title != nil {
		inputProto.Title = Name2NameProto(input.Title)
	}
	inputProto.Fold = &wrappers.BoolValue{Value: input.Fold}
	if input.ShowType == "" {
		inputProto.ShowType = "list"
	} else {
		inputProto.ShowType = input.ShowType
	}
	inputProto.Meta = input.Meta
	if input.SubInputs != nil {
		// 拼装subInputs
		var subInputsProto []*ui.Input
		for _, subInputId := range input.SubInputs {
			subInputsProto = append(subInputsProto, Input2InputProto(id2Data.InputKey2Input[subInputId], id2Data))
		}
		inputProto.SubInputs = subInputsProto
	} else if input.Fields != nil {
		// 拼装fields
		var fieldsProto []*ui.InputField
		for _, fieldId := range input.Fields {
			field := id2Data.FieldKey2Field[fieldId]
			fieldsProto = append(fieldsProto, Field2InputFieldProto(field, id2Data))
		}
		inputProto.Fields = fieldsProto
	}
	return &inputProto
}

// Field2InputFieldProto 将配置结构Field转化为渲染结构InputField
func Field2InputFieldProto(field *yaml.Field, id2Data *conform.Id2UIData) *ui.InputField {
	var fieldProto ui.InputField
	// 将驱动的名字作为前缀
	fieldProto.Id = id2Data.DriverID + "/" + field.ID
	fieldProto.Title = Name2NameProto(field.Title)
	if field.Edit == nil {
		field.Edit = true
	}
	fieldProto.Edit = &wrappers.BoolValue{Value: field.Edit.(bool)}
	fieldProto.Value = []string{""}
	if field.DefaultValue == nil {
		field.DefaultValue = []string{""}
	}

	fieldProto.Invisible = field.Invisible

	// 如果默认值中出现了${xxx}，如果xxx是某个InputField的ID，那么这个ID需要加上id2Data.DriverID前缀
	DefaultValueAddPrefix(field, id2Data)
	//elRegex := regexp.MustCompile(ELRegex)
	//for index, defaultV := range field.DefaultValue {
	//	field.DefaultValue[index] = string(elRegex.ReplaceAllFunc([]byte(defaultV), func(tarEl []byte) []byte {
	//		tarElString := string(tarEl)
	//		tarContentElString := tarElString[2 : len(tarElString)-1]
	//		if id2Data.FieldKey2Field[tarContentElString] != nil {
	//			return []byte("${" + id2Data.DriverID + "/" + tarContentElString + "}")
	//		}
	//		return tarEl
	//	}))
	//}

	fieldProto.DefaultValue = field.DefaultValue
	if field.InputType == "" {
		field.InputType = "input"
	}
	fieldProto.InputType = field.InputType
	fieldProto.Meta = field.Meta
	// validate
	if field.Validate != nil {
		if field.Validate.Require == nil {
			field.Validate.Require = true
		}
		switch field.InputType {
		case "input", "date", "time", "dateTime", "textArea", "password", "ip":
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateInput{
				ValidateDes: field.Validate.ValidateDes,
				Require:     &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:  field.Validate.RequireDes,
				Regex:       field.Validate.Regex,
				RegexDes:    field.Validate.RegexDes,
				PlaceHolder: field.Validate.PlaceHolder,
				Meta:        field.Validate.Meta,
			})
		case "inputs":
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateInputs{
				ValidateDes: field.Validate.ValidateDes,
				Require:     &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:  field.Validate.RequireDes,
				Regex:       field.Validate.Regex,
				RegexDes:    field.Validate.RegexDes,
				PlaceHolder: field.Validate.PlaceHolder,
				MaxCount:    field.Validate.MaxCount,
				MinCount:    field.Validate.MinCount,
				Meta:        field.Validate.Meta,
			})
		case "numberInput", "numberRange", "numberSlide":
			// TODO 根据field.Validate.Step
			//if field.Validate.Mid == "" {
			//	maxF, _ := strconv.ParseFloat(field.Validate.Max, 64)
			//	minF, _ := strconv.ParseFloat(field.Validate.Min, 64)
			//	midF := (maxF-minF)/2 + minF
			//	field.Validate.Mid = midF
			//}
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateNumber{
				Require:     &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:  field.Validate.RequireDes,
				ValidateDes: field.Validate.ValidateDes,
				Step:        field.Validate.Step,
				Max:         field.Validate.Max,
				Min:         field.Validate.Min,
				Mid:         field.Validate.Mid,
				Unit:        field.Validate.Unit,
				Meta:        field.Validate.Meta,
			})
		case "numberInputs":
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateNumbers{
				Require:     &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:  field.Validate.RequireDes,
				ValidateDes: field.Validate.ValidateDes,
				Step:        field.Validate.Step,
				Max:         field.Validate.Max,
				Min:         field.Validate.Min,
				Unit:        field.Validate.Unit,
				MaxCount:    field.Validate.MaxCount,
				MinCount:    field.Validate.MinCount,
				Meta:        field.Validate.Meta,
			})
		case "switch":
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateSwitch{
				Require:     &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:  field.Validate.RequireDes,
				ValidateDes: field.Validate.ValidateDes,
				Meta:        field.Validate.Meta,
			})
		case "select", "radio", "radioButton":
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateSelect{
				Require:     &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:  field.Validate.RequireDes,
				ValidateDes: field.Validate.ValidateDes,
				Options:     Names2NamesProto(field.Validate.Options),
				MaxCount:    field.Validate.MaxCount,
				MinCount:    field.Validate.MinCount,
				Meta:        field.Validate.Meta,
				PlaceHolder: field.Validate.PlaceHolder,
			})
		case "file":
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateFile{
				ValidateDes:      field.Validate.ValidateDes,
				Require:          &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:       field.Validate.RequireDes,
				Regex:            field.Validate.Regex,
				RegexDes:         field.Validate.RegexDes,
				PlaceHolder:      field.Validate.PlaceHolder,
				FileName:         field.Validate.FileName,
				FileNameRegex:    field.Validate.FileNameRegex,
				FileNameRegexDes: field.Validate.FileNameRegexDes,
				Meta:             field.Validate.Meta,
				Url:              field.Validate.URL,
				FileNameSuffix:   field.Validate.FileNameSuffix,
				Min:              field.Validate.Min,
				Max:              field.Validate.Max,
			})
		case "files":
			fieldProto.Validate, _ = ptypes.MarshalAny(&ui.ValidateFiles{
				ValidateDes:      field.Validate.ValidateDes,
				Require:          &wrappers.BoolValue{Value: field.Validate.Require.(bool)},
				RequireDes:       field.Validate.RequireDes,
				Regex:            field.Validate.Regex,
				RegexDes:         field.Validate.RegexDes,
				PlaceHolder:      field.Validate.PlaceHolder,
				FileName:         field.Validate.FileName,
				FileNameRegex:    field.Validate.FileNameRegex,
				FileNameRegexDes: field.Validate.FileNameRegexDes,
				MaxCount:         field.Validate.MaxCount,
				MinCount:         field.Validate.MinCount,
				Meta:             field.Validate.Meta,
				Url:              field.Validate.URL,
				FileNameSuffix:   field.Validate.FileNameSuffix,
				Min:              field.Validate.Min,
				Max:              field.Validate.Max,
			})
		default:
			fieldProto.Validate = nil
		}
	}
	// inputReacts
	if id2Data.FieldKey2Reactions[field.ID] != nil {
		reactions := id2Data.FieldKey2Reactions[field.ID]
		var inputReactsProto []*ui.InputReact
		for _, reaction := range reactions {
			inputReactsProto = append(inputReactsProto, Reaction2InputReact(reaction, id2Data))
		}
		fieldProto.InputReacts = inputReactsProto
	}
	// buttons
	if field.Buttons != nil {
		var buttonsProto []*ui.InputField
		for _, buttonFieldId := range field.Buttons {
			buttonField := id2Data.FieldKey2Field[buttonFieldId]
			buttonsProto = append(buttonsProto, Field2InputFieldProto(buttonField, id2Data))
		}
		fieldProto.Buttons = buttonsProto
	}
	// links
	if field.Links != nil {
		var linksProto []*ui.Name
		for _, link := range field.Links {
			linksProto = append(linksProto, Name2NameProto(link))
		}
		fieldProto.Links = linksProto
	}
	return &fieldProto
}

func Reaction2InputReact(reaction *yaml.Reaction, id2Data *conform.Id2UIData) *ui.InputReact {
	var inputReactProto ui.InputReact
	inputReactProto.TriggerRegex = reaction.TriggerRegex
	if reaction.ReactType == "" {
		reaction.ReactType = "showInput"
	}
	inputReactProto.ReactType = reaction.ReactType
	if reaction.InputId != "" {
		inputReactProto.Inputs = []*ui.Input{Input2InputProto(id2Data.InputKey2Input[reaction.InputId], id2Data)}
	}
	if reaction.UrlReact != nil {
		inputReactProto.UrlReact = UrlAction2UrlAction(reaction.UrlReact, id2Data)
	}
	if reaction.TargetInputId != "" {
		inputReactProto.TargetInputId = id2Data.DriverID + "/" + reaction.TargetInputId
	}
	inputReactProto.Meta = reaction.Meta
	return &inputReactProto
}

func UrlAction2UrlAction(urlAction *yaml.UrlAction, id2Data *conform.Id2UIData) *ui.UrlAction {
	var urlActionProto ui.UrlAction
	if urlAction.Name != nil {
		urlActionProto.Name = Name2NameProto(urlAction.Name)
	}
	urlActionProto.Method = urlAction.Method
	urlActionProto.Url = urlAction.Url
	if urlAction.Body != "" {
		urlActionProto.Body = Input2InputProto(id2Data.InputKey2Input[urlAction.Body], id2Data)
	}
	urlActionProto.Meta = urlAction.Meta
	return &urlActionProto
}

func Name2NameProto(name *yaml.Name) *ui.Name {
	if name != nil {
		var nameProto ui.Name
		nameProto.Id = name.ID
		nameProto.Text = name.Text
		nameProto.Description = name.Description
		nameProto.DocLink = name.DocLink
		nameProto.Meta = name.Meta
		var subNamesProto []*ui.Name
		for _, subName := range name.SubNames {
			subNamesProto = append(subNamesProto, Name2NameProto(subName))
		}
		nameProto.SubNames = subNamesProto
		return &nameProto
	}
	return nil
}

func Names2NamesProto(names []*yaml.Name) []*ui.Name {
	namesProto := make([]*ui.Name, len(names))
	for i, name := range names {
		nameProto := Name2NameProto(name)
		namesProto[i] = nameProto
	}
	return namesProto
}

func DefaultValueAddPrefix(field *yaml.Field, id2Data *conform.Id2UIData) {
	// 如果默认值中出现了${xxx}，如果xxx是某个InputField的ID，那么这个ID需要加上id2Data.DriverID前缀
	for in, defaultValue := range field.DefaultValue {
		L := len(defaultValue)

		var lInd []int
		var lInd2RInd map[int]int
		lInd2RInd = make(map[int]int)
		for i := 0; i < L; i++ {
			if defaultValue[i] == '$' {
				if i+1 >= L {
					break
				}
				i = i + 1
				if defaultValue[i] == '{' {
					if i+1 >= L {
						break
					}
					i = i + 1
					// 记录左括号下一个字符的下标
					lInd = append(lInd, i)
					continue
				}
			} else if defaultValue[i] == '}' {
				if len(lInd) <= 0 {
					break
				}
				// 于当前右括号对应的左括号下标 lI （其实是 左括号下标 + 1）
				lI := lInd[len(lInd)-1]
				// 记录匹配的左右括号
				lInd2RInd[lI] = i
				lInd = lInd[:len(lInd)-1]
			}
		}
		org2Rep := make(map[string]string)
		for lI, rI := range lInd2RInd {
			mayBeFieldName := defaultValue[lI:rI]
			if id2Data.FieldKey2Field[mayBeFieldName] != nil {
				org2Rep[mayBeFieldName] = id2Data.DriverID + "/" + mayBeFieldName
			}
		}
		newDefaultValue := defaultValue
		for org, rep := range org2Rep {
			newDefaultValue = strings.Replace(newDefaultValue, org, rep, -1)
		}
		field.DefaultValue[in] = newDefaultValue
	}
}
