package syslog

import (
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"net"
	"personal/model"
	"personal/transfer"
	"personal/utils"
	"strconv"
	"strings"
	"time"
)

// 限制goroutine数量
var limitChan = make(chan bool, 1000)
var sysLogDataChan = make(chan *model.MetricValue, 5000)
var severity = [...]string{"Emergency", "Alert", "Critical", "Error", "Warning", "Notice", "Info", "Debug"}

// UDP goroutine 实现并发读取UDP数据
func udpTdProcess(conn *net.UDPConn) {
	var timestamp int64
	// 最大读取数据大小
	data := make([]byte, 1024)
	n, clientip, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Info("udpTdProcess ReadFromUDP failed read udp msg, error: " + err.Error())
		<-limitChan
		return
	}
	message := string(data[:n])
	if message == "" {
		log.Infof("udpTdProcess end :>>>>>>>>>>>>> message: %v", message)
		<-limitChan
		return
	}
	a := strings.Index(message, "<") + 1
	b := strings.Index(message, ">")
	lev := message[a:b]
	pri, _ := strconv.Atoi(lev)
	level := pri & 7
	clientIP := clientip.IP.String()
	uuid, _ := utils.ReplaceStringByRegex(base64.RawStdEncoding.EncodeToString([]byte(clientIP)), "=", "")
	sysLogData := &model.MetricValue{
		Timestamp: timestamp,
		Endpoint:  "syslog_" + "test_name" + "_" + uuid,
		Metric:    clientIP,
		Desc:      message,
		Type:      severity[level],
		Index:     level,
		Step:      int64(pri), // 程序模块（Facility）、严重性（Severity）
	}
	log.Infof("udpTdProcess end :>>>>>>>>>>>>> data: %v", sysLogData)
	<-limitChan
	sysLogDataChan <- sysLogData
	return
}

func udpTdServer(address string) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Errorf("read from connect failed, err:" + err.Error())
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	if err != nil {
		log.Errorf("read from connect failed, err:" + err.Error())
		return
	}

	// 周期触发定时的计时器
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	var sysLogDataList = make([]*model.MetricValue, 0)
	for {
		limitChan <- true
		go udpTdProcess(conn)
		select {
		case SysLogData := <-sysLogDataChan:
			sysLogDataList = append(sysLogDataList, SysLogData)
			log.Infof("syslog 接受日志数据：%v", SysLogData)
		// 每隔10s 发送一次,数组置为空
		case <-ticker.C:
			transfer.SendToTransfer(sysLogDataList)
			sysLogDataList = make([]*model.MetricValue, 0)
		}
	}
}

func UdpSysLog() {
	address := "0.0.0.0:514"
	log.Info("SysLogMain start :>>>>>>>>>>>>>")
	udpTdServer(address)
}
