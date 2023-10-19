package consumer

import (
    "fmt"
    log "github.com/sirupsen/logrus"
)

type TestPro struct {
    msgContent string
}

// 实现接收者
func (t *TestPro) Consumer(dataByte []byte) error {
    log.Info(string(dataByte))
    return nil
}

func RunConsumer() {
    msg := fmt.Sprintf("这是消费者测试任务")
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
    mq.listenReceiver(t)
    mq.Start()
    log.Info("Consumer：消费者消费成功")
}
