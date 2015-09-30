package cpf

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/go-chat-bot/bot"
)

const (
	tamanhoCPF                      = 11
	msgParametroInvalido            = "Parâmetro inválido."
	msgQuantidadeParametrosInvalida = "Quantidade de parâmetros inválida."
	msgFmtCpfValido                 = "CPF %s é válido."
	msgFmtCpfInvalido               = "CPF %s é inválido."
)

func cpf(command *bot.Cmd) (string, error) {

	var param string
	switch len(command.Args) {
	case 0:
		param = "1"
	case 1:
		param = command.Args[0]
	default:
		return msgQuantidadeParametrosInvalida, nil

	}

	if len(param) > 2 {
		if valid(param) {
			return fmt.Sprintf(msgFmtCpfValido, command.Args[0]), nil
		}
		return fmt.Sprintf(msgFmtCpfInvalido, command.Args[0]), nil
	}

	qtCPF, err := strconv.Atoi(param)
	if err != nil {
		return msgParametroInvalido, nil
	}

	var cpf string
	for i := 0; i < qtCPF; i++ {
		cpf += gerarCPF() + " "
	}
	return cpf, nil
}

func gerarCPF() string {
	doc := rand.Perm(9)
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
	var calc float64
	for i, j := 2, len(doc)-1; j >= 0; i, j = i+1, j-1 {
		calc += float64(i * doc[j])
	}
	mod := int(math.Mod(calc*10, 11))
	if mod == 10 {
		return 0
	}
	return mod
}

func valid(cpf string) bool {
	if len(cpf) != tamanhoCPF {
		return false
	}

	for i := 0; i <= 9; i++ {
		if cpf == strings.Repeat(string(i), 11) {
			return false
		}
	}

	s := strings.Split(cpf, "")

	doc := make([]int, 9)
	for i := 0; i <= 8; i++ {
		digito, err := strconv.Atoi(s[i])
		if err != nil {
			return false
		}
		doc[i] = digito
	}

	dv1 := calcDV(doc)
	doc = append(doc, dv1)
	dv2 := calcDV(doc)

	dv1Valido := strconv.Itoa(dv1) == string(s[9])
	dv2Valido := strconv.Itoa(dv2) == string(s[10])
	return dv1Valido && dv2Valido
}

func init() {
	bot.RegisterCommand(
		"cpf",
		"Gerador/Validador de CPF.",
		"n para gerar n CPF e !cpf 12345678909 para validar um CPF",
		cpf)
}
