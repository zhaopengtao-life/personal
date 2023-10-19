package kafka

import (
    "fmt"
    "strings"

    "github.com/IBM/sarama"
    log "github.com/sirupsen/logrus"
)

const (
    brokerAddress = "172.31.12.5:9092"
    topicName     = "zpt_test_topic"
)

//
// CreateProducer
//  @Description: 建立生产者链接
//  @return sarama.SyncProducer
//  @return error
//
func CreateProducer() (sarama.SyncProducer, error) {
    config := sarama.NewConfig()
    // sasl认证
    config.Net.SASL.Enable = true
    config.Net.SASL.User = "admin"
    config.Net.SASL.Password = "admin"

    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true

    // 创建生产者
    producer, err := sarama.NewSyncProducer([]string{brokerAddress}, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create producer: %v", err)
    }

    return producer, nil
}

//
// CreateTopicAndPartitioner
//  @Description: 创建topic分区
//  @return error
//
func CreateTopicAndPartitioner() error {
    config := sarama.NewConfig()
    // sasl认证
    config.Net.SASL.Enable = true
    config.Net.SASL.User = "admin"
    config.Net.SASL.Password = "admin"

    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true

    client, err := sarama.NewClient([]string{brokerAddress}, config)
    if err != nil {
        log.Errorf("Failed to create client: %s ", err.Error())
        return err
    }
    defer client.Close()

    // Check if the topic exists
    topics, err := client.Topics()
    if err != nil {
        log.Errorf("Failed to retrieve topic list: %s ", err.Error())
        return err
    }

    // Check if the topic is in the list of topics
    if contains(topics, topicName) {
        log.Infof("Topic '%s' exists", topicName)
        return nil
    }

    admin, err := sarama.NewClusterAdmin([]string{brokerAddress}, config)
    if err != nil {
        log.Errorf("failed to create cluster admin: %v", err)
        return err
    }
    defer admin.Close()

    topicDetail := &sarama.TopicDetail{
        NumPartitions:     5, // 设置分区数量
        ReplicationFactor: 1, // 设置副本因子
    }

    err = admin.CreateTopic(topicName, topicDetail, false)
    if err != nil {
        log.Errorf("failed to create topic: %v", err)
        return err
    }
    log.Infof("Topic '%s' created successfully with %d partitions", topicName, 3)
    return nil
}

// Helper function to check if a string is in a list of strings
func contains(list []string, str string) bool {
    for _, item := range list {
        if strings.ToLower(item) == strings.ToLower(str) {
            return true
        }
    }
    return false
}

//
// DeleteTopic
//  @Description: 删除topic
//  @param topic
//
func DeleteTopic(topic string) {
    config := sarama.NewConfig()
    // sasl认证
    config.Net.SASL.Enable = true
    config.Net.SASL.User = "admin"
    config.Net.SASL.Password = "admin"

    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true
    admin, err := sarama.NewClusterAdmin([]string{brokerAddress}, config)
    if err != nil {
        log.Fatalf("Error creating cluster admin: %v", err)
    }
    defer admin.Close()

    // 删除topic
    err = admin.DeleteTopic(topic)
    if err != nil {
        log.Fatalf("Error deleting topic: %v", err)
    }
    log.Printf("Topic '%s' deleted successfully.\n", topic)
}

//
// ProduceMessage
//  @Description: 消费发送到指定topic
//  @param producer
//  @param topic
//  @param message
//  @param num
//  @return error
//
func ProduceMessage(producer sarama.SyncProducer, topic, message string) error {
    // 构建消息
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(message),
    }

    // 发送消息
    partitions, offset, err := producer.SendMessage(msg)
    log.Infof("Topic: %s, offset: %d Partitions: %v", topicName, offset, partitions)
    return err
}
