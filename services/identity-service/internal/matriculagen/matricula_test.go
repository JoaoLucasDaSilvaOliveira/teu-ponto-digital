package matriculagen

import (
	"bytes"
	"errors"
	"io"
	"testing"

	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

func TestGerarComFonte(t *testing.T) {
	t.Parallel()
	m, err := gerar(bytes.NewReader([]byte{0, 9, 10, 35, 36, 251}))
	if err != nil {
		t.Fatal(err)
	}
	if m.String() != "09AZ0Z" {
		t.Fatalf("matrícula = %q, esperado 09AZ0Z", m.String())
	}
	if err := m.Validar(); err != nil {
		t.Fatal(err)
	}
}

func TestGerarPropagaFalhaSemValorParcial(t *testing.T) {
	t.Parallel()
	falha := errors.New("fonte indisponível")
	for _, tc := range []struct {
		nome  string
		fonte io.Reader
		causa error
	}{
		{"imediata", leitorComErro{falha}, falha},
		{"apos prefixo", io.MultiReader(bytes.NewReader([]byte{0, 9, 35}), leitorComErro{falha}), falha},
		{"apos rejeicao", io.MultiReader(bytes.NewReader([]byte{252, 255}), leitorComErro{falha}), falha},
		{"fonte esgotada", bytes.NewReader([]byte{0, 1}), io.EOF},
	} {
		t.Run(tc.nome, func(t *testing.T) {
			t.Parallel()
			m, err := gerar(tc.fonte)
			if !errors.Is(err, tc.causa) {
				t.Fatalf("erro = %v, esperado causa %v", err, tc.causa)
			}
			if err.Error() == tc.causa.Error() {
				t.Fatal("erro deve incluir contexto de geração")
			}
			if m != (vo.Matricula{}) {
				t.Fatal("falha retornou matrícula parcial")
			}
		})
	}
}

type leitorComErro struct{ err error }

func (l leitorComErro) Read([]byte) (int, error) { return 0, l.err }

func TestGerar(t *testing.T) {
	t.Parallel()
	m, err := Gerar()
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validar(); err != nil {
		t.Fatalf("gerador padrão retornou matrícula inválida: %v", err)
	}
}

func TestGerarDescartaAmostrasComVies(t *testing.T) {
	t.Parallel()
	m, err := gerar(bytes.NewReader([]byte{252, 253, 254, 255, 0, 9, 10, 35, 36, 251}))
	if err != nil {
		t.Fatal(err)
	}
	if m.String() != "09AZ0Z" {
		t.Fatalf("matrícula = %q, esperado 09AZ0Z", m.String())
	}
}
