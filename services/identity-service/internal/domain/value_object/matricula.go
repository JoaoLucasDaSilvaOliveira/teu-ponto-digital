package value_object

import "fmt"

// Matricula identifica um cadastro de funcionário, sem representar a matrícula da folha.
// O valor zero é inválido. Unicidade por CNPJ exige persistência e coordenação.
type Matricula struct {
	valor string
}

// NovaMatricula representa uma matrícula existente e normaliza somente letras ASCII.
func NovaMatricula(entrada string) (Matricula, error) {
	valor := []byte(entrada)
	for i, c := range valor {
		if c >= 'a' && c <= 'z' {
			valor[i] = c - ('a' - 'A')
		}
	}
	m := Matricula{valor: string(valor)}
	if err := m.Validar(); err != nil {
		return Matricula{}, err
	}
	return m, nil
}

func (m Matricula) String() string { return m.valor }

// Validar verifica o formato canônico, inclusive em valores não construídos.
func (m Matricula) Validar() error {
	if len(m.valor) != 6 {
		return fmt.Errorf("matrícula deve conter seis caracteres: %w", ErrValorInvalido)
	}
	for i := range len(m.valor) {
		c := m.valor[i]
		if !(c >= '0' && c <= '9' || c >= 'A' && c <= 'Z') {
			return fmt.Errorf("matrícula deve conter somente dígitos e letras ASCII maiúsculas: %w", ErrValorInvalido)
		}
	}
	return nil
}
