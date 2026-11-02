package valueobject

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

// DataNascimento representa uma data civil, sem instante ou fuso horário.
type DataNascimento struct{ data civil.Date }

func NovaDataNascimento(ano int, mes time.Month, dia int, hoje civil.Date) (DataNascimento, error) {
	nascimento := DataNascimento{data: civil.Date{Year: ano, Month: mes, Day: dia}}
	if err := nascimento.ValidarEm(hoje); err != nil {
		return DataNascimento{}, err
	}
	return nascimento, nil
}

func (d DataNascimento) Data() civil.Date { return d.data }

// ValidarEm usa a data de referência recebida; nunca consulta o relógio global.
func (d DataNascimento) ValidarEm(hoje civil.Date) error {
	if !dataCivilValida(hoje) {
		return fmt.Errorf("nascimento: referência inválida: %w", ErrValorInvalido)
	}
	if !dataCivilValida(d.data) || d.data.After(hoje) {
		return fmt.Errorf("nascimento: data inválida ou futura: %w", ErrValorInvalido)
	}
	return nil
}

func dataCivilValida(data civil.Date) bool {
	return data.Year >= 1 && data.Year <= 9999 && data.IsValid()
}
