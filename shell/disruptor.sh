#!/bin/bash

logTime=$(date "+%Y-%m-%d %H:%M:%S")
server_path="/etc/init.d/disruptor"

while getopts o:u:l:p: option
do
 # shellcheck disable=SC2220
 case "${option}"
 in
   o) registerCode="$OPTARG";;
   u) logPath="$OPTARG";;
   l) installPath=${OPTARG};;
   p) params=${OPTARG};;
 esac
done

# 创建日志
create_file(){
    logFile="/data/agent_init.log"
    echo "$logTime 创建 Agent 安装日志" > $logFile
    echo "$logTime Linux 使用init管理系统服务，接受到的注册码：$registerCode" >> $logFile
}

# 停止进程
kill_process(){
    /bin/rm -rf $server_path
    items=$(ps -ef | grep '/data/aiops/disruptor' | grep -v grep | awk '{print $2}')
    for item in $items;
    do
        kill -9 "$item"
        echo "$logTime 杀死获取到的 Agent 程序进程Pid: $item" >> $logFile
    done
}

# 启动服务
run_agent(){
    # 启动服务，不存在升级，不需要带注册码，存在则安装需要带注册码
    if [ "$registerCode" == "" ]; then
      /data/aiops/disruptor &
    else
      /data/aiops/disruptor -key "$registerCode" &
    fi
    echo "$logTime Linux 使用init管理系统服务已启动" >> $logFile
}

case "$1" in
  start)
    # 创建日志
    create_file
    # 启动服务
    run_agent
    echo "Starting myservice"
    ;;
  stop)
    # 创建日志
    create_file
    # 停止进程
    kill_process
    echo "Stopping myservice"
    ;;
  restart)
    # 创建日志
    create_file
    # 停止进程
    kill_process
    sleep 1
    # 启动服务
    run_agent
    ;;
  *)
    code=$1
    code1=$2
    code2=$3
    echo "code: $code，code1: $code1，code2: $code2 验证数据接受"
    echo "Usage: $0 {start|stop|restart}"
    exit 1
    ;;
esac

exit 0