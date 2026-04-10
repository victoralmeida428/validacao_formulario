package pdf

import (
	"bytes"
	"fmt"
	"time"

	"validador/src/assets"
	"validador/src/validator"

	"github.com/jung-kurt/gofpdf"
)

// larguraTotal para A4 Paisagem: 277mm
const larguraTotal = 277.0

func gerarCabecalho(pdf *gofpdf.Fpdf, tr func(string) string, nome, codigo, revisao string) {
	opcoesImg := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	leitorDeBytes := bytes.NewReader(assets.LogoBytes)
	pdf.RegisterImageOptionsReader("logo_embed", opcoesImg, leitorDeBytes)
	pdf.ImageOptions("logo_embed", 10, 10, 40, 0, false, opcoesImg, 0, "")

	pdf.SetY(15)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(larguraTotal, 10, tr("Relatório de Auditoria de Fórmulas"), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(245, 245, 245)
	dataVisaoGeral := time.Now().Format("02/01/2006 15:04:05")

	wCol := larguraTotal / 2
	pdf.CellFormat(wCol, 8, tr(" Documento: "+nome), "1", 0, "L", true, 0, "")
	pdf.CellFormat(wCol, 8, tr(" Código: "+codigo), "1", 1, "L", true, 0, "")
	pdf.CellFormat(wCol, 8, tr(" Revisão: "+revisao), "1", 0, "L", true, 0, "")
	pdf.CellFormat(wCol, 8, tr(" Data da Validação: "+dataVisaoGeral), "1", 1, "L", true, 0, "")
	pdf.Ln(8)
}

func gerarResumoExecutivo(pdf *gofpdf.Fpdf, tr func(string) string, dados DadosRelatorio) {
	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(larguraTotal, 10, tr("RESUMO EXECUTIVO POR ABA"), "B", 1, "L", false, 0, "")
	pdf.Ln(5)

	wAba := 102.0
	wRest := 35.0

	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(wAba, 8, tr("Aba"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wRest, 8, tr("Total"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wRest, 8, tr("Passou"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wRest, 8, tr("Falha"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wRest, 8, tr("Rev. Man."), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wRest, 8, tr("% Aprov."), "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	for _, aba := range dados.AbasOrdem {
		res := dados.ResumoPorAba[aba]
		aprovacao := 0.0
		if res.Total > 0 {
			aprovacao = float64(res.Passou) / float64(res.Total) * 100
		}

		if res.Falha > 0 {
			pdf.SetFillColor(255, 235, 235)
		} else if res.Revisao > 0 {
			pdf.SetFillColor(255, 248, 225)
		} else {
			pdf.SetFillColor(235, 255, 235)
		}

		pdf.CellFormat(wAba, 7, tr(aba), "1", 0, "L", true, 0, "")
		pdf.CellFormat(wRest, 7, fmt.Sprintf("%d", res.Total), "1", 0, "C", true, 0, "")

		pdf.SetTextColor(0, 120, 0)
		pdf.CellFormat(wRest, 7, fmt.Sprintf("%d", res.Passou), "1", 0, "C", true, 0, "")
		pdf.SetTextColor(0, 0, 0)

		if res.Falha > 0 {
			pdf.SetTextColor(200, 0, 0)
		}
		pdf.CellFormat(wRest, 7, fmt.Sprintf("%d", res.Falha), "1", 0, "C", true, 0, "")
		pdf.SetTextColor(0, 0, 0)

		if res.Revisao > 0 {
			pdf.SetTextColor(220, 110, 0)
		}
		pdf.CellFormat(wRest, 7, fmt.Sprintf("%d", res.Revisao), "1", 0, "C", true, 0, "")
		pdf.SetTextColor(0, 0, 0)

		pdf.CellFormat(wRest, 7, fmt.Sprintf("%.1f%%", aprovacao), "1", 1, "C", true, 0, "")
	}

	aprovacaoGeral := 0.0
	if dados.TotalGeral.Total > 0 {
		aprovacaoGeral = float64(dados.TotalGeral.Passou) / float64(dados.TotalGeral.Total) * 100
	}

	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(210, 210, 210)
	pdf.CellFormat(wAba, 8, tr("TOTAL GERAL"), "1", 0, "L", true, 0, "")
	pdf.CellFormat(wRest, 8, fmt.Sprintf("%d", dados.TotalGeral.Total), "1", 0, "C", true, 0, "")
	pdf.SetTextColor(0, 120, 0)
	pdf.CellFormat(wRest, 8, fmt.Sprintf("%d", dados.TotalGeral.Passou), "1", 0, "C", true, 0, "")
	pdf.SetTextColor(200, 0, 0)
	pdf.CellFormat(wRest, 8, fmt.Sprintf("%d", dados.TotalGeral.Falha), "1", 0, "C", true, 0, "")
	pdf.SetTextColor(220, 110, 0)
	pdf.CellFormat(wRest, 8, fmt.Sprintf("%d", dados.TotalGeral.Revisao), "1", 0, "C", true, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(wRest, 8, fmt.Sprintf("%.1f%%", aprovacaoGeral), "1", 1, "C", true, 0, "")
	pdf.Ln(8)
}

func gerarSecaoInconsistencias(pdf *gofpdf.Fpdf, tr func(string) string, dados DadosRelatorio) {
	pdf.SetFont("Arial", "B", 13)
	pdf.SetTextColor(150, 0, 0)
	pdf.CellFormat(larguraTotal, 10, tr("1. INCONSISTÊNCIAS E REVISÕES MANUAIS"), "B", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(5)

	for _, aba := range dados.AbasOrdem {
		var inconsistencias []validator.ValidacaoFormula
		for _, r := range dados.DadosPorAba[aba] {
			if r.Status == "FALHA" || r.Status == "REVISÃO MANUAL" {
				inconsistencias = append(inconsistencias, r)
			}
		}

		if len(inconsistencias) > 0 {
			res := dados.ResumoPorAba[aba]
			pdf.SetFont("Arial", "B", 10)
			pdf.SetFillColor(255, 220, 220)
			header := fmt.Sprintf(" Aba: %s    [Falha: %d | Revisão Manual: %d | Total: %d]",
				aba, res.Falha, res.Revisao, res.Total)
			pdf.CellFormat(larguraTotal, 8, tr(header), "1", 1, "L", true, 0, "")
			gerarTabelaGrupo(pdf, inconsistencias, tr)
			pdf.Ln(5)
		}
	}
}

func gerarSecaoSucessos(pdf *gofpdf.Fpdf, tr func(string) string, dados DadosRelatorio) {
	pdf.SetFont("Arial", "B", 13)
	pdf.SetTextColor(0, 100, 0)
	pdf.CellFormat(larguraTotal, 10, tr("2. FÓRMULAS VALIDADAS COM SUCESSO"), "B", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(5)

	for _, aba := range dados.AbasOrdem {
		var sucessos []validator.ValidacaoFormula
		for _, r := range dados.DadosPorAba[aba] {
			if r.Status == "PASSOU" {
				sucessos = append(sucessos, r)
			}
		}

		if len(sucessos) > 0 {
			res := dados.ResumoPorAba[aba]
			pdf.SetFont("Arial", "B", 10)
			pdf.SetFillColor(220, 255, 220)
			header := fmt.Sprintf(" Aba: %s    [Sucessos: %d | Total: %d]", aba, res.Passou, res.Total)
			pdf.CellFormat(larguraTotal, 8, tr(header), "1", 1, "L", true, 0, "")
			gerarTabelaGrupo(pdf, sucessos, tr)
			pdf.Ln(5)
		}
	}
}

func gerarTabelaGrupo(pdf *gofpdf.Fpdf, dados []validator.ValidacaoFormula, tr func(string) string) {
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(230, 230, 230)

	wCoord := 12.0
	wForm := 70.0
	wInp := 90.0
	wOrig := 35.0
	wCalc := 35.0
	wStat := 35.0
	lineHeight := 4.0

	pdf.CellFormat(wCoord, 7, tr("Cél."), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wForm, 7, tr("Fórmula do Excel"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wInp, 7, tr("Inputs Utilizados"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wOrig, 7, tr("V. Salvo"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wCalc, 7, tr("Recalc."), "1", 0, "C", true, 0, "")
	pdf.CellFormat(wStat, 7, "Status", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 7)
	for _, d := range dados {
		linhasForm := pdf.SplitLines([]byte(d.Formula), wForm)
		linhasInp := pdf.SplitLines([]byte(d.Inputs), wInp)

		maxLinhas := len(linhasForm)
		if len(linhasInp) > maxLinhas {
			maxLinhas = len(linhasInp)
		}

		rowHeight := float64(maxLinhas) * lineHeight
		if rowHeight < 6 {
			rowHeight = 6
		}

		if pdf.GetY()+rowHeight > 190 {
			pdf.AddPage()
		}

		startX := pdf.GetX()
		pdf.SetFillColor(255, 255, 255)
		estiloRect := "D"

		switch d.Status {
		case "FALHA":
			pdf.SetFillColor(255, 245, 245)
			estiloRect = "FD"
		case "REVISÃO MANUAL":
			pdf.SetFillColor(255, 252, 235)
			estiloRect = "FD"
		}

		pdf.CellFormat(wCoord, rowHeight, tr(d.Referencia), "1", 0, "C", estiloRect == "FD", 0, "")

		x, y := pdf.GetX(), pdf.GetY()
		pdf.Rect(x, y, wForm, rowHeight, estiloRect)
		pdf.MultiCell(wForm, lineHeight, tr(d.Formula), "", "L", false)
		pdf.SetXY(x+wForm, y)

		x, y = pdf.GetX(), pdf.GetY()
		pdf.Rect(x, y, wInp, rowHeight, estiloRect)
		pdf.MultiCell(wInp, lineHeight, tr(d.Inputs), "", "L", false)
		pdf.SetXY(x+wInp, y)

		pdf.CellFormat(wOrig, rowHeight, tr(d.Original), "1", 0, "C", estiloRect == "FD", 0, "")
		pdf.CellFormat(wCalc, rowHeight, tr(d.Calculado), "1", 0, "C", estiloRect == "FD", 0, "")

		switch d.Status {
		case "PASSOU":
			pdf.SetTextColor(0, 100, 0)
		case "FALHA":
			pdf.SetTextColor(180, 0, 0)
		default:
			pdf.SetTextColor(200, 100, 0)
		}
		pdf.CellFormat(wStat, rowHeight, tr(d.Status), "1", 1, "C", estiloRect == "FD", 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(startX)
	}
}
