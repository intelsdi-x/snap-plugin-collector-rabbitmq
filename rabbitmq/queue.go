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
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
)

// queueMetrics is a map of available metrics for a queue with the type
// of the value returned by the API.
var queueMetrics = map[string]string{
	"consumers":                             "int",
	"disk_reads":                            "int",
	"disk_writes":                           "int",
	"memory":                                "int",
	"message_bytes":                         "int",
	"message_bytes_persistent":              "int",
	"message_bytes_ram":                     "int",
	"message_bytes_ready":                   "int",
	"message_bytes_unacknowledged":          "int",
	"message_stats/disk_reads":              "int",
	"message_stats/disk_reads_details/rate": "float64",
	"message_stats/publish":                 "int",
	"message_stats/publish_details/rate":    "float64",
	"messages":                              "int",
	"messages_details/rate":                 "float64",
	"messages_persistent":                   "int",
	"messages_ram":                          "int",
	"messages_ready":                        "int",
	"messages_ready_details/rate":           "float64",
	"messages_ready_ram":                    "int",
	"messages_unacknowledged":               "int",
	"messages_unacknowledged_details/rate":  "float64",
	"messages_unacknowledged_ram":           "int",
}

// getQueues returns a slice of named queues available on the node or cluster
// that includes all queues for all vhosts.
func (c *client) getQueues() ([]string, error) {
	queues, err := c.queryEndpoint(APIQueuesEndpoint)
	if err != nil {
		return nil, err
	}
	qArr := make([]string, len(queues))
	for i, e := range queues {
		qArr[i] = e["name"].(string)
	}
	return qArr, nil
}

// getQueue returns the API result for a single queue in a single vhost.
func (c *client) getQueue(vhost, name string) (result, error) {
	queue, err := c.queryEndpoint(APIQueuesEndpoint, vhost, name)
	if err != nil {
		return nil, err
	}
	return queue[0], nil
}

// getQueuesByVhost returns a slice of results of all queues for the specified
// vhost.
func (c *client) getQueuesByVhost(vhost string) ([]result, error) {
	queues, err := c.queryEndpoint(APIQueuesEndpoint, vhost)
	if err != nil {
		return nil, err
	}
	return queues, nil
}

// getQueueMetrics returns a slice of plugin.PluginMetricTypes made from the
// map of available metrics for a queue.
func getQueueMetrics() []plugin.PluginMetricType {
	var mts []plugin.PluginMetricType
	ns := []string{"intel", "rabbitmq", "queues", "*", "*"}
	for k := range queueMetrics {
		mts = append(mts, plugin.PluginMetricType{
			Namespace_: generateNamespace(ns, k),
			Labels_:    queueLabels(),
		})
	}
	return mts
}

func queueLabels() []core.Label {
	return []core.Label{
		{
			Index: 3,
			Name:  "vhost",
		},
		{
			Index: 4,
			Name:  "queue",
		}}
}
