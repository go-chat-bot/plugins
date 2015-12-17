package lula

import (
	"fmt"
	"math/rand"
	"regexp"

	"github.com/go-chat-bot/bot"
)

const (
	pattern = "(?i)\\b(lula)\\b"
)

var (
	re         = regexp.MustCompile(pattern)
	frasesLula = []string{
		"Lá, a crise é um tsunami. Aqui, se chegar, vai ser uma marolinha, que não dá nem para esquiar.",
		"É uma crise causada, fomentada, por comportamentos irracionais de gente branca, de olhos azuis, que antes da crise parecia que sabia tudo e que, agora, demonstra não saber nada.",
		"Crise? Que crise? Pergunta para o Bush.",
		"O caos aéreo é como uma metástase. A gente acha que está tudo bem, mas só descobre que o problema é bem maior quando ele (o câncer) surge.",
		"Um dia acordei invocado e liguei para o Bush.",
		"A polícia só bate em quem tem que bater.",
		"Uma mulher não pode ser submissa ao homem por causa de um prato de comida. Tem que ser submissa porque gosta dele.",
		"São privilegiados aqueles que podem pagar Imposto de Renda, porque ganham um pouco mais.",
		"Beliscão dói para cacete.",
		"Resolveram fazer um estudo para saber se a perereca estava em extinção. Aí teve que contratar gente para procurar perereca, e procure perereca, procure perereca, perereca, perereca...",
		"Você (como médico) diria ao paciente: 'Meu sifu'?",
		"Eu sei o que é greve de fome. Dá uma fome danada.",
		"Se você um dia for presidente da República, vai ver como é bom uma medida provisória.",
		"Tem que fazer uma reza profunda para que a gente deixe o otimismo (sic) no banheiro, dê descarga nele logo cedo e saia pensando em coisas boas.",
		"Política é olho no olho. É, como diria o povo brasileiro, tête à tête.",
		"Nunca fiz concessão política. Faço acordo... Se Jesus viesse para cá, e Judas tivesse a votação num partido qualquer, Jesus teria que chamar Judas para fazer coalizão.",
		"Eu e Palocci somos unha e carne. Tenho total confiança nele.",
		"[José Dirceu] é o capitão do time. Aquele que pode reclamar do juiz sem ser expulso de campo",
		"No ano que vem, estarei livre para andar por este país. Vou continuar fazendo política pelo Brasil, mas não vou dar pitaco no novo governo, porque eu quero que eles digam que nunca antes na história deste país um ex-presidente foi tão ex-presidente.",
		"Notícia é o que a gente quer esconder; o resto é propaganda.",
		"Desde que virei adulto, não vejo mais a Playboy.",
		"Sexo é uma coisa que quase todo mundo gosta e é uma necessidade orgânica.",
		"O dia que o mundo experimentar a boa cachaça brasileira, o uísque vai perder mercado.",
		"A dor da escravidão é como a de cálculo renal: não adianta dizer, tem que sentir.",
	}
)

func lula(command *bot.PassiveCmd) (string, error) {
	if re.MatchString(command.Raw) {
		return fmt.Sprintf(":lula: %s", frasesLula[rand.Intn(len(frasesLula))]), nil
	}
	return "", nil
}

func init() {
	bot.RegisterPassiveCommand(
		"lula",
		lula)
}
