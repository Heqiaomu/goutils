package yaml

// Fields 表单输入项数组
type Fields struct {
	Fields []*Field `yaml:"fields"`
}

// Field 表单输入项
type Field struct {
	ID           string            `yaml:"id"`      // 必填（自动填充，用于前端辨别处理用）
	Key          string            `yaml:"key"`     // v2新增，必填，需要保证唯一
	Tag          string            `yaml:"tag"`     // 标签，用于标记某些具有特定性质的输入项，比如标记port，表示该输入项是一个端口
	Title        *Name             `yaml:"title"`   // 必填
	Edit         interface{}       `yaml:"edit"`    // 缺省为true
	DefaultValue []string          `yaml:"default"` // 缺省为[""]
	InputType    string            `yaml:"type"`    // 缺省为"input"
	Validate     *Validate         `yaml:"validate"`
	Buttons      []string          `yaml:"buttons"`
	Links        []*Name           `yaml:"links"`
	Meta         map[string]string `yaml:"mate"`
	Invisible    bool              `yaml:"invisible"` // 缺省为false-可见
}

// Validate 表单校验
type Validate struct {
	ValidateDes      string            `yaml:"desc"`
	Require          interface{}       `yaml:"require"`     // 缺省为true
	RequireDes       string            `yaml:"requireDesc"` // 缺省为true
	Regex            string            `yaml:"regex"`
	RegexDes         string            `yaml:"regexDesc"`
	PlaceHolder      string            `yaml:"holder"`
	MaxCount         int64             `yaml:"maxCount"`
	MinCount         int64             `yaml:"minCount"`
	Options          []*Name           `yaml:"options"`
	Step             string            `yaml:"step"`
	Max              string            `yaml:"max"`
	Min              string            `yaml:"min"`
	Mid              string            `yaml:"mid"`
	Unit             string            `yaml:"unit"`
	FileName         []string          `yaml:"fileName"`
	FileNameRegex    string            `yaml:"fileNameRegex"`
	FileNameRegexDes string            `yaml:"fileNameDesc"`
	FileNameSuffix   string            `yaml:"fileNameSuffix"`
	URL              string            `yaml:"url"`
	Meta             map[string]string `yaml:"meta"`
}
