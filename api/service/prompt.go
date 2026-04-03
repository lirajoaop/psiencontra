package service

import (
	"fmt"
	"strings"

	"github.com/joaop/psiencontra/api/schemas"
)

func BuildPrompt(responses []schemas.Response) string {
	var sb strings.Builder

	sb.WriteString(`Voce e um especialista em Psicologia com profundo conhecimento das 8 abordagens teoricas e dos 9 campos de atuacao. Analise as respostas de um estudante de Psicologia e gere um ranking de afinidade.

As 8 abordagens teoricas sao:
1. Psicanalise (Freud, Lacan, Winnicott)
2. Fenomenologia-Existencial (Husserl, Heidegger, Rogers)
3. Analise do Comportamento (Skinner)
4. Terapia Cognitivo-Comportamental (Beck, Ellis)
5. Psicologia Analitica (Jung)
6. Gestalt-terapia (Perls)
7. Psicologia Socio-Historica (Vigotski)
8. Sistemica (Bateson, Minuchin)

Os 9 campos de atuacao sao:
1. Psicologia Clinica
2. Psicologia Organizacional
3. Psicologia Escolar/Educacional
4. Psicologia Social e Comunitaria
5. Psicologia da Saude/Hospitalar
6. Psicologia Juridica
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

	sb.WriteString(`
Com base nas respostas acima, gere EXCLUSIVAMENTE um JSON valido (sem markdown, sem crases, sem texto extra) no seguinte formato:
{
  "approach_scores": {
    "psicanalise": <0-100>,
    "fenomenologia": <0-100>,
    "comportamental": <0-100>,
    "tcc": <0-100>,
    "junguiana": <0-100>,
    "gestalt": <0-100>,
    "socio_historica": <0-100>,
    "sistemica": <0-100>
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
    "psicanalise": {"score": <0-100>, "description": "<2-3 frases explicando a afinidade>"},
    "fenomenologia": {"score": <0-100>, "description": "<2-3 frases>"},
    "comportamental": {"score": <0-100>, "description": "<2-3 frases>"},
    "tcc": {"score": <0-100>, "description": "<2-3 frases>"},
    "junguiana": {"score": <0-100>, "description": "<2-3 frases>"},
    "gestalt": {"score": <0-100>, "description": "<2-3 frases>"},
    "socio_historica": {"score": <0-100>, "description": "<2-3 frases>"},
    "sistemica": {"score": <0-100>, "description": "<2-3 frases>"}
  },
  "field_details": {
    "clinica": {"score": <0-100>, "description": "<2-3 frases explicando a afinidade>"},
    "organizacional": {"score": <0-100>, "description": "<2-3 frases>"},
    "escolar": {"score": <0-100>, "description": "<2-3 frases>"},
    "social": {"score": <0-100>, "description": "<2-3 frases>"},
    "saude": {"score": <0-100>, "description": "<2-3 frases>"},
    "juridica": {"score": <0-100>, "description": "<2-3 frases>"},
    "esporte": {"score": <0-100>, "description": "<2-3 frases>"},
    "neuropsicologia": {"score": <0-100>, "description": "<2-3 frases>"},
    "psicometria": {"score": <0-100>, "description": "<2-3 frases>"}
  },
  "summary": "<Paragrafo de 4-5 frases resumindo o perfil geral do estudante, suas inclinacoes teoricas predominantes e campos de atuacao mais alinhados>"
}

REGRAS IMPORTANTES:
- Os scores devem variar de 0 a 100 e devem ser coerentes entre si
- Considere tanto as respostas objetivas quanto as dissertativas
- Nas dissertativas, analise o vocabulario, as referencias e a forma de pensar
- Gere scores diferenciados — evite empates ou scores muito proximos
- Retorne APENAS o JSON, sem nenhum texto adicional`)

	return sb.String()
}
