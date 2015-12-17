package lula

import (
	"strings"
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLula(t *testing.T) {
	Convey("Given a text", t, func() {
		cmd := &bot.PassiveCmd{}

		Convey("When the text does not match lula", func() {
			cmd.Raw = "My name is go-bot, I am awesome."
			s, err := lula(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the text match lula", func() {
			cmd.Raw = "eu não votei na lula!"

			s, err := lula(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldNotEqual, "")
			So(strings.HasPrefix(s, ":lula: "), ShouldBeTrue)
		})
	})
}
