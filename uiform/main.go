package main

import (
	"flag"
	"fmt"
	v0 "github.com/Heqiaomu/goutil/uiform/generate/v0"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Heqiaomu/goutil/uiform/conform"
	"github.com/Heqiaomu/goutil/uiform/generate"
	"github.com/Heqiaomu/protocol/ui"
	"github.com/golang/protobuf/jsonpb"
)

const (
	helpDoc = "inst: help\n" +
		"command: ./uiform -help\n" +
		"desc: show help doc"

	initDoc = "inst: init\n" +
		"command: ./uiform -init [yourDriverPath]\n" +
		"args:\n" +
		"  [yourDriverPath] the root project path of driver\n" +
		"desc: generate README.md and info.yaml in [yourDriverPath]\n"
	inituiDoc = "inst: initui\n" +
		"command: ./uiform -initui [yourDriverPath]\n" +
		"args:\n" +
		"  [yourDriverPath] the root project path of driver\n" +
		"desc: generate a file fold named as 'uiform' in [yourDriverPath], and generate files named as [action.id] defined by info.yaml in 'uiform', and generate 3 yaml files named as 'fields.yaml', 'inputs.yaml', 'reactions.yaml' in each action file fold"
	genDoc = "inst: gen\n" +
		"command: ./uiform -gen [yourDriverPath]\n" +
		"args:\n" +
		"  [yourDriverPath] the root project path of driver\n" +
		"desc: generate 'input.json' file in each action file"
	doc = helpDoc + "\n\n" + initDoc + "\n\n" + inituiDoc + "\n\n" + genDoc

	infoYamlContent = "name:\n" +
		"  id: hpc\n" +
		"  text: 'Hyperchain'\n" +
		"version: v1.0.0\n" +
		"type: host\n" +
		"resources: [hosts,credentials]\n" +
		"logo: ./logo.png\n" +
		"actions:\n" +
		"  - name:\n" +
		"      id: create\n" +
		"      text: '创建'\n" +
		"      desc: '创建Hyperchain'\n" +
		"    target: chains\n" +
		"  - name:\n" +
		"      id: start\n" +
		"      text: '启动'\n" +
		"      desc: '启动Hyperchain'\n" +
		"    target: chain\n" +
		"    states: [unavailable]\n" +
		"  - name:\n" +
		"      id: stop\n" +
		"      text: '停止'\n" +
		"      desc: '停止Hyperchain'\n" +
		"    target: chain\n" +
		"    states: [available]\n" +
		"  - name:\n" +
		"      id: remove\n" +
		"      text: '删除'\n" +
		"      desc: '删除Hyperchain'\n" +
		"    target: chain\n" +
		"    states: [unavailable, available, stoped]\n" +
		"exepath: 'exe'"

	readmeContent = "# XXX说明文档\n" +
		"\n" +
		"## 简介\n" +
		"\n" +
		"这是一段驱动的简介\n" +
		"\n" +
		"## 支持版本\n" +
		"\n" +
		"+ Hyperchain\n" +
		"  + v1.8.5\n" +
		"  + v1.6.4\n" +
		"\n" +
		"## 使用说明\n" +
		"\n" +
		"这是一段驱动的使用说明\n" +
		"\n" +
		"\n"

	fieldsYamlContent = "fields:\n" +
		"  - id: 'chain/name'\n" +
		"    title:\n" +
		"      text: '链名'"

	inputsYamlContent = "inputs:\n" +
		"  - id: 'chain/Input'\n" +
		"    title:\n" +
		"      text: '联盟配置'\n" +
		"    fields: ['chain/name']\n"
)

// main
func main() {
	// 当前工作路径
	dir, _ := os.Getwd()
	// 指令
	var help bool
	var init bool
	var initui bool
	var gen bool
	//var auto bool
	flag.BoolVar(&help, "help", false, "show help doc")
	flag.BoolVar(&init, "init", false, "init the driver project, see help for more info")
	flag.BoolVar(&initui, "initui", false, "init uiform of the driver project, see help for more info")
	flag.BoolVar(&gen, "gen", false, "generate the input.json files for all actions, see help for more info")
	//flag.BoolVar(&auto, "auto", false, "")

	flag.Parse()

	if help {
		fmt.Println(doc)
		os.Exit(0)
		return
	} else if init {
		driverPath := flag.Arg(0)
		if driverPath == "" {
			driverPath = dir
		}
		err := Init(driverPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		} else {
			fmt.Println("init info.yaml and README.md success!")
			os.Exit(0)
		}
		return
	} else if initui {
		driverPath := flag.Arg(0)
		if driverPath == "" {
			driverPath = dir
		}
		err := InitUI(driverPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		} else {
			fmt.Println("init uiform and 3 yaml file in each action success!")
			os.Exit(0)
		}
		return
	} else if gen {
		driverPath := flag.Arg(0)
		if driverPath == "" {
			driverPath = dir
		}
		err := GenerateAllActionInputJsonFile(driverPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		} else {
			fmt.Println("gen input.json success!")
			os.Exit(0)
		}
		return
	} else {
		fmt.Println(doc)
		os.Exit(0)
		return
	}
}

// GenerateInputJsonFile 扫描驱动的根目录
// dir 表示驱动下info.yaml所在的路径 例如：drivers/drivers-name/
// 该文件夹下的uiform文件夹应该具有多个以info.yaml文件中定义的action.id为文件夹名的文件夹
func GenerateAllActionInputJsonFile(driverDir string) error {
	isExist, err := v0.PathExists(driverDir)
	if err != nil {
		return fmt.Errorf("finding driver path=%s err, err=%s", driverDir, err.Error())
	}
	if !isExist {
		return fmt.Errorf("cannot find driver path=%s ", driverDir)
	}
	dependsDir := "./depends"
	// info 文件提取
	info, err := generate.ReadInfo(driverDir)
	if err != nil {
		return err
	}
	if info.ID == "" {
		info.ID = info.Name.ID + info.Version
	} else {
		if info.ID != info.Name.ID+info.Version {
			return fmt.Errorf("info.yaml 'id'=%s is not equals to 'name.text'+'version' %s%s ", info.ID, info.Name.ID, info.Version)
		}
	}
	// 遍历所有功能
	for _, action := range info.Actions {
		actionDir := filepath.Join(driverDir, "uiform", action.ID)
		uiData, err := generate.ReadUiformAction(info.UIFormVersion, actionDir, dependsDir)
		if err != nil {
			return err
		}
		uiData.DriverID = info.ID
		uiData.Id2UIData.DriverID = info.ID
		// 将该功能下的fields、inputs、reactions转化为input的json格式文件，输出到该文件夹下
		// 第一个Input作为该功能入参
		mainInputProto := generateMainInput(info.UIFormVersion, uiData)
		// 转化为json
		m := jsonpb.Marshaler{}
		mainInputProtoJson, err := m.MarshalToString(mainInputProto)
		if err != nil {
			return err
		}
		// 写入input文件
		err = ioutil.WriteFile(filepath.Join(actionDir, "input.json"), []byte(mainInputProtoJson), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// generateMainInput 创建动作的表单
func generateMainInput(uiformVersion string, uiData *conform.UIData) *ui.Input {
	mainInput := uiData.Inputs[0]
	return generate.Input2InputProto(uiformVersion, mainInput, uiData.Id2UIData)
}

// Init 生成 README.md 和 info.yaml 模板
func Init(driverDir string) error {
	isExist, err := v0.PathExists(driverDir)
	if err != nil {
		return fmt.Errorf("finding driver path=%s err, err=%s", driverDir, err.Error())
	}
	if !isExist {
		return fmt.Errorf("cannot find driver path=%s ", driverDir)
	}
	infoFilePath := filepath.Join(driverDir, "info.yaml")
	isExist, err = v0.PathExists(infoFilePath)
	if err != nil {
		return fmt.Errorf("finding info.yaml=%s err, err=%s", infoFilePath, err.Error())
	}
	if isExist {
		fmt.Printf("info.yaml is existed in driver path=%s, skip create info.yaml\n", driverDir)
	} else {
		if err := ioutil.WriteFile(infoFilePath, []byte(infoYamlContent), 0755); err != nil {
			return err
		}
	}
	readmePath := filepath.Join(driverDir, "README.md")
	isExist, err = v0.PathExists(readmePath)
	if err != nil {
		return fmt.Errorf("finding README.md=%s err, err=%s", readmePath, err.Error())
	}
	if isExist {
		fmt.Printf("README.md is existed in driver path=%s, skip create README.md\n", driverDir)
	} else {
		if err := ioutil.WriteFile(readmePath, []byte(readmeContent), 0755); err != nil {
			return err
		}
	}
	return nil
}

// InitUIForm 根据 info.yaml 生成uiform下各个动作的 yaml 文件
func InitUI(driverDir string) error {
	isExist, err := v0.PathExists(driverDir)
	if err != nil {
		return fmt.Errorf("finding driver path=%s err, err=%s", driverDir, err.Error())
	}
	if !isExist {
		return fmt.Errorf("cannot find driver path=%s ", driverDir)
	}
	// info 文件提取
	info, err := generate.ReadInfo(driverDir)
	if err != nil {
		return err
	}
	uiformDir := filepath.Join(driverDir, "uiform")
	isExist, err = v0.PathExists(uiformDir)
	if err != nil {
		return fmt.Errorf("finding uiform=%s err, err=%s", uiformDir, err.Error())
	}
	if isExist {
		fmt.Printf("uiform is existed in driver path=%s, skip create uiform\n", driverDir)
	} else {
		if err := os.Mkdir(uiformDir, 0755); err != nil {
			return err
		}
	}
	// 遍历所有功能
	for _, action := range info.Actions {
		// uiform
		actionDir := filepath.Join(driverDir, "uiform", action.ID)
		isExist, err = v0.PathExists(actionDir)
		if err != nil {
			return fmt.Errorf("finding action=%s err, err=%s", actionDir, err.Error())
		}
		if isExist {
			fmt.Printf(action.ID+" is existed in path=%s, skip create action\n", actionDir)
		} else {
			if err := os.Mkdir(actionDir, 0755); err != nil {
				return err
			}
		}
		// fields.yaml
		fieldsDir := filepath.Join(driverDir, "uiform", action.ID, "fields.yaml")
		isExist, err = v0.PathExists(fieldsDir)
		if err != nil {
			return fmt.Errorf("finding fields.yaml=%s err, err=%s", fieldsDir, err.Error())
		}
		if isExist {
			fmt.Printf("fields.yaml is existed in path=%s, skip create fields.yaml\n", fieldsDir)
		} else {
			if err := ioutil.WriteFile(fieldsDir, []byte(fieldsYamlContent), 0755); err != nil {
				return err
			}
		}
		// inputs.yaml
		inputsDir := filepath.Join(driverDir, "uiform", action.ID, "inputs.yaml")
		isExist, err = v0.PathExists(inputsDir)
		if err != nil {
			return fmt.Errorf("finding inputs.yaml=%s err, err=%s", inputsDir, err.Error())
		}
		if isExist {
			fmt.Printf("inputs.yaml is existed in path=%s, skip create inputs.yaml\n", inputsDir)
		} else {
			if err := ioutil.WriteFile(inputsDir, []byte(inputsYamlContent), 0755); err != nil {
				return err
			}
		}
		// reactions.yaml
		reactionsDir := filepath.Join(driverDir, "uiform", action.ID, "reactions.yaml")
		isExist, err = v0.PathExists(reactionsDir)
		if err != nil {
			return fmt.Errorf("finding reactions.yaml=%s err, err=%s", reactionsDir, err.Error())
		}
		if isExist {
			fmt.Printf("reactions.yaml is existed in path=%s, skip create reactions.yaml\n", reactionsDir)
		} else {
			if err := ioutil.WriteFile(reactionsDir, []byte("fieldReactions:\n"), 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
