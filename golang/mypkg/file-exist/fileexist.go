package fileexist

import (
	"os"
)

// /etc/kubeedge/config/cloudcore.yaml
// 判断文件是否存在

// FileIsExist check file is exist
func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
