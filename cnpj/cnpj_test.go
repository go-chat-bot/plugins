package cnpj

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCPF(t *testing.T) {
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
			cnpjValido := "99999999000100"
			bot.Args = []string{cnpjValido}

			got, error := cnpj(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf(msgFmtCnpjInvalido, cnpjValido))
		})

		Convey("Quando não é passado parâmetro deve gerar apenas 1 CPF", func() {
			got, error := cnpj(bot)

			So(error, ShouldBeNil)

			So(quantidadeCnpjGerado(got), ShouldEqual, 1)
		})

		Convey("Quando é passado uma quantidade de CPF para gerar", func() {
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
