package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"personal_work/host_info/linux"

	"github.com/creack/pty"
	"github.com/olahol/melody"
)

func main() {
	// 初始化日志
	//initlog.Initlog()
	// 本机IP
	linux.GetLocalIp()

	WebTerminal()

	//web_terminal.WebTerminal()
	//web_terminal.WebTerminalVim()
}

//go:embed testdata/test.txt
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
			//fmt.Println("f.Read: ", string(buf[:read]))
			m.Broadcast(buf[:read]) // 将数据发送给网页
		}
	}()

	m.HandleMessage(func(s *melody.Session, msg []byte) { // 处理来自WebSocket的消息
		fmt.Println("m.HandleMessage: ", string(msg))
		f.Write(msg) // 将消息写到虚拟终端
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r) // 访问 /webterminal 时将转交给melody处理
	})

	http.Handle("/test.txt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := content.Open("testdata/test.txt")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		io.Copy(w, file)
	}))

	http.ListenAndServe("0.0.0.0:8080", nil) // 启动服务器，访问 http://本机(服务器)IP地址:22333/ 进行测试
}
