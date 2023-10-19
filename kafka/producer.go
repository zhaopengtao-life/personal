package kafka

import (
    log "github.com/sirupsen/logrus"
)

func Producer() {
    // 创建一个Kafka生产者
    producer, err := CreateProducer()
    if err != nil {
        log.Errorf("Error creating producer: %v", err)
    }
    defer producer.Close()

    // 建立topic和Partitioner
    if err := CreateTopicAndPartitioner(); err != nil {
        log.Errorf("Error creating topic and partitioner: %v", err)
    }

    for i := 1; i < 9; i++ {
        // 向Kafka发送消息
        message := "Hello, Kafka!"
        err = ProduceMessage(producer, topicName, message)
        if err != nil {
            log.Errorf("Error producing message: %v", err)
        }
        log.Infof("Message sent successfully!")
    }
}
