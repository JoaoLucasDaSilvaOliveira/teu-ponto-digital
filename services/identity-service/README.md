# identity-service

Identidades de acesso, cadastros de empresas e funcionários e vínculos de autorização do Teu Ponto.

Em desenvolvimento: módulo inicializado com testes de dependências. Domínio, API e persistência seguem o [plano de implementação](../../docs/PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md).

## Pré-requisitos

- Go 1.24.0 ou superior.

## Desenvolvimento

Execute em `services/identity-service/`.

### Testes

```bash
go test ./...
go test ./test
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
