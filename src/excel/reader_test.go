package excel

import (
	"reflect"
	"testing"
)

func TestExpandirIntervalo(t *testing.T) {
	tests := []struct {
		name    string
		inicio  string
		fim     string
		want    []string
		wantErr bool
	}{
		{
			name:    "Single column",
			inicio:  "A1",
			fim:     "A3",
			want:    []string{"A1", "A2", "A3"},
			wantErr: false,
		},
		{
			name:    "Single row",
			inicio:  "A1",
			fim:     "C1",
			want:    []string{"A1", "B1", "C1"},
			wantErr: false,
		},
		{
			name:    "Multiple columns and rows",
			inicio:  "A1",
			fim:     "B2",
			want:    []string{"A1", "A2", "B1", "B2"},
			wantErr: false,
		},
		{
			name:    "Inverted start and end (rows)",
			inicio:  "A3",
			fim:     "A1",
			want:    []string{"A1", "A2", "A3"},
			wantErr: false,
		},
		{
			name:    "Inverted start and end (cols and rows)",
			inicio:  "B2",
			fim:     "A1",
			want:    []string{"A1", "A2", "B1", "B2"},
			wantErr: false,
		},
		{
			name:    "Invalid cell start",
			inicio:  "A",
			fim:     "A3",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid cell end",
			inicio:  "A1",
			fim:     "3",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandirIntervalo(tt.inicio, tt.fim)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandirIntervalo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expandirIntervalo() = %v, want %v", got, tt.want)
			}
		})
	}
}
