package godoc

import (
	"fmt"
	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	validResults = `{
    	"results": [
	        {
	            "path": "github.com/go-chat-bot/bot",
	            "synopsis": "IRC bot written in go"
	        }
    	]
	}`

	emptyResults = `{"results":[]}`
)

func TestGoDoc(t *testing.T) {
	cmd := &bot.Cmd{}
	apiResult := ""

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, apiResult)
		}))

	godocSearchURL = ts.URL

	Convey("Given a search query text", t, func() {

		Reset(func() {
			cmd.RawArgs = ""
			apiResult = ""
		})

		Convey("When the result is empty", func() {
			cmd.RawArgs = "non existant package"
			apiResult = emptyResults

			s, err := search(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, noPackagesFound)
		})

		Convey("When the result is ok", func() {
			cmd.RawArgs = "go-bot"
			apiResult = validResults

			s, err := search(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "IRC bot written in go http://godoc.org/github.com/go-chat-bot/bot")
		})

		Convey("When the query is empty", func() {
			cmd.RawArgs = ""

			s, err := search(cmd)

			So(err, ShouldBeNil)
			So(s, ShouldEqual, "")
		})

		Convey("When the api is unreachable", func() {
			godocSearchURL = "127.0.0.1:0"
			cmd.RawArgs = "go-bot"

			_, err := search(cmd)

			So(err, ShouldNotBeNil)
		})
	})
}
