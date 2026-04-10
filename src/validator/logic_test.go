package validator

import (
	"math"
	"reflect"
	"testing"
	"validador/src/excel"
)

func TestParseLocalFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"Empty string", "", 0, false},
		{"Dot decimal", "1.23", 1.23, false},
		{"Comma decimal", "1,23", 1.23, false},
		{"Multiple decimal points", "1.23.45", 0, true},
		{"Negative value", "-5,55", -5.55, false},
		{"Invalid value", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLocalFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLocalFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("parseLocalFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidar(t *testing.T) {
	tests := []struct {
		name     string
		input    []excel.DetalhesFormula
		expected []ValidacaoFormula
	}{
		{
			name: "Pass validation with exact match",
			input: []excel.DetalhesFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "SUM(1,1)", ValorSalvo: "2", Recalculado: "2", Inputs: "1, 1"},
			},
			expected: []ValidacaoFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "SUM(1,1)", Original: "2", Calculado: "2", Status: "PASSOU", Inputs: "1, 1"},
			},
		},
		{
			name: "Pass validation within tolerance",
			input: []excel.DetalhesFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "1/3", ValorSalvo: "0.3333", Recalculado: "0.33333333"},
			},
			expected: []ValidacaoFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "1/3", Original: "0.3333", Calculado: "0.33333333", Status: "PASSOU"},
			},
		},
		{
			name: "Fail validation outside tolerance",
			input: []excel.DetalhesFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "2+2", ValorSalvo: "4", Recalculado: "5"},
			},
			expected: []ValidacaoFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "2+2", Original: "4", Calculado: "5", Status: "FALHA"},
			},
		},
		{
			name: "Manual review due to excel error",
			input: []excel.DetalhesFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "1/0", ValorSalvo: "0", Recalculado: "#DIV/0!"},
			},
			expected: []ValidacaoFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "1/0", Original: "0", Calculado: "#DIV/0!", Status: "REVISÃO MANUAL"},
			},
		},
		{
			name: "Textual equivalence of zeroes",
			input: []excel.DetalhesFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "0", ValorSalvo: "0.0000", Recalculado: ""},
			},
			expected: []ValidacaoFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "0", Original: "0.0000", Calculado: "", Status: "PASSOU"},
			},
		},
		{
			name: "AST Empty calculation with IF",
			input: []excel.DetalhesFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "IF(A2=1, 2, 3)", ValorSalvo: "2", Recalculado: ""},
			},
			expected: []ValidacaoFormula{
				{Aba: "Aba1", Referencia: "A1", Formula: "IF(A2=1, 2, 3)", Original: "2", Calculado: "", Status: "REVISÃO MANUAL"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := Validar(tt.input)
			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("Validar() = %+v, want %+v", results, tt.expected)
			}
		})
	}
}
