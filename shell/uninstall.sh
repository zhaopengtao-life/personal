#!/bin/bash
# 一键升级 agent(arm64/amd64 包版本)
# 使用 wget 安装部署 wget需要提前配置好,否则该脚本会失败
# 本脚本使用于 centos

logTime=$(date "+%Y-%m-%d %H:%M:%S")
install_package="/data/aiops/disruptor"

while getopts o:u:l:p: option
do
 # shellcheck disable=SC2220
 case "${option}"
 in
   o) filePath="$OPTARG";;
   u) install="$OPTARG";;
   l) logPath=${OPTARG};;
   p) params=${OPTARG};;
 esac
done

# Linux内核
get_linux_core(){
    linux_core="false"
    main_version=`uname -r |awk -F'.' '{print $1}'`
    minor_version=`uname -r |awk -F'.' '{print $2}'`
    if [[ ${main_version} -ge 3 && ${minor_version} -ge 1 ]];then
        linux_core="true"
    fi
    echo "$logTime 当前内核版本为：$main_version，次版本为：$minor_version"
}

# 获取当前系统为：ucpe/通用
get_ucpe_or_host(){
    ucpe_model="false"
    if [ ! -f "/usr/local/svxnetworks/7xcli_info" ]; then
        echo "$logTime 是host，不存在ucpe目录：/usr/local/svxnetworks/7xcli_info"
    else
        ucpe_model="true"
        echo "$logTime 是ucpe，存在ucpe目录：/usr/local/svxnetworks/7xcli_info"
    fi
}

# 停止进程
ucpe_kill_agent_process(){
    pid=$(ps -ef | grep '/data/aiops/disruptor' | grep -v grep | awk '{print $2}')
    echo "$logTime ucpe: disruptor process $pid"
    # 杀死进程
    if [[ $pid == "" ]];then
        echo "$logTime ucpe: 没有匹配到正在执行的 disruptor 进程"
    else
        kill -9 $pid
        echo "$logTime ucpe: disruptor kill -9 $pid"
    fi
}

# ucpe：停止系统服务
ucpe_stop_systemd_service(){
    # 是否存在系统服务
    serverPath="/lib/systemd/system/disruptor.service"
    exePath="ExecStart=/data/aiops/disruptor"
    systemInfo=$(grep $exePath $serverPath)
    echo "$logTime ucpe: 准备停止系统服务"

    if [[ ! -f "$serverPath" ]];then
        echo "$logTime ucpe: 不存在系统服务文件：$serverPath，卸载失败，自动退出"
        exit 1
    else
        if [[ "$systemInfo" != "" ]];then
            echo "$logTime ucpe: 系统服务校验路径：$exePath，存在系统服务配置文件：$systemInfo"
            # 禁用开机自启动
            systemctl disable disruptor.service
            echo "$logTime ucpe: 禁用开机自启动：systemctl disable disruptor.service"
            echo "$logTime ucpe: 停止系统服务：systemctl stop disruptor.service"
            # 停止系统服务
            systemctl stop disruptor.service
        fi
    fi
}

# 清空注册码文件
del_aiops_conf(){
    if [[ ! -f "/data/aiops/aiops.conf" ]];then
        echo "$logTime 不存在注册码文件" >> $logFile
    else
        # 删除最后一行
        sed -i '$d' "/data/aiops/aiops.conf"
        sleep 0.5s
        echo "registerCode=''" >> "/data/aiops/aiops.conf"
        echo "$logTime Agent：注册码：$registerCode，aiops.conf 文件注册码清空完成"
    fi
}

# 停止进程 > 3.0
host_kill_agent_process(){
    if [ -f $filePath"/disruptor" ]; then
        pid=$(ps -ef | grep $filePath"/disruptor" | grep -v grep | awk '{print $2}')
        echo "$logTime 内核大于3.0 $filePath/disruptor process $pid "
        # 杀死进程
        if [[ $pid == "" ]];then
            echo "$logTime 没有匹配到正在执行的 disruptor 进程"
        else
            kill -9 $pid
            echo "$logTime disruptor kill -9 $pid " >> $logFile
        fi
    fi

    if [ -f $filePath"/LinuxDisruptor" ]; then
        pid=$(ps -ef | grep $filePath"/LinuxDisruptor" | grep -v grep | awk '{print $2}')
        echo "$logTime 内核大于3.0 $filePath/LinuxDisruptor process $pid "
        # 杀死进程
        if [[ "$pid" == "" ]];then
            echo "$logTime 没有匹配到正在执行的 LinuxDisruptor 进程"
        else
            kill -9 $pid
            echo "$logTime LinuxDisruptor kill -9 $pid " >> $logFile
        fi
    fi
}

# host：停止系统服务
host_stop_systemd_service(){
    if [[ -f $filePath"/disruptor" ]]; then
        # 是否存在系统服务
        serverPath="/lib/systemd/system/disruptor.service"
        exePath="ExecStart=$filePath/disruptor"
        systemInfo=$(grep $exePath $serverPath)
        echo "$logTime ucpe: 准备停止系统服务: $filePath/disruptor"
        if [[ ! -f "$serverPath" ]];then
              echo "$logTime 不存在系统服务文件：$serverPath"
        else
            if [[ "$systemInfo" != "" ]];then
                echo "$logTime 系统服务校验路径：$exePath，存在系统服务配置文件：$systemInfo"
                # 禁用开机自启动
                systemctl disable disruptor.service
                echo "$logTime 禁用开机自启动：systemctl disable disruptor.service"
                echo "$logTime 停止系统服务：systemctl stop disruptor.service"
                # 停止系统服务
                systemctl stop disruptor.service
            fi
        fi
    fi

    if [[ -f $filePath"/LinuxDisruptor" ]]; then
        # 是否存在系统服务
        serverPath="/lib/systemd/system/LinuxDisruptor.service"
        exePath="ExecStart=$filePath/LinuxDisruptor"
        systemInfo=$(grep $exePath $serverPath)
        echo "$logTime ucpe: 准备停止系统服务: $filePath/LinuxDisruptor"
        if [[ ! -f "$serverPath" ]];then
            echo "$logTime 不存在系统服务文件：$serverPath"
        else
            if [[ "$systemInfo" != "" ]];then
                echo "$logTime 系统服务校验路径：$exePath，存在系统服务配置文件：$systemInfo"
                # 禁用开机自启动
                systemctl disable LinuxDisruptor.service
                echo "$logTime 禁用开机自启动：systemctl disable LinuxDisruptor.service"
                echo "$logTime 停止系统服务：systemctl stop LinuxDisruptor.service"
                # 停止系统服务
                systemctl stop LinuxDisruptor.service

            fi
        fi
    fi
}

# 获取对应路径的程序pid
host_core_stop_systemd_service(){
    if [[ -f $filePath"/disruptor" ]]; then
        echo "$logTime ucpe: 准备停止系统服务：$filePath/disruptor"
        serverPath="/etc/init.d/disruptor.sh"
        if [[ ! -f "$serverPath" ]]; then
            echo "$logTime 不存在系统服务文件：$serverPath"
        else
            #删除系统服务
            rm -rf $serverPath
            echo "$logTime 删除系统服务：$serverPath"
        fi
    fi
    if [[ -f $filePath"/LinuxDisruptor" ]]; then
        echo "$logTime ucpe: 准备停止系统服务:$filePath/LinuxDisruptor"
        # 是否存在系统服务
        serverPath="/etc/init.d/LinuxDisruptor.sh"
        if [[ ! -f "$serverPath" ]]; then
            echo "$logTime 不存在系统服务文件：$serverPath"
        else
            # 删除系统服务
            rm -rf $serverPath
            echo "$logTime 删除系统服务：$serverPath"
        fi
    fi
}

run(){
    # 创建日志
    logFile="$filePath/agent_uninstall.log"
    echo "$logTime 创建 Agent 停止/卸载 日志" > $logFile
    echo "$logTime 接受到传参注册码：$filePath"
    # 系统内核
    get_linux_core
    # ucpe/通用
    get_ucpe_or_host
    #ucpe
    if [[ "$ucpe_model" == "true" ]];then
        # 停止系统服务
        ucpe_stop_systemd_service
        # 停止进程
        ucpe_kill_agent_process
    fi
    # 主机
    if [[ "$ucpe_model" == "false" && "$filePath" != "" ]];then
        if [[ "$linux_core" == "true" ]];then
            # 内核 > 3,0
            host_stop_systemd_service
        else
            # 内核 < 3,0
            host_core_stop_systemd_service
        fi
        # 停止进程
        host_kill_agent_process
    else
        echo "$logTime 主机卸载，卸载路径不能为空：$filePath，自动退出"
        exit 1
    fi

    del_aiops_conf
}

run













