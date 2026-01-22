package locations

import (
	"testing"
)

func TestNormalizeAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple address",
			input:    "Rua Teste, 123",
			expected: "rua teste 123",
		},
		{
			name:     "with trailing comma",
			input:    "Avenida Paulista, 1000,",
			expected: "avenida paulista 1000",
		},
		{
			name:     "abbreviated street name",
			input:    "R. das Flores, 45",
			expected: "rua das flores 45",
		},
		{
			name:     "abbreviated avenue",
			input:    "Av. Brasil, 200",
			expected: "avenida brasil 200",
		},
		{
			name:     "multiple abbreviations",
			input:    "R. Dr. Silva, Av. das Nações",
			expected: "rua doutor silva avenida das nações",
		},
		{
			name:     "uppercase and lowercase mix",
			input:    "RUA TESTE, 123",
			expected: "rua teste 123",
		},
		{
			name:     "extra whitespace",
			input:    "  Rua Teste   ,  123  ",
			expected: "rua teste 123",
		},
		{
			name:     "saint abbreviation",
			input:    "R. Sta. Maria, 50",
			expected: "rua santa maria 50",
		},
		{
			name:     "doctor abbreviation",
			input:    "Dr. João Silva",
			expected: "doutor joão silva",
		},
		{
			name:     "doutora abbreviation",
			input:    "Dra. Ana Costa",
			expected: "doutora ana costa",
		},
		{
			name:     "park abbreviation",
			input:    "Pq. Ibirapuera",
			expected: "parque ibirapuera",
		},
		{
			name:     "square abbreviation",
			input:    "Pç. da Sé",
			expected: "praca da sé",
		},
		{
			name:     "neighborhood abbreviation",
			input:    "Jd. Botânico",
			expected: "jardim botânico",
		},
		{
			name:     "village abbreviation",
			input:    "Vl. Madalena",
			expected: "vila madalena",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAddress(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeAddress(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
