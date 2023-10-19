package model

type SysLogData struct {
    Timestamp int64  `json:"ts"`
    Endpoint  string `json:"endpoint"`  // GetBaseEncode(ip)
    Metric    string `json:"metric"`    // clientIP
    Message   string `json:"message"`   // 接受到的数据，长度截取
    LevelInfo string `json:"levelInfo"` // 日志级别
    Level     string `json:"level"`     // 日志级别
    Pri       string `json:"pri"`       // 程序模块（Facility）、严重性（Severity）
}
