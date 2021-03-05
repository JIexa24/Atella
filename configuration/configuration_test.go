package configuration

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig_Init(t *testing.T) {
	c := GetDefault()
	Convey("Load config", t, func() {
    err := ReadConfig("./testdata/atella.yml", "", &c)
		So(err, ShouldBeNil)
	})

}
