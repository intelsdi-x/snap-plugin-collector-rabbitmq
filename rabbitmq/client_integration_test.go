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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPIConnection(t *testing.T) {
	c := newClient("http://localhost:15672", "guest", "guest")
	Convey("Testing API nodes endpoint", t, func() {
		_, e := c.get(joinUrl(c.url, APINodesEndpoint))
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})

	Convey("Testing API queues endpoint", t, func() {
		_, e := c.get(joinUrl(c.url, APIQueuesEndpoint))
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})

	Convey("Testing API queues endpoint for default vhost", t, func() {
		_, e := c.get(joinUrl(c.url, APIQueuesEndpoint, DefaultVhost))
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})
	Convey("Testing API vhosts endpoint", t, func() {
		_, e := c.get(joinUrl(c.url, APIVhostsEndpoint))
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})

	Convey("Testing API vhosts endpoint for default vhost", t, func() {
		_, e := c.get(joinUrl(c.url, APIVhostsEndpoint, DefaultVhost))
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})
}

func TestNodesAPI(t *testing.T) {
	c := newClient("http://localhost:15672", "guest", "guest")
	nodes, e := c.getNodes()
	Convey("Get all nodes from RabbitMQ cluster", t, func() {
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
		Convey("should not be empty", func() {
			So(len(nodes), ShouldBeGreaterThan, 0)
		})
	})

	_, e = c.getNode("rabbit@localhost")
	Convey("Get rabbit@localhost node from RabbitMQ cluster", t, func() {
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})

	metrics := getNodeMetrics()
	Convey("Get metrics that are available to be collected from a node", t, func() {
		Convey("should return an array of metrics with length equal to nodeMetrics length", func() {
			So(len(metrics), ShouldEqual, len(nodeMetrics))
		})
	})
}

func TestQueuesAPI(t *testing.T) {
	c := newClient("http://localhost:15672", "guest", "guest")
	queues, e := c.getQueues()
	Convey("Get all queues from RabbitMQ cluster", t, func() {
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
		Convey("should not be empty", func() {
			So(len(queues), ShouldBeGreaterThan, 0)
		})
	})

	_, e = c.getQueue("%2f", "snapq")
	Convey("Get test queue from RabbitMQ cluster", t, func() {
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})

	metrics := getQueueMetrics()
	Convey("Get metrics that are available to be collected from a queue", t, func() {
		Convey("should return an array of metrics with length equal to queueMetrics length", func() {
			So(len(metrics), ShouldEqual, len(queueMetrics))
		})
	})
}

func TestExchangesAPI(t *testing.T) {
	c := newClient("http://localhost:15672", "guest", "guest")
	exchanges, e := c.getExchanges()
	Convey("Get all exchanges from RabbitMQ cluster", t, func() {
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
		Convey("should not be empty", func() {
			So(len(exchanges), ShouldBeGreaterThan, 0)
		})
	})

	_, e = c.getExchange("%2f", "snap")
	Convey("Get pulse exchange from RabbitMQ cluster", t, func() {
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
	})

	metrics := getExchangeMetrics()
	Convey("Get metrics available for an exchange", t, func() {
		Convey("should return an array with length same as length of exchangeMetrics", func() {
			So(len(metrics), ShouldEqual, len(exchangeMetrics))
		})
	})
}

func TestVhostsAPI(t *testing.T) {
	c := newClient("http://localhost:15672", "guest", "guest")
	vhosts, e := c.getVhosts()
	Convey("Get all vhosts from RabbitMQ cluster", t, func() {
		Convey("should not receive an error", func() {
			So(e, ShouldBeNil)
		})
		Convey("should not be empty", func() {
			So(len(vhosts), ShouldBeGreaterThan, 0)
		})
	})

	metrics := getVhostMetrics()
	Convey("Get metrics available for a vhost", t, func() {
		Convey("should return an array with length same as length of vhostMetrics", func() {
			So(len(metrics), ShouldEqual, len(vhostMetrics))
		})
	})
}
