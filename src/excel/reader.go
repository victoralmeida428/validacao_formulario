package excel

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

// reRefs agora captura referências individuais (ex: A1) E intervalos completos (ex: P24:P28)
var reRefs = regexp.MustCompile(`\b[A-Z]+\$?[0-9]+(?::[A-Z]+\$?[0-9]+)?\b`)

// expandirIntervalo converte notações de range ("A3:A5") em uma lista de células
func expandirIntervalo(inicio, fim string) ([]string, error) {
	colInicio, linhaInicio, err := excelize.CellNameToCoordinates(inicio)
	if err != nil {
		return nil, err
	}
	colFim, linhaFim, err := excelize.CellNameToCoordinates(fim)
	if err != nil {
		return nil, err
	}

	// Garante que a iteração vá sempre da menor para a maior coordenada
	if colInicio > colFim {
		colInicio, colFim = colFim, colInicio
	}
	if linhaInicio > linhaFim {
		linhaInicio, linhaFim = linhaFim, linhaInicio
	}

	var celulas []string
	for col := colInicio; col <= colFim; col++ {
		for linha := linhaInicio; linha <= linhaFim; linha++ {
			nome, err := excelize.CoordinatesToCellName(col, linha)
			if err == nil {
				celulas = append(celulas, nome)
			}
		}
	}
	return celulas, nil
}

// extrairInputs lê os valores reais de todas as células referenciadas por uma fórmula.
func extrairInputs(f *excelize.File, aba, formula string) string {
	formula = strings.TrimPrefix(strings.TrimSpace(formula), "=")
	formula = strings.ReplaceAll(formula, "$", "")
	formula = strings.ToUpper(formula)

	matches := reRefs.FindAllString(formula, -1)
	
	// SINTAXE CORRIGIDA: Mapa com tipo definido
	seen := make(map[string]bool)
	var inputValues []string

	for _, match := range matches {
		var celulasAlvo []string

		if strings.Contains(match, ":") {
			partesIntervalo := strings.Split(match, ":")
			if len(partesIntervalo) == 2 {
				// SINTAXE CORRIGIDA: Acesso aos índices do slice
				expandidas, err := expandirIntervalo(partesIntervalo[0], partesIntervalo[1])
				if err == nil {
					celulasAlvo = append(celulasAlvo, expandidas...)
				}
			}
		} else {
			celulasAlvo = append(celulasAlvo, match)
		}

		for _, ref := range celulasAlvo {
			// SINTAXE CORRIGIDA: Acesso à chave do mapa
			if seen[ref] {
				continue
			}
			seen[ref] = true

			val, err := f.GetCellValue(aba, ref)
			if err != nil {
				continue
			}

			val = strings.TrimSpace(val)
			if val == "" {
				val = "0"
			}

			inputValues = append(inputValues, fmt.Sprintf("%s=%s", ref, val))
		}
	}

	return strings.Join(inputValues, "; ")
}

func BuscarTodasFormulas(caminho string, onProgress func(float64)) ([]DetalhesFormula, error) {
	fTemp, err := excelize.OpenFile(caminho)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir arquivo: %v", err)
	}
	abas := fTemp.GetSheetList()
	fTemp.Close()

	var resultadosGlobais []DetalhesFormula
	var wg sync.WaitGroup
	var mu sync.Mutex

	totalAbas := len(abas)
	abasConcluidas := 0

	limiteConcorrencia := 3
	sem := make(chan struct{}, limiteConcorrencia)

	for _, nomeAba := range abas {
		wg.Add(1)

		go func(aba string) {
			sem <- struct{}{}
			defer func() {
				<-sem       // Libera a vaga
				wg.Done()
			}()

			fAba, errAba := excelize.OpenFile(caminho)
			if errAba != nil {
				return
			}
			defer fAba.Close()

			var backups []backupFormula
			rows, errRows := fAba.GetRows(aba)

			if errRows == nil {
				for rIndex, row := range rows {
					for cIndex := range row {
						colName, errCol := excelize.ColumnNumberToName(cIndex + 1)
						if errCol != nil {
							continue
						}
						coord := fmt.Sprintf("%s%d", colName, rIndex+1)

						styleID, _ := fAba.GetCellStyle(aba, coord)
						formulaOriginal, _ := fAba.GetCellFormula(aba, coord)

						fAba.SetCellStyle(aba, coord, coord, 0)
						valorSalvo, _ := fAba.GetCellValue(aba, coord)

						b := backupFormula{
							coord:       coord,
							valorSalvo:  valorSalvo,
							formulaOrig: formulaOriginal,
							styleID:     styleID,
						}

						if formulaOriginal != "" {
							b.formulaS1 = descascarIF(formulaOriginal)
							if b.formulaS1 != formulaOriginal {
								fAba.SetCellFormula(aba, coord, b.formulaS1)
							}
						}

						backups = append(backups, b)
					}
				}

				for i, b := range backups {
					if b.formulaOrig != "" {
						formulaAtual := b.formulaS1
						if strings.Contains(strings.ToUpper(formulaAtual), "AVERAGE") ||
							strings.Contains(strings.ToUpper(formulaAtual), "STDEV") ||
							strings.Contains(strings.ToUpper(formulaAtual), "SUMSQ") {

							formulaResolvida := simplificarFormula(fAba, aba, formulaAtual)
							if formulaResolvida != formulaAtual {
								fAba.SetCellFormula(aba, b.coord, formulaResolvida)
								// SINTAXE CORRIGIDA: Acesso ao índice i do slice backups
								backups[i].formulaS1 = formulaResolvida
							}
						}
					}
				}

				var resultadosAba []DetalhesFormula

				for _, b := range backups {
					if b.formulaOrig != "" {
						recalculado, errCalc := fAba.CalcCellValue(aba, b.coord)
						if errCalc != nil {
							recalculado = "Erro: " + errCalc.Error()
						}

						inputs := extrairInputs(fAba, aba, b.formulaOrig)

						resultadosAba = append(resultadosAba, DetalhesFormula{
							Aba:         aba,
							Referencia:  b.coord,
							ValorSalvo:  b.valorSalvo,
							Formula:     b.formulaOrig,
							Recalculado: recalculado,
							Inputs:      inputs,
						})
					}
				}

				for _, b := range backups {
					if b.formulaOrig != "" {
						fAba.SetCellFormula(aba, b.coord, b.formulaOrig)
					}
					fAba.SetCellStyle(aba, b.coord, b.coord, b.styleID)
				}

				mu.Lock()
				resultadosGlobais = append(resultadosGlobais, resultadosAba...)
				abasConcluidas++
				if onProgress != nil {
					onProgress(float64(abasConcluidas) / float64(totalAbas))
				}
				mu.Unlock()
			}
		}(nomeAba)
	}

	wg.Wait()
	return resultadosGlobais, nil
}
