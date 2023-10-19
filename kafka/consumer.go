package kafka

import (
    "github.com/IBM/sarama"
    log "github.com/sirupsen/logrus"
    "personal_work/redis"
    "sync"
)

var ConsumerInfo *Consumer

type Consumer struct {
    Addr          []string
    Topic         string
    PartitionList []int32
    Config        *sarama.Config
    Consumer      sarama.Consumer
}

func KafkaConsumer() {
    // 初始化
    NewKafkaConsumer()
    var wg sync.WaitGroup
    wg.Add(len(ConsumerInfo.PartitionList))
    // 然后每个分区开一个 goroutine 来消费
    for _, partitionId := range ConsumerInfo.PartitionList { // 遍历所有的分区
        go ConsumePartition(partitionId, ConsumerInfo.Topic, ConsumerInfo.Consumer, &wg)
    }
    wg.Wait()
    return
}

func NewKafkaConsumer() {
    c := Consumer{}
    c.Addr = []string{"172.31.12.5:9092"}
    c.Topic = "zpt_test_topic"

    config := sarama.NewConfig()
    config.Net.SASL.Enable = true
    config.Net.SASL.User = "admin"
    config.Net.SASL.Password = "admin"
    c.Config = config

    consumer, err := sarama.NewConsumer(c.Addr, config)
    if err != nil {
        log.Printf("fail to get consumer of NewConsumer err: %v\n", err)
    }
    c.Consumer = consumer

    // 根据topic取到所有的分区
    partitionList, err := consumer.Partitions(c.Topic)
    if err != nil {
        log.Printf("fail to get list of partition err: %v\n", err)
        return
    }
    c.PartitionList = partitionList
    ConsumerInfo = &c
}

//
// ConsumePartition
//  @Description: 数据消费，起到分流分压的作用，不需要协程并发消费
//  @param partitionId
//  @param topic
//  @param consumer
//  @param wg
//
func ConsumePartition(partitionId int32, topic string, consumer sarama.Consumer, wg *sync.WaitGroup) {
    defer wg.Done()
    offset := redis.GetInt(topic)
    if offset == 0 {
        offset = sarama.OffsetOldest
    }
    partitionConsumer, err := consumer.ConsumePartition(topic, partitionId, offset)
    if err != nil {
        log.Fatal("ConsumePartition err: ", err)
    }
    defer partitionConsumer.Close()
    // 存储的为原始数据，需要进行数据转换
    for message := range partitionConsumer.Messages() {
        value := string(message.Value)
        log.Infof("ConsumePartition partitionId: %v， offset: %v， value: %v", partitionId, offset, value)
        // 缓存offset，确保断开恢复，不再消费历史数据
        redis.SetInt(topic, message.Offset)
    }
    return
}
