package valueobject_test

import (
	"errors"
	"testing"

	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

func TestNovoPerfil(t *testing.T) {
	t.Parallel()
	for _, codigo := range []string{"GESTOR", "SUPERVISOR", "COLABORADOR", "CONTABILIDADE_DP"} {
		t.Run(codigo, func(t *testing.T) {
			t.Parallel()
			perfil, err := vo.NovoPerfil(codigo)
			if err != nil {
				t.Fatal(err)
			}
			if perfil.Codigo() != codigo {
				t.Fatal("código alterado")
			}
			if err := perfil.Validar(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNovoPerfilRejeitaCodigosDesconhecidos(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, codigo string }{
		{"zero", ""}, {"master_inexistente", "GESTOR_MASTER"}, {"rh_inexistente", "RH"},
		{"minusculo", "gestor"}, {"espaco", " GESTOR"}, {"multiplos", "GESTOR,COLABORADOR"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			perfil, err := vo.NovoPerfil(tc.codigo)
			if !errors.Is(err, vo.ErrValorInvalido) || perfil != (vo.Perfil{}) {
				t.Fatalf("esperado erro e valor zero: %v", err)
			}
		})
	}
	if !errors.Is((vo.Perfil{}).Validar(), vo.ErrValorInvalido) {
		t.Fatal("valor zero válido")
	}
}
