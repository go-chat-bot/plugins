package guid

import (
	"strings"
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	guidSize = 36
)

func TestGUID(t *testing.T) {
	Convey("GUID", t, func() {
		bot := &bot.Cmd{
			Command: "guid",
		}

		Convey("Should return a valid GUID", func() {
			got, error := guid(bot)

			So(error, ShouldBeNil)
			So(len(got), ShouldEqual, guidSize)
		})

		Convey("Should return a GUID version 4", func() {
			got, error := guid(bot)

			So(error, ShouldBeNil)
			So(strings.Split(got, "")[14], ShouldEqual, "4")
		})

		Convey("Should return a upper GUID", func() {
			bot.Args = []string{"upper"}
			got, error := guid(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, strings.ToUpper(got))
		})
	})
}
