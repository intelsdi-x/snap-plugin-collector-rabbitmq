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
	"strings"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
)

// exchangeMetrics is a map of available metrics for an exchange
// with the type for the value returned by the API.
var exchangeMetrics = map[string]string{
	"message_stats/confirm":                  "int",
	"message_stats/confirm_details/rate":     "float64",
	"message_stats/publish_in":               "int",
	"message_stats/publish_in_details/rate":  "float64",
	"message_stats/publish_out":              "int",
	"message_stats/publish_out_details/rate": "float64",
}

// getExchanges returns a slice of exchanges names available on the node or cluster.
// This is a list of all exchanges on the system for all vhosts.
func (c *client) getExchanges() ([]string, error) {
	exchanges, err := c.queryEndpoint(APIExchangesEndpoint)
	if err != nil {
		return nil, err
	}
	exchArr := make([]string, len(exchanges))
	for i, e := range exchanges {
		exchArr[i] = e["name"].(string)
	}
	return exchArr, nil
}

// getExchange returns an API result for an exchange for a vhost.
func (c *client) getExchange(vhost, name string) (result, error) {
	exchange, err := c.queryEndpoint(APIExchangesEndpoint, vhost, name)
	if err != nil {
		return nil, err
	}
	return exchange[0], nil
}

// getExchangesByVhost returns a slice of API results for all exchanges for a specific
// vhost.
func (c *client) getExchangesByVhost(vhost string) ([]result, error) {
	exchanges, err := c.queryEndpoint(APIExchangesEndpoint, vhost)
	if err != nil {
		return nil, err
	}
	return exchanges, nil
}

// getExchangeMetrics returns a slice of plugin.MetricType created from
// the map of available metrics for an exchange.
func getExchangeMetrics() []plugin.MetricType {
	var mts []plugin.MetricType
	nsPrefix := []string{"intel", "rabbitmq", "exchanges"}
	for k := range exchangeMetrics {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace(nsPrefix...).
				AddDynamicElement("vhost", "Virtual host name").
				AddDynamicElement("name", "Exchange name").
				AddStaticElements(strings.Split(k, "/")...),
		})
	}
	return mts
}
