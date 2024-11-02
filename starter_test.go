package gorabbitmq_starter_test

import (
	goframework_rabbitmq "github.com/kordar/goframework-rabbitmq"
	logger "github.com/kordar/gologger"
	rabbitmq "github.com/kordar/gorabbitmq"
	gorabbitmq_starter "github.com/kordar/gorabbitmq-starter"
	"testing"
)

type DemoConsumer struct {
	*rabbitmq.BaseConsumer
}

func (b *DemoConsumer) QueueName() string {
	return "demo"
}

func (b *DemoConsumer) RoutingKey() string {
	return "com.demo"
}

func (b *DemoConsumer) OnReceive(bytes []byte) bool {
	logger.Info("=================%v", string(bytes))
	return true
}

func TestRabbitmq(t *testing.T) {
	rabbitmqModule := gorabbitmq_starter.NewRabbitmqModule("test", func(moduleName string, itemId string, item map[string]string) {
		goframework_rabbitmq.Subscribe(itemId, item["dsn"], &DemoConsumer{})
	}, nil)
	rabbitmqModule.Load(map[string]interface{}{
		"id":                   "demo",
		"exchange_name":        "demo",
		"exchange_type":        "topic",
		"exchange_durable":     "1",
		"exchange_auto_delete": "0",
		"exchange_internal":    "0",
		"exchange_no_wait":     "0",
		"dsn":                  "amqp://admin:admin@192.168.0.190:5672/%2f",
	})

}

func TestNewRabbitmqPublishModule(t *testing.T) {
	rabbitmqPublishModule := gorabbitmq_starter.NewRabbitmqPublishModule("test", func(moduleName string, itemId string, item map[string]string) {
		logger.Info(moduleName, item)
		goframework_rabbitmq.AddChannelObject(itemId, &rabbitmq.ChannelObject{
			Name:         "topic1",
			ExchangeName: "demo",
			RoutingKey:   "com.demo",
		})
	})
	rabbitmqPublishModule.Load(map[string]interface{}{
		"id":  "demo",
		"dsn": "amqp://admin:admin@192.168.0.190:5672/%2f",
	})
	goframework_rabbitmq.Publish("demo", "topic1", []byte("AAAAAAAAAAA"))
	goframework_rabbitmq.GetPublishClient("demo")
}
