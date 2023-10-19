package web_terminal

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"

    log "github.com/sirupsen/logrus"
    "golang.org/x/crypto/ssh"
)

func getClientConfig(privateKeyFile, user, passwd string) *ssh.ClientConfig {
    var auth ssh.AuthMethod
    if privateKeyFile != "" {
        // true: privateKey specified
        auth = loadPrivateKeyFile(privateKeyFile)
    } else {
        auth = ssh.Password(passwd)
    }
    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            auth,
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    return config
}

func ClientSSH() {
    // SSH连接配置
    config := getClientConfig("", "root", "Xnetworks.c0M")

    // 连接SSH服务器
    conn, err := ssh.Dial("tcp", "10.168.1.227:7722", config)
    if err != nil {
        log.Errorf("ClientSSH Failed to dial Error: %v", err)
    }
    defer conn.Close()

    // 创建SSH会话
    session, err := conn.NewSession()
    if err != nil {
        log.Errorf("ClientSSH Failed to create session Error: %v", err)
    }
    defer session.Close()

    // 硬编码要执行的命令
    command := "ls -l"

    output, err := session.CombinedOutput(command)
    if err != nil {
        log.Fatalf("Failed to execute command: %v", err)
    }
    log.Infof("session.CombinedOutput: %v", string(output))

    // 等待SSH会话结束
    if err := session.Wait(); err != nil {
        log.Errorf("ClientSSH session Wait Error: %v", err)
    }
}

func ClientSSHStdin() {
    // SSH连接配置
    config := getClientConfig("", "root", "Xnetworks.c0M")

    // 连接SSH服务器
    conn, err := ssh.Dial("tcp", "10.168.1.227:7722", config)
    if err != nil {
        log.Errorf("ClientSSH Failed to dial Error: %v", err)
    }
    defer conn.Close()

    // 创建SSH会话
    session, err := conn.NewSession()
    if err != nil {
        log.Errorf("ClientSSH Failed to create session Error: %v", err)
    }
    defer session.Close()

    // 将SSH会话的标准输入输出连接到标准输入输出
    session.Stdout = os.Stdout
    session.Stderr = os.Stderr
    stdin, err := session.StdinPipe()
    if err != nil {
        log.Errorf("ClientSSH Failed to create stdin pipe Error: %v", err)
    }

    // 启动SSH会话
    err = session.Shell()
    if err != nil {
        log.Errorf("ClientSSH Failed to start shell Error: %v", err)
    }

    // 读取用户输入并发送到SSH会话
    reader := bufio.NewReader(os.Stdin)
    for {
        text, _ := reader.ReadString('\n')
        if text == "exit\n" {
            break
        }
        _, err := fmt.Fprint(stdin, text)
        if err != nil {
            log.Errorf("ClientSSH 执行输出结果 input：%v", err)
        }
    }

    // 等待SSH会话结束
    if err := session.Wait(); err != nil {
        log.Errorf("ClientSSH session Wait Error: %v", err)
    }
}

func loadPrivateKeyFile(dir string) ssh.AuthMethod {
    buffer, err := ioutil.ReadFile(dir)
    if err != nil {
        return nil
    }
    key, err := ssh.ParsePrivateKey(buffer)
    if err != nil {
        return nil
    }
    return ssh.PublicKeys(key)
}
