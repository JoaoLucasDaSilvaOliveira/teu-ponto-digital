package valueobject_test

import (
	"errors"
	"testing"

	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

func TestNovaModalidadeColeta(t *testing.T) {
	t.Parallel()
	for _, codigo := range []string{"APARELHO_PROPRIO", "APARELHO_COLETIVO"} {
		t.Run(codigo, func(t *testing.T) {
			t.Parallel()
			modalidade, err := vo.NovaModalidadeColeta(codigo)
			if err != nil {
				t.Fatal(err)
			}
			if modalidade.Codigo() != codigo {
				t.Fatal("código alterado")
			}
			if err := modalidade.Validar(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNovaModalidadeColetaRejeitaCodigosDesconhecidos(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, codigo string }{
		{"zero", ""}, {"desconhecido", "FACIAL"}, {"minusculo", "aparelho_proprio"},
		{"espaco", "APARELHO_COLETIVO "}, {"multiplos", "APARELHO_PROPRIO,APARELHO_COLETIVO"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			modalidade, err := vo.NovaModalidadeColeta(tc.codigo)
			if !errors.Is(err, vo.ErrValorInvalido) || modalidade != (vo.ModalidadeColeta{}) {
				t.Fatalf("esperado erro e valor zero: %v", err)
			}
		})
	}
	if !errors.Is((vo.ModalidadeColeta{}).Validar(), vo.ErrValorInvalido) {
		t.Fatal("valor zero válido")
	}
}
