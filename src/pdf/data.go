package pdf

import "validador/src/validator"

func processarDados(resultados []validator.ValidacaoFormula) DadosRelatorio {
	resumoPorAba := make(map[string]*ResumoAba)                 
	dadosPorAba := make(map[string][]validator.ValidacaoFormula)
	var abasOrdem []string
	var totalGeral ResumoAba

	for _, r := range resultados {
		if _, existe := resumoPorAba[r.Aba]; !existe {
			abasOrdem = append(abasOrdem, r.Aba)
			resumoPorAba[r.Aba] = &ResumoAba{}        
			dadosPorAba[r.Aba] = []validator.ValidacaoFormula{}
		}

		dadosPorAba[r.Aba] = append(dadosPorAba[r.Aba], r)

		res := resumoPorAba[r.Aba]
		res.Total++
		totalGeral.Total++

		switch r.Status {
		case "PASSOU":
			res.Passou++
			totalGeral.Passou++
		case "FALHA":
			res.Falha++
			totalGeral.Falha++
		case "REVISÃO MANUAL":
			res.Revisao++
			totalGeral.Revisao++
		}
	}

	return DadosRelatorio{
		ResumoPorAba: resumoPorAba,
		DadosPorAba:  dadosPorAba,
		AbasOrdem:    abasOrdem,
		TotalGeral:   totalGeral,
	}
}