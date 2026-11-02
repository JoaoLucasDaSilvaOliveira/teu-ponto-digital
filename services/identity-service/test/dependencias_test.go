package identity_test

import (
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

var errFonte = errors.New("fonte indisponível para avaliação")

type fonteComFalha struct{}

func (fonteComFalha) Read([]byte) (int, error) { return 0, errFonte }

func TestUUID(t *testing.T) {
	t.Parallel()
	t.Run("geracao", func(t *testing.T) {
		value, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("gerar uuid: %v", err)
		}
		if value == uuid.Nil || value.Version() != 4 || value.Variant() != uuid.RFC4122 {
			t.Fatal("UUID gerado fora do formato esperado")
		}
	})
	t.Run("falha_da_fonte", func(t *testing.T) {
		_, err := uuid.NewRandomFromReader(fonteComFalha{})
		if !errors.Is(err, errFonte) {
			t.Fatalf("causa não preservada: %v", err)
		}
	})
	t.Run("parse_nao_impoe_nil", func(t *testing.T) {
		value, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
		if err != nil || value != uuid.Nil {
			t.Fatal("resultado inesperado ao interpretar UUID nulo")
		}
		t.Log("LIMITAÇÃO: domínio deve rejeitar uuid.Nil explicitamente")
	})
}

func TestCivil(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		date  civil.Date
		valid bool
	}{
		{"bissexto", civil.Date{Year: 2024, Month: time.February, Day: 29}, true},
		{"dia_inexistente", civil.Date{Year: 2023, Month: time.February, Day: 29}, false},
		{"valor_zero", civil.Date{}, false},
		{"ano_zero", civil.Date{Year: 0, Month: time.January, Day: 1}, true},
		{"ano_acima_9999", civil.Date{Year: 10000, Month: time.January, Day: 1}, true},
		{"data_futura_para_nascimento", civil.Date{Year: 2099, Month: time.January, Day: 1}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.date.IsValid() != tc.valid {
				t.Fatal("validade civil diferente da esperada")
			}
		})
	}
	t.Run("somente_data_sem_fuso", func(t *testing.T) {
		date, err := civil.ParseDate("2024-02-29")
		if err != nil || date.String() != "2024-02-29" {
			t.Fatalf("round trip inválido: %v", err)
		}
	})
}
