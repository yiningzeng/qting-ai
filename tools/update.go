package tools

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func CheckUpdate() (string, error) {
	_, err := os.Stat("update.tar.gz")
	if err != nil {
		return "当前已是最新版", err
	} else {
		return "存在更新", nil
	}
}

// 更新之后需要删除update.tar.gz文件
func Update() (string, error) {
	updateFile := "update.tar.gz"
	_, err := os.Stat(updateFile)
	cmd := exec.Command("/bin/bash", "-c", "tar -xzf " + updateFile)
	bytes,err := cmd.Output()
	if err != nil {
		logrus.Error(err)
		return "更新失败", err
	}
	resp := string(bytes)
	logrus.Info(resp)

	if _, err = os.Stat("powerAi.deb"); err == nil {
		// 这里需要更新下powerAi
	}

	// 更新完成后
	err = os.Remove(updateFile)
	if err != nil {
		// 删除失败

	} else {
		// 删除成功

	}
	return "更新完成", nil
}
