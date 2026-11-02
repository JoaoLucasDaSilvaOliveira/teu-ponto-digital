# identity-service

Módulo Go das identidades de acesso, cadastros de empresas e funcionários e vínculos de autorização do Teu Ponto. Inicialização da [issue #7](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/7), conforme o [plano aprovado](../../docs/PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md) e a [arquitetura do backend](../docs/adr/0001-arquitetura-backend.md).

## Estado atual

O módulo é `github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/services/identity-service`, com Go 1.24.0 ou superior. As dependências diretas fixadas são `github.com/google/uuid v1.6.0` e `cloud.google.com/go v0.123.0` (somente o pacote de valores `civil`). Versões e hashes estão em [go.mod](go.mod) e [go.sum](go.sum); a [avaliação de dependências](../../docs/DEPENDENCIAS_DOMINIO_IDENTITY_SERVICE.md) registra as limitações observadas.

O único pacote atual é `test/`, com testes de caracterização de UUID e data civil. Esses testes não implementam as invariantes dos futuros objetos de valor. Os pacotes de domínio serão criados nas respectivas issues, sem dependências de HTTP, ORM ou broker.

O bootstrap vazio originalmente em `services/cmd/main/main.go` foi movido para este módulo na rodada anterior. Na #7, `cmd/main/main.go`, que continha apenas `package main`, foi removido: passava em teste e análise estática, mas falhava no build por não possuir função `main`. Não há executável nesta etapa. Quando houver comportamento de aplicação executável, seu bootstrap será criado em `cmd/main/main.go`, compondo dependências concretas conforme o ADR. API, persistência e autenticação completa permanecem nas etapas posteriores.

## Comandos

Execute a partir da raiz deste repositório:

```bash
cd services/identity-service
```

Todos os comandos seguintes partem dessa pasta. Não há módulo ou workspace Go na raiz do repositório.

```bash
# Executar todos os testes, sem reutilizar resultados em cache
go test -count=1 ./...

# Executar somente os testes do pacote atual
go test -count=1 ./test

# Selecionar um teste e exibir seus cenários
go test -count=1 -v ./test -run '^TestCivil$'

# Análise estática e compilação dos pacotes existentes
go vet ./...
go build ./...

# Verificar formatação (saída vazia significa arquivos formatados)
gofmt -l .

# Aplicar formatação
gofmt -w .

# Conferir integridade e consistência das dependências
go mod verify
go mod tidy -diff
```

`./...` inclui recursivamente os pacotes do módulo, inclusive os futuros pacotes de domínio. Para testar um pacote específico, substitua `./test` por seu caminho relativo ao módulo. `go build ./...` verifica os pacotes existentes; não gera um servidor. Não há comando `go run` nesta etapa. Para atualizar os manifestos após uma alteração intencional de imports, execute `go mod tidy` e revise o diff.

## Verificação da inicialização

Em 02/11/2026, no ambiente `go version go1.27.1-X:nodwarf5 linux/amd64`, passaram os três comandos de teste acima, `go vet ./...`, `go build ./...`, `go mod verify` e `go mod tidy -diff`; `gofmt -l .` não listou arquivos. A compilação foi verificada após a remoção do bootstrap vazio. O mínimo declarado é Go 1.24.0, conforme a dependência de datas; esta rodada não executou verificações nesse toolchain nem em outras plataformas. Revisão e homologação continuam pendentes.
