package pdf

import (
	"reflect"
	"testing"
	"validador/src/validator"
)

func TestProcessarDados(t *testing.T) {
	tests := []struct {
		name       string
		resultados []validator.ValidacaoFormula
		wantAbas   []string
		wantTotal  int
		wantPassou int
		wantFalha  int
		wantRev    int
	}{
		{
			name:       "Empty list",
			resultados: []validator.ValidacaoFormula{},
			wantAbas:   nil,
			wantTotal:  0,
		},
		{
			name: "Mixed results",
			resultados: []validator.ValidacaoFormula{
				{Aba: "Aba1", Status: "PASSOU"},
				{Aba: "Aba1", Status: "FALHA"},
				{Aba: "Aba2", Status: "PASSOU"},
				{Aba: "Aba3", Status: "REVISÃO MANUAL"},
			},
			wantAbas:   []string{"Aba1", "Aba2", "Aba3"},
			wantTotal:  4,
			wantPassou: 2,
			wantFalha:  1,
			wantRev:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processarDados(tt.resultados)
			
			if len(got.AbasOrdem) == 0 && len(tt.wantAbas) == 0 {
				// both empty, OK
			} else if !reflect.DeepEqual(got.AbasOrdem, tt.wantAbas) {
				t.Errorf("processarDados().AbasOrdem = %v, want %v", got.AbasOrdem, tt.wantAbas)
			}
			if got.TotalGeral.Total != tt.wantTotal {
				t.Errorf("processarDados().TotalGeral.Total = %v, want %v", got.TotalGeral.Total, tt.wantTotal)
			}
			if got.TotalGeral.Passou != tt.wantPassou {
				t.Errorf("processarDados().TotalGeral.Passou = %v, want %v", got.TotalGeral.Passou, tt.wantPassou)
			}
			if got.TotalGeral.Falha != tt.wantFalha {
				t.Errorf("processarDados().TotalGeral.Falha = %v, want %v", got.TotalGeral.Falha, tt.wantFalha)
			}
			if got.TotalGeral.Revisao != tt.wantRev {
				t.Errorf("processarDados().TotalGeral.Revisao = %v, want %v", got.TotalGeral.Revisao, tt.wantRev)
			}
			
			for _, aba := range tt.wantAbas {
				if _, ok := got.ResumoPorAba[aba]; !ok {
					t.Errorf("processarDados().ResumoPorAba missing key %v", aba)
				}
				if _, ok := got.DadosPorAba[aba]; !ok {
					t.Errorf("processarDados().DadosPorAba missing key %v", aba)
				}
			}
		})
	}
}
