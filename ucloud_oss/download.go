package ucloud_oss

import (
    "os"
    "time"

    log "github.com/sirupsen/logrus"
)

func DownLoad(url, scriptName string) error {
    log.Println("准备下载文件: ", url)
    reqUrl := UFileClient.GetPrivateURL(url, 60*time.Minute)
    err := UFileClient.Download(reqUrl)
    if err != nil {
        log.Fatalln(string(UFileClient.DumpResponse(true)))
        return err
    }
    // 保存到本地
    f, err := os.OpenFile(scriptName, os.O_CREATE|os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0777)
    defer f.Close()
    if err != nil {
        log.Errorf("创建文件失败，错误信息为：%v", err)
        return err
    }
    data := string(UFileClient.LastResponseBody)
    _, err = f.WriteString(data)
    if err != nil {
        log.Errorf("下载数据保存到本地文件失败：%v", err)
        return err
    }
    log.Infof("下载数据保存到本地文件成功：%v", scriptName)
    return nil
}
