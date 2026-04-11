package service

import (
	"fmt"
	"strings"

	"github.com/joaop/psiencontra/api/schemas"
)

func BuildPrompt(responses []schemas.Response) string {
	return buildSimplePrompt(responses)
}

func buildSimplePrompt(responses []schemas.Response) string {
	var sb strings.Builder

	sb.WriteString(`Você é um especialista em Psicologia com profundo conhecimento das 8 abordagens teóricas e dos 9 campos de atuação. Analise as respostas de um estudante de Psicologia e gere um ranking de afinidade.

As 8 abordagens teóricas são:
1. Psicanálise (Freud, Lacan, Winnicott)
2. Fenomenologia-Existencial (Husserl, Heidegger, Sartre)
3. Análise do Comportamento (Skinner)
4. Terapia Cognitivo-Comportamental (Beck, Ellis)
5. Psicologia Analítica (Jung)
6. Gestalt-terapia (Perls)
7. Psicologia Sócio-Histórica (Vigotski)
8. Humanismo / Abordagem Centrada na Pessoa (Rogers, Maslow)

Os 9 campos de atuação são:
1. Psicologia Clínica
2. Psicologia Organizacional
3. Psicologia Escolar/Educacional
4. Psicologia Social e Comunitária
5. Psicologia da Saúde/Hospitalar
6. Psicologia Jurídica
7. Psicologia do Esporte
8. Neuropsicologia
9. Psicometria

Respostas do estudante:
`)

	for _, r := range responses {
		sb.WriteString(fmt.Sprintf("\nPergunta %d: %s\n", r.QuestionID, r.QuestionText))
		if r.AnswerType == "multiple_choice" {
			sb.WriteString(fmt.Sprintf("Resposta (objetiva): %s\n", r.AnswerValue))
		} else {
			sb.WriteString(fmt.Sprintf("Resposta (dissertativa): %s\n", r.AnswerValue))
		}
	}

	sb.WriteString(jsonInstructions("2-3 frases resumindo o perfil do estudante"))
	return sb.String()
}

func BuildDetailedPrompt(responses []schemas.Response) string {
	var sb strings.Builder

	sb.WriteString(`Você é um psicólogo especialista em orientação profissional com profundo conhecimento das 8 abordagens teóricas e dos 9 campos de atuação em Psicologia. Analise as respostas de um questionário detalhado (escala Likert + respostas dissertativas) e gere um perfil de afinidade preciso e diferenciado.

As 8 abordagens teóricas são:
1. Psicanálise (Freud, Lacan, Winnicott)
2. Fenomenologia-Existencial (Husserl, Heidegger, Sartre)
3. Análise do Comportamento (Skinner)
4. Terapia Cognitivo-Comportamental (Beck, Ellis)
5. Psicologia Analítica (Jung)
6. Gestalt-terapia (Perls)
7. Psicologia Sócio-Histórica (Vigotski)
8. Humanismo / Abordagem Centrada na Pessoa (Rogers, Maslow)

Os 9 campos de atuação são:
1. Psicologia Clínica
2. Psicologia Organizacional
3. Psicologia Escolar/Educacional
4. Psicologia Social e Comunitária
5. Psicologia da Saúde/Hospitalar
6. Psicologia Jurídica
7. Psicologia do Esporte
8. Neuropsicologia
9. Psicometria

O questionário tem 4 blocos:
- Bloco 1 (Abordagens): 32 afirmações Likert (1-5), 4 por abordagem, cobrindo visão de sujeito, concepção de mudança, postura clínica e epistemologia.
- Bloco 2 (Campos): 27 afirmações Likert (1-5), 3 por campo, cobrindo público/contexto, atividade e motivação.
- Bloco 3 (Vinhetas Clínicas): 3 cenários dissertativos que revelam abordagem e campo simultaneamente.
- Bloco 4 (Reflexão Pessoal): 3 perguntas abertas sobre influências, motivação e valores.

Respostas do estudante:
`)

	for _, r := range responses {
		sb.WriteString(fmt.Sprintf("\nPergunta %d: %s\n", r.QuestionID, r.QuestionText))
		switch r.AnswerType {
		case "likert":
			sb.WriteString(fmt.Sprintf("Resposta (Likert 1-5): %s\n", r.AnswerValue))
		case "multiple_choice":
			sb.WriteString(fmt.Sprintf("Resposta (objetiva): %s\n", r.AnswerValue))
		default:
			sb.WriteString(fmt.Sprintf("Resposta (dissertativa): %s\n", r.AnswerValue))
		}
	}

	sb.WriteString(`
INSTRUÇÕES DE ANÁLISE:
1. Para itens Likert: calcule a média ponderada por abordagem/campo (4 itens por abordagem, 3 por campo). Converta a média 1-5 para escala 0-100.
2. Para vinhetas clínicas: analise o vocabulário, o raciocínio clínico e a postura implícita. Uma vinheta pode indicar afinidade com múltiplas abordagens e campos simultaneamente — distribua o peso proporcionalmente.
3. Para reflexões pessoais: identifique referências intelectuais, valores e motivações que reforcem ou contradigam o perfil Likert.
4. O score final deve integrar as 3 fontes (Likert + vinhetas + reflexões), priorizando o Likert como base quantitativa e usando as respostas abertas como ajuste qualitativo.
5. Gere scores genuinamente diferenciados — o perfil deve ter picos e vales claros, não uma linha plana.

`)
	sb.WriteString(jsonInstructions("3-5 frases descrevendo o perfil do estudante com mais profundidade, mencionando conexões entre abordagem e campo quando houver"))
	return sb.String()
}

func jsonInstructions(summaryDesc string) string {
	return fmt.Sprintf(`
Gere EXCLUSIVAMENTE um JSON válido (sem markdown, sem crases, sem texto extra) no seguinte formato:
{
  "approach_scores": {
    "psicanalise": <0-100>,
    "fenomenologia": <0-100>,
    "comportamental": <0-100>,
    "tcc": <0-100>,
    "junguiana": <0-100>,
    "gestalt": <0-100>,
    "socio_historica": <0-100>,
    "humanismo": <0-100>
  },
  "field_scores": {
    "clinica": <0-100>,
    "organizacional": <0-100>,
    "escolar": <0-100>,
    "social": <0-100>,
    "saude": <0-100>,
    "juridica": <0-100>,
    "esporte": <0-100>,
    "neuropsicologia": <0-100>,
    "psicometria": <0-100>
  },
  "approach_details": {
    "psicanalise": {"score": <0-100>, "description": "<1 frase curta>"},
    "fenomenologia": {"score": <0-100>, "description": "<1 frase curta>"},
    "comportamental": {"score": <0-100>, "description": "<1 frase curta>"},
    "tcc": {"score": <0-100>, "description": "<1 frase curta>"},
    "junguiana": {"score": <0-100>, "description": "<1 frase curta>"},
    "gestalt": {"score": <0-100>, "description": "<1 frase curta>"},
    "socio_historica": {"score": <0-100>, "description": "<1 frase curta>"},
    "humanismo": {"score": <0-100>, "description": "<1 frase curta>"}
  },
  "field_details": {
    "clinica": {"score": <0-100>, "description": "<1 frase curta>"},
    "organizacional": {"score": <0-100>, "description": "<1 frase curta>"},
    "escolar": {"score": <0-100>, "description": "<1 frase curta>"},
    "social": {"score": <0-100>, "description": "<1 frase curta>"},
    "saude": {"score": <0-100>, "description": "<1 frase curta>"},
    "juridica": {"score": <0-100>, "description": "<1 frase curta>"},
    "esporte": {"score": <0-100>, "description": "<1 frase curta>"},
    "neuropsicologia": {"score": <0-100>, "description": "<1 frase curta>"},
    "psicometria": {"score": <0-100>, "description": "<1 frase curta>"}
  },
  "summary": "<%s>"
}

REGRAS IMPORTANTES:
- Os scores devem variar de 0 a 100 e devem ser coerentes entre si
- Considere tanto as respostas objetivas quanto as dissertativas
- Nas dissertativas, analise o vocabulário, as referências e a forma de pensar
- Gere scores diferenciados — evite empates ou scores muito próximos
- Retorne APENAS o JSON, sem nenhum texto adicional`, summaryDesc)
}
