package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestGenerateAllActionInputJsonFileHpcv100(t *testing.T) {
	err := GenerateAllActionInputJsonFile("./test_data/drivers/hpcv1.0.0")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGenerateAllActionInputJsonFileFabricv100(t *testing.T) {
	err := GenerateAllActionInputJsonFile("./test_data/drivers/fabricv1.0.0")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGenerateAllActionInputJsonFileFabricv101(t *testing.T) {
	err := GenerateAllActionInputJsonFile("./test_data/drivers/fabricv1.0.1")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGenerateAllActionInputJsonFileHaweiv100(t *testing.T) {
	err := GenerateAllActionInputJsonFile("./test_data/drivers/huaweiv1.0.0")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestInitUI(t *testing.T) {
	err := InitUI("./test_data/drivers/newv1.0.0")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGenerateAllActionInputJsonFileNewv100(t *testing.T) {
	err := GenerateAllActionInputJsonFile("./test_data/drivers/newv1.0.0")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestDefaultValueEpl(t *testing.T) {
	// 如果默认值中出现了${xxx}，如果xxx是某个InputField的ID，那么这个ID需要加上id2Data.DriverID前缀
	defaultValue := "${nodes/${index}/type}${index}.${chain/organizationDomain}.${chain/implementation}"
	elRegex := regexp.MustCompile(`\$\{([a-zA-Z/${}-]+)\}`)
	defaultValue = string(elRegex.ReplaceAllFunc([]byte(defaultValue), func(tarEl []byte) []byte {
		tarElString := string(tarEl)
		//tarContentElString := tarElString[2 : len(tarElString)-1]
		fmt.Println(tarElString)
		return tarEl
	}))
}

func TestDefaultValueEpl2(t *testing.T) {
	// 如果默认值中出现了${xxx}，如果xxx是某个InputField的ID，那么这个ID需要加上id2Data.DriverID前缀
	defaultValue := "${nodes/${index}/type}${index}.${chain/organizationDomain}.${chain/implementation}"
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
	for lI, rI := range lInd2RInd {
		fmt.Println(defaultValue[lI:rI])
	}
}
