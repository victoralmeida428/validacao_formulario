package excel

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

var reAverage = regexp.MustCompile(`(?i)AVERAGE\(([^)]+)\)`)
var reStdev = regexp.MustCompile(`(?i)STDEV\(([^)]+)\)`)
var reSumSq = regexp.MustCompile(`(?i)SUMSQ\(([^)]+)\)`)

func simplificarFormula(f *excelize.File, aba, formula string) string {
	formula = processarFuncaoComplexa(f, aba, formula, reAverage, calcularMedia)
	formula = processarFuncaoComplexa(f, aba, formula, reStdev, calcularDesvioPadrao)
	formula = processarFuncaoComplexa(f, aba, formula, reSumSq, calcularSomaQuadrados)
	return formula
}

func processarFuncaoComplexa(f *excelize.File, aba, formula string, re *regexp.Regexp, calculo func([]float64) float64) string {
	matches := re.FindAllStringSubmatch(formula, -1)
	for _, match := range matches {
		argumentos := strings.FieldsFunc(match[1], func(r rune) bool {
			return r == ',' || r == ';'
		})

		var todosValores []float64
		for _, arg := range argumentos {
			arg = strings.TrimSpace(arg)
			if strings.Contains(arg, ":") {
				partes := strings.Split(arg, ":")
				if len(partes) == 2 {
					todosValores = append(todosValores, mapearIntervalo(f, aba, partes[0], partes[1])...)
				}
			} else {
				formCell, _ := f.GetCellFormula(aba, arg)
				var valStr string
				if formCell != "" {
					valStr, _ = f.CalcCellValue(aba, arg)
				}
				if valStr == "" {
					valStr, _ = f.GetCellValue(aba, arg)
				}

				valStr = strings.TrimSpace(valStr)
				if valStr == "-" {
					valStr = "0"
				}

				valLimpo := strings.ReplaceAll(valStr, ",", ".")
				if v, err := strconv.ParseFloat(valLimpo, 64); err == nil {
					todosValores = append(todosValores, v)
				}
			}
		}

		resultado := calculo(todosValores)
		// ATENÇÃO: O parâmetro '-1' garante que NÃO HAVERÁ ARREDONDAMENTO
		strResultado := strconv.FormatFloat(resultado, 'f', -1, 64)
		formula = strings.Replace(formula, match[0], strResultado, 1)
	}
	return formula
}

func mapearIntervalo(f *excelize.File, aba, celulaInicio, celulaFim string) []float64 {
	var valores []float64

	colInicio, linInicio, err1 := excelize.CellNameToCoordinates(celulaInicio)
	colFim, linFim, err2 := excelize.CellNameToCoordinates(celulaFim)

	if err1 != nil || err2 != nil {
		return valores
	}

	for lin := linInicio; lin <= linFim; lin++ {
		for col := colInicio; col <= colFim; col++ {
			nomeCelula, _ := excelize.CoordinatesToCellName(col, lin)

			formula, _ := f.GetCellFormula(aba, nomeCelula)
			var valorStr string
			if formula != "" {
				valorStr, _ = f.CalcCellValue(aba, nomeCelula)
			}
			if valorStr == "" {
				valorStr, _ = f.GetCellValue(aba, nomeCelula)
			}

			valorStr = strings.TrimSpace(valorStr)
			if valorStr == "-" {
				valorStr = "0"
			}

			valLimpo := strings.ReplaceAll(valorStr, ",", ".")
			if val, err := strconv.ParseFloat(valLimpo, 64); err == nil {
				valores = append(valores, val)
			}
		}
	}
	return valores
}