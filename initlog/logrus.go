package initlog

import (
    "github.com/natefinch/lumberjack"
    "github.com/sirupsen/logrus"
    "io"
    "io/ioutil"
    "os"
    "personal_work/utils"
    "strings"
    "time"
)

func Initlog() {
    defer func() {
        if err := recover(); err != nil {
            logrus.Errorf("Initlog Start Init Panic: %v", err)
        }
    }()
    //设置日志级别为debug
    logrus.SetLevel(logrus.DebugLevel)
    logSize := 10
    pathDir := utils.GetFilePath()
    // 删除多余日志，只保留2次日志
    DeleteLog(pathDir + "/")
    lumberJackLogger := &lumberjack.Logger{
        Filename:   "./agent-" + time.Now().Format("2006-01-02") + ".log", //日志文件位置
        MaxSize:    logSize,                                               // 单文件最大容量,单位是MB
        MaxBackups: 1,                                                     // 最大保留过期文件个数
        //MaxAge:   1,     // 保留过期文件的最大时间间隔,单位是天
        Compress: false, // 是否需要压缩滚动日志, 使用的 gzip 压缩
    }
    mw := io.MultiWriter(os.Stdout, lumberJackLogger)
    logrus.SetOutput(mw)
    logrus.SetFormatter(new(LogFormatter))
}

//
// DeleteLog
//  @Description: 删除多余日志只保留2个
//  @param logPath
//
func DeleteLog(logPath string) {
    //获取文件或目录相关信息
    fileInfoList, err := ioutil.ReadDir(logPath)
    if err != nil {
        logrus.Errorf("DeleteLog ioutil.ReadDir: %v", err)
        return
    }
    if len(fileInfoList) < 3 {
        return
    }
    // 保留最近两个文件
    for i := 0; i < len(fileInfoList)-2; i++ {
        // 不是agent的日志不进行删除
        if !strings.Contains(fileInfoList[i].Name(), "agent-") {
            continue
        }
        // 删除文件
        err := os.Remove(logPath + fileInfoList[i].Name())
        if err != nil {
            logrus.Errorf("DeleteLog os.Remove: %v", err)
        }
    }
}
