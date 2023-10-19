package ucloud_oss

import (
    log "github.com/sirupsen/logrus"
    "github.com/ufilesdk-dev/ufile-gosdk/example/helper"
    "os"
)

func UpLoad(dirPath, uploadPath string, stderr, stdout string) error {
    log.Infof("准备上传文件: %v ", uploadPath)
    f, err := os.OpenFile(dirPath, os.O_CREATE|os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0777)
    defer f.Close()
    if err != nil {
        log.Errorf("UpLoad Open File：%v Lose，Error：%v", dirPath, err)
        return err
    }
    if stdout != "" {
        _, err = f.WriteString(stdout)
        if err != nil {
            log.Errorf("UpLoad Save File Lose，Error：%v", err)
            return err
        }
    }
    if stderr != "" {
        _, err = f.WriteString(stderr)
        if err != nil {
            log.Errorf("UpLoad Save File Lose，Error：%v", err)
            return err
        }
    }
    if _, err := os.Stat(dirPath); os.IsNotExist(err) {
        helper.GenerateFakefile(dirPath, 1024)
    }
    err = UFileClient.PutFile(dirPath, uploadPath, "")
    if err != nil {
        log.Error("UpLoad File Lose，Error：", string(UFileClient.DumpResponse(true)))
        return err
    }
    log.Infof("文件: %v , 上传成功", dirPath)
    return nil
}
