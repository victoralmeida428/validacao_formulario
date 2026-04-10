package pdf

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"validador/src/validator"

	"github.com/jung-kurt/gofpdf"
)

// GerarRelatorio orquestra a criação do relatório PDF, delegando
// a lógica de negócio e o layout de seções para funções dedicadas.
func GerarRelatorio(resultados []validator.ValidacaoFormula, caminho, nome, revisao, codigo string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// 1. Organização dos Dados
	dadosProcessados := processarDados(resultados)

	// 2. Construção Visual do PDF
	gerarCabecalho(pdf, tr, nome, codigo, revisao)
	gerarResumoExecutivo(pdf, tr, dadosProcessados)
	gerarSecaoInconsistencias(pdf, tr, dadosProcessados)

	pdf.AddPage()
	gerarSecaoSucessos(pdf, tr, dadosProcessados)

	// 3. Salvamento do Arquivo
	horaArquivo := time.Now().Format("20060102_150405")
	ext := filepath.Ext(caminho)
	base := strings.TrimSuffix(caminho, ext)
	caminhoComHora := fmt.Sprintf("%s_%s%s", base, horaArquivo, ext)

	return pdf.OutputFileAndClose(caminhoComHora)
}