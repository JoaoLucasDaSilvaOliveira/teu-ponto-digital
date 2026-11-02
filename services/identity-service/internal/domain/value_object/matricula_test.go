package valueobject_test

import (
	"errors"
	"testing"

	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

func TestNovaMatricula(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ nome, entrada, esperada string }{
		{"existente", "0AZ9BY", "0AZ9BY"},
		{"minusculas ASCII", "a0z9bY", "A0Z9BY"},
		{"zeros preservados", "000000", "000000"},
		{"letras", "ZZZZZZ", "ZZZZZZ"},
	} {
		t.Run(tc.nome, func(t *testing.T) {
			t.Parallel()
			m, err := vo.NovaMatricula(tc.entrada)
			if err != nil {
				t.Fatal(err)
			}
			if m.String() != tc.esperada {
				t.Fatalf("String() = %q, esperado %q", m.String(), tc.esperada)
			}
			if err := m.Validar(); err != nil {
				t.Fatal(err)
			}
			existente, err := vo.NovaMatricula(tc.esperada)
			if err != nil {
				t.Fatal(err)
			}
			if m != existente {
				t.Fatal("matrículas de mesmo conteúdo devem ser iguais")
			}
		})
	}
}

func TestNovaMatriculaRejeitaFormatoInvalido(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ nome, entrada string }{
		{"vazia", ""},
		{"curta", "12345"},
		{"longa", "1234567"},
		{"espaco externo", " ABC123"},
		{"espaco interno", "ABC 12"},
		{"espaco final", "ABC12 "},
		{"hifen", "ABC-12"},
		{"controle", "ABC12\n"},
		{"nao ASCII", "ſ1234"},
		{"acentuada", "Á1234"},
		{"utf8 invalido", "ABCDE\xff"},
	} {
		t.Run(tc.nome, func(t *testing.T) {
			t.Parallel()
			m, err := vo.NovaMatricula(tc.entrada)
			if !errors.Is(err, vo.ErrValorInvalido) {
				t.Fatalf("erro = %v, esperado ErrValorInvalido", err)
			}
			if m != (vo.Matricula{}) {
				t.Fatal("entrada inválida retornou matrícula parcial")
			}
		})
	}
}

func TestMatriculaZeroInvalida(t *testing.T) {
	t.Parallel()
	var m vo.Matricula
	if err := m.Validar(); !errors.Is(err, vo.ErrValorInvalido) {
		t.Fatalf("erro = %v, esperado ErrValorInvalido", err)
	}
}
