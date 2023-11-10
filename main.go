package main

import (
	"embed"
	"fmt"
	"github.com/creack/pty"
	"github.com/olahol/melody"
	"net/http"
	"os/exec"
	"strconv"
)

func main() {
	// 初始化日志
	//initlog.Initlog()
	// 本机IP
	//linux.GetLocalIp()
	//tcp.OpsServer()

	//WebTerminal()
	//web_terminal.Demo()
	//web_terminal.WebTerminalDemo()

	var ifHCOutOctets_start, ifHCOutOctets_end uint64
	ifHCOutOctets_start = 172252719015828
	ifHCOutOctets_end = 172249605304388
	// 假设1: ifHCOutOctets_start 大，ifHCOutOctets_end 小     结果为正
	// 假设2: ifHCOutOctets_start 小，ifHCOutOctets_end 大     结果为负

	first := int64(ifHCOutOctets_start - ifHCOutOctets_end)
	fmt.Println("first: ", first)

	fmt.Println("000000000000", float64(int64(first))/float64(61))
	IfHCOutFloatValue := 8 * (float64(int64(first)) / float64(61))
	fmt.Println("1111111111", IfHCOutFloatValue)
}

func FloatFomatStr(receiver float64) float64 {
	IfHCInOctets := (strconv.FormatFloat(receiver, 'f', 2, 64))
	value, _ := strconv.ParseFloat(IfHCInOctets, 64)
	return value
}

//go:embed index.html node_modules/xterm/css/xterm.css node_modules/xterm/lib/xterm.js
var content embed.FS

func WebTerminal() {
	c := exec.Command("sh") // 系统默认shell交互程序
	f, err := pty.Start(c)  // pty用于调用系统自带的虚拟终端
	if err != nil {
		panic(err)
	}

	m := melody.New() // melody用于实现WebSocket功能

	go func() { // 处理来自虚拟终端的消息
		for {
			buf := make([]byte, 1024)
			read, err := f.Read(buf)
			if err != nil {
				return
			}
			fmt.Println("f.Read: ", string(buf[:read]))
			m.Broadcast(buf[:read]) // 将数据发送给网页
		}
	}()

	m.HandleMessage(func(s *melody.Session, msg []byte) { // 处理来自WebSocket的消息
		fmt.Println("m.HandleMessage: ", string(msg))
		f.Write(msg) // 将消息写到虚拟终端
	})

	http.HandleFunc("/webterminal", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r) // 访问 /webterminal 时将转交给melody处理
	})

	fs := http.FileServer(http.FS(content))
	http.Handle("/", http.StripPrefix("/", fs)) // 设置静态文件服务

	http.ListenAndServe("0.0.0.0:8080", nil) // 启动服务器，访问 http://本机(服务器)IP地址:8080/ 进行测试
}
