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
		PublicKey:       "*",
		PrivateKey:      "*",
		BucketName:      "dev-*-*-*",
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
