package web_terminal

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage {
			// Handle the received message here
			// For simplicity, we'll just print it
			println(string(p))
		}
	}
}
func WebTerminal() {
	http.HandleFunc("/ws", handleConnection)
	http.ListenAndServe(":8080", nil)
}

func WebTerminalVim() {
	// 启动一个虚拟终端会话（在这个示例中使用bash）
	cmd := exec.Command("bash")

	// 创建命令的标准输入、标准输出和标准错误管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("无法获取标准输入管道: %v\n", err)
		return
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("无法获取标准输出管道: %v\n", err)
		return
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("无法获取标准错误管道: %v\n", err)
		return
	}
	defer stderr.Close()

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("无法启动命令: %v\n", err)
		return
	}

	// 创建一个Scanner来读取终端的输出
	scanner := bufio.NewScanner(stdout)

	// 启动一个goroutine来读取并打印终端的输出
	go func() {
		for scanner.Scan() {
			fmt.Println("终端输出:", scanner.Text())
		}
	}()

	// 从标准输入读取用户输入并发送到终端
	fmt.Println("请输入命令，输入exit退出：")
	for {
		reader := bufio.NewReader(os.Stdin)
		userInput, _ := reader.ReadString('\n')
		_, err := stdin.Write([]byte(userInput))
		if err != nil {
			fmt.Printf("无法发送命令到终端: %v\n", err)
			break
		}
	}

	// 等待命令执行完毕
	if err := cmd.Wait(); err != nil {
		fmt.Printf("命令执行出错: %v\n", err)
	}
}
