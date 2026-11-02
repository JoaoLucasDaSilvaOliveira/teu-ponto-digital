package valueobject_test

import (
	"errors"
	"strings"
	"testing"

	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

func TestNovoEmail(t *testing.T) {
	t.Parallel()
	email, err := vo.NovoEmail("  Ana.Silva+Ponto@EXAMPLE.COM  ")
	if err != nil {
		t.Fatal(err)
	}
	if email.String() != "Ana.Silva+Ponto@example.com" || email.Chave() != "ana.silva+ponto@example.com" {
		t.Fatal("normalização incorreta")
	}
	if err := email.Validar(); err != nil {
		t.Fatal(err)
	}
}

func TestNovoEmailRejeitaEntradasInvalidas(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, entrada string }{
		{"vazio", ""}, {"espacos", "   "}, {"sem_dominio", "ana@"},
		{"sem_local", "@example.com"}, {"sem_arroba", "ana.example.com"},
		{"nome", "Ana <ana@example.com>"}, {"delimitadores", "<ana@example.com>"},
		{"lista", "ana@example.com,bia@example.com"}, {"comentario", "ana@example.com (Ana)"},
		{"crlf", "ana@example.com\r\nBcc: bia@example.com"},
		{"quebra_externa", "\nana@example.com\r"}, {"espaco_interno", "ana @example.com"},
		{"controle_na_parte_local", "\"ana\t\"@example.com"},
		{"grupo", "equipe: ana@example.com;"},
		{"nome_vazio", `"" <ana@example.com>`},
		{"utf8_invalido", "ana\xff@example.com"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := vo.NovoEmail(tc.entrada)
			if !errors.Is(err, vo.ErrValorInvalido) || got != (vo.Email{}) {
				t.Fatalf("esperado erro e valor zero: %v", err)
			}
			if tc.entrada != "" && strings.Contains(err.Error(), tc.entrada) {
				t.Fatal("erro expõe entrada")
			}
		})
	}
	if !errors.Is((vo.Email{}).Validar(), vo.ErrValorInvalido) {
		t.Fatal("valor zero válido")
	}
}

func TestEmailChave(t *testing.T) {
	t.Parallel()
	base, err := vo.NovoEmail("Ana.Silva+Ponto@EXAMPLE.COM")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, entrada string
		igual         bool
	}{
		{"caixa", "ana.silva+ponto@example.com", true},
		{"pontos_preservados", "AnaSilva+Ponto@example.com", false},
		{"sufixo_preservado", "Ana.Silva@example.com", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			outro, err := vo.NovoEmail(tc.entrada)
			if err != nil {
				t.Fatal(err)
			}
			if (base.Chave() == outro.Chave()) != tc.igual {
				t.Fatal("chave incorreta")
			}
		})
	}
}

func TestNovoEmailPreservaLocalEntreAspas(t *testing.T) {
	t.Parallel()
	email, err := vo.NovoEmail(`"Ana Silva+Ponto"@EXAMPLE.COM`)
	if err != nil {
		t.Fatal(err)
	}
	if email.String() != `"Ana Silva+Ponto"@example.com` || email.Chave() != `"ana silva+ponto"@example.com` {
		t.Fatal("parte local alterada")
	}
	if err := email.Validar(); err != nil {
		t.Fatal(err)
	}
}

func TestNovoEmailLocalComCaracteresEscapados(t *testing.T) {
	t.Parallel()
	entrada := `"Ana\"(Ponto)"@EXAMPLE.COM`
	email, err := vo.NovoEmail(entrada)
	if err != nil {
		t.Fatal(err)
	}
	if email.String() != `"Ana\"(Ponto)"@example.com` {
		t.Fatal("grafia da parte local alterada")
	}
}
