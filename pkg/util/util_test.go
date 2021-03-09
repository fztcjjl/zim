package util

import (
	"github.com/fztcjjl/zim/pkg/idgen"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIdCodec(t *testing.T) {
	Convey("should equal IdEncode IdDecode", t, func() {
		id := idgen.Next()
		code := IdEncode(id)
		So(IdDecode(code), ShouldEqual, id)
	})
}
