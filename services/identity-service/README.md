# identity-service

Identidades de acesso, cadastros de empresas e funcionários e vínculos de autorização do Teu Ponto.

Objetos de valor implementados: e-mail, nome, data de nascimento, perfil e modalidade de coleta. As demais entregas seguem o [plano de implementação](../../docs/PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md).

Os contratos dos valores estão na [especificação do domínio](../../docs/ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md) e na [autorização](../../docs/AUTORIZACAO_IDENTITY_SERVICE.md). E-mail usa `Chave()` para comparação sem distinção de caixa e `String()` para apresentação; sua construção não confirma titularidade. Nascimento recebe uma referência `civil.Date` explícita e não consulta o relógio.

## Pré-requisitos

- Go 1.24.0 ou superior.

## Desenvolvimento

Execute em `services/identity-service/`.

### Testes

```bash
go test ./...
go test ./test
go test -race ./...
go test -cover ./internal/domain/value_object
go test ./internal/domain/value_object
go test -v ./test -run '^TestCivil$'
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
