package valueobject_test

import (
	"errors"
	"testing"

	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

func TestNovoNome(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, entrada, esperado string }{
		{"acentos", "  João da Conceição  ", "João da Conceição"},
		{"sem_sobrenome", "Ana", "Ana"}, {"grafia", "Setor  A - Açúcar", "Setor  A - Açúcar"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			nome, err := vo.NovoNome(tc.entrada)
			if err != nil {
				t.Fatal(err)
			}
			if nome.String() != tc.esperado {
				t.Fatalf("nome = %q", nome.String())
			}
			if err := nome.Validar(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNovoNomeRejeitaEntradasInvalidas(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, entrada string }{
		{"vazio", ""}, {"espacos", "  "}, {"utf8", "Jo\xffão"},
		{"controle", "Ana\x00"}, {"tab", "\tAna"}, {"quebra", "Ana\n"}, {"del", "Ana\x7f"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			nome, err := vo.NovoNome(tc.entrada)
			if !errors.Is(err, vo.ErrValorInvalido) || nome != (vo.Nome{}) {
				t.Fatalf("esperado erro e valor zero: %v", err)
			}
		})
	}
	if !errors.Is((vo.Nome{}).Validar(), vo.ErrValorInvalido) {
		t.Fatal("valor zero válido")
	}
}
