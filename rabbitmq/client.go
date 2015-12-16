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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	APINodesEndpoint     = "nodes"
	APIQueuesEndpoint    = "queues"
	APIVhostsEndpoint    = "vhosts"
	APIExchangesEndpoint = "exchanges"
	DefaultVhost         = "%2f" // The default vhost on a system is just "/" but we use HEX format
	DefaultPrefix        = "api"
)

type result map[string]interface{}

type client struct {
	url      string
	source   string
	user     string
	password string
	http     *http.Client
}

// newClient returns a client to the RabbitMQ management REST API
// specified with the url, user, and password passed into the method.
func newClient(_url, user, password string) *client {
	c := &client{
		url:      joinUrl(_url, DefaultPrefix),
		user:     user,
		password: password,
	}
	u, err := url.Parse(_url)
	if err != nil {
		panic(err)
	}
	c.source = splitHost(u.Host)
	c.http = &http.Client{}
	return c
}

// get makes a get request to the RabbitMQ management REST API
// and returns body of the response as a []byte array
func (c *client) get(u string) ([]byte, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	default:
		return b, nil
	case 401, 404:
		var e map[string]string
		err := json.Unmarshal(b, &e)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error calling %s. Error: %s Reason: %s", u, e["error"], e["reason"])
	}
}

// queryEndpoint makes a get request for the specified query and returns any result received
// as a result slice. Not all queries returns a slice, but this made it simpler to always return
// return a slice.
func (c *client) queryEndpoint(ep ...string) ([]result, error) {
	if len(ep) == 0 {
		return nil, fmt.Errorf("Must pass an endpoint to query")
	}
	b, err := c.get(joinUrl(c.url, ep...))
	if err != nil {
		return nil, err
	}
	var r result
	err = json.Unmarshal(b, &r)
	if err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal array") {
			var r []result
			err = json.Unmarshal(b, &r)
			if err != nil {
				return nil, err
			}
			return r, nil
		}
		return nil, err
	}
	arr := make([]result, 1)
	arr[0] = r
	return arr, nil
}
