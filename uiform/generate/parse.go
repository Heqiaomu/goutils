package generate

import (
	"fmt"
	v0 "github.com/Heqiaomu/goutil/uiform/generate/v0"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Heqiaomu/goutil/uiform/yaml"
	"github.com/golang/protobuf/jsonpb"

	ui "github.com/Heqiaomu/protocol/ui"
)

type InputKVItem struct {
	KeyText string
	Value   []string
	KeyType string
}

// Input2Map parse
func Input2Map(input *ui.Input, result map[string]*InputKVItem) (err error) {
	err = parseCreateChainFormLoop(input, result)
	if err != nil {
		return
	}
	return
}

// parseCreateChainFormLoop parse create chain form loop
func parseCreateChainFormLoop(input *ui.Input, result map[string]*InputKVItem) (err error) {
	if input.ShowType == "info" {
		for k, v := range input.Meta {
			if result[k] == nil {
				result[k] = &InputKVItem{}
			}
			result[k].Value = strings.Split(v, ",")
			result[k].KeyText = k
		}
	}
	for _, sif := range input.Fields {
		parseCreateChainFormLoopField(sif, result)
	}

	for _, si := range input.SubInputs {
		if err = parseCreateChainFormLoop(si, result); err != nil {
			return err
		}
	}
	return
}
func parseCreateChainFormLoopField(sif *ui.InputField, result map[string]*InputKVItem) (err error) {
	resKey := sif.Id
	// TODO: 参数校验
	if result[resKey] == nil {
		result[resKey] = &InputKVItem{}
	}
	result[resKey].Value = sif.Value
	result[resKey].KeyText = sif.Title.Text
	result[resKey].KeyType = sif.InputType
	for _, sifa := range sif.InputReacts {
		// 触发
		for _, sifai := range sifa.Inputs {
			if err = parseCreateChainFormLoop(sifai, result); err != nil {
				return err
			}
		}
	}
	for _, button := range sif.Buttons {
		parseCreateChainFormLoopField(button, result)
	}
	return nil
}

// func getKey(parentKey, id string) string {
// 	return strings.Trim(parentKey+"."+id, ".")
// }

// GetActionInputsFromDriverDir 将驱动路径下的所有动作中的input.json提取出来
func GetActionInputsFromDriverDir(driverDir string) (map[string]*ui.Input, error) {
	// info 文件提取
	info, err := ReadInfo(driverDir)
	if err != nil {
		return nil, err
	}
	action2Input := make(map[string]*ui.Input)
	// 遍历所有功能
	for _, action := range info.Actions {
		uiInput, err := GetActionInput(driverDir, action.ID)
		if err != nil {
			return nil, err
		}
		action2Input[action.ID] = uiInput
	}
	return action2Input, nil
}

// GetActionInputJsonsFromDriverDir 将驱动路径下的所有动作中的input.json提取出来
func GetActionInputJsonsFromDriverDir(driverDir string) (map[string]string, error) {
	// info 文件提取
	info, err := ReadInfo(driverDir)
	if err != nil {
		return nil, err
	}
	action2Input := make(map[string]string)
	// 遍历所有功能
	for _, action := range info.Actions {
		uiInputJsonBytes, err := GetActionInputJsonBytes(driverDir, action.ID)
		if err != nil {
			return nil, err
		}
		action2Input[action.ID] = string(uiInputJsonBytes)
	}
	return action2Input, nil
}

func GetActionInput(driverDir string, actionId string) (*ui.Input, error) {
	inputJsonBytes, err := GetActionInputJsonBytes(driverDir, actionId)
	if err != nil {
		return nil, err
	}
	// 转化为json
	var uiInput ui.Input
	err = jsonpb.UnmarshalString(string(inputJsonBytes), &uiInput)
	if err != nil {
		return nil, fmt.Errorf("trans input(json)=%s to ui.proto fail, err=%s", string(inputJsonBytes), err.Error())
	}
	return &uiInput, nil
}

func GetActionInputJsonBytes(driverDir string, actionId string) ([]byte, error) {
	actionDir := filepath.Join(driverDir, "uiform", actionId)
	// 读取input.json
	inputJsonFilePath := filepath.Join(actionDir, "input.json")
	isExist, err := v0.PathExists(inputJsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("finding input.json=%s err, err=%s", inputJsonFilePath, err.Error())
	}
	if !isExist {
		return nil, fmt.Errorf("cannot find input.json=%s ", inputJsonFilePath)
	}
	inputJsonBytes, err := ioutil.ReadFile(inputJsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("read input.json=%s err, err=%s", inputJsonFilePath, err.Error())
	}
	return inputJsonBytes, nil
}

// InputJson2Input 将 InputJson 转成 ui.Input 对象
func InputJson2Input(inputJson []byte) (*ui.Input, error) {
	var uiInput ui.Input
	err := jsonpb.UnmarshalString(string(inputJson), &uiInput)
	if err != nil {
		return nil, fmt.Errorf("trans input(json)=%s to ui.proto fail, err=%s", string(inputJson), err.Error())
	}
	return &uiInput, nil
}

// Input2InputJson 将 InputJson 转成 ui.Input 对象
func Input2InputJson(input *ui.Input) (string, error) {
	m := jsonpb.Marshaler{}
	inputJson, err := m.MarshalToString(input)
	if err != nil {
		return "", fmt.Errorf("trans input to inputJson fail, err=%s", err.Error())
	}
	return inputJson, nil
}

// GetDriverInfosFromDriverRootPath 从驱动安装的根目录中，扫描所有驱动的info.yaml文件提取驱动信息
func GetDriverInfosFromDriverRootPath(driverRootPath string) []*yaml.Info {
	var infos []*yaml.Info
	_ = filepath.Walk(driverRootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "info.yaml" {
			infoYamlFilePath, err := filepath.Abs(path)
			if err != nil {
				return nil
			}
			infoYaml, err := ReadInfoFromInfoYaml(infoYamlFilePath)
			if err != nil {
				return nil
			}
			infos = append(infos, infoYaml)
		}
		return nil
	})
	return infos
}
