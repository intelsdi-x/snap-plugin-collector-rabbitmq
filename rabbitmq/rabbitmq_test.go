//
// +build small

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

	"github.com/intelsdi-x/snap/control/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRabbitMQGetMetricTypes(t *testing.T) {
	r := &rabbitMQCollector{}
	cfg := plugin.NewPluginConfigType()
	mts, err := r.GetMetricTypes(cfg)
	Convey("Getting metric types", t, func() {
		Convey("should not error", func() {
			So(err, ShouldBeNil)
		})
		Convey("should return a list of metrics greater than 0", func() {
			So(len(mts), ShouldBeGreaterThan, 0)
		})
	})
}
