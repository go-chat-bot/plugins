package cnpj

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/go-chat-bot/bot"
)

const (
	tamanhoCNPJ                     = 14
	msgParametroInvalido            = "Parâmetro inválido."
	msgQuantidadeParametrosInvalida = "Quantidade de parâmetros inválida."
	msgFmtCnpjValido                = "CNPJ %s é válido."
	msgFmtCnpjInvalido              = "CNPJ %s é inválido."
)

func cnpj(command *bot.Cmd) (string, error) {

	var param string
	if len(command.Args) == 0 {
		param = "1"
	} else if len(command.Args) == 1 {
		param = command.Args[0]
	} else {
		return msgQuantidadeParametrosInvalida, nil
	}

	if len(param) > 2 {
		if valid(param) {
			return fmt.Sprintf(msgFmtCnpjValido, command.Args[0]), nil
		}
		return fmt.Sprintf(msgFmtCnpjInvalido, command.Args[0]), nil
	}

	qtCNPJ, err := strconv.Atoi(param)
	if err != nil {
		return msgParametroInvalido, nil
	}

	var cpf string
	for i := 0; i < qtCNPJ; i++ {
		cpf += gerarCNPJ() + " "
	}
	return cpf, nil
}

func gerarCNPJ() string {
	doc := rand.Perm(12)
	dv1 := calcDV(doc)
	doc = append(doc, dv1)
	dv2 := calcDV(doc)
	doc = append(doc, dv2)

	var str string
	for _, value := range doc {
		str += strconv.Itoa(value)
	}
	return str
}

func calcDV(doc []int) int {
	multiplicadores := []int{2, 3, 4, 5, 6, 7, 8, 9}

	var calc float64
	m := 0
	for i := len(doc) - 1; i >= 0; i-- {
		calc += float64(multiplicadores[m] * doc[i])
		m++
		if m >= len(multiplicadores) {
			m = 0
		}
	}
	mod := int(math.Mod(calc*10, 11))
	if mod == 10 {
		return 0
	}
	return mod
}

func valid(cnpj string) bool {
	if len(cnpj) != tamanhoCNPJ {
		return false
	}

	if cnpj == "00000000000000" {
		return false
	}

	s := strings.Split(cnpj, "")

	doc := make([]int, 12)
	for i := 0; i <= 11; i++ {
		digito, err := strconv.Atoi(s[i])
		if err != nil {
			return false
		}
		doc[i] = digito
	}

	dv1 := calcDV(doc)
	doc = append(doc, dv1)
	dv2 := calcDV(doc)

	dv1Valido := strconv.Itoa(dv1) == string(s[12])
	dv2Valido := strconv.Itoa(dv2) == string(s[13])
	return dv1Valido && dv2Valido
}

func init() {
	bot.RegisterCommand(
		"cnpj",
		"Gerador/Validador de CNPJ.",
		"n para gerar n CNPJ e !cnpj 11111111111 para validar um CNPJ",
		cnpj)
}
