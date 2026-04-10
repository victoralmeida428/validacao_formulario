package excel

// backupFormula guarda o estado de cada célula durante a compilação
type backupFormula struct {
	coord       string
	valorSalvo  string
	formulaOrig string
	formulaS1   string // Fórmula após a poda do IF
	styleID     int    // NOVO: Guarda a máscara visual (casas decimais)
}

// No pacote excel (models.go)
type DetalhesFormula struct {
	Aba         string
	Referencia  string
	ValorSalvo  string
	Formula     string
	Recalculado string
	Inputs      string // NOVO: Guarda os valores base encontrados
}

// No pacote validator (logic.go)
type ValidacaoFormula struct {
	Aba        string
	Referencia string
	Formula    string
	Original   string
	Calculado  string
	Status     string
	Inputs     string // NOVO: Repassa os inputs para o PDF
}