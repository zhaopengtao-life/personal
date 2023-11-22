#!/bin/bash

old_pid=""
old_value=""
new_value=""
system_type=""
old_version=""
new_version=""
package_type=""
today=$(date "+%Y-%m-%d")
logTime=$(date "+%Y-%m-%d %H:%M:%S")
install_package="/data/aiops/disruptor"
install_path="/data/aiops/install/service"
update_package_bak="/data/aiops/disruptor_"$today
server_path="/lib/systemd/system/disruptor.service"
server_bak_path="/lib/systemd/system/disruptor_bak.service"
amd_install_package="/data/aiops/install/service/install_agent_amd64"
arm_install_package="/data/aiops/install/service/install_agent_arm64"

while getopts o:u:l:p: option
do
   # shellcheck disable=SC2220
   case "${option}"
   in
       o) registerCode="$OPTARG";;
       u) install="$OPTARG";;
       l) logPath="$OPTARG";;
       p) params="$OPTARG";;
   esac
done

# 主机：安装/升级，添加用户权限
add_user_auth(){
    # NOPASSWD后面带有冒号:，表示执行sudo时可以不需要输入密码
    echo "aiops ALL=(ALL) NOPASSWD:/data/aiops/script/7xops_*,/data/aiops/ops_script/job/*,/data/aiops/ops_script/inspection/*,/data/aiops/ops_script/executor_inspection/*,/data/aiops/ops_script/executor_job/*,/usr/local/svxnetworks/7xcli_info,/etc/init.d/disruptor.sh" >> /etc/sudoers
    echo "$logTime HOST：aiops 角色控制权限：/etc/sudoers 添加完成"

    #  检查是否添加aiops用户权限
    SUDO_NOPASS=$(cat /etc/sudoers | grep "aiops ALL=(ALL) NOPASSWD")
    if [[ "$SUDO_NOPASS" == "" ]]; then
        echo "$logTime HOST：aiops用户：无密钥sudo权限添加失败，请重新添加"
        exit 1
    else
        echo "$logTime HOST：aiops用户：无密钥sudo权限添加成功"
    fi
}

# 主机：安装/升级，创建用户
create_user(){
    # 检查/etc/passwd文件的权限
    if [ -w "/etc/passwd" ]; then
        echo "$logTime HOST：/etc/passwd文件具有可修改权限"
    else
        chattr -i /etc/passwd
        echo "$logTime HOST：/etc/passwd文件没有可修改权限，临时分配可修改权限"
    fi

    # 检查/etc/passwd文件的权限
    if [ -w "/etc/shadow" ]; then
        echo "$logTime HOST：/etc/shadow文件具有可修改权限"
    else
        chattr -i /etc/shadow
        echo "$logTime HOST：/etc/shadow文件没有可修改权限，临时分配可修改权限"
    fi

    group_user=$(grep aiops /etc/group)
    echo "$logTime 当前用户组：$group_user"
    if [[ $group_user == "" ]];then
        adduser --disabled-password --gecos aiops aiops
        echo "$logTime 添加用户组aiops"
    fi
    sleep 1
    usermod -s /usr/sbin/nologin aiops
    sleep 1
    echo "$logTime Ubuntu 存在家目录，声明无登录权限"
    chattr +i /etc/passwd
    chattr +i /etc/shadow

    SUDO_NOPASS=$(cat /etc/sudoers | grep "aiops ALL=(ALL) NOPASSWD")
    if [[ "$SUDO_NOPASS" == "" ]]; then
        add_user_auth
    fi
    chown -R aiops:aiops /data/aiops
}

create_user