package web_terminal

import (
    "bufio"
    "io"
    "net"
    "os"
    "strings"
    "time"

    log "github.com/sirupsen/logrus"
)

type TelnetClient struct {
    IP               string
    Port             string
    IsAuthentication bool
    UserName         string
    Password         string
}

var g_WriteChan chan string

func ClientTelnet() {
    g_WriteChan = make(chan string)
    telnetClientObj := new(TelnetClient)
    telnetClientObj.IP = "10.168.1.227"
    telnetClientObj.Port = "23"
    telnetClientObj.IsAuthentication = true
    telnetClientObj.UserName = "root"
    telnetClientObj.Password = "Xnetworks.c0M"
    go telnetClientObj.Telnet(50)

    for {
        // 接受传递进去的参数
        line := readLine()
        g_WriteChan <- line
    }
}

func readLine() string {
    line, err := bufio.NewReader(os.Stdin).ReadString('\n')
    if err != nil && err != io.EOF {
        log.Errorf("Telnet ReadLine Error: %v", err)
    }
    return strings.TrimSpace(line)
}

func (this *TelnetClient) Telnet(timeout int) (err error) {
    raddr := this.IP + ":" + this.Port
    conn, err := net.DialTimeout("tcp", raddr, time.Duration(timeout)*time.Second)
    if nil != err {
        log.Errorf("Telnet, method: net.DialTimeout, errInfo: %v", err)
        return
    }
    defer conn.Close()
    if false == this.telnetProtocolHandshake(conn) {
        log.Error("Telnet, method: this.telnetProtocolHandshake, errInfo: telnet protocol handshake failed!!!")
        return
    }
    go func() {
        for {
            data := make([]byte, 1024)
            _, err := conn.Read(data)
            if err != nil {
                // 链接超时，抛出异常
                log.Errorf("Telnet, method: for conn.Read, errInfo: %v", err)
                break
            }
            log.Infof("Telnet 传递参数输出结果: %v", string(data))
        }
    }()

    // 控制链接时长
    conn.SetReadDeadline(time.Now().Add(time.Minute * 6))

    for {
        select {
        case cmd, _ := <-g_WriteChan:
            _, err = conn.Write([]byte(cmd + "\n"))
            if nil != err {
                log.Errorf("Telnet, method: for conn.Write, errInfo: %v", err)
                return
            }
        default:
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func (this *TelnetClient) telnetProtocolHandshake(conn net.Conn) bool {
    var buf [4096]byte
    n, err := conn.Read(buf[0:])
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Read, errInfo: %v", err)
        return false
    }
    buf[1] = 252
    buf[4] = 252
    buf[7] = 252
    buf[10] = 252
    n, err = conn.Write(buf[0:n])
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Write, errInfo: %v", err)
        return false
    }
    n, err = conn.Read(buf[0:])
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Read, errInfo: %v", err)
        return false
    }
    buf[1] = 252
    buf[4] = 251
    buf[7] = 252
    buf[10] = 254
    buf[13] = 252
    log.Infof("telnetProtocolHandshake init buf: %v", (buf[0:n]))
    n, err = conn.Write(buf[0:n])
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Write, errInfo: %v", err)
        return false
    }

    n, err = conn.Read(buf[0:])
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Read, errInfo: %v", err)
        return false
    }
    _, err = conn.Write([]byte(this.UserName + "\n"))
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Write, errInfo: %v", err)
        return false
    }
    time.Sleep(time.Second * 1)
    _, err = conn.Read(buf[0:])
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Read, errInfo: %v", err)
        return false
    }
    _, err = conn.Write([]byte(this.Password + "\n"))
    if nil != err {
        log.Errorf("telnetProtocolHandshake, method: conn.Write, errInfo: %v", err)
        return false
    }
    time.Sleep(time.Second * 1)
    _, err = conn.Read(buf[0:])
    if nil != err {
        log.Println("telnetProtocolHandshake, method: conn.Read, errInfo: %v", err)
        return false
    }
    return true
}
