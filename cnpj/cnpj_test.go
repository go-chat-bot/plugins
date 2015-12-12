package cnpj

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCNPJ(t *testing.T) {
	Convey("CNPJ", t, func() {
		bot := &bot.Cmd{
			Command: "cnpj",
		}

		Convey("Quando é passado um CNPJ válido para validação", func() {
			cnpjValido := "99999999000191"
			bot.Args = []string{cnpjValido}

			got, error := cnpj(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf(msgFmtCnpjValido, cnpjValido))
		})

		Convey("Quando é passado um CNPJ inválido para validação", func() {
			cnpjInvalido := "99999999000100"
			bot.Args = []string{cnpjInvalido}

			got, error := cnpj(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf(msgFmtCnpjInvalido, cnpjInvalido))
		})

		Convey("Quando não é passado parâmetro deve gerar apenas 1 CNPJ", func() {
			got, error := cnpj(bot)

			So(error, ShouldBeNil)

			So(quantidadeCnpjGerado(got), ShouldEqual, 1)
			So(valid(strings.Trim(got, " ")), ShouldEqual, true)
		})

		Convey("Quando é passado uma quantidade de CNPJ para gerar", func() {
			bot.Args = []string{"3"}

			got, error := cnpj(bot)

			So(error, ShouldBeNil)

			So(quantidadeCnpjGerado(got), ShouldEqual, 3)
		})

		Convey("Quando é passado um parâmetro inválido", func() {
			bot.Args = []string{"123"}

			got, error := cnpj(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf(msgFmtCnpjInvalido, "123"))
		})
		Convey("Quando é passado o CNPJ com números repetidos deve invalidar", func() {
			for i := 0; i <= 9; i++ {
				cnpjInvalido := strings.Repeat(string(i), 14)

				bot.Args = []string{cnpjInvalido}
				got, error := cnpj(bot)

				So(error, ShouldBeNil)
				So(got, ShouldEqual, fmt.Sprintf(msgFmtCnpjInvalido, cnpjInvalido))
			}
		})
	})
}

func quantidadeCnpjGerado(r string) int {
	return len(strings.Split(strings.Trim(r, " "), " "))
}
