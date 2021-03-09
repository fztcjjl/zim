package util

import (
	"fmt"
	"testing"
)

func TestGenUserSig(t *testing.T) {
	//Convey("should equal IdEncode IdDecode", t, func() {
	//	id := idgen.Next()
	//	code := IdEncode(id)
	//	So(IdDecode(code), ShouldEqual, id)
	//})

	sig, _ := GenUserSig(1, "123456", "test", 86400)
	fmt.Println(CheckSig("123456", sig))
}
