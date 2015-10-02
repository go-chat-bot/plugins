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
	})
}
