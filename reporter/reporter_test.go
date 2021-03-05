package reporter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig_RandomHex(t *testing.T) {
	Convey("Test string generation with positive len", t, func() {
		r, err := RandomHex(10)
		So(err, ShouldBeNil)
		So(len(r), ShouldBeGreaterThan, 0)
	})

	Convey("Test string generation with zero len", t, func() {
		r, err := RandomHex(0)
		So(err, ShouldBeNil)
		So(len(r), ShouldBeZeroValue)
	})

	Convey("Test string generation with negative len", t, func() {
		r, err := RandomHex(-1)
		So(err, ShouldBeError)
		So(len(r), ShouldBeZeroValue)
	})
}

