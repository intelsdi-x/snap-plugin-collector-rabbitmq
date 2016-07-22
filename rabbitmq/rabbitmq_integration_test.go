//
// +build integration

/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/streadway/amqp"
)

func TestRabbitMQCollectMetrics(t *testing.T) {
	// Connect to RabbitMQ
	c, err := NewConsumer(
		"localhost:5672",
		"snap",
		"fanout",
		"snapq",
		"metrics",
		"conTag",
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		err := publishAmqp(
			"localhost:5672",
			"snap",
			"fanout",
			"metrics",
			"test"+string(i),
		)
		if err != nil {
			panic(err)
		}
	}

	if err := c.Shutdown(); err != nil {
		panic(err)
	}

	cfg := plugin.NewPluginConfigType()
	cfg.ConfigDataNode.AddItem("url", ctypes.ConfigValueStr{Value: "http://localhost:15672"})
	cfg.ConfigDataNode.AddItem("user", ctypes.ConfigValueStr{Value: "guest"})
	cfg.ConfigDataNode.AddItem("password", ctypes.ConfigValueStr{Value: "guest"})

	r := &rabbitMQCollector{}

	mts := []plugin.MetricType{
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "nodes", "rabbit@localhost", "disk_free"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "nodes", "rabbit@localhost", "memory", "total"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "nodes", "*", "sockets_used"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "nodes", "*", "memory", "mnesia"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "vhosts", "%2f", "message_stats", "deliver_get"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "vhosts", "*", "messages"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "vhosts", "*", "messages_details", "rate"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "exchanges", "*", "*", "message_stats", "publish_in_details", "rate"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "exchanges", "%2f", "snap", "message_stats", "publish_in"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "exchanges", "%2f", "*", "message_stats", "publish_out"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "exchanges", "*", "*", "message_stats", "confirm"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "queues", "%2f", "snapq", "messages_unacknowledged"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "queues", "%2f", "snapq", "messages_details", "rate"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "queues", "%2f", "*", "message_stats", "publish"),
		},
		plugin.MetricType{
			Config_:    cfg.ConfigDataNode,
			Namespace_: core.NewNamespace("intel", "rabbitmq", "queues", "*", "*", "messages"),
		},
	}
	metrics, err := r.CollectMetrics(mts)
	Convey("Collecting metrics from RabbitMQ node or cluster", t, func() {
		Convey("Should not return an error", func() {
			So(err, ShouldBeNil)
		})
		Convey("And metrics collected should be greater than 0", func() {
			So(len(metrics), ShouldBeGreaterThan, 0)
		})
	})
}

// We need to connect to Rabbitmq instance and send and receive messages so we can
// make sure the api has the available metrics to return

type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	tag        string
	quit       chan error
}

func NewConsumer(host, exchange, exchangeType, queueName, key, tag string) (*Consumer, error) {
	c := &Consumer{
		connection: nil,
		channel:    nil,
		tag:        tag,
		quit:       make(chan error),
	}

	var err error

	c.connection, err = amqp.Dial("amqp://" + host)
	if err != nil {
		return nil, err
	}

	go func() {
		<-c.connection.NotifyClose(make(chan *amqp.Error))
	}()

	c.channel, err = c.connection.Channel()
	if err != nil {
		return nil, err
	}

	err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type of exchange
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          //arguments
	)
	if err != nil {
		return nil, err
	}

	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     //delete when complete
		false,     //exclusive
		false,     //noWait
		nil,       //arguments
	)
	if err != nil {
		return nil, err
	}

	err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // exchange
		false,      // noWait
		nil,        // arguments
	)

	if err != nil {
		return nil, err
	}

	msgs, err := c.channel.Consume(
		queue.Name, // name of the queue
		c.tag,      // We don't care about a tag for this
		true,       // auto-ack messages
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	go receive(msgs, c.quit)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return err
	}

	if err := c.connection.Close(); err != nil {
		return err
	}

	return <-c.quit
}

func receive(msgs <-chan amqp.Delivery, quit chan error) {
	for m := range msgs {
		fmt.Println(string(m.Body))
	}
	quit <- nil
}

func publishAmqp(host, exchange, exchangeType, routingKey, body string) error {
	connection, err := amqp.Dial("amqp://" + host)
	if err != nil {
		return err
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(
		exchange,     // name of exchange
		exchangeType, // type of exchange
		true,         // durable
		false,        // auto-delete
		false,        //internal
		false,        //noWait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		exchange,   // name of exchange
		routingKey, // routing to queue
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
	if err != nil {
		return err
	}

	return nil

}
