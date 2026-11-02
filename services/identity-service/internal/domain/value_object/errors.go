// Package valueobject contém valores imutáveis do domínio de identidades.
package valueobject

import "errors"

// ErrValorInvalido identifica uma violação das invariantes de um valor.
var ErrValorInvalido = errors.New("valor inválido")
