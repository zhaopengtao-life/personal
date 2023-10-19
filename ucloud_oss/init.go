package ucloud_oss

import (
    "time"

    log "github.com/sirupsen/logrus"
    ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

var (
    UFileClient *ufsdk.UFileRequest
)

func Init() {
    config := &ufsdk.Config{
        PublicKey:       "4eU4QThQx6JGbp8MHU3u9wEBesuAZN7hE",
        PrivateKey:      "I9aDEX6RXG8E9ggL2zAWMv1M6gxrE4PovCVIExLn43DG",
        BucketName:      "dev-ops-7x-networks",
        FileHost:        "cn-sh2.ufileos.com",
        VerifyUploadMD5: false,
    }
    client, err := ufsdk.NewFileRequest(config, nil)
    if err != nil {
        log.Errorf("oss文件链接初始化失败：%v", err)
        time.Sleep(1 * time.Second)
    }
    UFileClient = client
    return
}
