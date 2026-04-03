package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/joaop/psiencontra/api/schemas"
	"github.com/jung-kurt/gofpdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

type scoreEntry struct {
	Key   string
	Score float64
	Desc  string
}

var approachLabels = map[string]string{
	"psicanalise":    "Psicanalise",
	"fenomenologia":  "Fenomenologia-Existencial",
	"comportamental": "Analise do Comportamento",
	"tcc":            "Terapia Cognitivo-Comportamental",
	"junguiana":      "Psicologia Analitica (Jung)",
	"gestalt":        "Gestalt-terapia",
	"socio_historica": "Psicologia Socio-Historica",
	"sistemica":      "Sistemica",
}

var fieldLabels = map[string]string{
	"clinica":         "Psicologia Clinica",
	"organizacional":  "Psicologia Organizacional",
	"escolar":         "Psicologia Escolar/Educacional",
	"social":          "Psicologia Social e Comunitaria",
	"saude":           "Psicologia da Saude/Hospitalar",
	"juridica":        "Psicologia Juridica",
	"esporte":         "Psicologia do Esporte",
	"neuropsicologia": "Neuropsicologia",
	"psicometria":     "Psicometria",
}

func (s *PDFService) Generate(result *schemas.Result) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	// Title
	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(88, 28, 135) // violet-900
	pdf.CellFormat(0, 15, "PsiEncontra", "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 8, "Seu Perfil de Afinidade em Psicologia", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Summary
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(88, 28, 135)
	pdf.CellFormat(0, 10, "Resumo Geral", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(50, 50, 50)
	pdf.MultiCell(0, 5, result.Explanation, "", "L", false)
	pdf.Ln(8)

	// Approach Scores
	approachEntries := parseScoresWithDetails(result.ApproachScores, result.ApproachDetails, approachLabels)
	s.renderSection(pdf, "Abordagens Teoricas", approachEntries)

	// Field Scores
	pdf.AddPage()
	fieldEntries := parseScoresWithDetails(result.FieldScores, result.FieldDetails, fieldLabels)
	s.renderSection(pdf, "Campos de Atuacao", fieldEntries)

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(0, 5, "Gerado por PsiEncontra - Este resultado e apenas orientativo e nao substitui acompanhamento profissional.", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *PDFService) renderSection(pdf *gofpdf.Fpdf, title string, entries []scoreEntry) {
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(88, 28, 135)
	pdf.CellFormat(0, 10, title, "", 1, "L", false, 0, "")
	pdf.Ln(2)

	for i, e := range entries {
		// Score bar
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(50, 50, 50)
		label := fmt.Sprintf("%d. %s - %d%%", i+1, e.Key, int(e.Score))
		pdf.CellFormat(0, 7, label, "", 1, "L", false, 0, "")

		// Bar
		barWidth := 170.0
		filledWidth := barWidth * e.Score / 100
		pdf.SetFillColor(139, 92, 246) // violet-500
		pdf.Rect(20, pdf.GetY(), filledWidth, 4, "F")
		pdf.SetFillColor(229, 231, 235) // gray-200
		pdf.Rect(20+filledWidth, pdf.GetY(), barWidth-filledWidth, 4, "F")
		pdf.Ln(6)

		// Description
		if e.Desc != "" {
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(80, 80, 80)
			pdf.MultiCell(0, 4, e.Desc, "", "L", false)
		}
		pdf.Ln(3)
	}
}

func parseScoresWithDetails(scoresJSON, detailsJSON json.RawMessage, labels map[string]string) []scoreEntry {
	var scores map[string]float64
	json.Unmarshal(scoresJSON, &scores)

	var details map[string]struct {
		Score       float64 `json:"score"`
		Description string  `json:"description"`
	}
	json.Unmarshal(detailsJSON, &details)

	var entries []scoreEntry
	for key, score := range scores {
		label := labels[key]
		if label == "" {
			label = key
		}
		desc := ""
		if d, ok := details[key]; ok {
			desc = d.Description
		}
		entries = append(entries, scoreEntry{Key: label, Score: score, Desc: desc})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	return entries
}
