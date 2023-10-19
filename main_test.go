package main

import (
    log "github.com/sirupsen/logrus"
    "personal_work/calculate"
    "personal_work/command"
    "personal_work/excel_pdf/pdf"
    "personal_work/host_info/linux"
    "personal_work/host_info/windows"
    "personal_work/kafka"
    "personal_work/rabbitmq/consumer"
    "personal_work/rabbitmq/producer"
    "personal_work/rabbitmq/push"
    "personal_work/redis"
    "personal_work/syslog"
    "personal_work/ucloud_oss"
    "personal_work/utils"
    "personal_work/web_terminal"
    "testing"
)

// 计算
func TestExec(t *testing.T) {
    calculate.Exec()
}

// 本机信息
func TestLinuxInfo(t *testing.T) {
    linux.LinuxInfo()
}
func TestLocalInfo(t *testing.T) {
    linux.LocalInfo()
}
func TestWinInfo(t *testing.T) {
    windows.WinInfo()
}

// syslog数据监听，并转发
func TestSysLog(t *testing.T) {
    syslog.UdpSysLog()
}

// 执行shell命令
func TestCommand(t *testing.T) {
    command.GetCommandData("ls -a", 60)
}

// 执行shell命令
func TestDataAes(t *testing.T) {
    utils.DataAes("aiops_agent")
}
func TestMarshalAes(t *testing.T) {
    utils.MarshalAes("aiops")
}

// web 终端：test测试类无法测试，需切换到main执行验证
func TestClientTelnet(t *testing.T) {
    web_terminal.ClientTelnet()
}
func TestClientSSHStdin(t *testing.T) {
    web_terminal.ClientSSHStdin()
}
func TestWebTerminal(t *testing.T) {
    web_terminal.WebTerminal()
}

// 数组/切片
func TestSliceSort(t *testing.T) {
    utils.SliceSort()
}
func TestSliceCopy(t *testing.T) {
    utils.SliceCopy()
}
func TestGoSliceSort(t *testing.T) {
    utils.GoSliceSort()
}

// Kafka 生产/消费
func TestProducer(t *testing.T) {
    //kafka.DeleteTopic("zpt_test_topic")
    kafka.Producer()
}
func TestConsume(t *testing.T) {
    redis.Init()
    kafka.KafkaConsumer()
}

// rabbitmq 生产/消费
func TestRunProducer(t *testing.T) {
    producer.RunProducer()
}
func TestPush(t *testing.T) {
    push.PushScriptRunStatus(true, "123213213213", "111111111")
}
func TestRunConsumer(t *testing.T) {
    consumer.RunConsumer()
}

// oss 上传/下载
func TestDownLoad(t *testing.T) {
    ucloud_oss.Init()
    url := "下载链接"
    dirPath := "上传路径"
    err := ucloud_oss.DownLoad(url, dirPath)
    if err != nil {
        log.Errorf("CsvDownLoad DownLoad url: %v Error: %v", url, err)
    }
}
func TestUpload(t *testing.T) {
    ucloud_oss.Init()
    dirPath := "本地文件"
    uploadPath := "上传路径"
    err := ucloud_oss.UpLoad(dirPath, uploadPath, "error", "123")
    if err != nil {
        log.Errorf("UpLoad Error: %v", err)
    }
}

// pdf 读/写
func TestGeneratePdf(t *testing.T) {
    pdf.GeneratePdf()
}
