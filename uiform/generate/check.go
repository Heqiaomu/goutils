package generate

import (
	"fmt"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	"github.com/go-playground/validator/v10"
	"path/filepath"
	"strings"
)

// checkInfoYaml 验证info格式是否正确
func checkInfoYaml(infoYamlFilePath string, info *yaml.Info) error {
	oldLogo := info.Logo
	info.Logo = filepath.Join(filepath.Dir(infoYamlFilePath), info.Logo)
	defer func() { info.Logo = oldLogo }()
	v := validator.New()
	verr := v.Struct(info)
	if verr != nil {
		var errs strings.Builder
		errs.WriteString("invalid form of info.yaml:\n")
		for _, e := range verr.(validator.ValidationErrors) {
			errs.WriteString(fmt.Sprintf("%v\n", e))
		}
		return fmt.Errorf(errs.String())
	}
	return nil
}
