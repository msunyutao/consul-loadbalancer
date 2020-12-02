package util_test

import (
	"testing"

	"github.com/msunyutao/consul-loadbalancer/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIntRandomGenerator(t *testing.T) {
	Convey("Test IntRandomGenerator", t, func() {
		Convey("Given min: 1, max: 100, IntPseudoRandom return random in [1, 100]", func() {
			for i := 0; i < 20; i++ {
				t.Logf("random: %d", util.IntPseudoRandom(1, 100))
			}
		})
		Convey("Given min: 1, max: 100, IntGenuineRandom return random in [1, 100]", func() {
			for i := 0; i < 20; i++ {
				t.Logf("random: %d", util.IntGenuineRandom(1, 100))
			}
		})
		Convey("FloatPseudoRandom return random in [0, 1.0)", func() {
			for i := 0; i < 50; i++ {
				t.Logf("random: %f", util.FloatPseudoRandom())
			}
		})
	})
}
