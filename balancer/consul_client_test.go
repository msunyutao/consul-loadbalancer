package balancer_test

import (
	"testing"

	"github.com/msunyutao/consul-loadbalancer/balancer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConsulKVClient(t *testing.T) {
	Convey("Test NewConsulKVClient", t, func() {
		Convey("Given consul address", func() {
			cc, err := balancer.NewConsulClient("127.0.0.1:8500")
			So(err, ShouldBeNil)
			So(cc, ShouldNotBeNil)
			serviceNodes, err := cc.GetServiceNodes("hb-aerospike-1", "", true)
			So(err, ShouldNotBeNil)
			So(serviceNodes, ShouldBeNil)
			serviceNodes, err = cc.GetServiceNodes("hb-aerospike", "", true)
			So(err, ShouldBeNil)
			So(serviceNodes, ShouldNotBeNil)
			t.Logf("%d, %+v", len(serviceNodes), serviceNodes)
		})
	})
}
