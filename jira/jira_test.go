package jira

import (
	"fmt"
	gojira "github.com/andygrunwald/go-jira"
	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.URL.String(), "/")
			fname := parts[len(parts)-1]
			dat, err := ioutil.ReadFile("mocks/" + fname + ".json")
			if err != nil {
				fmt.Printf("No mock file %s.json", fname)
				// provide empty data when file does not exist
				return
			}

			fmt.Fprintf(w, "%s", dat)
		},
	))
	baseURL = ts.URL
	err := initJIRAClient()
	if err != nil {
		fmt.Print(err.Error())
	}
	return ts
}

func TestJira(t *testing.T) {
	ts := setup()
	defer ts.Close()
	url = "https://example.atlassian.net/browse/"
	projects["BOT"] = gojira.Project{}
	projects["JENKINS"] = gojira.Project{}
	projects["MON"] = gojira.Project{}
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
			cmd.Raw = "My name is go-bot, I am awesome. JENKINS-33149"

			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual,
				"JENKINS-33149 (ndeloof, Closed): Images that specify an "+
					"entrypoint can not be used as a build environment - "+
					"https://example.atlassian.net/browse/JENKINS-33149")
			So(<-s.Done, ShouldEqual, true)
			So(s.Message, ShouldBeEmpty)
		})

		Convey("When the text has a jira issue in the midle of a word", func() {
			cmd.Raw = "My name is goJENKINS-3314"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual,
				"JENKINS-3314 (no assignee, Closed): <import file=\"...\"/>"+
					" to inherit portions of configurations - "+
					"https://example.atlassian.net/browse/JENKINS-3314")
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
			cmd.Raw = "::JENKINS-3314,JENKINS-33149 and BOT-321"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(<-s.Message, ShouldEqual,
				"JENKINS-3314 (no assignee, Closed): <import file=\"...\"/>"+
					" to inherit portions of configurations - "+
					"https://example.atlassian.net/browse/JENKINS-3314")
			So(<-s.Message, ShouldEqual,
				"JENKINS-33149 (ndeloof, Closed): Images that specify an "+
					"entrypoint can not be used as a build environment - "+
					"https://example.atlassian.net/browse/JENKINS-33149")
			So(<-s.Message, ShouldEqual, fmt.Sprintf("%s%s", url, "BOT-321"))
			So(s.Message, ShouldBeEmpty)
			So(<-s.Done, ShouldEqual, true)
		})

		Convey("When jira from non-existing project is mentioned", func() {
			cmd.Raw = "I saw this NON-123 issue once!"
			s, err := jira(cmd)

			So(err, ShouldBeNil)
			So(s.Message, ShouldBeEmpty)
			So(<-s.Done, ShouldEqual, true)
		})
	})
}

func TestChannelConfig(t *testing.T) {

	Convey("Given environment variables", t, func() {
		Convey("When there is correct channel template config", func() {
			loadChannelConfigs("mocks/config1.json")

			So(len(channelConfigs), ShouldEqual, 1)

			conf, ok := channelConfigs["#chan1"]
			So(ok, ShouldEqual, true)
			So(conf.Template, ShouldEqual, "{{.Self}}")
		})

		Convey("When there are more channel configurations", func() {
			loadChannelConfigs("mocks/config2.json")

			So(len(channelConfigs), ShouldEqual, 2)

			conf, ok := channelConfigs["#chan1"]
			So(ok, ShouldEqual, true)
			So(conf.Template, ShouldEqual, "{{.Self}} - 1")

			conf, ok = channelConfigs["#chan2"]
			So(ok, ShouldEqual, true)
			So(conf.Template, ShouldEqual, "{{.Self}} - 2")
		})

		Convey("When there is channel notification config", func() {
			loadChannelConfigs("mocks/config3.json")

			So(notifyNewConfig, ShouldHaveLength, 1)

			conf, ok := notifyNewConfig["PROJ1"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan1")
		})

		Convey("When there is channel notification config with many projects", func() {
			loadChannelConfigs("mocks/config4.json")

			So(notifyNewConfig, ShouldHaveLength, 2)

			conf, ok := notifyNewConfig["PROJ1"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan1")

			conf, ok = notifyNewConfig["PROJ2"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan1")
		})

		Convey("When there is multiple channel notifications with many projects", func() {
			loadChannelConfigs("mocks/config5.json")

			So(notifyNewConfig, ShouldHaveLength, 3)

			conf, ok := notifyNewConfig["PROJ1"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 2)
			So(conf, ShouldContain, "#chan1")
			So(conf, ShouldContain, "#chan2")

			conf, ok = notifyNewConfig["PROJ2"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan1")

			conf, ok = notifyNewConfig["PROJ3"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan2")

			So(notifyResConfig, ShouldHaveLength, 3)

			conf, ok = notifyResConfig["PROJ1"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan1")

			conf, ok = notifyResConfig["PROJ2"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan2")

			conf, ok = notifyResConfig["PROJ3"]
			So(ok, ShouldEqual, true)
			So(conf, ShouldHaveLength, 1)
			So(conf, ShouldContain, "#chan2")

		})
	})

}
