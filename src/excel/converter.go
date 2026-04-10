package excel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// obterComandoLibreOffice detecta o sistema operacional e busca o executável correto
func obterComandoLibreOffice() string {
	if runtime.GOOS == "windows" {
		// No Windows, busca nos caminhos padrões de instalação
		caminhosComuns := []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
			`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
		}
		for _, caminho := range caminhosComuns {
			if _, err := os.Stat(caminho); err == nil {
				return caminho
			}
		}
		// Fallback: se não achar nas pastas, tenta pelo PATH do Windows
		return "soffice.exe"
	}

	// No Linux/macOS o binário é registrado globalmente como 'soffice'
	return "soffice"
}

// PrepararArquivo atua como o interceptador primário de formatos.
func PrepararArquivo(caminhoOriginal string) (string, func(), error) {
	ext := strings.ToLower(filepath.Ext(caminhoOriginal))

	// Se for Excel nativo, libera a passagem direto com um cleanup vazio
	if ext == ".xlsx" || ext == ".xls" {
		return caminhoOriginal, func() {}, nil
	}

	// Se for LibreOffice, inicia o protocolo de conversão em vácuo
	if ext == ".ods" {
		tempDir := os.TempDir()
		comando := obterComandoLibreOffice()

		// Aciona o LibreOffice nativo do sistema operacional em modo invisível
		cmd := exec.Command(comando, "--headless", "--convert-to", "xlsx", "--outdir", tempDir, caminhoOriginal)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", func() {}, fmt.Errorf("falha ao converter ODS. Verifique se o LibreOffice está instalado.\nErro: %v\nOutput: %s", err, string(output))
		}

		// Identifica o nome do arquivo recém-criado na pasta temporária
		nomeBase := strings.TrimSuffix(filepath.Base(caminhoOriginal), ext)
		caminhoConvertido := filepath.Join(tempDir, nomeBase+".xlsx")

		// Verifica se a conversão realmente gerou o arquivo
		if _, err := os.Stat(caminhoConvertido); os.IsNotExist(err) {
			return "", func() {}, fmt.Errorf("arquivo convertido não foi encontrado no diretório temporário: %s", caminhoConvertido)
		}

		// Protocolo de destruição (limpa a memória após o uso)
		cleanup := func() {
			os.Remove(caminhoConvertido)
		}

		return caminhoConvertido, cleanup, nil
	}

	return "", func() {}, fmt.Errorf("formato de arquivo não suportado pelo sistema: %s", ext)
}
