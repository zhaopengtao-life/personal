package ssh

import (
	"fmt"
	"github.com/creack/pty"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

var ExitTcp = make(chan bool, 1)

func WebTerminal(sessionId uint32, charWrit chan string, charRea chan string) {
	cmd := exec.Command("sh") // 系统默认shell交互程序
	f, err := pty.Start(cmd)  // pty用于调用系统自带的虚拟终端
	if err != nil {
		log.Errorf("ServerTcp Failed to Start data Error: %v", err)
		ExitTcp <- true
		return
	}
	// 处理来自虚拟终端的消息
	go func(sessionId uint32) {
		for {
			log.Info("处理来自虚拟终端的消息sessionId：", sessionId)
			buf := make([]byte, 1024)
			read, err := f.Read(buf)
			if err != nil {
				log.Errorf("ServerTcp Failed to Read data Error: %v", err)
				ExitTcp <- true
				continue
			}
			// 执行完的结果
			fmt.Println("WebTerminal f.Read: ", string(buf[:read]))
			//lock.NewTerminalWriteMap().Set(sessionId, string(buf[:read]))
		}
	}(sessionId)
	//  写入执行脚本命令
EXIT:
	for {
		log.Info("处理执行脚本命令 sessionId：", sessionId)
		select {
		case <-ExitTcp:
			log.Infof("Exit the current function WriteTcp")
			//ExitTcp <- true
			break EXIT
		case commandTcp := <-charRea:
			log.Info("写入执行脚本命令: ", commandTcp)
			// 将消息写到虚拟终端
			_, err := f.Write([]byte(commandTcp))
			if err != nil {
				log.Errorf("ServerTcp Failed to Write data Error: %v", err)
				//ExitTcp <- true
				continue
			}
		}
	}
}

func WebTerminalBak() {
	cmd := exec.Command("sh") // 系统默认shell交互程序
	f, err := pty.Start(cmd)  // pty用于调用系统自带的虚拟终端
	if err != nil {
		log.Errorf("ServerTcp Failed to Start data Error: %v", err)
		ExitTcp <- true
		return
	}
	// 处理来自虚拟终端的消息
	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := f.Read(buf)
			if err != nil {
				log.Errorf("ServerTcp Failed to Read data Error: %v", err)
				ExitTcp <- true
				continue
			}
			// 执行完的结果
			fmt.Println("WebTerminal f.Read: ", string(buf[:read]))
			//lock.NewTerminalWriteMap().Set(sessionId, string(buf[:read]))
		}
	}()
	//  写入执行脚本命令
	for {
		command := ""
		log.Info("写入执行脚本命令: ", command)
		// 将消息写到虚拟终端
		_, err := f.Write([]byte(command))
		if err != nil {
			log.Errorf("ServerTcp Failed to Write data Error: %v", err)
			ExitTcp <- true
			continue
		}
	}
}
