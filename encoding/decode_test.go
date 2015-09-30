package encoding

import (
	"strings"
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

		Convey("Should return a error message when pass a invalid hash", func() {
			bot.Args = []string{"base64", "R28gaXMgYW4gb3BlbiBzb3VyY2Ugc", "HJvZ3JhbW1pbmcgbGFuZ3VhZ2U="}
			got, error := decode(bot)

			want := 0
			So(error, ShouldBeNil)
			So(strings.Index(got, "Error: "), ShouldEqual, want)
		})

		Convey("Should return a error message when pass correct amount of params but invalid param", func() {
			bot.Args = []string{"invalid_code", "R28gaXMgYW4gb3BlbiBzb3VyY2UgcHJvZ3JhbW1pbmcgbGFuZ3VhZ2U="}
			got, error := decode(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidParams)
		})

		Convey("Should return a error message when don't pass any params", func() {
			got, error := decode(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidAmountOfParams)
		})
	})
}
