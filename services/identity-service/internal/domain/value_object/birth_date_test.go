package valueobject_test

import (
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

func TestNovaDataNascimento(t *testing.T) {
	t.Parallel()
	hoje := civil.Date{Year: 2026, Month: time.November, Day: 2}
	for _, tc := range []struct {
		name string
		data civil.Date
	}{
		{"bissexto", civil.Date{Year: 2024, Month: time.February, Day: 29}},
		{"hoje_sem_idade_minima", hoje},
		{"ano_minimo", civil.Date{Year: 1, Month: time.January, Day: 1}},
		{"secular_bissexto", civil.Date{Year: 2000, Month: time.February, Day: 29}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := vo.NovaDataNascimento(tc.data.Year, tc.data.Month, tc.data.Day, hoje)
			if err != nil {
				t.Fatal(err)
			}
			if data.Data() != tc.data {
				t.Fatal("data alterada")
			}
			if err := data.ValidarEm(hoje); err != nil {
				t.Fatal(err)
			}
			copia := data.Data()
			copia.Day = 15
			if data.Data() != tc.data {
				t.Fatal("getter permite mutação")
			}
		})
	}
	max := civil.Date{Year: 9999, Month: time.December, Day: 31}
	if _, err := vo.NovaDataNascimento(9999, time.December, 31, max); err != nil {
		t.Fatal(err)
	}
}

func TestNovaDataNascimentoRejeitaDatasInvalidas(t *testing.T) {
	t.Parallel()
	hoje := civil.Date{Year: 2026, Month: time.November, Day: 2}
	for _, tc := range []struct {
		name string
		ano  int
		mes  time.Month
		dia  int
	}{
		{"zero", 0, 0, 0}, {"ano_zero", 0, time.January, 1},
		{"ano_negativo", -1, time.January, 1}, {"ano_excedido", 10000, time.January, 1},
		{"mes_zero", 2024, 0, 1}, {"mes_excedido", 2024, 13, 1},
		{"dia_zero", 2024, time.January, 0}, {"dia_negativo", 2024, time.January, -1},
		{"dia_excedido", 2024, time.April, 31}, {"nao_bissexto", 2023, time.February, 29},
		{"secular_nao_bissexto", 1900, time.February, 29}, {"amanha", 2026, time.November, 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := vo.NovaDataNascimento(tc.ano, tc.mes, tc.dia, hoje)
			if !errors.Is(err, vo.ErrValorInvalido) || data != (vo.DataNascimento{}) {
				t.Fatalf("esperado erro e valor zero: %v", err)
			}
		})
	}
}

func TestDataNascimentoValidarEm(t *testing.T) {
	t.Parallel()
	nascimento := civil.Date{Year: 2024, Month: time.February, Day: 29}
	data, err := vo.NovaDataNascimento(2024, time.February, 29, nascimento)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		hoje civil.Date
	}{
		{"referencia_zero", civil.Date{}},
		{"referencia_inexistente", civil.Date{Year: 2023, Month: time.February, Day: 29}},
		{"referencia_ano_zero", civil.Date{Year: 0, Month: time.January, Day: 1}},
		{"referencia_ano_excedido", civil.Date{Year: 10000, Month: time.January, Day: 1}},
		{"anterior_ao_nascimento", civil.Date{Year: 2024, Month: time.February, Day: 28}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if !errors.Is(data.ValidarEm(tc.hoje), vo.ErrValorInvalido) {
				t.Fatal("referência inválida aceita")
			}
			got, err := vo.NovaDataNascimento(2024, time.February, 29, tc.hoje)
			if !errors.Is(err, vo.ErrValorInvalido) || got != (vo.DataNascimento{}) {
				t.Fatalf("construtor aceitou referência: %v", err)
			}
			if data.Data() != nascimento {
				t.Fatal("validação alterou data")
			}
		})
	}
	if !errors.Is((vo.DataNascimento{}).ValidarEm(nascimento), vo.ErrValorInvalido) {
		t.Fatal("valor zero válido")
	}
}
