package generate

import (
	"fmt"
	"testing"
)

func TestGetActionInputsFromDriverDir(t *testing.T) {
	m, err := GetActionInputsFromDriverDir("../test_data/drivers/hpcv1.0.0")
	fmt.Println(err)
	fmt.Println(m)
}

func TestParseCreateChainForm(t *testing.T) {
	m, err := GetActionInputsFromDriverDir("../test_data/drivers/huaweiv1.0.0")
	fmt.Println(err)
	fmt.Println(m)
	for _, input := range m {
		kv := make(map[string]*InputKVItem)
		err = Input2Map(input, kv)
		fmt.Println(err)
		fmt.Println(kv)
	}
}

func TestGetDriverInfosFromDriverRootPath(t *testing.T) {
	infos := GetDriverInfosFromDriverRootPath("../test_data/drivers")
	fmt.Println(len(infos))
}

func TestParseCreateK8SForm(t *testing.T) {
	m, err := GetActionInputsFromDriverDir("../test_data/drivers/k8sv1_0_1")
	fmt.Println(err)
	fmt.Println(m)
	for _, input := range m {
		kv := make(map[string]*InputKVItem)
		err = Input2Map(input, kv)
		fmt.Println(err)
		fmt.Println(kv)
	}
}
