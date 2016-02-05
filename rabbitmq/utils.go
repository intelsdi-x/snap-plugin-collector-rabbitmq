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
)

// joinUrl will take a url and an array of endpoints and join
// them with a '/' and return the resulting string
func joinUrl(url string, endpoint ...string) string {
	p := url
	for _, e := range endpoint {
		p += "/" + e
	}
	return p
}

// metricString is used to join a metric for lookup in a metric
// list or cache
func metricString(ns []string) string {
	return strings.Join(ns, "/")
}

// createKey returns a string of two strings joined by ':'
// Used as a key for cache
func createKey(a, b string) string {
	return a + ":" + b
}

// walkResult walks through a result returned from the API and returns
// the value to the metric requested
func walkResult(res result, keys []string) (interface{}, error) {
	if len(keys) == 1 {
		if val, ok := res[keys[0]]; ok {
			return val, nil
		}
	} else {
		if val, ok := res[keys[0]]; ok {
			if v, ok := (val).(map[string]interface{}); ok {
				return walkResult(v, keys[1:])
			}
		}
	}
	return nil, fmt.Errorf("Metric not found")
}

// typeValue returns the correct type for a value returned by the API
func typeValue(valueType string, value interface{}) interface{} {
	switch valueType {
	default:
		return 0
	case "int":
		return int(value.(float64))
	case "float64":
		return value.(float64)
	}
}

// generateNamespace is used to append a metric a namespace and return
// the entire namespace to be returned with a metric.
func generateNamespace(ns []string, metric string) []string {
	return append(ns, strings.Split(metric, "/")...)
}

// splitHost takes a url.Host string and splits the string into the host
// and port, returning just the host portion
func splitHost(host string) string {
	return strings.Split(host, ":")[0]
}

// copyNamespace copies the namespace provided to a new slice
// and returns it
func copyNamespace(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
