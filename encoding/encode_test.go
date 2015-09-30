package encoding

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEncode(t *testing.T) {
	Convey("Encode", t, func() {
		bot := &bot.Cmd{
			Command: "encode",
		}

		Convey("Should encode a value", func() {
			bot.Args = []string{"base64", "The Go Programming Language"}
			got, error := encode(bot)

			want := "VGhlIEdvIFByb2dyYW1taW5nIExhbmd1YWdl"
			So(error, ShouldBeNil)
			So(got, ShouldEqual, want)
		})

		Convey("Should encode multiple words", func() {
			bot.Args = []string{"base64", "The", "Go", "Programming", "Language"}
			got, error := encode(bot)

			want := "VGhlIEdvIFByb2dyYW1taW5nIExhbmd1YWdl"
			So(error, ShouldBeNil)
			So(got, ShouldEqual, want)
		})

		Convey("Should return a error message when pass correct amount of params but invalid param", func() {
			bot.Args = []string{"invalid_code", "R28gaXMgYW4gb3BlbiBzb3VyY2UgcHJvZ3JhbW1pbmcgbGFuZ3VhZ2U="}
			got, error := encode(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidParams)
		})

		Convey("Should return a error message when don't pass any params", func() {
			got, error := encode(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidAmountOfParams)
		})

	})
}
