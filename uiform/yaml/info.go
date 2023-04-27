package yaml

type Info struct {
	UIFormVersion string    `yaml:"uiform-version" validate:"oneof='' 'v1' 'v2'"` // 选填，不填默认为v1
	ID            string    `yaml:"id"`                                           // 选填
	Name          *Name     `yaml:"name" validate:"required"`                     // 必填
	Version       string    `yaml:"version" validate:"required"`                  // 必填
	Type          string    `yaml:"type" validate:"required"`                     // 必填
	Resources     []string  `yaml:"resources" validate:"required"`                // 必填
	Logo          string    `yaml:"logo" validate:"required,file"`                // 必填
	Actions       []*Action `yaml:"actions" validate:"required"`                  // 必填
	ExePath       string    `yaml:"exePath" validate:"required"`                  // 必填
	Tag           string    `yaml:"tag"`                                          //  逗号隔开
}

type Action struct {
	ID              string   `yaml:"id"`                         // 选填，如果 ID 为空， 但 Name.ID 和 Target 存在，则 ID 自动填充为 Name.ID-Target
	Name            *Name    `yaml:"name" validate:"required"`   // 必填
	Target          string   `yaml:"target" validate:"required"` // 必填
	AvailableStates []string `yaml:"states"`
	Time            int      `yaml:"time"`     // 选填，单位秒
	IsSystem        bool     `yaml:"isSystem"` // true表示系统发起调用，不需要显示在前端页面
	Tag             string   `yaml:"tag"`      // 例如 'k8s,virtualbox'
}
