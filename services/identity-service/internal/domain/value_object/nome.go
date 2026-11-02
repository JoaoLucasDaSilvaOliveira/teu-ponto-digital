package valueobject

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Nome preserva a grafia de nomes pessoais e identificações da empresa/setor.
type Nome struct{ texto string }

func NovoNome(entrada string) (Nome, error) {
	// Controles são rejeitados antes de aparar espaços, inclusive nas extremidades.
	if !utf8.ValidString(entrada) || strings.ContainsFunc(entrada, unicode.IsControl) {
		return Nome{}, fmt.Errorf("nome: texto inválido: %w", ErrValorInvalido)
	}
	nome := Nome{texto: strings.TrimSpace(entrada)}
	if err := nome.Validar(); err != nil {
		return Nome{}, err
	}
	return nome, nil
}

func (n Nome) String() string { return n.texto }
func (n Nome) Validar() error {
	if n.texto == "" || !utf8.ValidString(n.texto) || strings.ContainsFunc(n.texto, unicode.IsControl) {
		return fmt.Errorf("nome: texto obrigatório sem controles: %w", ErrValorInvalido)
	}
	return nil
}
