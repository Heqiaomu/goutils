package yaml

type Reactions struct {
	FieldReactions []*FieldReactions `yaml:"fieldReactions"`
}
type FieldReactions struct {
	Field     string      `yaml:"field"`
	Reactions []*Reaction `yaml:"reactions"`
}
type Reaction struct {
	TriggerRegex  string            `yaml:"regex"` // 缺省为"*"
	ReactType     string            `yaml:"type"`  // 缺省为"showInput"
	UrlReact      *UrlAction        `yaml:"urlReact"`
	InputId       string            `yaml:"input"`
	TargetInputId string            `yaml:"target"`
	Meta          map[string]string `yaml:"meta"`
}
type UrlAction struct {
	Name   *Name             `yaml:"name"`
	Method string            `yaml:"method"`
	Url    string            `yaml:"url"`
	Body   string            `yaml:"body"`
	Meta   map[string]string `yaml:"meta"`
}
