package valueobject

import "fmt"

// ModalidadeColeta representa a atribuição da empresa, sem executar a coleta.
type ModalidadeColeta struct{ codigo string }

func NovaModalidadeColeta(codigo string) (ModalidadeColeta, error) {
	modalidade := ModalidadeColeta{codigo: codigo}
	if err := modalidade.Validar(); err != nil {
		return ModalidadeColeta{}, err
	}
	return modalidade, nil
}

func (m ModalidadeColeta) Codigo() string { return m.codigo }
func (m ModalidadeColeta) Validar() error {
	switch m.codigo {
	case "APARELHO_PROPRIO", "APARELHO_COLETIVO":
		return nil
	default:
		return fmt.Errorf("modalidade de coleta: código desconhecido: %w", ErrValorInvalido)
	}
}
