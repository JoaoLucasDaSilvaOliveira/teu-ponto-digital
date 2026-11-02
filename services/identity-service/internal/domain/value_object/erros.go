package value_object

import "errors"

// ErrValorInvalido identifica valores que não satisfazem as regras de domínio.
var ErrValorInvalido = errors.New("valor inválido")
