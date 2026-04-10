package pdf

import "validador/src/validator"

type ResumoAba struct {
	Total   int
	Passou  int
	Falha   int
	Revisao int
}

type DadosRelatorio struct {
	ResumoPorAba map[string]*ResumoAba                  // CORRIGIDO
	DadosPorAba  map[string][]validator.ValidacaoFormula // CORRIGIDO
	AbasOrdem    []string
	TotalGeral   ResumoAba
}