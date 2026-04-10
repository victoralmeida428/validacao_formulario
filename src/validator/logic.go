package validator

import (
	"math"
	"strconv"
	"strings"

	"validador/src/excel"
)

type ValidacaoFormula struct {
	Aba        string
	Referencia string
	Formula    string
	Original   string
	Calculado  string
	Status     string
	Inputs     string
}

func Validar(dados []excel.DetalhesFormula) []ValidacaoFormula {
	var resultados []ValidacaoFormula

	const tolAbsoluta = 0.0005
	const tolRelativa = 0.001

	for _, d := range dados {
		status := "PASSOU"

		sOrig := strings.TrimSpace(d.ValorSalvo)
		sCalc := strings.TrimSpace(d.Recalculado)

		// 1. Ignorar falhas nativas do motor do Excel
		if strings.Contains(sCalc, "IF accepts") ||
			strings.Contains(sCalc, "formula not") ||
			strings.Contains(sCalc, "Erro:") ||
			strings.Contains(sCalc, "#DIV/0!") ||
			strings.Contains(sCalc, "#VALUE!") ||
			strings.Contains(sCalc, "#N/A") ||
			strings.Contains(sCalc, "#NAME?") {
			status = "REVISÃO MANUAL"
		} else {
			valOrig, err1 := parseLocalFloat(sOrig)
			valCalc, err2 := parseLocalFloat(sCalc)

			isCalcEmptyOrZero := (sCalc == "")
			isOrigValidNumber := (err1 == nil)

			// 2. Regra de Limitação de AST
			if isCalcEmptyOrZero && isOrigValidNumber && strings.Contains(strings.ToUpper(d.Formula), "IF(") {
				status = "REVISÃO MANUAL"
			} else if err1 == nil && err2 == nil {

				// 3. ANÁLISE METROLÓGICA DINÂMICA
				diferencaAbsoluta := math.Abs(valOrig - valCalc)
				maiorValor := math.Max(math.Abs(valOrig), math.Abs(valCalc))
				limiteDinamico := math.Max(tolAbsoluta, tolRelativa*maiorValor)

				if diferencaAbsoluta > limiteDinamico {
					status = "FALHA"
				}

			} else {
				// 4. Equivalência Textual de Nulos
				isOrigZero := (sOrig == "" || sOrig == "0" || sOrig == "0.0000" || sOrig == "-0.000001")
				isCalcZero := (sCalc == "" || sCalc == "0" || sCalc == "0.0000")

				if isOrigZero && isCalcZero {
					status = "PASSOU"
				} else if sOrig != sCalc {
					status = "FALHA"
				}
			}
		}

		resultados = append(resultados, ValidacaoFormula{
			Aba:        d.Aba,
			Referencia: d.Referencia,
			Formula:    d.Formula,
			Original:   d.ValorSalvo,
			Calculado:  d.Recalculado,
			Status:     status,
			Inputs:     d.Inputs, // ← repassa os inputs extraídos pelo reader
		})
	}
	return resultados
}

func parseLocalFloat(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	sLimpo := strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(sLimpo, 64)
}
