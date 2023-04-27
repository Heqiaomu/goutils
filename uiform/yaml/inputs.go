package yaml

// Inputs 表单数组
type Inputs struct {
	Inputs []*Input `yaml:"inputs"`
}

// Input 表单
type Input struct {
	ID            string            `yaml:"id"` // v1必填 v2选填
	Title         *Name             `yaml:"title"`
	Fold          bool              `yaml:"fold"`
	ShowType      string            `yaml:"type"` // 缺省为"list"
	Fields        []string          `yaml:"fields"`
	SubInputs     []string          `yaml:"sub"`
	Meta          map[string]string `yaml:"meta"`
	Depend        string            `yaml:"depend"`
	Key           string            `yaml:"key"` // v2新增字段，必填，需保证全局唯一
	Tag           string            `yaml:"tag"`
	UIFormVersion string            `yaml:"version"`
}
