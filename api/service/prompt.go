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

// BuildDetailedPrompt builds a prompt for the AI to generate ONLY the textual
// descriptions and summary for a detailed questionnaire result. The numeric
// scores themselves are computed deterministically via ipsative normalization
// in scoring.go — the AI does not touch them. This guarantees psychometric
// reproducibility: the same Likert answers always yield the same scores.
func BuildDetailedPrompt(approachScores, fieldScores map[string]float64) string {
	var sb strings.Builder

	sb.WriteString(`Você é um psicólogo especialista em orientação profissional. Os scores de afinidade JÁ FORAM CALCULADOS de forma determinística a partir das respostas Likert do estudante, usando normalização ipsativa (desvio em relação à média pessoal do respondente). Seu papel é APENAS gerar descrições textuais e o resumo narrativo — NÃO recalcule nem altere os scores.

As 8 abordagens teóricas:
1. Psicanálise (Freud, Lacan, Winnicott)
2. Fenomenologia-Existencial (Husserl, Heidegger, Sartre)
3. Análise do Comportamento (Skinner)
4. Terapia Cognitivo-Comportamental (Beck, Ellis)
5. Psicologia Analítica (Jung)
6. Gestalt-terapia (Perls)
7. Psicologia Sócio-Histórica (Vigotski)
8. Humanismo / Abordagem Centrada na Pessoa (Rogers, Maslow)

Os 9 campos de atuação:
1. Psicologia Clínica
2. Psicologia Organizacional
3. Psicologia Escolar/Educacional
4. Psicologia Social e Comunitária
5. Psicologia da Saúde/Hospitalar
6. Psicologia Jurídica
7. Psicologia do Esporte
8. Neuropsicologia
9. Psicometria

SCORES PRÉ-CALCULADOS (use exatamente estes valores — não recalcule):

Abordagens:
`)
	for _, key := range approachOrder {
		sb.WriteString(fmt.Sprintf("- %s: %.0f\n", key, approachScores[key]))
	}
	sb.WriteString("\nCampos:\n")
	for _, key := range fieldOrder {
		sb.WriteString(fmt.Sprintf("- %s: %.0f\n", key, fieldScores[key]))
	}

	sb.WriteString(`
SUA TAREFA:
1. Para cada abordagem e cada campo, escreva uma descrição de 1 frase curta que reflita o score e o que ele indica sobre o estudante (ex: score alto = afinidade forte, score baixo = pouca identificação).
2. Escreva um summary de 3-5 frases descrevendo o perfil geral do estudante, destacando as abordagens e campos de maior afinidade, possíveis tensões entre eles, e conexões observáveis.
3. Use EXATAMENTE os scores fornecidos acima. NÃO invente, NÃO recalcule.

`)
	sb.WriteString(detailedJSONInstructions(approachScores, fieldScores))
	return sb.String()
}

var approachOrder = []string{
	"psicanalise", "fenomenologia", "comportamental", "tcc",
	"junguiana", "gestalt", "socio_historica", "humanismo",
}

var fieldOrder = []string{
	"clinica", "organizacional", "escolar", "social", "saude",
	"juridica", "esporte", "neuropsicologia", "psicometria",
}

func detailedJSONInstructions(approachScores, fieldScores map[string]float64) string {
	var sb strings.Builder
	sb.WriteString(`Gere EXCLUSIVAMENTE um JSON válido (sem markdown, sem crases, sem texto extra) no seguinte formato. Os campos "score" devem ser EXATAMENTE os valores pré-calculados:
{
  "approach_scores": {
`)
	for i, key := range approachOrder {
		comma := ","
		if i == len(approachOrder)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    \"%s\": %.0f%s\n", key, approachScores[key], comma))
	}
	sb.WriteString(`  },
  "field_scores": {
`)
	for i, key := range fieldOrder {
		comma := ","
		if i == len(fieldOrder)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    \"%s\": %.0f%s\n", key, fieldScores[key], comma))
	}
	sb.WriteString(`  },
  "approach_details": {
`)
	for i, key := range approachOrder {
		comma := ","
		if i == len(approachOrder)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    \"%s\": {\"score\": %.0f, \"description\": \"<1 frase curta refletindo o score>\"}%s\n", key, approachScores[key], comma))
	}
	sb.WriteString(`  },
  "field_details": {
`)
	for i, key := range fieldOrder {
		comma := ","
		if i == len(fieldOrder)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    \"%s\": {\"score\": %.0f, \"description\": \"<1 frase curta refletindo o score>\"}%s\n", key, fieldScores[key], comma))
	}
	sb.WriteString(`  },
  "summary": "<3-5 frases descrevendo o perfil do estudante, destacando maiores afinidades e conexões>"
}

REGRAS:
- Use EXATAMENTE os scores fornecidos acima. NÃO altere nenhum valor numérico.
- As descrições devem ser coerentes com o score: scores altos indicam afinidade forte, scores baixos indicam pouca identificação.
- Retorne APENAS o JSON, sem nenhum texto adicional.`)
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
