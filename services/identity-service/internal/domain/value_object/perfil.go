package valueobject

import "fmt"

// Perfil identifica um dos quatro perfis aprovados; não concede ações sozinho.
type Perfil struct{ codigo string }

func NovoPerfil(codigo string) (Perfil, error) {
	perfil := Perfil{codigo: codigo}
	if err := perfil.Validar(); err != nil {
		return Perfil{}, err
	}
	return perfil, nil
}

func (p Perfil) Codigo() string { return p.codigo }
func (p Perfil) Validar() error {
	switch p.codigo {
	case "GESTOR", "SUPERVISOR", "COLABORADOR", "CONTABILIDADE_DP":
		return nil
	default:
		return fmt.Errorf("perfil: código desconhecido: %w", ErrValorInvalido)
	}
}
