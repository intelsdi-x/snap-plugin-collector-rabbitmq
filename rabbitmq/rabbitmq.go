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
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	name       = "rabbitmq-collector"
	version    = 3
	pluginType = plugin.CollectorPluginType
)

// Meta returns required plugin meta details to the snap daemon. This includes
// the name, version, and type of plugin along with what content encoding (GOB or JSONRPC)
// should be used for communication between the plugin and the snap daemon.
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

type rabbitMQCollector struct {
	client *client
}

// NewRabbitMQCollector returns a pointer to an instance of the rabbitMQCollector
// that will be used to call functions the instance supports
func NewRabbitMQCollector() *rabbitMQCollector {
	return &rabbitMQCollector{}
}

// GetConfigPolicy returns the config policy required for collecting metrics
// from this plugin. In order to collect metrics, the plugin needs to know
// the RabbitMQ host to collect from and the credentials (user and password) for
// communicating the API of the RabbitMQ server.
//
// Note: For remote instances of RabbitMQ, the credentials of guest/guest are disabled
// by the RabbitMQ server config.
func (r *rabbitMQCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule, _ := cpolicy.NewStringRule("url", false, "http://localhost:15672")
	rule2, _ := cpolicy.NewStringRule("user", false, "guest")
	rule3, _ := cpolicy.NewStringRule("password", false, "guest")
	config := cpolicy.NewPolicyNode()
	config.Add(rule, rule2, rule3)
	c.Add([]string{""}, config)
	return c, nil
}

// GetMetricTypes returns the metrics available to be collected from this plugin. The available metrics
// are broken down between nodes, exchanges, queues, and vhosts. For a full list of available
// metrics, please see the README or the other source files pertaining to those topics.
func (r *rabbitMQCollector) GetMetricTypes(_ plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}

	mts = append(mts, getNodeMetrics()...)
	mts = append(mts, getExchangeMetrics()...)
	mts = append(mts, getVhostMetrics()...)
	mts = append(mts, getQueueMetrics()...)

	return mts, nil
}

// CollectMetrics takes a list of metrics from the snap daemon, collects the metrics from the required
// hosts and returns the collected metrics to the snap daemon.
func (r *rabbitMQCollector) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	// Cache nodes, vhosts to reduce API calls on a collection
	// Exchanges and queues are multi level by vhost. Map may better serve them.
	var nodes []string
	var vhosts []string

	// Cache API results
	var nodeResCache map[string]result
	var exchResCache map[string]result
	var queueResCache map[string]result
	var vhostResCache map[string]result
	var metrics []plugin.MetricType

	timestamp := time.Now()
	for _, mt := range mts {
		var err error
		if r.client == nil {
			if err := r.init(mt.Config().Table()); err != nil {
				return nil, err
			}
		}
		tags := mt.Tags()
		if tags == nil {
			tags = map[string]string{}
		}
		tags["hostname"] = r.client.source

		ns := mt.Namespace().Strings()
		metType := ns[2]
		switch metType {
		default:
			return nil, fmt.Errorf("Unknown type, %s, for requested metric", metType)
		case "nodes":
			// Let's first verify the metric given is part of our available metrics
			metric := ns[4:]
			metStr := metricString(metric)
			if _, ok := nodeMetrics[metStr]; !ok {
				return nil, fmt.Errorf("Unsupported node metric: %s", metStr)
			}
			node := ns[3]

			// Check for wildcard to see if we are collecting for a single node or all nodes
			// A cache is used so that additional node metrics that be collected on a single pass
			// will not result in multiple API calls
			if node == "*" {
				if len(nodes) == 0 {
					nodes, err = r.client.getNodes()
					if err != nil {
						return nil, err
					}
				}
				for _, n := range nodes {
					if _, ok := nodeResCache[n]; !ok {
						res, err := r.client.getNode(n)
						if err != nil {
							return nil, err
						}
						nodeResCache = map[string]result{
							n: res,
						}
					}
					if v, err := walkResult(nodeResCache[n], metric); err == nil {
						ns_ := copyNamespace(mt.Namespace())
						ns_[3].Value = n
						metrics = append(metrics, plugin.MetricType{
							Namespace_: ns_,
							Tags_:      tags,
							Data_:      typeValue(nodeMetrics[metStr], v),
							Timestamp_: timestamp,
						})
					}
				}
			} else {
				if _, ok := nodeResCache[node]; !ok {
					res, err := r.client.getNode(node)
					if err != nil {
						return nil, err
					}
					nodeResCache = map[string]result{
						node: res,
					}
				}
				v, err := walkResult(nodeResCache[node], metric)
				if err != nil {
					return nil, err
				}
				ns_ := copyNamespace(mt.Namespace())
				metrics = append(metrics, plugin.MetricType{
					Namespace_: ns_,
					Tags_:      tags,
					Data_:      typeValue(nodeMetrics[metStr], v),
					Timestamp_: timestamp,
				})

			}
		case "exchanges":
			metric := ns[5:]
			metStr := metricString(metric)
			if _, ok := exchangeMetrics[metStr]; !ok {
				return nil, fmt.Errorf("Unsupported exchange metric: %s", metStr)
			}

			vhost := ns[3]
			exchange := ns[4]

			if vhost == "*" {
				if len(vhosts) == 0 {
					vhosts, err = r.client.getVhosts()
					if err != nil {
						return nil, err
					}
				}
				if exchange == "*" {
					for _, vh := range vhosts {
						exchanges, err := r.client.getExchangesByVhost(vh)
						if err != nil {
							return nil, err
						}
						// Need an error check here as default exchange is "", which is a bit idiotic for pulling stats
						// and a filter to filter out all default amq. exchanges when a vhost is created.
						for _, exch := range exchanges {
							if exch["name"] == "" || strings.Contains(exch["name"].(string), "amq.") {
								continue
							}
							key := createKey(vh, exch["name"].(string))
							if _, ok := exchResCache[key]; !ok {
								res, err := r.client.getExchange(vh, exch["name"].(string))
								if err != nil {
									return nil, err
								}
								exchResCache = map[string]result{
									key: res,
								}
							}
							if v, err := walkResult(exchResCache[key], metric); err == nil {
								ns_ := copyNamespace(mt.Namespace())
								ns_[3].Value = vh
								ns_[4].Value = exch["name"].(string)
								metrics = append(metrics, plugin.MetricType{
									Namespace_: ns_,
									Tags_:      tags,
									Data_:      typeValue(exchangeMetrics[metStr], v),
									Timestamp_: timestamp,
								})
							}
						}
					}
				} else {
					return nil, fmt.Errorf("We do not currently .../exchanges/*/<some name>/<some metric>")
				}

			} else {
				if exchange == "*" {
					// need a better to cache these results as they are vhost:exchange
					exchanges, err := r.client.getExchangesByVhost(vhost)
					if err != nil {
						return nil, err
					}
					for _, exch := range exchanges {
						if exch["name"] == "" || strings.Contains(exch["name"].(string), "amq.") {
							continue
						}
						key := createKey(vhost, exch["name"].(string))
						if _, ok := exchResCache[key]; !ok {
							res, err := r.client.getExchange(vhost, exch["name"].(string))
							if err != nil {
								return nil, err
							}
							exchResCache = map[string]result{
								key: res,
							}
						}
						if v, err := walkResult(exchResCache[key], metric); err == nil {
							ns_ := copyNamespace(mt.Namespace())
							ns_[4].Value = exch["name"].(string)
							metrics = append(metrics, plugin.MetricType{
								Namespace_: ns_,
								Tags_:      tags,
								Data_:      typeValue(exchangeMetrics[metStr], v),
								Timestamp_: timestamp,
							})
						}
					}
				} else {
					key := createKey(vhost, exchange)
					if _, ok := exchResCache[key]; !ok {
						res, err := r.client.getExchange(vhost, exchange)
						if err != nil {
							return nil, err
						}
						exchResCache = map[string]result{
							key: res,
						}
					}
					v, err := walkResult(exchResCache[key], metric)
					if err != nil {
						return nil, err
					}
					ns_ := copyNamespace(mt.Namespace())
					metrics = append(metrics, plugin.MetricType{
						Namespace_: ns_,
						Tags_:      tags,
						Data_:      typeValue(exchangeMetrics[metStr], v),
						Timestamp_: timestamp,
					})
				}
			}
		case "queues":
			metric := ns[5:]
			metStr := metricString(metric)
			if _, ok := queueMetrics[metStr]; !ok {
				return nil, fmt.Errorf("Unsupported queue metric: %s", metStr)
			}

			vhost := ns[3]
			queue := ns[4]

			if vhost == "*" {
				if len(vhosts) == 0 {
					vhosts, err = r.client.getVhosts()
					if err != nil {
						return nil, err
					}
				}
				if queue == "*" {
					for _, vh := range vhosts {
						queues, err := r.client.getQueuesByVhost(vh)
						if err != nil {
							return nil, err
						}
						for _, q := range queues {
							key := createKey(vh, q["name"].(string))
							if _, ok := queueResCache[key]; !ok {
								res, err := r.client.getQueue(vh, q["name"].(string))
								if err != nil {
									return nil, err
								}
								queueResCache = map[string]result{
									key: res,
								}
							}
							if v, err := walkResult(queueResCache[key], metric); err == nil {
								ns_ := copyNamespace(mt.Namespace())
								ns_[3].Value = vh
								ns_[4].Value = q["name"].(string)
								metrics = append(metrics, plugin.MetricType{
									Namespace_: ns_,
									Tags_:      tags,
									Data_:      typeValue(queueMetrics[metStr], v),
									Timestamp_: timestamp,
								})
							}
						}
					}
				} else {
					return nil, fmt.Errorf("We do not currently .../queues/*/<some name>/<some metric>")
				}

			} else {
				if queue == "*" {
					// need a better to cache these results as they are vhost:exchange
					queues, err := r.client.getQueuesByVhost(vhost)
					if err != nil {
						return nil, err
					}
					for _, q := range queues {
						key := createKey(vhost, q["name"].(string))
						if _, ok := queueResCache[key]; !ok {
							res, err := r.client.getQueue(vhost, q["name"].(string))
							if err != nil {
								return nil, err
							}
							queueResCache = map[string]result{
								key: res,
							}
						}
						if v, err := walkResult(queueResCache[key], metric); err == nil {
							ns_ := copyNamespace(mt.Namespace())
							ns_[4].Value = q["name"].(string)
							metrics = append(metrics, plugin.MetricType{
								Namespace_: ns_,
								Tags_:      tags,
								Data_:      typeValue(queueMetrics[metStr], v),
								Timestamp_: timestamp,
							})
						}
					}
				} else {
					key := createKey(vhost, queue)
					if _, ok := queueResCache[key]; !ok {
						res, err := r.client.getQueue(vhost, queue)
						if err != nil {
							return nil, err
						}
						queueResCache = map[string]result{
							key: res,
						}
					}
					v, err := walkResult(queueResCache[key], metric)
					if err != nil {
						return nil, err
					}
					ns_ := copyNamespace(mt.Namespace())
					metrics = append(metrics, plugin.MetricType{
						Namespace_: ns_,
						Tags_:      tags,
						Data_:      typeValue(queueMetrics[metStr], v),
						Timestamp_: timestamp,
					})
				}
			}
		case "vhosts":
			// Let's first verify the metric given is part of our available metrics
			metric := ns[4:]
			metStr := metricString(metric)
			if _, ok := vhostMetrics[metStr]; !ok {
				return nil, fmt.Errorf("Unsupported node metric: %s", metStr)
			}
			vhost := ns[3]

			// Check for wildcard to see if we are collecting for a single node or all nodes
			// A cache is used so that additional node metrics that be collected on a single pass
			// will not result in multiple API calls
			if vhost == "*" {
				if len(vhosts) == 0 {
					vhosts, err = r.client.getVhosts()
					if err != nil {
						return nil, err
					}
				}
				for _, vh := range vhosts {
					if _, ok := vhostResCache[vh]; !ok {
						res, err := r.client.getVhost(vh)
						if err != nil {
							return nil, err
						}
						vhostResCache = map[string]result{
							vh: res,
						}
					}
					if v, err := walkResult(vhostResCache[vh], metric); err == nil {
						ns_ := copyNamespace(mt.Namespace())
						ns_[3].Value = vh
						metrics = append(metrics, plugin.MetricType{
							Namespace_: ns_,
							Tags_:      tags,
							Data_:      typeValue(vhostMetrics[metStr], v),
							Timestamp_: timestamp,
						})
					}
				}
			} else {
				if _, ok := vhostResCache[vhost]; !ok {
					res, err := r.client.getVhost(vhost)
					if err != nil {
						return nil, err
					}
					vhostResCache = map[string]result{
						vhost: res,
					}
				}
				v, err := walkResult(vhostResCache[vhost], metric)
				if err != nil {
					return nil, err
				}
				ns_ := copyNamespace(mt.Namespace())
				metrics = append(metrics, plugin.MetricType{
					Namespace_: ns_,
					Tags_:      tags,
					Data_:      typeValue(vhostMetrics[metStr], v),
					Timestamp_: timestamp,
				})

			}
		}

	}

	return metrics, nil
}

func (r *rabbitMQCollector) init(cfg map[string]ctypes.ConfigValue) error {
	if _, ok := cfg["url"]; !ok {
		return fmt.Errorf("Plugin config must contain url for RabbitMQ instance to connect to")
	}

	if _, ok := cfg["user"]; !ok {
		return fmt.Errorf("Plugin config must contain user for accessing RabbitMQ instance")
	}

	if _, ok := cfg["password"]; !ok {
		return fmt.Errorf("Plugin config must contain password for accessing RabbitMQ instance")
	}

	r.client = newClient(
		cfg["url"].(ctypes.ConfigValueStr).Value,
		cfg["user"].(ctypes.ConfigValueStr).Value,
		cfg["password"].(ctypes.ConfigValueStr).Value,
	)
	return nil
}
