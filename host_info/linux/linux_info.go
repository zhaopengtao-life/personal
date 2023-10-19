package linux

import (
    "io/ioutil"
    "net"
    "os"
    "os/exec"
    "runtime"
    "strings"

    log "github.com/sirupsen/logrus"
)

type DMI struct {
    ProductName    string
    ProductSerial  string
    ProductUuid    string
    ProductVersion string
}

var SysinfoDmi DMI

func init() {
    SysinfoDmi.ProductName = getStringFromFile("/sys/class/dmi/id/product_name")
    SysinfoDmi.ProductSerial = getStringFromFile("/sys/class/dmi/id/product_serial")
    SysinfoDmi.ProductUuid = getStringFromFile("/sys/class/dmi/id/product_uuid")
    SysinfoDmi.ProductVersion = getStringFromFile("/sys/class/dmi/id/product_version")
}

func getStringFromFile(path string) string {
    //读取文件全部内容
    b, err := ioutil.ReadFile(path)
    if err != nil {
        return ""
    }
    return strings.Replace(string(b), "\n", "", -1)
}

func GetLocalIp() string {
    localIp, err := ExternalIP()
    if nil != err {
        log.Errorf("getLocalIp get local Ip Error: %v", err)
    }
    log.Info("配置文件localIP为空，获取到的本地IP：", localIp.String())
    return localIp.String()
}

func GetLocalHostName() string {
    var err error
    var systemType, system, hostname, agentName string
    systemType = runtime.GOOS
    system = runtime.GOARCH
    hostname, err = os.Hostname()
    if err != nil {
        log.Errorf(" getLocalHostName get Hostname Error: %v", err)
    }
    agentName = systemType + "-" + system + "-" + hostname
    log.Info("配置文件agentName为空，获取到的本地agentName：", agentName)
    return agentName
}

func GetMacAddrs() (macAddrs []string) {
    netInterfaces, err := net.Interfaces()
    if err != nil {
        log.Errorf("fail to get net interfaces: %v", err)
        return macAddrs
    }

    for _, netInterface := range netInterfaces {
        macAddr := netInterface.HardwareAddr.String()
        if len(macAddr) == 0 {
            continue
        }
        macAddrs = append(macAddrs, macAddr)
    }
    return macAddrs
}

type CommPack struct {
    bytesData []byte
}

func getCpuId() (CommPack, error) {
    cmd := exec.Command("/bin/sh", "-c", `dmidecode -t 4 | grep ID`)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Errorf("StdoutPipe: %v", err.Error())
        return CommPack{nil}, err
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        log.Errorf("StderrPipe: %v", err.Error())
        return CommPack{nil}, err
    }

    if err := cmd.Start(); err != nil {
        log.Errorf("Start: %v", err.Error())
        return CommPack{nil}, err
    }

    bytesErr, err := ioutil.ReadAll(stderr)
    if err != nil {
        log.Errorf("ReadAll stderr: %v", err.Error())
        return CommPack{nil}, err
    }

    if len(bytesErr) != 0 {
        log.Errorf("stderr is not nil: %v", bytesErr)
        return CommPack{nil}, err
    }

    bytes, err := ioutil.ReadAll(stdout)
    if err != nil {
        log.Errorf("ReadAll stdout: %v", err.Error())
        return CommPack{nil}, err
    }

    if err := cmd.Wait(); err != nil {
        log.Errorf("Wait: %v", err.Error())
        return CommPack{nil}, err
    }

    return CommPack{bytes}, err
}

func ExternalIP() (net.IP, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        log.Errorf(" externalIP get Interfaces Error: %v", err)
        return nil, err
    }
    for _, inter := range interfaces {
        if inter.Flags&net.FlagUp == 0 {
            continue // interface down
        }
        if inter.Flags&net.FlagLoopback != 0 {
            continue // loopback interface
        }
        addrs, err := inter.Addrs()
        if err != nil {
            return nil, err
        }
        for _, addr := range addrs {
            ip := GetIpFromAddr(addr)
            if ip == nil {
                continue
            }
            return ip, nil
        }
    }
    return nil, nil
}

func GetIpFromAddr(addr net.Addr) net.IP {
    var ip net.IP
    switch v := addr.(type) {
    case *net.IPNet:
        ip = v.IP
    case *net.IPAddr:
        ip = v.IP
    }
    if ip == nil || ip.IsLoopback() {
        return nil
    }
    ip = ip.To4()
    if ip == nil {
        return nil // not an ipv4 address
    }
    return ip
}

func LinuxInfo() []string {
    var cpuData string
    var metricItems = make([]string, 0)
    cpuId, _ := getCpuId()
    cpuDatas := strings.Split(string(cpuId.bytesData), "ID:")
    if len(cpuDatas) > 1 {
        cpuida := strings.ReplaceAll(cpuDatas[1], " ", "")
        cpuData = strings.ReplaceAll(cpuida, "\n", "")
    }
    metricItems = append(metricItems, cpuData)
    // ip（ip4)
    metricItems = append(metricItems, GetLocalIp())
    // mac(集合拼接)
    metricItems = append(metricItems, strings.Join(GetMacAddrs(), ","))
    // 机器名称（系统类型+系统架构+名称）
    metricItems = append(metricItems, GetLocalHostName())
    // 名称
    metricItems = append(metricItems, SysinfoDmi.ProductName)
    // 序列号
    metricItems = append(metricItems, SysinfoDmi.ProductSerial)
    // UUID
    metricItems = append(metricItems, SysinfoDmi.ProductUuid)
    // 版本
    metricItems = append(metricItems, SysinfoDmi.ProductVersion)
    return metricItems
}
