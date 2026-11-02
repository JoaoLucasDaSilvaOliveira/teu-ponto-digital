// Package matriculagen gera matrículas internas com aleatoriedade criptográfica.
// Não consulta cadastros nem garante unicidade: a persistência futura deve reservar
// a matrícula por CNPJ, inclusive no histórico, e coordenar colisões.
package matriculagen

import (
	"crypto/rand"
	"fmt"
	"io"

	vo "github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service/internal/domain/value_object"
)

// Gerar retorna uma matrícula pronta para o construtor do funcionário.
// Erros da fonte são propagados preservando a causa, sem retornar valor parcial.
func Gerar() (vo.Matricula, error) {
	return gerar(rand.Reader)
}

// gerar mantém a fonte controlável sem expor configuração insegura ao consumidor.
func gerar(fonte io.Reader) (vo.Matricula, error) {
	const alfabeto = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 252 é o maior múltiplo de 36 abaixo de 256: sete amostras por símbolo.
	const limite = 256 - 256%len(alfabeto)
	var valor [6]byte
	var amostra [1]byte
	for i := 0; i < len(valor); {
		if _, err := io.ReadFull(fonte, amostra[:]); err != nil {
			return vo.Matricula{}, fmt.Errorf("gerar matrícula: obter aleatoriedade: %w", err)
		}
		if int(amostra[0]) >= limite {
			continue
		}
		valor[i] = alfabeto[int(amostra[0])%len(alfabeto)]
		i++
	}
	return vo.NovaMatricula(string(valor[:]))
}
