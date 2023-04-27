package yaml

// Name 定义性名字，ID和Text都必须填写
type Name struct {
	ID          string            `yaml:"id" validate:"required"`   // 必填
	Text        string            `yaml:"text" validate:"required"` // 必填
	Description string            `yaml:"desc"`
	DocLink     string            `yaml:"link"`
	SubNames    []*Name           `yaml:"sub"`
	Meta        map[string]string `yaml:"meta"`
}

// DesName 描述性名字，ID字段不需要必填
type DesName struct {
	ID          string            `yaml:"id"`
	Text        string            `yaml:"text" validate:"required"` // 必填
	Description string            `yaml:"desc"`
	DocLink     string            `yaml:"link"`
	SubNames    []*Name           `yaml:"sub"`
	Meta        map[string]string `yaml:"meta"`
}
