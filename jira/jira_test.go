package jira

import (
	"fmt"
	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJira(t *testing.T) {
	url = "https://example.atlassian.net/browse/"
	Convey("Given a text", t, func() {
		cmd := &bot.PassiveCmd{}
		Convey("When the text does not match a jira issue syntax", func() {
			cmd.Raw = "My name is go-bot, I am awesome."
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Done, ShouldEqual, true)
			So(s.Message, ShouldBeEmpty)
		})

		Convey("When the text match a jira issue syntax", func() {
			cmd.Raw = "My name is go-bot, I am awesome. MON-965"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "MON-965"))
			So(<-s.Done, ShouldEqual, true)
			So(s.Message, ShouldBeEmpty)
		})

		Convey("When the text has a jira issue in the midle of a word", func() {
			cmd.Raw = "My name is goBOT-123"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-123"))
			So(<-s.Done, ShouldEqual, true)
			So(s.Message, ShouldBeEmpty)
		})

		Convey("When the text has a jira issue syntax with only two numbers", func() {
			cmd.Raw = "BOT-12"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-12"))
			So(<-s.Done, ShouldEqual, true)
			So(s.Message, ShouldBeEmpty)
		})

		Convey("When the jira issue isn't preceeded by space", func() {
			cmd.Raw = "::BOT-122"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-122"))
			So(<-s.Done, ShouldEqual, true)
			So(s.Message, ShouldBeEmpty)
		})

		Convey("When multiple jiras are referenced", func() {
			cmd.Raw = "::BOT-122,BOT-234 and BOT-321"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-122"))
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-234"))
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-321"))
			So(s.Message, ShouldBeEmpty)
			So(<- s.Done, ShouldEqual, true)
		})
	})
}
