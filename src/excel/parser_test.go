package excel

import (
	"reflect"
	"testing"
)

func TestDescascarIF(t *testing.T) {
	tests := []struct {
		name    string
		formula string
		want    string
	}{
		{
			name:    "No IF formula",
			formula: "=SUM(A1:B2)",
			want:    "=SUM(A1:B2)",
		},
		{
			name:    "Proteção visual (empty true)",
			formula: `=IF(A1="", "", SUM(A1:B2))`,
			want:    "SUM(A1:B2)",
		},
		{
			name:    "Proteção visual (dash true)",
			formula: `=IF(A1="-", "-", SUM(A1:B2))`,
			want:    "SUM(A1:B2)",
		},
		{
			name:    "Proteção visual (empty false)",
			formula: `=IF(A1="", SUM(A1:B2), "")`,
			want:    "SUM(A1:B2)",
		},
		{
			name:    "Proteção visual (zero false)",
			formula: `=IF(A1="", SUM(A1:B2), 0)`,
			want:    "SUM(A1:B2)",
		},
		{
			name:    "Lógica de negócio real",
			formula: `=IF(A1>10, B1, C1)`,
			want:    `=IF(A1>10, B1, C1)`, // retém original com =
		},
		{
			name:    "Tratar sinais de dólar",
			formula: `=$IF($A$1="", "", SUM($B$1:$B$2))`,
			want:    `SUM(B1:B2)`,
		},
		{
			name:    "IF(cond, A, no_false)",
			formula: `IF(A1="", "")`,
			want:    `IF(A1="", "")`, // doesn't fit the clear replace paths fully, actually descascar returns formula since false is empty and true is placeholder -> would return descascar("")
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := descascarIF(tt.formula); got != tt.want {
				t.Errorf("descascarIF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapturarArgumentosIF(t *testing.T) {
	tests := []struct {
		name string
		f    string
		want []string
	}{
		{
			name: "Simple IF",
			f:    "IF(A1>10, 1, 0)",
			want: []string{"A1>10", "1", "0"},
		},
		{
			name: "Semicolon separated",
			f:    "IF(A1>10; 1; 0)",
			want: []string{"A1>10", "1", "0"},
		},
		{
			name: "Quotes with comma",
			f:    `IF(A1="O,K", "Yes", "No")`,
			want: []string{`A1="O,K"`, `"Yes"`, `"No"`},
		},
		{
			name: "Nested parenthesis",
			f:    `IF((A1+B1)>10, SUM(C1, C2), "Zero")`,
			want: []string{`(A1+B1)>10`, `SUM(C1, C2)`, `"Zero"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := capturarArgumentosIF(tt.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("capturarArgumentosIF() = %v, want %v", got, tt.want)
			}
		})
	}
}
