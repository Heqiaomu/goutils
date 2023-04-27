package generate

import (
	"encoding/base64"
	"fmt"
	"github.com/Heqiaomu/goutil/uiform/yaml"
	"io/ioutil"
	"path/filepath"
)

func ModifyInfo(infoYamlFilePath string, info *yaml.Info) error {
	if info.UIFormVersion == "" {
		info.UIFormVersion = "v1"
	}

	info.ID = info.Name.ID + info.Version

	for _, action := range info.Actions {
		if action.ID == "" {
			if action.Name.ID != "" && action.Target != "" {
				action.ID = action.Name.ID + "-" + action.Target
			} else {
				return fmt.Errorf("fail to gen action ID, because action ID and name ID and target are all empty. please make your actions in info.yaml have ID or name ID and target")
			}
		}
	}

	// 将logo的路径转化为字符流
	logoPath := filepath.Join(filepath.Dir(infoYamlFilePath), info.Logo)
	logo, err := ioutil.ReadFile(logoPath)
	if err != nil {
		return fmt.Errorf("open logo failure, logoPath=%s, err=%s", logoPath, err.Error())
	}
	info.Logo = base64.StdEncoding.EncodeToString(logo)

	return nil
}
