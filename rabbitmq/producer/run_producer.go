package producer

import (
    "fmt"
)

type TestPro struct {
    msgContent string
}

// 实现发送者
func (t *TestPro) MsgContent() string {
    return t.msgContent
}

func RunProducer() {
    msg := fmt.Sprintf("这是生产者测试任务")
    t := &TestPro{
        msg,
    }
    queueExchange := &QueueExchange{
        "test_rabbit",
        "",
        "test_rabbit",
        "direct",
    }
    mq := New(queueExchange.QuName, queueExchange.ExName, queueExchange.ExType, "")
    mq.RegisterProducer(t)
    mq.Start()
    fmt.Println("Producer：生产者发送成功")
}
