package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"validador/src/gui"

	"fyne.io/fyne/v2/app"
)

// mostrarAlertaGUI aciona um alerta visual nativo, contornando a ausência do terminal (-H=windowsgui)
func mostrarAlertaGUI(mensagem string) {
	if runtime.GOOS == "windows" {
		// Trata a string para o formato VBScript (Escapa aspas duplas e insere quebras de linha nativas)
		msgLimpa := strings.ReplaceAll(mensagem, `"`, `""`)
		msgLimpa = strings.ReplaceAll(msgLimpa, "\n", `" & vbCrLf & "`)

		// Cria o código nativo do Windows para Popup Crítico (Ícone 16 = X Vermelho)
		codigoVBS := fmt.Sprintf(`MsgBox "%s", 16, "Erro Fatal - Validador ISO 17025"`, msgLimpa)

		// Salva temporariamente na pasta %TEMP% do Windows
		caminhoTemp := filepath.Join(os.TempDir(), "panic_alert.vbs")
		err := os.WriteFile(caminhoTemp, []byte(codigoVBS), 0644)
		
		if err == nil {
			// Executa a caixa de diálogo nativa e trava a execução até o usuário clicar em "OK"
			cmd := exec.Command("wscript", caminhoTemp)
			cmd.Run()
			
			// Autodestruição do arquivo temporário
			os.Remove(caminhoTemp)
		}
	} else {
		// Fallback silencioso caso alguém rode em Linux/Mac (onde o windowsgui não se aplica)
		fmt.Println("Erro Fatal:", mensagem)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			mensagemErro := fmt.Sprintf("ALERTA CRÍTICO!\n\nO sistema encontrou uma falha fatal e não pôde continuar.\n\nDetalhes do Erro:\n%v\n\nPor favor, tire um print desta tela e envie ao suporte técnico.", r)
			mostrarAlertaGUI(mensagemErro)
			os.Exit(1)
		}
	}()


	myApp := app.New()
	mainWindow := myApp.NewWindow("Auditoria de Fórmulas - Excel")

	// Inicia a interface construída no pacote gui
	gui.SetupUI(mainWindow)
	mainWindow.ShowAndRun()
}