package ssh

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	lock "personal/utils/locak"
)

var WebTerminalAddr string

func InitTcp(sessionId int32) error {
	// 初始化链接
	WebTerminalAddr = "192.168.10.96:33332"
	address, err := net.ResolveTCPAddr("tcp4", WebTerminalAddr)
	if err != nil {
		log.Printf("InitTcp ResolveTCPAddr Failed Addr: %v  Error: %v", WebTerminalAddr, err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		log.Printf("InitTcp DialTCP Failed Addr: %v  Error: %v", WebTerminalAddr, err)
		return err
	}
	lock.TcpClientMap.Set(sessionId, conn)
	return nil
}

// Encode 编码：发送
func Encode(sessionId int32, data []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(fmt.Sprintf("webTerminal ssh  Encode Recover: %v", err))
			//ExitTcp <- true
		}
	}()
	var pkg = new(bytes.Buffer)
	// 写入消息头
	// 小端排列，排列方式从左至右。详情搜索大小端排列
	err = binary.Write(pkg, binary.BigEndian, byte(10))
	if err != nil {
		log.Errorf("Encode Type Write Error: %v", err)
		//ExitTcp <- true
		return
	}
	// 写入消息头
	err = binary.Write(pkg, binary.BigEndian, sessionId)
	if err != nil {
		log.Errorf("Encode SessionId Write Error: %v", err)
		//ExitTcp <- true
		return
	}
	// 写入消息头
	err = binary.Write(pkg, binary.BigEndian, uint32(len(data)))
	if err != nil {
		log.Errorf("Encode Length Write Error: %v", err)
		//ExitTcp <- true
		return
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.BigEndian, data)
	if err != nil {
		log.Errorf("Encode Body Write Error: %v", err)
		//ExitTcp <- true
		return
	}
	for i := 0; i < 4; i++ {
		if i == 3 {
			//ExitTcp <- true
			break
		}
		conn := lock.TcpClientMap.Get(sessionId)
		_, err = conn.Write(pkg.Bytes())
		if err != nil {
			log.Errorf("Encode Send Write Error: %v", err)
			InitTcp(sessionId)
			continue
		}
		break
	}
	return nil
}

// Decode 解码：读取
func Decode(sessionId int32) (data string, err error) {
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(fmt.Sprintf("webTerminal ssh Decode Recover: %v", err))
		}
	}()
	// 读取消息头
	header := make([]byte, 9) // 4字节消息类型 + 4字节消息长度
	conn := lock.TcpClientMap.Get(sessionId)
	_, err = conn.Read(header)
	if err != nil {
		log.Errorf("Decode reading message header Error: %v", err)
		return
	}
	types := header[0]
	//log.Info("header: ", header)
	headers := header[1:]
	// 会话ID
	messageSession := binary.BigEndian.Uint32(headers[:4])
	if types == 0 {
		return
	}
	//log.Info("messageSession: ", messageSession)
	// 数据长度
	messageLength := binary.BigEndian.Uint32(headers[4:])
	//log.Info("messageLength: ", messageLength)
	// 读取消息数据
	messageData := make([]byte, messageLength)
	err = readExactly(conn, messageData)
	if err != nil {
		log.Errorf("Decode reading messageData header Error: %v", err)
		return
	}
	log.Infof("接受到ops长连接解析后的消息,会话ID: %v, 数据长度: %v, 消息数据: %v,", messageSession, messageLength, string(messageData))
	return string(messageData), err
}

func readExactly(conn net.Conn, data []byte) error {
	toRead := len(data)
	for toRead > 0 {
		n, err := conn.Read(data[len(data)-toRead:])
		if err != nil {
			return err
		}
		toRead -= n
	}
	return nil
}
