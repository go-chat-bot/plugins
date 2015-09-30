package encoding

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDecode(t *testing.T) {
	Convey("Decode", t, func() {
		bot := &bot.Cmd{
			Command: "decode",
		}

		Convey("Should decode a value", func() {
			bot.Args = []string{"base64", "R28gaXMgYW4gb3BlbiBzb3VyY2UgcHJvZ3JhbW1pbmcgbGFuZ3VhZ2U="}
			got, error := decode(bot)

			want := "Go is an open source programming language"
			So(error, ShouldBeNil)
			So(got, ShouldEqual, want)
		})

		Convey("Should return a error message when pass correct amount of params but invalid param", func() {
			bot.Args = []string{"invalid_code", "R28gaXMgYW4gb3BlbiBzb3VyY2UgcHJvZ3JhbW1pbmcgbGFuZ3VhZ2U="}
			got, error := decode(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidParam)
		})

		Convey("Should return a error message when don't pass any params", func() {
			got, error := decode(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidAmountOfParams)
		})

		Convey("Should return a error message when pass invalid amount of params", func() {
			bot.Args = []string{"param1", "param2", "param3"}
			got, error := decode(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidAmountOfParams)
		})
	})
}
