package models

// 引入amqo包
import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	url = beego.AppConfig.DefaultString("amqpconn", "amqp://baymin:baymin1024@localhost:5672/")
	exchangeName = beego.AppConfig.DefaultString("exchange_name", "ai.train.topic")
	exchangeDurable = beego.AppConfig.DefaultBool("exchange_durable", true)
	queueName = beego.AppConfig.DefaultString("queue_name", "ai.train.topic-queue")
	queueDurable = beego.AppConfig.DefaultBool("queue_durable", true)
	exchangeRoutingKey = beego.AppConfig.DefaultString("train_routing_key", "train.start.#")
)

func Publish(data string) {
	// 创建链接
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	// 打开一个通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	// 生成一个交换机（交换机不存在的情况下）
	err = ch.ExchangeDeclare(exchangeName,"topic", exchangeDurable,false,false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	// 生成一个队列队列（队列不存在的情况下）
	_, err = ch.QueueDeclare(queueName, queueDurable, false, false, false, nil)
	failOnError(err, "Failed to declare an queue")
	// 列队与交换机绑定
	err = ch.QueueBind(queueName, exchangeRoutingKey, exchangeName, false, nil)
	failOnError(err, "Bind queue to exchange failure")
	// 指定交换机发布消息
	err = ch.Publish(exchangeName, exchangeRoutingKey, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			DeliveryMode: 2,
			Body: []byte(data),
		})
	failOnError(err, "Message publish failure")
}

func GetMsg(ack bool) (bool, []byte) {
	// 创建链接
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	// 打开一个通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	// 指定队列获取消息
	msg, ok, err := ch.Get(queueName, false)
	failOnError(err, "Message empty")
	if ok && ack {
		_ = msg.Ack(false)
	}
	return ok, msg.Body
}

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Error("%s: %s", msg, err)
	}
}
