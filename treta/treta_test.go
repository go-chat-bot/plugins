package treta

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTreta(t *testing.T) {
	Convey("treta", t, func() {
		bot := &bot.Cmd{
			Command: "treta",
		}

		Convey("Should return a random treta", func() {
			got, error := treta(bot)

			So(error, ShouldBeNil)
			So(got, ShouldNotBeBlank)
		})

		Convey("Should return a VIM treta", func() {
			bot.Args = []string{"vim"}
			got, error := treta(bot)

			So(error, ShouldBeNil)
			So(quotes["VIM"], ShouldContain, got)
		})

		Convey("Should return a Ruby treta", func() {
			bot.Args = []string{"ruby"}
			got, error := treta(bot)

			So(error, ShouldBeNil)
			So(quotes["RUBY"], ShouldContain, got)
		})

		Convey("Should return a error message when pass a invalid param", func() {
			bot.Args = []string{"kkk"}
			got, error := treta(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, msgInvalidParam)
		})
	})
}
