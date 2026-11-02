# identity-service

Identidades de acesso, cadastros de empresas e funcionários e vínculos de autorização do Teu Ponto.

Em desenvolvimento: módulo inicializado, com geração e validação de matrículas internas. Demais objetos do domínio, API e persistência seguem o [plano de implementação](../../docs/PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md).

## Matrículas internas

`internal/matriculagen.Gerar()` retorna `(vo.Matricula, error)` com seis caracteres de `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`, usando `crypto/rand.Reader` e descarte de amostras para evitar viés. Falhas de leitura preservam a causa via `errors.Is` e retornam matrícula zero, inválida.

`internal/domain/value_object.NovaMatricula(string)` representa também matrículas existentes, normaliza somente letras ASCII minúsculas e rejeita formato inválido com `ErrValorInvalido`. `String()` devolve o valor canônico; `Validar()` permite rejeitar o valor zero ao recebê-lo em construtores. O construtor futuro de Funcionário recebe a matrícula já gerada.

A geração não garante unicidade. A persistência futura deve impor restrição por CNPJ, incluindo cadastros encerrados; somente colisões de matrícula permitem nova geração, até dez tentativas por operação, sem cadastro parcial. A matrícula da folha permanece uma referência externa separada. Consulte os [contratos de domínio](../../docs/ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md#biblioteca-de-matrícula).

## Pré-requisitos

- Go 1.24.0 ou superior.

## Desenvolvimento

Execute em `services/identity-service/`.

### Testes

```bash
go test ./...
go test ./test
go test -v ./test -run '^TestCivil$'
go test ./internal/domain/value_object ./internal/matriculagen
go test -race ./...
```

### Build e análise estática

```bash
go build ./...
go vet ./...
```

### Formatação

```bash
gofmt -w .
```

### Dependências

```bash
go mod verify
go mod tidy -diff
```

Consulte a [avaliação de dependências](../../docs/DEPENDENCIAS_DOMINIO_IDENTITY_SERVICE.md) e a [arquitetura do backend](../docs/adr/0001-arquitetura-backend.md).
