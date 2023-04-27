package generate

import (
	"github.com/Heqiaomu/goutil/uiform/conform"
	v1 "github.com/Heqiaomu/goutil/uiform/generate/v1"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	"github.com/Heqiaomu/protocol/ui"
)

const (
	ELRegex = `\$\{([a-zA-Z/-]+)\}`
)

// Input2InputProto 将配置结构Input转化为渲染结构Input
func Input2InputProto(uiformVersion string, input *yaml.Input, id2Data *conform.Id2UIData) *ui.Input {
	switch uiformVersion {
	case "v1":
		return v1.Input2InputProto(input, id2Data)
	case "v2":
	}
	return nil
}
