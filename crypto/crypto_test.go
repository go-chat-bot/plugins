package crypto

import (
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCrypto(t *testing.T) {
	Convey("Crypto", t, func() {
		bot := &bot.Cmd{
			Command: "crypto",
		}

		Convey("Should return a error message when don't pass any params", func() {
			got, error := crypto(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidAmountOfParams)
		})

		Convey("Should return a error message when pass an invalid algorithm", func() {
			bot.Args = []string{"invalidAlgorithm", "input data"}
			got, error := crypto(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, invalidParams)
		})

		Convey("using MD5 algorithm", func() {

			Convey("Should encrypt a value", func() {
				bot.Args = []string{"md5", "go-chat-bot"}
				got, error := crypto(bot)
				want := "1120d1df84fec8a0557e8737ac021651"

				So(error, ShouldBeNil)
				So(got, ShouldEqual, want)
			})

			Convey("Should encrypt multiple words", func() {
				bot.Args = []string{"md5", "The", "Go", "Programming", "Language"}
				got, error := crypto(bot)
				want := "adb505803d3502f2f00c88365ab85bf0"

				So(error, ShouldBeNil)
				So(got, ShouldEqual, want)
			})
		})

		Convey("using SHA-1 algorithm", func() {

			Convey("Should encrypt a value", func() {
				bot.Args = []string{"sha1", "go-chat-bot"}
				got, error := crypto(bot)
				want := "385ca248ffebb5ed7f62d1ea2b0545cff80ac18e"

				So(error, ShouldBeNil)
				So(got, ShouldEqual, want)
			})

			Convey("Should encrypt multiple words", func() {
				bot.Args = []string{"sha-1", "The", "Go", "Programming", "Language"}
				got, error := crypto(bot)
				want := "88a93e668044877a845097aaf620532a232bfd34"

				So(error, ShouldBeNil)
				So(got, ShouldEqual, want)
			})
		})
	})
}
