package push

import (
    "github.com/sirupsen/logrus"
    "github.com/streadway/amqp"
)

// 获取rabbitmq连接
type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    //队列名称
    QueueName string
    //交换机
    Exchange string
    //key
    key string
    //连接信息
    Url string
}

//创建RabbitMQ结构体实例
func NewRabbitMQ(url, queueName, exchange, key string) (*RabbitMQ, error) {
    rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, key: key, Url: url}
    var err error
    //创建RabbitMQ连接
    rabbitmq.conn, err = amqp.Dial(rabbitmq.Url)
    if err != nil {
        return nil, err
    }

    rabbitmq.channel, err = rabbitmq.conn.Channel()
    if err != nil {
        return nil, err
    }
    return rabbitmq, nil
}

//断开channel和connection
func (r *RabbitMQ) Destroy() {
    _ = r.channel.Close()
    _ = r.conn.Close()
}

//简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string) error {
    //1.申请队列,如果队列不存在会自动创建,如果存在则跳过创建
    //保证队列存在,消息队列能发送到队列中
    queue, err := r.channel.QueueDeclare(
        r.QueueName,
        true,  //是否持久化
        false, //是否为自动删除
        false, //是否具有排他性
        false, //是否阻塞
        nil,   //额外属性

    )
    if err != nil {
        logrus.Info("QueueDeclare error:", err)
        return err
    }
    value := amqp.Table{}
    value["bytes"] = "java.lang.String"

    values := amqp.Table{}
    values["kind"] = 'x'
    values["value"] = value

    entry := amqp.Table{}
    entry["key"] = "__TypeId__(string)"
    entry["value"] = values

    Header := amqp.Table{}
    Header["num_entry"] = 1
    Header["entries"] = entry
    //2.发送消息到队列中
    err = r.channel.Publish(
        r.Exchange,
        queue.Name,
        // 如果为true,根据exchange类型和routekey规则,如果无法找到符合条件的队列那么会把发送的消息返回给发送者
        false,
        // 如果为true,当exchange发送消息队列到队列后发现队列上没有绑定消费者,则会把消息发还给发送者
        false,
        amqp.Publishing{
            Headers:      Header,
            ContentType:  "application/json",
            Priority:     0,
            DeliveryMode: 2,
            Expiration:   "5000",
            Body:         []byte(message),
        })
    if err != nil {
        logrus.Info("Publish error:", err, ";message:", message)
        return err
    }
    return nil
}
