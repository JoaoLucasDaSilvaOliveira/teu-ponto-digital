package valueobject

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

// Email é um endereço simples. Validade sintática não comprova titularidade.
// A comparação na aplicação usa Chave, não a igualdade da apresentação.
type Email struct{ endereco string }

func NovoEmail(entrada string) (Email, error) {
	if strings.ContainsAny(entrada, "\r\n") {
		return Email{}, fmt.Errorf("email: quebra de linha: %w", ErrValorInvalido)
	}
	entrada = strings.TrimSpace(entrada)
	if err := validarEndereco(entrada); err != nil {
		return Email{}, err
	}
	separador := strings.LastIndexByte(entrada, '@')
	return Email{endereco: entrada[:separador+1] + strings.ToLower(entrada[separador+1:])}, nil
}

func (e Email) String() string { return e.endereco }
func (e Email) Chave() string  { return strings.ToLower(e.endereco) }
func (e Email) Validar() error { return validarEndereco(e.endereco) }

func validarEndereco(endereco string) error {
	parsed, err := mail.ParseAddress(endereco)
	// ParseAddress também aceita cabeçalhos com nome, comentários e delimitadores.
	// O valor de domínio recebe somente o endereço simples, sem essas formas.
	if err != nil || parsed.Name != "" || !enderecoSimples(endereco) {
		return fmt.Errorf("email: endereço simples obrigatório: %w", ErrValorInvalido)
	}
	return nil
}

// Mantém a grafia da parte local, inclusive aspas que ParseAddress remove.
// A gramática do endereço continua sendo validada pela biblioteca padrão.
func enderecoSimples(endereco string) bool {
	aspas, escape := false, false
	for _, r := range endereco {
		if unicode.IsControl(r) {
			return false
		}
		if escape {
			escape = false
			continue
		}
		if aspas && r == '\\' {
			escape = true
			continue
		}
		if r == '"' {
			aspas = !aspas
			continue
		}
		if !aspas && (unicode.IsSpace(r) || strings.ContainsRune("<>(),;", r)) {
			return false
		}
	}
	return !aspas && !escape
}
