// package pdf

// import (
// 	"bytes"
// 	"fmt"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"validador/src/assets"
// 	"validador/src/validator"

// 	"github.com/jung-kurt/gofpdf"
// )

// // ResumoAba armazena as estatísticas de validação para o quadro gerencial
// type ResumoAba struct {
// 	Total   int
// 	Passou  int
// 	Falha   int
// 	Revisao int
// }

// func GerarRelatorio(resultados []validator.ValidacaoFormula, caminho, nome, revisao, codigo string) error {
// 	pdf := gofpdf.New("L", "mm", "A4", "")
// 	pdf.AddPage()

// 	tr := pdf.UnicodeTranslatorFromDescriptor("")

// 	// CABEÇALHO PADRÃO
// 	opcoesImg := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
// 	leitorDeBytes := bytes.NewReader(assets.LogoBytes)
// 	pdf.RegisterImageOptionsReader("logo_embed", opcoesImg, leitorDeBytes)
// 	pdf.ImageOptions("logo_embed", 10, 10, 40, 0, false, opcoesImg, 0, "")

// 	pdf.SetY(15)
// 	pdf.SetFont("Arial", "B", 18)
// 	pdf.CellFormat(0, 10, tr("Relatório de Auditoria de Fórmulas"), "", 1, "C", false, 0, "")
// 	pdf.Ln(15)

// 	// QUADRO DE IDENTIFICAÇÃO
// 	pdf.SetFont("Arial", "B", 11)
// 	pdf.SetFillColor(245, 245, 245)
// 	dataVisaoGeral := time.Now().Format("02/01/2006 15:04:05")

// 	pdf.CellFormat(138.5, 8, tr(" Documento: "+nome), "1", 0, "L", true, 0, "")
// 	pdf.CellFormat(138.5, 8, tr(" Código: "+codigo), "1", 1, "L", true, 0, "")
// 	pdf.CellFormat(138.5, 8, tr(" Revisão: "+revisao), "1", 0, "L", true, 0, "")
// 	pdf.CellFormat(138.5, 8, tr(" Data da Validação: "+dataVisaoGeral), "1", 1, "L", true, 0, "")
// 	pdf.Ln(10)

// 	// -----------------------------------------------------------------------
// 	// ESTRUTURA DE DADOS: agrupa resultados por aba preservando a ordem
// 	// -----------------------------------------------------------------------
// 	resumoPorAba := make(map[string]*ResumoAba)
// 	dadosPorAba := make(map[string][]validator.ValidacaoFormula)
// 	var abasOrdem []string

// 	for _, r := range resultados {
// 		if _, existe := resumoPorAba[r.Aba]; !existe {
// 			abasOrdem = append(abasOrdem, r.Aba)
// 			resumoPorAba[r.Aba] = &ResumoAba{}
// 			dadosPorAba[r.Aba] = []validator.ValidacaoFormula{}
// 		}

// 		dadosPorAba[r.Aba] = append(dadosPorAba[r.Aba], r)

// 		res := resumoPorAba[r.Aba]
// 		res.Total++
// 		switch r.Status {
// 		case "PASSOU":
// 			res.Passou++
// 		case "FALHA":
// 			res.Falha++
// 		case "REVISÃO MANUAL":
// 			res.Revisao++
// 		}
// 	}

// 	// -----------------------------------------------------------------------
// 	// SEÇÃO 0: RESUMO EXECUTIVO POR ABA
// 	// -----------------------------------------------------------------------
// 	pdf.SetFont("Arial", "B", 14)
// 	pdf.CellFormat(277, 10, tr("RESUMO EXECUTIVO POR ABA"), "B", 1, "L", false, 0, "")
// 	pdf.Ln(5)

// 	// Cabeçalho
// 	pdf.SetFont("Arial", "B", 10)
// 	pdf.SetFillColor(230, 230, 230)
// 	pdf.CellFormat(87, 8, tr("Aba"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(38, 8, tr("Total"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(38, 8, tr("Passou"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(38, 8, tr("Falha"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(38, 8, tr("Revisão Manual"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(38, 8, tr("% Aprovação"), "1", 1, "C", true, 0, "")

// 	// Linhas por aba
// 	pdf.SetFont("Arial", "", 10)
// 	totalGeral := ResumoAba{}
// 	for _, aba := range abasOrdem {
// 		res := resumoPorAba[aba]
// 		totalGeral.Total += res.Total
// 		totalGeral.Passou += res.Passou
// 		totalGeral.Falha += res.Falha
// 		totalGeral.Revisao += res.Revisao

// 		aprovacao := 0.0
// 		if res.Total > 0 {
// 			aprovacao = float64(res.Passou) / float64(res.Total) * 100
// 		}

// 		// Cor de fundo da linha baseada no resultado
// 		if res.Falha > 0 {
// 			pdf.SetFillColor(255, 235, 235)
// 		} else if res.Revisao > 0 {
// 			pdf.SetFillColor(255, 248, 225)
// 		} else {
// 			pdf.SetFillColor(235, 255, 235)
// 		}

// 		pdf.CellFormat(87, 7, tr(aba), "1", 0, "L", true, 0, "")
// 		pdf.CellFormat(38, 7, fmt.Sprintf("%d", res.Total), "1", 0, "C", true, 0, "")

// 		// Passou em verde
// 		pdf.SetTextColor(0, 120, 0)
// 		pdf.CellFormat(38, 7, fmt.Sprintf("%d", res.Passou), "1", 0, "C", true, 0, "")
// 		pdf.SetTextColor(0, 0, 0)

// 		// Falha em vermelho
// 		if res.Falha > 0 {
// 			pdf.SetTextColor(200, 0, 0)
// 		}
// 		pdf.CellFormat(38, 7, fmt.Sprintf("%d", res.Falha), "1", 0, "C", true, 0, "")
// 		pdf.SetTextColor(0, 0, 0)

// 		// Revisão em laranja
// 		if res.Revisao > 0 {
// 			pdf.SetTextColor(220, 110, 0)
// 		}
// 		pdf.CellFormat(38, 7, fmt.Sprintf("%d", res.Revisao), "1", 0, "C", true, 0, "")
// 		pdf.SetTextColor(0, 0, 0)

// 		pdf.CellFormat(38, 7, fmt.Sprintf("%.1f%%", aprovacao), "1", 1, "C", true, 0, "")
// 	}

// 	// Linha de TOTAL GERAL
// 	aprovacaoGeral := 0.0
// 	if totalGeral.Total > 0 {
// 		aprovacaoGeral = float64(totalGeral.Passou) / float64(totalGeral.Total) * 100
// 	}
// 	pdf.SetFont("Arial", "B", 10)
// 	pdf.SetFillColor(210, 210, 210)
// 	pdf.CellFormat(87, 8, tr("TOTAL GERAL"), "1", 0, "L", true, 0, "")
// 	pdf.CellFormat(38, 8, fmt.Sprintf("%d", totalGeral.Total), "1", 0, "C", true, 0, "")
// 	pdf.SetTextColor(0, 120, 0)
// 	pdf.CellFormat(38, 8, fmt.Sprintf("%d", totalGeral.Passou), "1", 0, "C", true, 0, "")
// 	pdf.SetTextColor(0, 0, 0)
// 	pdf.SetTextColor(200, 0, 0)
// 	pdf.CellFormat(38, 8, fmt.Sprintf("%d", totalGeral.Falha), "1", 0, "C", true, 0, "")
// 	pdf.SetTextColor(220, 110, 0)
// 	pdf.CellFormat(38, 8, fmt.Sprintf("%d", totalGeral.Revisao), "1", 0, "C", true, 0, "")
// 	pdf.SetTextColor(0, 0, 0)
// 	pdf.CellFormat(38, 8, fmt.Sprintf("%.1f%%", aprovacaoGeral), "1", 1, "C", true, 0, "")

// 	pdf.Ln(10)

// 	// -----------------------------------------------------------------------
// 	// SEÇÃO 1: INCONSISTÊNCIAS E REVISÕES MANUAIS (agrupadas por aba)
// 	// -----------------------------------------------------------------------
// 	pdf.SetFont("Arial", "B", 14)
// 	pdf.SetTextColor(150, 0, 0)
// 	pdf.CellFormat(277, 10, tr("1. INCONSISTÊNCIAS E REVISÕES MANUAIS"), "B", 1, "L", false, 0, "")
// 	pdf.SetTextColor(0, 0, 0)
// 	pdf.Ln(5)

// 	temErrosGlobais := false
// 	for _, aba := range abasOrdem {
// 		var inconsistencias []validator.ValidacaoFormula
// 		for _, r := range dadosPorAba[aba] {
// 			if r.Status == "FALHA" || r.Status == "REVISÃO MANUAL" {
// 				inconsistencias = append(inconsistencias, r)
// 			}
// 		}

// 		if len(inconsistencias) > 0 {
// 			temErrosGlobais = true
// 			res := resumoPorAba[aba]

// 			pdf.SetFont("Arial", "B", 11)
// 			pdf.SetFillColor(255, 220, 220)
// 			cabecalhoAba := fmt.Sprintf(" Aba: %s    [Falha: %d  |  Revisão Manual: %d  |  de %d fórmulas]",
// 				aba, res.Falha, res.Revisao, res.Total)
// 			pdf.CellFormat(277, 8, tr(cabecalhoAba), "1", 1, "L", true, 0, "")
// 			gerarTabelaGrupo(pdf, inconsistencias, tr)
// 			pdf.Ln(5)
// 		}
// 	}

// 	if !temErrosGlobais {
// 		pdf.SetFont("Arial", "I", 11)
// 		pdf.CellFormat(277, 10, tr("Nenhuma inconsistência encontrada em nenhuma aba."), "", 1, "L", false, 0, "")
// 	}

// 	pdf.Ln(10)
// 	pdf.AddPage()

// 	// -----------------------------------------------------------------------
// 	// SEÇÃO 2: VALIDAÇÕES COM SUCESSO (agrupadas por aba)
// 	// -----------------------------------------------------------------------
// 	pdf.SetFont("Arial", "B", 14)
// 	pdf.SetTextColor(0, 100, 0)
// 	pdf.CellFormat(277, 10, tr("2. FÓRMULAS VALIDADAS COM SUCESSO (PASSOU)"), "B", 1, "L", false, 0, "")
// 	pdf.SetTextColor(0, 0, 0)
// 	pdf.Ln(5)

// 	for _, aba := range abasOrdem {
// 		var sucessos []validator.ValidacaoFormula
// 		for _, r := range dadosPorAba[aba] {
// 			if r.Status == "PASSOU" {
// 				sucessos = append(sucessos, r)
// 			}
// 		}

// 		if len(sucessos) > 0 {
// 			res := resumoPorAba[aba]

// 			pdf.SetFont("Arial", "B", 11)
// 			pdf.SetFillColor(220, 255, 220)
// 			cabecalhoAba := fmt.Sprintf(" Aba: %s    [Passou: %d  |  de %d fórmulas]",
// 				aba, res.Passou, res.Total)
// 			pdf.CellFormat(277, 8, tr(cabecalhoAba), "1", 1, "L", true, 0, "")
// 			gerarTabelaGrupo(pdf, sucessos, tr)
// 			pdf.Ln(5)
// 		}
// 	}

// 	// SALVAMENTO DO ARQUIVO
// 	horaArquivo := time.Now().Format("20060102_150405")
// 	ext := filepath.Ext(caminho)
// 	base := strings.TrimSuffix(caminho, ext)
// 	caminhoComHora := fmt.Sprintf("%s_%s%s", base, horaArquivo, ext)

// 	return pdf.OutputFileAndClose(caminhoComHora)
// }

// func gerarTabelaGrupo(pdf *gofpdf.Fpdf, dados []validator.ValidacaoFormula, tr func(string) string) {
// 	pdf.SetFont("Arial", "B", 8)
// 	pdf.SetFillColor(230, 230, 230)
// 	pdf.SetTextColor(0, 0, 0)

// 	wCoord := 15.0
// 	wForm := 75.0
// 	wInp := 77.0
// 	wOrig := 35.0
// 	wCalc := 35.0
// 	wStat := 40.0

// 	pdf.CellFormat(wCoord, 7, tr("Célula"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(wForm, 7, tr("Fórmula do Excel"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(wInp, 7, tr("Inputs Utilizados"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(wOrig, 7, tr("Valor Salvo"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(wCalc, 7, tr("Recalculado"), "1", 0, "C", true, 0, "")
// 	pdf.CellFormat(wStat, 7, "Status", "1", 1, "C", true, 0, "")

// 	pdf.SetFont("Arial", "", 7)
// 	for _, d := range dados {
// 		formulaStr := d.Formula
// 		if len(formulaStr) > 55 {
// 			formulaStr = formulaStr[:55] + "..."
// 		}

// 		inputStr := d.Inputs
// 		if len(inputStr) > 60 {
// 			inputStr = inputStr[:60] + "..."
// 		}

// 		// Fundo alternado por status para facilitar leitura visual
// 		switch d.Status {
// 		case "FALHA":
// 			pdf.SetFillColor(255, 245, 245)
// 		case "REVISÃO MANUAL":
// 			pdf.SetFillColor(255, 252, 235)
// 		default:
// 			pdf.SetFillColor(245, 255, 245)
// 		}

// 		pdf.CellFormat(wCoord, 6, tr(d.Referencia), "1", 0, "C", true, 0, "")
// 		pdf.CellFormat(wForm, 6, tr(formulaStr), "1", 0, "L", true, 0, "")
// 		pdf.CellFormat(wInp, 6, tr(inputStr), "1", 0, "L", true, 0, "")
// 		pdf.CellFormat(wOrig, 6, tr(d.Original), "1", 0, "C", true, 0, "")
// 		pdf.CellFormat(wCalc, 6, tr(d.Calculado), "1", 0, "C", true, 0, "")

// 		switch d.Status {
// 		case "FALHA":
// 			pdf.SetTextColor(200, 0, 0)
// 		case "REVISÃO MANUAL":
// 			pdf.SetTextColor(220, 110, 0)
// 		default:
// 			pdf.SetTextColor(0, 120, 0)
// 		}

// 		pdf.CellFormat(wStat, 6, tr(d.Status), "1", 1, "C", true, 0, "")
// 		pdf.SetTextColor(0, 0, 0)
// 	}
// }
package pdf