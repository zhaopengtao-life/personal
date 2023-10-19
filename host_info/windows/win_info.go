package windows

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "os/exec"
    "personal_work/host_info/linux"
    "strings"
)

//GetCPUID 获取cpuid
func GetWinCpuId() string {
    var cpuid string
    cmd := exec.Command("wmic", "cpu", "get", "processorid")
    b, e := cmd.CombinedOutput()

    if e == nil {
        cpuid = string(b)
        cpuid = cpuid[12 : len(cpuid)-2]
        cpuid = strings.ReplaceAll(cpuid, "\n", "")
    } else {
        fmt.Printf("%v", e)
    }
    return cpuid
}

//系统版本
//func GetSystemVersion() string {
//	version, err := syscall.GetVersion()
//	if err != nil {
//		return ""
//	}
//	return fmt.Sprintf("%d.%d (%d)", byte(version), uint8(version>>8), version>>16)
//}

type MotherboardInfo struct {
    ProductSerial string // 序列号
    ProductUuid   string // uuid
    ProductName   string // 名称
}

//主板信息
func GetMotherboardInfo() *MotherboardInfo {
    command := exec.Command("wmic", "baseboard", "get", "serialnumber")
    output, err := command.CombinedOutput()
    if err != nil {
        log.Errorf("output err: %v", err)
    }
    motherboardInfo := &MotherboardInfo{}
    if len(strings.Split(string(output), "\n")) > 2 {
        motherboardInfo.ProductSerial = strings.Split(string(output), "\n")[1]
    }
    commanda := exec.Command("wmic", "csproduct", "list", "full")
    outputa, err := commanda.CombinedOutput()
    if err != nil {
        log.Errorf("output err: %v", err)
    }
    datas := strings.Split(string(outputa), "\n")
    for _, v := range datas {
        if strings.Contains(v, "Name=") {
            motherboardInfo.ProductName = strings.Split(v, "Name=")[1]
        }
        if strings.Contains(v, "UUID=") {
            motherboardInfo.ProductUuid = strings.Split(v, "UUID=")[1]
        }
    }
    return motherboardInfo
}

func WinInfo() []string {
    var metricItems = make([]string, 0)
    // cpuid
    cpuData := strings.ReplaceAll(GetWinCpuId(), " ", "")
    metricItems = append(metricItems, strings.ReplaceAll(cpuData, "\r", ""))
    // ip（ip4)
    metricItems = append(metricItems, linux.GetLocalIp())
    // mac(集合拼接)
    metricItems = append(metricItems, strings.Join(linux.GetMacAddrs(), ","))
    // 机器名称（系统类型+系统架构+名称）
    metricItems = append(metricItems, linux.GetLocalHostName())
    datas := GetMotherboardInfo()
    // 名称
    if datas != nil {
        metricItems = append(metricItems, datas.ProductName)
        // 序列号
        metricItems = append(metricItems, datas.ProductSerial)
        // UUID
        metricItems = append(metricItems, datas.ProductUuid)
    }
    // 版本
    //metricItems = append(metricItems, GetSystemVersion())
    return metricItems
}
