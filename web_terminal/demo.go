package web_terminal

import (
	"github.com/creack/pty"
	log "github.com/sirupsen/logrus"
	"os/exec"
	lock "personal/utils/locak"
	"personal/web_terminal/ssh"
)

func Demo() {
	var sessionId int32 = -1070720063
	lock.TcpClientMap = lock.NewTcpClientMap()
	cmd := exec.Command("sh") // 系统默认shell交互程序
	f, err := pty.Start(cmd)  // pty用于调用系统自带的虚拟终端
	if err != nil {
		log.Errorf("ServerTcp Failed to Start data Error: %v", err)
		return
	}
	ssh.InitTcp(sessionId)
	// 数据发送
	err = ssh.Encode(sessionId, []byte("Start Client"))
	if err != nil {
		log.Errorf("准备运维链接窗口Error: %v", err)
	}

	// 处理来自虚拟终端的消息
	go func() {
		for {
			log.Info("开始执行虚拟终端接受到消息sessionId：", sessionId)
			buf := make([]byte, 1024)
			read, err := f.Read(buf)
			if err != nil {
				log.Errorf("ServerTcp Failed to Read data Error: %v", err)
				continue
			}
			// 执行完的结果
			log.Info("虚拟终端结果输出: ", string(buf[:read]))
			// 数据发送
			err = ssh.Encode(sessionId, buf[:read])
			if err != nil {
				log.Errorf("准备运维链接窗口Error: %v", err)
			}
		}
	}()
	for {
		// 数据解析
		data, err := ssh.Decode(sessionId)
		if err != nil {
			log.Errorf("ssh.Decode Error: %v", err)
		}
		if data == "" {
			continue
		}
		// 解析消息写到虚拟终端
		_, err = f.Write([]byte(data))
		if err != nil {
			log.Errorf("ServerTcp Failed to Write data Error: %v", err)
			continue
		}
	}
}

//func WebTerminalDemo() {
//    // 主程序阻塞
//    lock.TcpClientMap = lock.NewTcpClientMap()
//    lock.TerminalReadMap = lock.NewTerminalReadMap()
//    lock.TerminalWriteMap = lock.NewTerminalWriteMap()
//    var sessionId uint32 = 988396339
//    ssh.InitTcp(sessionId)
//    // 数据发送
//    err := ssh.Encode(sessionId, "开始准备运维链接窗口")
//    if err != nil {
//        log.Errorf("准备运维链接窗口Error: %v", err)
//    }
//
//    charRead := make(chan string, 200)
//    lock.NewTerminalReadMap().Set(sessionId, charRead)
//    lock.NewTerminalReadMap().Get(sessionId)
//    charWrite := make(chan string, 200)
//    lock.NewTerminalWriteMap().Set(sessionId, charWrite)
//    lock.NewTerminalWriteMap().Get(sessionId)
//
//    // 读取虚拟终端，发送ops
//    go func(sessionId uint32) {
//        for {
//            resultTcp := lock.NewTerminalWriteMap().Get(sessionId)
//            log.Info("读取虚拟终端，发送ops resultTcp: ", resultTcp)
//            select {
//            case result := <-resultTcp:
//                err := ssh.Encode(sessionId, result)
//                if err != nil {
//                    continue
//                }
//            }
//        }
//    }(sessionId)
//    // 接受的数据，并输入到虚拟终端
//    go func(sessionId uint32) {
//        for {
//            log.Infof("进入解析")
//            data, err := ssh.Decode(sessionId)
//            if err != nil {
//                log.Errorf("ssh.Decode Error: %v", err)
//                continue
//            }
//            log.Infof("解析后数据data: %v", data)
//        }
//    }(sessionId)
//    time.Sleep(3 * time.Second)
//    //ssh.WebTerminal(sessionId)
//}
