package cpf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-chat-bot/bot"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCPF(t *testing.T) {
	Convey("CPF", t, func() {
		bot := &bot.Cmd{
			Command: "cpf",
		}

		Convey("Quando é passado um CPF válido para validação", func() {
			cpfValido := "52998224725"
			bot.Args = []string{cpfValido}

			got, error := cpf(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf(msgFmtCpfValido, cpfValido))
		})

		Convey("Quando é passado um CPF inválido para validação", func() {
			cpfValido := "52998224700"
			bot.Args = []string{cpfValido}

			got, error := cpf(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf(msgFmtCpfInvalido, cpfValido))
		})

		Convey("Quando não é passado parâmetro deve gerar apenas 1 CPF", func() {
			got, error := cpf(bot)

			So(error, ShouldBeNil)

			So(quantidadeCpfGerado(got), ShouldEqual, 1)
		})

		Convey("Quando é passado uma quantidade de CPF para gerar", func() {
			bot.Args = []string{"3"}

			got, error := cpf(bot)

			So(error, ShouldBeNil)

			So(quantidadeCpfGerado(got), ShouldEqual, 3)
		})

		Convey("Quando é passado um parâmetro inválido", func() {
			bot.Args = []string{"123"}

			got, error := cpf(bot)

			So(error, ShouldBeNil)
			So(got, ShouldEqual, fmt.Sprintf(msgFmtCpfInvalido, "123"))
		})

		Convey("Quando é passado o CPF com números repetidos deve invalidar", func() {
			for i := 0; i <= 9; i++ {
				cpfInvalido := strings.Repeat(string(i), 11)

				bot.Args = []string{cpfInvalido}
				got, error := cpf(bot)

				So(error, ShouldBeNil)
				So(got, ShouldEqual, fmt.Sprintf(msgFmtCpfInvalido, cpfInvalido))
			}
		})
	})
}

func quantidadeCpfGerado(r string) int {
	return len(strings.Split(strings.Trim(r, " "), " "))
}
