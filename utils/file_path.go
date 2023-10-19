package utils

import (
    "github.com/sirupsen/logrus"
    "os"
    "path"
    "path/filepath"
    "runtime"
    "strings"
)

func GetFilePath() (pathDir string) {
    systemType := runtime.GOOS
    if systemType == "linux" {
        dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            logrus.Info("GetFromGlobal filepath Error", err)
        }
        pathDir = strings.Replace(dir, "\\", "/", -1)
    } else {
        str, _ := os.Getwd()
        pathDir = str
    }
    return
}

// 最终方案-全兼容
func GetCurrentAbPath() string {
    dir := getCurrentAbPathByExecutable()
    if strings.Contains(dir, getTmpDir()) {
        return getCurrentAbPathByCaller()
    }
    return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() (res string) {
    dir := os.Getenv("TEMP")
    if dir == "" {
        dir = os.Getenv("TMP")
        res, _ = filepath.EvalSymlinks(dir)
    }
    return res
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() (res string) {
    exePath, err := os.Executable()
    if err != nil {
        res, _ = filepath.EvalSymlinks(filepath.Dir(exePath))
    }
    return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() (res string) {
    _, filename, _, ok := runtime.Caller(0)
    if ok {
        res = path.Dir(filename)
    }
    return res
}
