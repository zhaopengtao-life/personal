package command

import (
    "bytes"
    "context"
    log "github.com/sirupsen/logrus"
    "os/exec"
    "syscall"
    "time"
)

func GetCommandData(command string, timeOut int64) {
    // 参数3只要在linux系统下才会生效,只写，不存在创建，存在则清空
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
    defer cancel()
    stdout, stderr, exitCode := ExecRunScriptCommand(ctx, command)
    if exitCode != 0 {
        log.Errorf("GetCommandData Error: %v", stderr)
    }
    log.Infof("GetCommandData Data: %v", stdout)
}

func ExecRunScriptCommand(ctx context.Context, args string) (stdout string, stderr string, exitCode int) {
    var outBuf, errBuf bytes.Buffer
    cmd := exec.CommandContext(ctx, "/bin/bash", "-c", args) // mac linux
    cmd.Stdout = &outBuf
    cmd.Stderr = &errBuf
    err := cmd.Run()
    stdout = outBuf.String()
    stderr = errBuf.String()
    log.Infof("ExecRunScriptCommand Accept Command:%v", cmd.Args)
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            ws := exitError.Sys().(syscall.WaitStatus)
            exitCode = ws.ExitStatus()
        } else {
            log.Infof("ExecRunScriptCommand Could not get exit code for failed program:  %v", args)
            exitCode = 1
            if stderr == "" {
                stderr = err.Error()
            }
        }
    } else {
        ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
        exitCode = ws.ExitStatus()
    }
    return
}
