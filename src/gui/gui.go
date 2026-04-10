package gui

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"validador/src/assets"
	"validador/src/excel"
	"validador/src/pdf"
	"validador/src/validator"
)

func SetupUI(win fyne.Window) {
	win.Resize(fyne.NewSize(550, 700))

	resLogo := fyne.NewStaticResource("logo.png", assets.LogoBytes)
	imgLogo := canvas.NewImageFromResource(resLogo)
	imgLogo.FillMode = canvas.ImageFillContain
	imgLogo.SetMinSize(fyne.NewSize(200, 100))

	entryNome := widget.NewEntry()
	entryNome.SetPlaceHolder("Ex: ANÁLISE DE ÍONS EM ÁGUA")
	entryCodigo := widget.NewEntry()
	entryCodigo.SetPlaceHolder("Ex: FLE 215")
	entryRevisao := widget.NewEntry()
	entryRevisao.SetPlaceHolder("Ex: 07")

	var caminhoArquivoSelecionado string
	lblArquivo := widget.NewLabel("Nenhum arquivo selecionado")

	var btnScan *widget.Button
	var btnSelect *widget.Button

	// -----------------------------------------------------------------
	// SISTEMA DE BINDING (MURALHA THREAD-SAFE DO FYNE)
	// -----------------------------------------------------------------
	statusBind := binding.NewString()
	_ = statusBind.Set("")

	progressoBind := binding.NewFloat()
	_ = progressoBind.Set(0.0)

	// O Fyne não tem binding nativo para .Enable()/.Disable(),
	// então criamos um listener para capturar as mudanças de forma segura.
	botoesHabilitados := binding.NewBool()
	_ = botoesHabilitados.Set(true)
	botoesHabilitados.AddListener(binding.NewDataListener(func() {
		ativo, _ := botoesHabilitados.Get()
		if ativo {
			if btnScan != nil {
				btnScan.Enable()
			}
			if btnSelect != nil {
				btnSelect.Enable()
			}
		} else {
			if btnScan != nil {
				btnScan.Disable()
			}
			if btnSelect != nil {
				btnSelect.Disable()
			}
		}
	}))

	// Criamos os widgets "conectados" aos Bindings em vez de criar vazios
	lblStatus := widget.NewLabelWithData(statusBind)
	lblStatus.Alignment = fyne.TextAlignCenter
	lblStatus.TextStyle = fyne.TextStyle{Bold: true}

	barraProgresso := widget.NewProgressBarWithData(progressoBind)
	// -----------------------------------------------------------------

	btnSelect = widget.NewButton("Escolher Arquivo .xlsx", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			caminhoArquivoSelecionado = reader.URI().Path()
			lblArquivo.SetText("Arquivo: " + reader.URI().Name())
		}, win)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx", ".ods"}))

		caminhoExe, err := os.Executable()
		if err == nil {
			diretorioExe := filepath.Dir(caminhoExe)
			uriDir := storage.NewFileURI(diretorioExe)
			listaDir, errList := storage.ListerForURI(uriDir)
			if errList == nil {
				fd.SetLocation(listaDir)
			}
		}

		fd.Show()
	})

	btnScan = widget.NewButton("Executar Auditoria de Fórmulas", func() {
		if caminhoArquivoSelecionado == "" {
			dialog.ShowInformation("Atenção", "Selecione um arquivo Excel primeiro.", win)
			return
		}

		nomeTxt := entryNome.Text
		revisaoTxt := entryRevisao.Text
		codigoTxt := entryCodigo.Text

		// 1. Mudamos as variáveis Binding. A UI atualiza sozinha e sem falhar.
		_ = botoesHabilitados.Set(false)
		_ = statusBind.Set("⏳ Analisando múltiplas abas simultaneamente. Aguarde...")
		_ = progressoBind.Set(0.0)

		go func() {
			// 2. No final, destravamos os botões via Binding. O crash desaparece.
			defer func() { _ = botoesHabilitados.Set(true) }()

			caminhoPronto, limparTemp, err := excel.PrepararArquivo(caminhoArquivoSelecionado)
			if err != nil {
				_ = statusBind.Set("❌ Erro na preparação do arquivo: " + err.Error())
				return
			}
			defer limparTemp()

			// A biblioteca de excel injeta progresso no Bind e não quebra mais a thread principal
			dadosFormulas, err := excel.BuscarTodasFormulas(caminhoPronto, func(progresso float64) {
				_ = progressoBind.Set(progresso)
			})

			if err != nil {
				_ = statusBind.Set("❌ Erro na leitura da planilha: " + err.Error())
				return
			}

			if len(dadosFormulas) == 0 {
				_ = statusBind.Set("⚠️ Nenhuma fórmula matemática encontrada neste arquivo.")
				_ = progressoBind.Set(0.0)
				return
			}

			_ = statusBind.Set("⏳ Montando dados e estruturando validação...")
			resultados := validator.Validar(dadosFormulas)

			nomeArquivoOriginal := filepath.Base(caminhoArquivoSelecionado)
			nomeBase := strings.TrimSuffix(nomeArquivoOriginal, filepath.Ext(nomeArquivoOriginal))
			nomePDF := nomeBase + ".pdf"

			_ = statusBind.Set("⏳ Desenhando e gerando o arquivo PDF...")
			err = pdf.GerarRelatorio(resultados, nomePDF, nomeTxt, revisaoTxt, codigoTxt)
			if err != nil {
				_ = statusBind.Set("❌ Erro ao gerar PDF: " + err.Error())
				return
			}

			_ = statusBind.Set("✅ Auditoria concluída! Relatório salvo: " + nomePDF)
			_ = progressoBind.Set(1.0)
		}()
	})

	formInputs := container.NewVBox(
		widget.NewLabel("Nome do Documento:"),
		entryNome,
		widget.NewLabel("Código do Documento:"),
		entryCodigo,
		widget.NewLabel("Revisão:"),
		entryRevisao,
	)

	content := container.NewVScroll(container.NewVBox(
		container.NewCenter(imgLogo),
		widget.NewLabelWithStyle("Varredura Automática Metrológica", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		formInputs,
		widget.NewLabel(""),
		btnSelect,
		lblArquivo,
		widget.NewLabel(""),
		barraProgresso,
		btnScan,
		lblStatus,
	))

	win.SetContent(content)
}
