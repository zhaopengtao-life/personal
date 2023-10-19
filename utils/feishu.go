package utils

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "runtime"
    "runtime/debug"
    "strconv"
    "strings"
)

type FSResponse struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

type TextItem struct {
    Tag  string `json:"tag"`
    Text string `json:"text"`
}

type NotifyData struct {
    MsgType string `json:"msg_type"`
    Content `json:"content"`
}

type Content struct {
    Post PostData `json:"post"`
}

type PostData struct {
    ZhCn zhCn `json:"zh_cn"`
}

type zhCn struct {
    Title   string       `json:"title"`
    Content [][]TextItem `json:"content"`
}

const WebhookUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/66e2ca1b-3472-4012-b601-6bf0ebf58ef3" // 飞书报警地址

func TmsFeiShuNotify(title string, errParam error, message interface{}) error {
    // 参数
    msgStr, _ := json.Marshal(message)
    pc, file, line, _ := runtime.Caller(1)
    f := runtime.FuncForPC(pc)

    var lineStr = strconv.Itoa(line)
    var errMsg = "状态"
    var msgItem = []TextItem{
        {
            Tag:  "text",
            Text: "环境：" + os.Getenv("RUN_TIME") + "\n",
        },
        {
            Tag:  "text",
            Text: "位置：" + file + ":" + lineStr + "\n",
        },
        {
            Tag:  "text",
            Text: "函数：" + f.Name() + "\n",
        },
        {
            Tag:  "text",
            Text: "类型：" + errMsg + "\n",
        },
        {
            Tag:  "text",
            Text: "调用链：" + string(debug.Stack()) + "\n",
        },
        {
            Tag:  "text",
            Text: "参数：" + string(msgStr) + "\n",
        },
    }

    var msg = NotifyData{
        MsgType: "post",
        Content: Content{PostData{ZhCn: zhCn{
            Title: title,
            Content: [][]TextItem{
                msgItem,
            },
        }}},
    }

    params, err := json.Marshal(msg)

    req, err := http.NewRequest("POST", WebhookUrl, strings.NewReader(string(params)))
    if err != nil {
        return err
    }
    req.Header.Add("Content-Type", "application/json")
    result, err := new(http.Client).Do(req)
    if err != nil {
        return err
    }

    defer result.Body.Close()
    body, err := ioutil.ReadAll(result.Body)
    if err != nil {
        return err
    }

    var response = FSResponse{}
    err = json.Unmarshal(body, &response)
    if err != nil {
        return err
    }
    if response.Code != 0 {
        return errors.New(fmt.Sprintf("GeoNotifyRobot StatusCode != 0; StatusCode ==%v ", response.Code))
    }

    return errors.New("200")
}
