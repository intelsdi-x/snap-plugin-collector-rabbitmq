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

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
)

// nodeMetrics is a map of metrics available to be collected per
// node along type of value that is returned.
var nodeMetrics = map[string]string{
	"disk_free":                            "int",
	"disk_free_details/rate":               "float64",
	"fd_used":                              "int",
	"fd_used_details/rate":                 "float64",
	"io_read_avg_time":                     "float64",
	"io_read_avg_time_details/rate":        "float64",
	"io_read_bytes":                        "int",
	"io_read_bytes_details/rate":           "float64",
	"io_read_count":                        "int",
	"io_read_count_details/rate":           "float64",
	"io_seek_avg_time":                     "float64",
	"io_seek_avg_time_details/rate":        "float64",
	"io_seek_count":                        "int",
	"io_seek_count_details/rate":           "float64",
	"io_sync_avg_time":                     "float64",
	"io_sync_avg_time_details/rate":        "float64",
	"io_sync_count":                        "int",
	"io_sync_count_details/rate":           "float64",
	"io_write_avg_time":                    "float64",
	"io_write_avg_time_details/rate":       "float64",
	"io_write_bytes":                       "int",
	"io_write_bytes_details/rate":          "float64",
	"io_write_count":                       "int",
	"io_write_count_details/rate":          "float64",
	"mem_used":                             "int",
	"mem_used_details/rate":                "float64",
	"memory/atom":                          "int",
	"memory/binary":                        "int",
	"memory/code":                          "int",
	"memory/connection_channels":           "int",
	"memory/connection_other":              "int",
	"memory/connection_readers":            "int",
	"memory/connection_writers":            "int",
	"memory/mgmt_db":                       "int",
	"memory/mnesia":                        "int",
	"memory/msg_index":                     "int",
	"memory/other_ets":                     "int",
	"memory/other_proc":                    "int",
	"memory/other_system":                  "int",
	"memory/plugins":                       "int",
	"memory/queue_procs":                   "int",
	"memory/queue_slave_procs":             "int",
	"memory/total":                         "int",
	"mnesia_disk_tx_count":                 "int",
	"mnesia_disk_tx_count_details/rate":    "float64",
	"mnesia_ram_tx_count":                  "int",
	"mnesia_ram_tx_count_details/rate":     "float64",
	"proc_used":                            "int",
	"proc_used_details/rate":               "float64",
	"queue_index_read_count":               "int",
	"queue_index_read_count_details/rate":  "float64",
	"queue_index_write_count":              "int",
	"queue_index_write_count_details/rate": "float64",
	"sockets_used":                         "int",
	"sockets_used_details/rate":            "float64",
}

// getNodes returns a slice of node names available from the node or cluster
func (c *client) getNodes() ([]string, error) {
	nodes, err := c.queryEndpoint(APINodesEndpoint)
	if err != nil {
		return nil, err
	}
	nodeArr := make([]string, len(nodes))
	for i, e := range nodes {
		nodeArr[i] = e["name"].(string)
	}
	return nodeArr, nil
}

// getNode returns the API result for a single node. The API request is made
// with '?memory=true' to return all memory metrics for the node as well
func (c *client) getNode(name string) (result, error) {
	node, err := c.queryEndpoint(APINodesEndpoint, name+"?memory=true")
	if err != nil {
		return nil, err
	}
	if len(node) > 1 {
		return nil, fmt.Errorf("Error: more than one result returned")
	}
	return node[0], nil
}

// getNodeMetrics returns a slice of available plugin.MetricTypes created
// from the nodeMetrics map
func getNodeMetrics() []plugin.MetricType {
	var mts []plugin.MetricType
	nsPrefix := []string{"intel", "rabbitmq", "nodes"}
	for k := range nodeMetrics {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace(nsPrefix...).
				AddDynamicElement("node_name", "Name of the node").
				AddStaticElements(strings.Split(k, "/")...),
		})
	}
	return mts
}
