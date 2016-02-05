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

// vhostMetrics is a map of available metrics for each vhost on the node
// or cluster with the value type returned by the API for the metric.
var vhostMetrics = map[string]string{
	"messages":                               "int",
	"messages_details/rate":                  "float64",
	"messages_acknowledged":                  "int",
	"messages_acknowledged_details/rate":     "float64",
	"messages_ready":                         "int",
	"messages_ready_details/rate":            "float64",
	"messages_unacknowledged":                "int",
	"messages_unacknowledged_details/rate":   "float64",
	"message_stats/confirm":                  "int",
	"message_stats/confirm_details/rate":     "float64",
	"message_stats/deliver_get":              "int",
	"message_stats/deliver_get_details/rate": "float64",
	"message_stats/get_no_ack":               "int",
	"message_stats/get_no_ack_details/rate":  "float64",
	"message_stats/publish":                  "int",
	"message_stats/publish_details/rate":     "float64",
}

// getVhosts returns a slice of available vhost names on the node or cluster
func (c *client) getVhosts() ([]string, error) {
	vhosts, err := c.queryEndpoint(APIVhostsEndpoint)
	if err != nil {
		return nil, err
	}
	vArr := make([]string, len(vhosts))
	for i, e := range vhosts {
		if e["name"] == "/" {
			vArr[i] = DefaultVhost
			continue
		}
		vArr[i] = e["name"].(string)
	}
	return vArr, nil
}

// getVhost returns the API result for a single vhost
func (c *client) getVhost(name string) (result, error) {
	vhost, err := c.queryEndpoint(APIVhostsEndpoint, name)
	if err != nil {
		return nil, err
	}
	return vhost[0], nil
}

// getVhostMetrics returns a slice of plugin.PluginMetricType that is
// created from the available metrics in the vhostMetrics map.
func getVhostMetrics() []plugin.PluginMetricType {
	var mts []plugin.PluginMetricType
	ns := []string{"intel", "rabbitmq", "vhosts", "*"}
	for k := range vhostMetrics {
		mts = append(mts, plugin.PluginMetricType{
			Namespace_: generateNamespace(ns, k),
			Labels_:    vhostLabels(),
		})
	}
	return mts
}

func vhostLabels() []core.Label {
	return []core.Label{{
		Index: 3,
		Name:  "vhost",
	}}
}
