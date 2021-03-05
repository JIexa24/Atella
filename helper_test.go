package atella

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig_Pause(t *testing.T) {
	Convey("Test with interrupt", t, func() {
		r := true
		Pause(1, &r)
	})

	Convey("Test without interrupt", t, func() {
		r := false
		Pause(1, &r)
	})
}

func TestConfig_ElExistsString(t *testing.T) {
	data := []string{"a", "ab", "abc", "abcd"}
	emptyData := []string{}

	Convey("Test has string", t, func() {
		result := ElExistsString(data, "ab")
		So(result, ShouldBeTrue)
	})

	Convey("Test hasn't string", t, func() {
		result := ElExistsString(data, "abcde")
		So(result, ShouldBeFalse)
		result = ElExistsString(emptyData, "abcde")
		So(result, ShouldBeFalse)
	})
}

func TestConfig_ElExistsInt64(t *testing.T) {
	data := []int64{0, 1, 2, 3, 4, -1}
	emptyData := []int64{}

	Convey("Test has int64", t, func() {
		result := ElExistsInt64(data, 1)
		So(result, ShouldBeTrue)
	})

	Convey("Test hasn't int64", t, func() {
		result := ElExistsInt64(data, 5)
		So(result, ShouldBeFalse)
		result = ElExistsInt64(emptyData, 5)
		So(result, ShouldBeFalse)
	})
}
