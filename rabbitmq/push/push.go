package push

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

type CommonResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Info    string      `json:"info"`
}

func NewProduce(queueMessage *CommonResult, queueName string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("NewProduce RabbitMQ Push error: %v", err)
			return
		}
	}()

	var mqUrl string
	mqUrl = "amqp://admin:admin@47.120.43.149:5672/vhost"

	rabbitmq, err := NewRabbitMQ(mqUrl, queueName, "", "")
	if err != nil {
		logrus.Errorf("NewProduce NewRabbitMQ Client Lose Error：%v", err)
		return
	}

	defer rabbitmq.Destroy()

	var sendMsg []byte
	// 序列化
	if queueMessage.Success == false && queueMessage.Info == "" {
		sendMsg, err = json.Marshal(&CommonResult{false, nil, "未知错误"})
	} else {
		sendMsg, err = json.Marshal(queueMessage)
	}
	if err != nil {
		logrus.Errorf("NewProduce Message Json Marshal Error：%v", err)
		return
	}

	// 发布队列
	err = rabbitmq.PublishSimple(string(sendMsg)) // 订阅模式发布
	if err != nil {
		logrus.Errorf("NewProduce RabbitMQ Push Lose Message：%v ; error: %v", string(sendMsg), err)
		return
	}
	logrus.Info("NewProduce RabbitMQ Push End And Success", string(sendMsg))
	return
}

func PushScriptRunStatus(flag bool, data interface{}, info string) {
	// 开始执行脚本通知
	commonResult := &CommonResult{
		Success: flag,
		Data:    data,
		Info:    info,
	}
	for i := 0; i < 3; i++ {
		err := NewProduce(commonResult, "test_rabbit")
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			logrus.Errorf("PushScriptRunStatus Task Status Retry Send :%v, Push Data, :%v Error, : %v ", i, data, err)
			continue
		}
		break
	}
}
