package conform

import (
	"git.hyperchain.cn/blocface/goutil/uiform/yaml"
)

type UIData struct {
	DriverID       string
	Fields         []*yaml.Field
	Inputs         []*yaml.Input
	FieldReactions []*yaml.FieldReactions
	Id2UIData      *Id2UIData
}

type Id2UIData struct {
	DriverID           string
	FieldKey2Field     map[string]*yaml.Field
	InputKey2Input     map[string]*yaml.Input
	FieldKey2Reactions map[string][]*yaml.Reaction
}
