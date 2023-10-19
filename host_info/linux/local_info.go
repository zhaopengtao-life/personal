package linux

import (
    "github.com/denisbrodbeck/machineid"
    "github.com/shirou/gopsutil/host"
    log "github.com/sirupsen/logrus"
    "net"
)

func LocalInfo() {
    hostInfo, _ := host.Info()
    log.Infoln("Hostname:", hostInfo.Hostname)
    log.Infoln("HostID:", hostInfo.HostID)
    log.Infoln("Platform:", hostInfo.Platform)
    log.Infoln("Platform version:", hostInfo.PlatformVersion)
    log.Infoln("Kernel version:", hostInfo.KernelVersion)

    // 本地IP4
    localIp, err := ExternalIP()
    if nil != err {
        log.Errorf("getLocalIp get local Ip Error: %v", err)
    }
    log.Info("localIP:", localIp.String())

    // 获取 IP 地址和网卡的 MAC 地址
    interfaces, err := net.Interfaces()
    if err != nil {
        log.Errorf("Error: %v", err)
        return
    }

    // 遍历网络接口，找到非零 MAC 地址（通常 eth0、en0）
    for _, iface := range interfaces {
        if iface.HardwareAddr != nil && len(iface.HardwareAddr) != 0 {
            macAddr := iface.HardwareAddr.String()
            log.Infoln("macAddr:", macAddr)
            break
        }
    }

    // 获取设备的 UUID
    uuid, err := machineid.ProtectedID("myAppName")
    if err != nil {
        log.Errorf("Error: %v", err)
    } else {
        log.Infoln("UUID:", uuid)
    }
}
