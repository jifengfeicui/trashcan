package initialize

import (
	"os"
	"path/filepath"
)

func CreateMkdirall() {
	_ = os.MkdirAll("./out", 0755)

}

func ChangeWorkingDir() {
	var err error
	executable, err := os.Executable()
	if err != nil {
		return
	}
	_ = os.Chdir(filepath.Dir(executable))
}
