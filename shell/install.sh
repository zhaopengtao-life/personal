#!/bin/bash

old_pid=""
md5_value=""
old_md5_value=""
new_md5_value=""
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

while getopts o:m:r:t: option
do
   # shellcheck disable=SC2220
   case "${option}"
   in
       o) registerCode="$OPTARG";;
       m) amd_value="$OPTARG";;
       r) arm_value="$OPTARG";;
       t) type="$OPTARG";;
   esac
done

# 删除文件，安装包
delete_install_package(){
    /bin/rm -rf '/data/aiops/install'
    /bin/rm -rf '/tmp/7xAgent_Install.tar.gz'
}

# 删除目录，安装包
delete_mkdir_package(){
    /bin/rm -rf '/data/aiops'
    /bin/rm -rf '/tmp/7xAgent_Install.tar.gz'
}

# 删除agent安装包
del_agent_install(){
    # 安装
    if [[ "$package_type" == "install" ]];then
        # host：rm /data/aiops目录，ucpe：rm /data/aiops/install
        if [[ "$ucpe_model" == "false" ]];then
            delete_mkdir_package
        else
            delete_install_package
        fi
        return
    fi
    # 升级
    if [[ "$package_type" == "update" ]];then
        delete_install_package
    fi
}

# 退出：文件备份还原
recover_bak_file(){
    if [[ "$package_type" == "install" ]];then
        # 如果是主机直接删除/data/aiops目录
        if [[ "$ucpe_model" == "false" ]];then
            /bin/rm -rf '/data/aiops'
        else
            delete_install_package
            # 安装包
            /bin/rm -rf $install_package
        fi
        return
    fi
    if [[ "$package_type" == "update" ]];then
        delete_install_package
        /bin/rm -rf $install_package
        # 还原备份的旧包
        /bin/mv "$update_package_bak" $install_package
    fi
}

# agent：进程pid
get_agent_process(){
    agent_pid=$(ps -ef | grep '/data/aiops/disruptor' | grep -v grep | awk '{print $2}')
}

# agent：重启系统服务
systemd_stop_or_run(){
    # 停止服务
    systemctl stop disruptor.service
    echo "$logTime Local：停止系统服务" >> $logFile
    sleep 1
    # 获取路径的进程 pid
    get_agent_process
    # 没有停止成功则主动kill
    if [[ "$agent_pid" != "" ]];then
        kill -9 "$agent_pid"
        echo "$logTime Local：kill -9 进程 $agent_pid" >> $logFile
    fi
    systemctl daemon-reload
    echo "$logTime Local：系统服务reload加载" >> $logFile
    sleep 1
    systemctl start disruptor.service
    echo "$logTime Local：升级已完成：系统服务已重启" >> $logFile
}

# 退出：系统服务备份还原
recover_systemd_service(){
    if [[ "$package_type" == "install" ]];then
        # 如果是主机直接删除/data/aiops目录
        if [[ "$ucpe_model" == "false" ]];then
            /bin/rm -rf '/data/aiops'
        else
            delete_install_package
            # 安装包
            /bin/rm -rf $install_package
        fi
        # 停止服务，无进程触发，不需要再次kill
        systemctl stop disruptor.service
        sleep 1
        # 删除服务
        /bin/rm -rf $server_path
    fi

    if [[ "$package_type" == "update" ]];then
        delete_install_package
        # 安装包
        /bin/rm -rf $install_package
         # 还原备份的旧包
        /bin/mv "$update_package_bak" $install_package
        # 系统服务
        /bin/rm -rf $server_path
        # 还原备份的系统服务
        /bin/mv $server_bak_path $server_path
        systemd_stop_or_run
    fi
}

# 系统类型：ucpe/host
get_ucpe_or_host(){
    ucpe_model="false"
    if [ ! -f "/usr/local/svxnetworks/7xcli_info" ]; then
        echo "$logTime Local 服务器是：host，不存在ucpe目录" >> $log
        return
    else
        ucpe_model="true"
        echo "$logTime Local 服务器是：ucpe，存在ucpe目录" >> $log
    fi
}

# 系统架构校验：AMD64 = x86/x86-64，ARMv8-A = aarch64
check_system_type(){
    system_build=$(uname -m)
    if [[ "$system_build" == "x86_64" || "$system_build" == "x86" || "$system_build" == "aarch64" || "$system_build" == *"arm"* || "$system_build" == *"ARM"* ]];then
        echo "$logTime Local 目前支持系统架构：$system_build" >> $log
        return
    else
        del_agent_install
        echo "$logTime Local 暂不支持系统架构：$system_build" >> $log
        exit 1
    fi
}

# 获取当前待比对的md5
get_agent_md5(){
    #  安装/升级：选择正确的md5
    if [[ "$system_build" == "x86_64" || "$system_build" == "x86" ]];then
        md5_value=$amd_value
        echo "$logTime Local：安装包系统架构amd：$md5_value" >> $log
        return
    fi
    if [[ "$system_build" == "aarch64" || "$system_build" == *"arm"* || "$system_build" == *"ARM"* ]];then
        md5_value=$arm_value
        echo "$logTime Local：安装包系统架构arm：$md5_value" >> $log
    fi
}

# 执行安装包
get_agent_package(){
    #  安装/升级：选择正确的包
    if [[ "$system_build" == "x86_64" || "$system_build" == "x86" ]];then
        /bin/cp $amd_install_package $install_package
        echo "$logTime Local：安装包系统架构：$system_build，安装包路径命名：$install_package"  >> $logFile
    fi
    if [[ "$system_build" == "aarch64" || "$system_build" == *"arm"* || "$system_build" == *"ARM"* ]];then
        /bin/cp $arm_install_package $install_package
        echo "$logTime Local：安装包系统架构：$system_build，安装包路径命名：$install_package"  >> $logFile
    fi
}

# 执行：安装/升级
get_install_or_update(){
    # 安装：注册码校验，安装包校验，md5校验
    if [[ ! -f "$install_package" ]];then
        package_type="install"
        logFile="/data/agent_install.log"
        echo "$logTime Local 创建 Agent 安装日志" > $logFile
        echo "$logTime Local 接受到传参注册码：$registerCode" >> $logFile

        # 注册码校验
        if [[ "$registerCode" == "" ]];then
            echo "$logTime Local 安装：注册码不能为空，自动退出，请携带注册码进行安装" >> $logFile
            # host：rm /data/aiops，ucpe：rm /data/aiops/install
            if [[ "$ucpe_model" == "false" ]];then
                /bin/rm -rf '/data/aiops'
            else
                /bin/rm -rf '/data/aiops/install'
            fi
            exit 1
        fi
        return
    fi

    # 升级：当前md5和升级md5做比较
    old_md5_value=$(md5sum $install_package | awk '{print $1}')
    if [[ "$old_md5_value" == "$md5_value" ]];then
        echo "$logTime Local Agent 已是最新的包，无需升级，$new_md5_value"
        exit 1
    fi
    package_type="update"
    logFile="/data/agent_update.log"
    echo "$logTime Local 创建 Agent 升级日志" > $logFile
    echo "$logTime Local 接受到传参注册码：$registerCode" >> $logFile
}


# 操作系统校验：主机安装/升级
check_os_info() {
    if [ -f /etc/os-release ]; then
        source /etc/os-release
        OS_NAME=$NAME
        OS_VERSION=$VERSION
    elif [ -f /etc/lsb-release ]; then
        source /etc/lsb-release
        OS_NAME=$DISTRIB_ID
        OS_VERSION=$DISTRIB_RELEASE
    else
        OS_NAME=$(uname -s)
        OS_VERSION=$(uname -r)
    fi
    system_type="true"
    #  类型判断
    if [[ "$OS_NAME" == *"CentOS"* || "$OS_NAME" == *"Ubuntu"* || "$OS_NAME" == *"RedHat"* || "$OS_NAME" == *"Red Hat"* ]]; then
        if [[ "$OS_NAME" == *"Ubuntu"* ]]; then
            system_type="false"
        fi
        echo "$logTime LocalHOST：操作系统: $OS_NAME ，当前版本：$OS_VERSION" >> $logFile
    else
        OS_NAME=$(cat /etc/redhat-release)
        #  类型判断
        if [[ "$OS_NAME" == *"CentOS"* || "$OS_NAME" == *"Ubuntu"* || "$OS_NAME" == *"RedHat"* || "$OS_NAME" == *"Red Hat"* ]]; then
            if [[ "$OS_NAME" == *"Ubuntu"* ]]; then
                system_type="false"
            fi
            echo "$logTime LocalHOST：操作系统: $OS_NAME" >> $logFile
        else
            del_agent_install
            echo "$logTime LocalHOST：暂不支持操作系统: $OS_NAME，自动退出" >> $logFile
            exit 1
        fi
    fi
}

# SELinux校验：主机安装/升级
check_selinux_status() {
    # Ubuntu，先执行dpkg -l来判断有没有安装 SELINUX
    if [[ "$system_type" == "false" ]]; then
        selinux_value=$(dpkg -l | grep policycoreutils)
        if [[ "$selinux_value" == "" ]]; then
            echo "$logTime LocalHOST：Ubuntu 没有使用 SELinux，通常使用AppArmor作为默认的强制访问控制（MAC）系统" >> $logFile
            return
        fi
    fi
    # CentOS/RedHat 判断有没有安装 SELINUX
    selinux_value=$(rpm -qa | grep policycoreutils)
    if [[ "$selinux_value" == "" ]]; then
        echo "$logTime LocalHOST：CentOS/RedHat 没有使用 SELinux" >> $logFile
        return
    fi
    # 检查SELINUX状态，Disabled/Permissive，进行下一步
    SELINUX_STATUS=$(getenforce)
    if [[ "$SELINUX_STATUS" == "disabled" || "$SELINUX_STATUS" == "Disabled" || "$SELINUX_STATUS" == "Permissive" || "$SELINUX_STATUS" == "permissive" ]]; then
        echo "$logTime LocalHOST：SELinux 当前设置状态: $SELINUX_STATUS" >> $logFile
        return
    fi
    echo "$logTime LocalHOST：SELinux 当前设置状态: $SELINUX_STATUS，自动退出，请手动谨慎进行设置Disabled/Permissive" >> $logFile
    del_agent_install
    exit 1

}

# 云端通讯校验：主机安装/升级
check_website() {
    WEBSITE_URL="https://aiops.7x-networks.net"
    RESPONSE_CODE=$(curl -Is $WEBSITE_URL | head -n 1 | cut -d ' ' -f 2)
    if [ "$RESPONSE_CODE" == "200" ]; then
        echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，返回状态码: $RESPONSE_CODE" >> $logFile
    else
        echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，请求失败：(Response Code: $RESPONSE_CODE)，自动退出，检查网络" >> $logFile
        del_agent_install
        exit 1
    fi

#    WEBSITE_URL="https://watsons-7xops-receiver.7x-networks.net:6061"
#    RESPONSE_CODE=$(curl -Is $WEBSITE_URL | head -n 1 | cut -d ' ' -f 2)
#    if [ "$RESPONSE_CODE" == "404" ]; then
#       echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，返回状态码: $RESPONSE_CODE" >> $logFile
#    else
#       echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，请求失败：(Response Code: $RESPONSE_CODE)，自动退出，检查网络" >> $logFile
#       del_agent_install
#       exit 1
#    fi

#    WEBSITE_URL="https://watsons-7xops-mq.7x-networks.net:5671"
#    RESPONSE_CODE=$(curl -Is $WEBSITE_URL)
#    if [[ $RESPONSE_CODE == *"AMQP"* ]]; then
#        echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL 返回值: $RESPONSE_CODE" >> $logFile
#    else
#        echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，请求失败：(Response Code: $RESPONSE_CODE)，自动退出，检查网络" >> $logFile
#        del_agent_install
#        exit 1
#    fi

#    WEBSITE_URL="https://receiver-aiops.7x-networks.net:6061"
#    RESPONSE_CODE=$(curl -Is $WEBSITE_URL | head -n 1 | cut -d ' ' -f 2)
#    if [ "$RESPONSE_CODE" == "404" ]; then
#      echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，返回状态码: $RESPONSE_CODE" >> $logFile
#    else
#      echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，请求失败：Response Code: $RESPONSE_CODE)，自动退出，检查网络" >> $logFile
#      del_agent_install
#      exit 1
#    fi

#    WEBSITE_URL="https://mq-aiops.7x-networks.net:5671"
#    RESPONSE_CODE=$(curl -Is $WEBSITE_URL)
#    if [[ $RESPONSE_CODE == *"AMQP"* ]]; then
#        echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL 返回值: $RESPONSE_CODE" >> $logFile
#    else
#        echo "$logTime LocalHOST：访问域名地址: $WEBSITE_URL ，请求失败：(Response Code: $RESPONSE_CODE)，自动退出，检查网络" >> $logFile
#        del_agent_install
#        exit 1
#    fi
}

# 校验是否安装成功
check_install_success(){
    # 获取路径的进程 pid
    get_agent_process
    if [[ "$agent_pid" == "" ]];then
        recover_systemd_service
        echo "$logTime LocalAgent：安装失败，自动退出，请重新安装"
        echo "$logTime LocalAgent：安装失败，自动退出，请重新安装" >> $logFile
        exit 1
    fi
    /bin/rm -rf "/lib/systemd/system/disruptor_bak.service"
    delete_install_package
    echo "$logTime LocalAgent：安装程序进程: $agent_pid" >> $logFile
    echo "$logTime LocalAgent：安装成功"
    echo "$logTime LocalAgent：安装成功" >> $logFile
    exit 0
}

# 校验是否升级成功
check_update_success(){
    # 获取路径的进程 pid
    get_agent_process
    if [[ "$agent_pid" == "" || "$agent_pid" == "$old_pid" ]];then
        recover_systemd_service
        echo "$logTime LocalAgent：升级失败，自动退出，请重新升级"
        echo "$logTime LocalAgent：升级失败，自动退出，请重新升级" >> $logFile
        exit 1
    fi
    /bin/rm -rf "/lib/systemd/system/disruptor_bak.service"
    delete_install_package
    echo "$logTime LocalAgent：升级程序进程: $agent_pid" >> $logFile
    echo "$logTime LocalAgent：升级成功"
    echo "$logTime LocalAgent：升级成功" >> $logFile
    exit 0
}

# Linux内核：主机安装/升级
#get_linux_core(){
#    linux_core="false"
#    Main_version=`uname -r |awk -F'.' '{print $1}'`
#    Minor_version=`uname -r |awk -F'.' '{print $2}'`
#    if [[ ${Main_version} -ge 3 && ${Minor_version} -ge 1 ]];then
#        linux_core="true"
#    fi
#    echo "$logTime Local：内核版本为：$Main_version，次版本为：$Minor_version" >> $logFile
#}

# 添加用户权限：主机安装/升级

# 创建用户权限：主机安装/升级
add_user_auth(){
    # NOPASSWD后面带有冒号:，表示执行sudo时可以不需要输入密码
    echo "aiops ALL=(ALL) NOPASSWD:/data/aiops/script/7xops_*,/data/aiops/ops_script/job/*,/data/aiops/ops_script/inspection/*,/data/aiops/ops_script/executor_inspection/*,/data/aiops/ops_script/executor_job/*,/usr/local/svxnetworks/7xcli_info,/etc/init.d/disruptor.sh" >> /etc/sudoers
    echo "$logTime LocalHOST：aiops 角色控制权限：/etc/sudoers 添加完成" >> $logFile

    #  检查是否添加aiops用户权限
    SUDO_NOPASS=$(cat /etc/sudoers | grep "aiops ALL=(ALL) NOPASSWD")
    if [[ "$SUDO_NOPASS" == "" ]]; then
        echo "$logTime LocalHOST：aiops用户无密钥sudo权限添加失败，请重新添加" >> $logFile
        del_agent_install
        exit 1
    else
        echo "$logTime LocalHOST：aiops用户无密钥sudo权限添加成功" >> $logFile
    fi
}

# 创建用户：主机安装/升级
create_user(){
    passwd_flag=true
    # 检查/etc/passwd文件的权限
    if [[ ! -w "/etc/passwd" ]]; then
        passwd_flag=false
        chattr -i /etc/passwd
        echo "$logTime LocalHOST：/etc/passwd文件没有可修改权限，临时分配可修改权限" >> $logFile
    fi
    shadow_flag=true
    # 检查/etc/passwd文件的权限
    if [[ ! -w "/etc/shadow" ]]; then
        shadow_flag=false
        chattr -i /etc/shadow
        echo "$logTime LocalHOST：/etc/shadow文件没有可修改权限，临时分配可修改权限" >> $logFile
    fi

    # CentOS，RedHat：添加用户，用户组
    if [[ "$system_type" == "true" ]]; then
        useradd  -m aiops -s /sbin/nologin
        echo "$logTime LocalHOST: CentOS，RedHat 声明无登录权限" >> $logFile
    fi
    # Ubuntu：添加用户，用户组
    if [[ "$system_type" == "false" ]]; then
        group_user=$(grep aiops /etc/group)
        if [[ $group_user == "" ]];then
            adduser --disabled-password --gecos aiops aiops
            sleep 1
            echo "$logTime LocalHOST: Ubuntu 添加用户组，用户：aiops"
            usermod -s /usr/sbin/nologin aiops
            sleep 1
            echo "$logTime LocalHOST: Ubuntu 声明无登录权限" >> $logFile
        else
            user=$(grep '^aiops:' /etc/passwd)
            if [[ $user == "" ]];then
                echo "$logTime LocalHOST: Ubuntu 建立用户，用户组：aiops，已存在，需先解决用户组存在问题" >> $logFile
                del_agent_install
                exit 1
            fi
        fi
    fi
    # 关闭临时权限
    if [[ "$passwd_flag" == "false" ]]; then
        chattr +i /etc/passwd
    fi
    if [[ "$shadow_flag" == "false" ]]; then
        chattr +i /etc/shadow
    fi
    # 用户无权限则进行赋予添加
    user_auth=$(cat /etc/sudoers | grep "aiops ALL=(ALL) NOPASSWD")
    if [[ "$user_auth" == "" ]]; then
        add_user_auth
    fi
}

# 用户权限校验：主机升级
check_user_auth() {
    uname=$(id aiops)
    # 不存在用户，创建用户并添加权限
    if [[ "$uname" == "" ]]; then
        # 添加用户权限
        create_user
        # 退出当前函数
        return
    fi
    # 存在用户，不存在用户权限进行添加，存在进行校验
    line_number=$(grep -n "aiops ALL=(ALL) NOPASSWD:" /etc/sudoers | cut -d ':' -f 1)
    if [[ "$line_number" == "" ]]; then
        echo "$logTime LocalHOST：不存在aiops用户权限，进行添加" >> $logFile
        add_user_auth
        return
    else
        echo "$logTime LocalHOST：存在aiops用户权限，进行校验" >> $logFile
    fi

    # 内置脚本
    ops=$(sudo -l -U "aiops" | grep -o "/data/aiops/script/7xops_" )
    if [[ "$ops" == "" ]]; then
        sed -i "${line_number}s|$|,/data/aiops/script/7xops_*|" /etc/sudoers
        echo "$logTime LocalHOST：aiops用户：添加/data/aiops/script/7xops_*，权限" >> $logFile
    fi
    #  作业
    job=$(sudo -l -U "aiops" | grep -o "/data/aiops/ops_script/job/" )
    if [[ "$job" == "" ]]; then
        sed -i "${line_number}s|$|,/data/aiops/ops_script/job/*|" /etc/sudoers
        echo "$logTime LocalHOST：aiops用户：添加/data/aiops/ops_script/job/*，权限" >> $logFile
    fi
    #  第三方作业
    executor_job=$(sudo -l -U "aiops" | grep -o "/data/aiops/ops_script/executor_job/")
    if [[ "$executor_job" == "" ]]; then
        sed -i "${line_number}s|$|,/data/aiops/ops_script/executor_job/*|" /etc/sudoers
        echo "$logTime LocalHOST：aiops用户：添加/data/aiops/ops_script/executor_job/*，权限" >> $logFile
    fi
    #  巡检
    inspection=$(sudo -l -U "aiops" | grep -o "/data/aiops/ops_script/inspection/")
    if [[ "$inspection" == "" ]]; then
        sed -i "${line_number}s|$|,/data/aiops/ops_script/inspection/*|" /etc/sudoers
        echo "$logTime LocalHOST：aiops用户：添加/data/aiops/ops_script/inspection/*，权限" >> $logFile
    fi
    #  第三方巡检
    executor_inspection=$(sudo -l -U "aiops" | grep -o "/data/aiops/ops_script/executor_inspection/")
    if [[ "$executor_inspection" == "" ]]; then
        sed -i "${line_number}s|$|,/data/aiops/ops_script/executor_inspection/*|" /etc/sudoers
        echo "$logTime LocalHOST：aiops用户：添加/data/aiops/ops_script/executor_inspection/*，权限" >> $logFile
    fi
}

# agent初始化：安装/升级
get_init_agent(){
    if [[ "$package_type" == "install" ]];then
        # 选择正确的安装/升级包
        get_agent_package
        # 获取待安装包的md5
        new_md5_value=$(md5sum $install_package | awk '{print $1}')
        if [[ "$md5_value" != "$new_md5_value" ]];then
            recover_bak_file
            echo "$logTime Local：正确的安装/升级包获取失败，自动退出" >> $logFile
            exit 1
        fi
        return
    fi

    # 升级，删除旧的安装包，替换新的安装包
    if [[ "$package_type" == "update" ]];then
        # 获取未升级的版本号
        old_version=$($install_package Version)
        get_agent_process
        echo "$logTime Local：安装/升级前版本： $old_version, 进程：$agent_pid, MD5： $old_md5_value" >> $logFile

        /bin/mv $install_package "$update_package_bak"
        echo "$logTime Local：旧的安装包备份，重命名为：$update_package_bak" >> $logFile
        # 选择正确的安装/升级包
        get_agent_package

        # 获取升级后的版本号，md5
        new_md5_value=$(md5sum $install_package | awk '{print $1}')
        new_version=$($install_package Version)
        echo "$logTime Local：安装/升级后版本号： $new_version, MD5： $new_md5_value" >> $logFile
        if [[ "$old_version" == "$new_version" || "$old_md5_value" == "$new_md5_value" ]];then
            recover_bak_file
            echo "$logTime Local：已是最高版本无需升级 $new_version" >> $logFile
            exit 1
        fi
    fi
}

# conf文件: ucpe更新，脚本权限修改
get_ucpe_chown(){
    # 更新注册码文件
    if [[ $registerCode != "" ]];then 那边
        # 删除最后一行
        sed -i '/registerCode/d' "/data/aiops/aiops.conf"
        sleep 1
        echo "registerCode=$registerCode" >> "/data/aiops/aiops.conf"
        echo "$logTime UCPE：注册码：$registerCode，aiops.conf 文件更新完成" >> $logFile
    fi
    # 补充脚本权限
    chown aiops:aiops /data/aiops/disruptor
    # 添加syslog权限
    sudo /usr/sbin/setcap 'cap_net_bind_service=+ep' /data/aiops/disruptor
    echo "$logTime UCPE：安装/升级包添加syslog权限完成，进行权限校验" >> $logFile
    sleep 1
    inspection=$(getcap /data/aiops/disruptor)
    if [[ "$inspection" == "" ]];then
        recover_bak_file
        echo "$logTime UCPE：syslog权限添加失败，自动退出请重新安装/升级" >> $logFile
        exit 1
    fi
}

# 系统服务：ucpe启动
run_ucpe_systemctl(){
    # 安装：启动系统服务
    if [[ "$package_type" == "install" ]];then
        echo "$logTime UCPE：启动系统服务" >> $logFile
        systemctl start disruptor.service
        echo "$logTime UCPE：系统服务已启动，待校验是否安装成功" >> $logFile
        sleep 1
        check_install_success
    fi

    # 升级：重启系统服务
    if [[ "$package_type" == "update" ]];then
        echo "$logTime UCPE：重启系统服务" >> $logFile
        # 停止服务再重启
        systemd_stop_or_run
        echo "$logTime UCPE：系统服务已重启，待校验是否升级成功" >> $logFile
        sleep 1
        check_update_success
    fi
}

# conf文件: host补充/更新，脚本权限修改
get_host_chown(){
    # 配置注册码文件
    if [[ ! -f "/data/aiops/aiops.conf" ]];then
        echo "$logTime HOST：不存在注册码文件，重新进行配置" >> $logFile
        if [[ $registerCode == "" ]];then
            recover_bak_file
            echo "$logTime ：HOST：不存在注册码文件且不存在注册码，退出请补充注册码，再重新安装/升级"
            exit 1
        fi
        /bin/cp $install_path'/aiops.conf' '/data/aiops/aiops.conf'
        echo "$logTime HOST：注册码文件：$install_path/aiops.conf，拷贝至：/data/aiops/aiops.conf" >> $logFile
    fi
    # 更新注册码文件
    if [[ $registerCode != "" ]];then
        sleep 1
        # 删除包含registerCode字段这一行
        sed -i '/registerCode/d' "/data/aiops/aiops.conf"
        sleep 1
        echo "registerCode=$registerCode" >> "/data/aiops/aiops.conf"
        echo "$logTime HOST：注册码：$registerCode，aiops.conf 文件更新完成" >> $logFile
    fi
    # 补充目录权限
    chown -R aiops:aiops /data/aiops
    # 添加syslog权限
    setcap 'cap_net_bind_service=+ep' /data/aiops/disruptor
    sleep 1
    # 添加icmp权限
    echo "net.ipv4.ping_group_range = 0  2147483647" >> /etc/sysctl.conf
    sleep 1
    echo "$logTime HOST：安装/升级包进行权限补充添加完成，进行权限校验" >> $logFile
    inspection=$(getcap /data/aiops/disruptor)
    if [[ "$inspection" == "" ]];then
        recover_bak_file
        echo "$logTime HOST：syslog权限添加失败，退出请重新安装/升级" >> $logFile
        exit 1
    fi
}

# 系统服务：主机内核 > 3.0
run_host_systemctl_service(){
    # 待安装/升级替换前系统服务路径
    system_service=$install_path"/disruptor.service"
    # 是否存在系统服务
    if [ ! -f "$server_path" ]; then
        /bin/cp "$system_service" "/lib/systemd/system/"
        echo "$logTime HOST：不存在系统服务，安装系统服务" >> $logFile
        echo "$logTime HOST：将 $system_service，系统服务文件，拷贝到指定目录下：/lib/systemd/system/" >> $logFile

        systemctl enable disruptor.service
        echo "$logTime HOST：配置开机自启动" >> $logFile
        sleep 1
        systemctl start disruptor.service
        echo "$logTime HOST：系统服务已启动，待校验是否安装成功" >> $logFile
    else
        /bin/mv $server_path $server_bak_path
        sleep 1
        echo "$logTime HOST：存在系统服务，系统服务：$server_path，重命名为：$server_bak_path" >> $logFile
        /bin/cp "$system_service" "/lib/systemd/system/"
        echo "$logTime HOST：将 $system_service，系统服务文件，拷贝到指定目录下：/lib/systemd/system/" >> $logFile

        # 停止服务再重启
        systemd_stop_or_run
        echo "$logTime HOST：系统服务已重启，待校验是否升级成功" >> $logFile
    fi
}

# 系统服务：主机启动
run_host_systemctl(){
    value=$(ps -p 1 -o comm=)
    # systemctl: /lib/systemd/system/
    if [[ "$value" == "systemd" ]];then
        run_host_systemctl_service
    fi
    # init: /etc/init.d/
    if [[ "$value" == "init" ]];then
        # 待安装/升级替换前系统服务路径
        system_service=$install_path"/disruptor"
        # 待安装/升级替换后系统服务路径
        server_path="/etc/init.d/disruptor"

        if [ ! -f "$server_path" ]; then
            /bin/cp "$system_service" "/etc/init.d/"
            echo "$logTime HOST：不存在系统服务：$server_path，安装系统服务" >> $logFile
        else
            /bin/rm -rf $server_path
            echo "$logTime HOST：存在系统服务：$server_path，删除并重新安装系统服务" >> $logFile

            service disruptor -u stop
            /bin/cp "$system_service" "/etc/init.d/"
            echo "$logTime HOST：将 $system_service，系统服务文件，拷贝到指定目录下：/etc/init.d/" >> $logFile
        fi
        # 2，3，5启动级别开on
        # CentOS，RedHat
        if [[ "$system_type" == "true" ]]; then
            chkconfig --level 235 disruptor on
            echo "$logTime chkconfig：2，3，5 中启用 disruptor 服务" >> $logFile
        fi
        # Ubuntu
        if [[ "$system_type" == "false" ]]; then
            update-rc.d disruptor defaults
            echo "$logTime update-rc.d：2，3，4，5 中启用 disruptor 服务" >> $logFile
        fi
        /etc/init.d/disruptor -o "$registerCode" -u "start"
        echo "$logTime HOST：安装已完成：系统服务已启动" >> $logFile
    fi

    if [[ "$ucpe_model" == "false" ]];then
        chown -R aiops:aiops /data/aiops
    fi
}

# 主机安装/升级
run_host_shell(){
    # 操作系统
    check_os_info
    # selinux状态
    check_selinux_status
    # 云端通讯
    check_website
    # 用户权限校验
    check_user_auth
    # 安装/升级初始化
    get_init_agent
    # 脚本权限
    get_host_chown
    # 启动系统服务
    run_host_systemctl
    sleep 1
    # 安装：成功校验
    if [[ "$package_type" == "install" ]];then
        check_install_success
    fi
    # 升级：成功校验
    if [[ "$package_type" == "update" ]];then
        check_update_success
    fi
}

run(){
    log="/data/agent_check_install.log"
    echo "$logTime Local 创建 Agent 校验日志" > $log

    # 校验：是否存在/data/aiops 目录
    if [[ ! -d "/data/aiops" ]]; then
        echo "$logTime Local： 目录/data/aiops不存在，自动退出，请先创建该目录" >> $log
        exit 1
    fi
    # 校验：下载安装包是否解压成功
    if [[ ! -f "$amd_install_package" || ! -f "$arm_install_package" ]];then
        # 如果是主机直接删除/data/aiops目录
        if [[ "$ucpe_model" == "false" ]];then
            delete_install_package
        fi
        echo "$logTime Local：一键安装/升级解压失败，退出请重新安装/升级" >> $log
        exit 1
    fi
    # 校验：下载安装包是否是最新的
    md_amd=$(md5sum $amd_install_package | awk '{print $1}')
    md_arm=$(md5sum $arm_install_package | awk '{print $1}')
    if [[ "$md_amd" != "$amd_value" || "$md_arm" != "$arm_value" ]];then
        echo "$logTime Local：md_amd: $md_amd, amd_value: $amd_value，md_arm: $md_arm, arm_value: $arm_value" >> $log
        echo "$logTime Local：一键安装/升级执行失败，不是最新的下载链接，退出请重新安装/升级" >> $log
        exit 1
    fi

    # ucpe/host
    get_ucpe_or_host
    # 系统架构
    check_system_type
    # 获取当前待比对的md5
    get_agent_md5
    # 安装/升级
    get_install_or_update

    if [[ "$ucpe_model" == "true" ]];then
        echo "$logTime UCPE：准备开始安装/升级" >> $logFile
        # 安装包
        get_init_agent
        # 脚本权限
        get_ucpe_chown
        # 启动/重启 系统服务
        run_ucpe_systemctl
    else
        echo "$logTime HOST：准备开始安装/升级" >> $logFile
        # 开始主机脚本
        run_host_shell
    fi
}

run



