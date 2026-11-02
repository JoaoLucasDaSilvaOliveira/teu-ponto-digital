# Dependências do domínio — identity-service

Avaliação da [issue #6](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/6), realizada em 02/11/2026. Por orientação do usuário, o módulo foi inicializado diretamente em [services/identity-service](../services/identity-service/go.mod), antecipando essa parte da #7. Fontes internas: [plano](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md), [contratos](ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md), [arquitetura do backend](../services/docs/adr/0001-arquitetura-backend.md). Pesquisa de regras e exemplos oficiais em [Validação estrutural de CPF e CNPJ](PESQUISA_VALIDACAO_CPF_CNPJ.md); testes das dependências adotadas em [test/dependencias_test.go](../services/identity-service/test/dependencias_test.go).

## Decisões para as próximas implementações

| Uso | Módulo e versão fixada | Licença | Decisão |
|---|---|---|---|
| IDs | `github.com/google/uuid v1.6.0` | BSD-3-Clause | Adotar a versão avaliada; preservar tipo UUID e rejeitar `uuid.Nil` no domínio. |
| Datas civis | `cloud.google.com/go v0.123.0`, pacote `cloud.google.com/go/civil` | Apache-2.0 | Adotar somente o pacote de valores, encapsulando `civil.Date` em DataNascimento. Não há módulo separado `cloud.google.com/go/civil` nesta seleção. |
| CPF/CNPJ | `github.com/cll0ud/cpfcnpj v0.1.0` | MIT | Avaliada e rejeitada para adoção nesta rodada. A biblioteca completa entradas curtas e descarta caracteres indevidos, contrariando o contrato. |
| Validação de documentos | Biblioteca padrão Go, implementação própria dentro dos valores CPF/CNPJ | Parte do projeto | Implementar na #8, a partir das regras oficiais e testes, sem dependência externa de CPF/CNPJ. Não copiar implementação/licença do candidato. |
| E-mail, nome e testes | Biblioteca padrão Go | Licença da distribuição Go | Manter a decisão do plano; sem pacote adicional. |
| Comparação de testes upstream | `github.com/google/go-cmp v0.7.0` | BSD-3-Clause | Auxiliar dos testes upstream de `civil`; seus hashes são mantidos por `go mod tidy`, sem dependência direta no `go.mod` do serviço. |

As versões foram resolvidas por `go list -m -json <módulo>@latest` com `GOPROXY=https://proxy.golang.org,direct` e depois avaliadas explicitamente nas versões acima. O `go.mod` do serviço fixa somente UUID e o módulo que fornece `civil`; `go mod tidy` gerou o `go.sum`. A candidata rejeitada não é dependência do serviço. Não usar `@latest` como receita reproduzível.

## Compatibilidade e procedência

| Módulo | Declaração do `go.mod` inspecionado | Commit da tag |
|---|---|---|
| `github.com/google/uuid v1.6.0` | Sem diretiva `go`; não deduzir versão mínima efetiva a partir dessa ausência | `0f11ee6918f41a04c201eceeadf612a377bc7fbc` |
| `cloud.google.com/go v0.123.0` | `go 1.24.0` | `4e8373586a5e48c18fbfd4bb0a3e259184e49a91` |
| `github.com/cll0ud/cpfcnpj v0.1.0` | `go 1.18` | `368530337fbbd03fefb43df980f6b2ccb1185ed4` |
| `github.com/google/go-cmp v0.7.0` | `go 1.21` | `9b12f366a942ebc7254abc7f32ca05068b455fb7` |

Ambiente efetivamente medido: `go version go1.27.1-X:nodwarf5 linux/amd64`. O módulo do serviço declara `go 1.24.0`, mínimo exigido pelo módulo selecionado de datas; a execução ocorreu com o toolchain acima. Não houve execução com Go 1.24 nem em outras plataformas.

Embora o módulo raiz `cloud.google.com/go` liste dependências de produtos cloud, a árvore efetiva de imports dos três pacotes avaliados mostrou apenas os próprios pacotes como não padrão; seus imports de produção transitivos são da biblioteca padrão. `civil` não abre conexão nem exige credenciais cloud. Não importar clientes de produtos para obter uma data civil. `database/sql/driver` aparece em `civil` como biblioteca padrão para conversões de valores, não como ORM/conexão do domínio. Essa verificação é por pacote/imports, não a afirmação de que o grafo inteiro do módulo raiz seja pequeno.

Os arquivos de licença foram lidos nas versões baixadas. BSD-3-Clause exige preservar os avisos/condições nas redistribuições e não usar nomes dos autores como endosso. Apache-2.0 exige preservar a licença e avisos aplicáveis, inclusive NOTICE quando fornecido; revisar os avisos da distribuição efetiva ao empacotar. A candidata é MIT, não copyleft, mas licença permissiva não resolve incompatibilidade de comportamento. Links versionados: [UUID: licença](https://github.com/google/uuid/blob/v1.6.0/LICENSE), [civil: licença do módulo](https://github.com/googleapis/google-cloud-go/blob/v0.123.0/LICENSE), [cpfcnpj: licença](https://github.com/cll0ud/cpfcnpj/blob/v0.1.0/LICENSE), [go-cmp: licença](https://github.com/google/go-cmp/blob/v0.7.0/LICENSE).

Os hashes das dependências estão em [go.sum](../services/identity-service/go.sum). `go mod verify` verifica a integridade do material baixado/cache, não certifica segurança ou correção de negócio.

## UUID

O pacote expõe `[16]byte` como `uuid.UUID`; a geração `NewRandom()` retorna valor e erro. No bootstrap/caso de uso, preferir essa forma para propagar falha de entropia; o construtor de entidade recebe UUID já gerado. `uuid.New()` é uma conveniência que pode entrar em panic se a geração falhar e não deve substituir o tratamento explícito dessa falha. A avaliação verificou versão 4, variante RFC4122, não nulidade e preservação da causa com fonte que falha. Não houve teste probabilístico de unicidade. [Código da versão: geração](https://github.com/google/uuid/blob/v1.6.0/version4.go).

Interpretar com sucesso um UUID não comprova tenant, autorização ou existência. O parser aceita UUID nulo, que continua inválido para os IDs das entidades. Não tratar `uuid.Parse` como validação completa do contrato de cadastro. Formatos de transporte e verificação autoritativa continuam nas camadas próprias. [Código da versão: tipo e parsing](https://github.com/google/uuid/blob/v1.6.0/uuid.go).

## Data civil

`civil.Date` guarda ano, mês e dia, sem fuso ou instante. Como seus campos são públicos, o valor próprio do projeto o guarda em campo privado e só devolve cópias. A avaliação confirmou 29/02/2024, rejeição de 29/02/2023 e do valor zero, e round trip de `ParseDate`. `IsValid()` também aceita ano 0, ano 10000 e datas futuras como datas civis: não é validador de nascimento. [Código versionado de civil](https://github.com/googleapis/google-cloud-go/blob/v0.123.0/civil/civil.go).

O objeto DataNascimento impõe adicionalmente ano de 1 a 9999 e data não futura em relação a `hoje` válido recebido explicitamente. Não usar relógio global nem converter nascimento automaticamente para meia-noite de um fuso. Instantes de confirmação/encerramento permanecem em `time.Time`, conforme os contratos.

## Avaliação da candidata CPF/CNPJ

Código, testes upstream e changelog de `v0.1.0` foram inspecionados. Os próprios testes aprovam completar zeros em CPF/CNPJ truncados; isso é comportamento intencional dessa versão, não apenas ausência de teste. As rotinas removem caracteres arbitrários antes de calcular DV. Portanto, usar diretamente seus construtores ou `IsValid` não cumpre as regras do Teu Ponto. [CPF](https://github.com/cll0ud/cpfcnpj/blob/v0.1.0/cpf.go), [CNPJ](https://github.com/cll0ud/cpfcnpj/blob/v0.1.0/cnpj.go), [changelog](https://github.com/cll0ud/cpfcnpj/blob/v0.1.0/CHANGELOG.md).

| Entrada | Resultado observado na candidata | Contrato Teu Ponto |
|---|---|---|
| CPF oficial `28001238938` e sua máscara | Aceita, DV correto | Aceitar |
| CPF sintético `01234567890` | Preserva zero inicial | Aceitar |
| CPF curto `1234567890` | Completa para `01234567890` e aceita | Rejeitar, não recuperar dígito ausente |
| CPF `abc28001238938` ou máscara `280-012-389.38` | Aceita ao desconsiderar caracteres | Rejeitar |
| CNPJ oficial `00000000000191` | Aceita com todos os zeros presentes | Aceitar sem conversão para número |
| CNPJs oficiais `12ABC34501DE35`, `ABCDEFGHIJKL80` | Aceita e calcula os DVs | Aceitar |
| CNPJ `191` | Completa para `00000000000191` e aceita | Rejeitar |
| CNPJ `12#ABC34501DE35` | Aceita descartando `#` para validar | Rejeitar |
| CNPJ vazio ou `!!!` | Retorna `00000000000000`, sem erro | Rejeitar |
| CNPJ `00000000000000` | `IsValid` retorna true | Rejeitar, conforme os testes de referência da Receita |
| CNPJ com DV errado | Construtor retorna erro e valor vazio; `IsValid` desse valor vazio retorna true | Sempre tratar erro; valor zero do domínio deve ser inválido |

Na avaliação exploratória foram caracterizados nove casos de CPF, quinze de CNPJ e o valor vazio retornado após erro. Nove entradas divergiram do contrato estrito, além da limitação do retorno após erro. Esses resultados fundamentam a rejeição; os testes da candidata não foram incorporados ao serviço. DVs alterados e letras no DV foram rejeitados corretamente. Minúsculas em CNPJ foram convertidas explicitamente para maiúsculas na entrada, comportamento adotado pelo valor próprio; o validador oficial trabalha sobre representação canônica em maiúsculas.

Uma fachada estrita poderia bloquear esses casos antes de chamar a candidata. Foi escolhida implementação própria dos pequenos algoritmos na #8 para evitar depender das normalizações permissivas e duplicar verificações em torno dos tipos públicos baseados em string, cujos métodos de formatação alteram o valor. Isso não antecipa código de produção nem altera os contratos externos de CPF/CNPJ. Uma versão futura da biblioteca exigirá nova avaliação explícita, não substituição silenciosa por `latest`.

## Contrato de entrada e encapsulamento para a #8

- CPF: aceitar exatamente 11 dígitos ASCII ou a máscara exata `DDD.DDD.DDD-DD`. CNPJ: aceitar exatamente `[A-Z0-9]{12}[0-9]{2}` ou a máscara exata `XX.XXX.XXX/XXXX-DD`, com as 12 primeiras posições alfanuméricas. Converter somente letras ASCII `a-z` para maiúsculas antes de validar CNPJ. Todas as letras A–Z são possíveis, inclusive nas duas primeiras posições.
- Validar texto inteiro e máscara antes de remover exclusivamente os separadores nos lugares previstos. Rejeitar espaços, sinais, símbolos extras, caracteres Unicode semelhantes, letras em CPF/DV e tamanhos diferentes. Não aparar nem completar documento; preservar cada zero recebido. Guardar representação canônica como string privada, nunca inteiro.
- Aplicar cálculo de DV após estrutura/alfabeto, com regras oficiais identificadas na pesquisa. Rejeitar CPF composto por um único dígito repetido e CNPJ zerado. Não inventar lista adicional de exceções de outros sistemas. Testar ambos os DVs separadamente e preservar estado em falha.
- Tipos próprios expõem o contrato do projeto; bibliotecas de valor não vazam por retorno de CPF/CNPJ, erros ou DTOs. Usar erros distinguíveis com contexto, sem incluir o documento nas mensagens/logs. Não produzir matrícula a partir do documento.
- Checksum/formato não comprovam existência, titularidade, situação fiscal, vínculo com empresa ou autorização. CPF didático foi obtido de fonte oficial indexada; os PDFs Sped tiveram timeout e a limitação permanece registrada na pesquisa. Exemplos sintéticos são identificados e não são enviados a cadastros de produção.

## Verificações e próximos passos

Comandos diretos em `services/identity-service/`:

```bash
go mod tidy
go test ./...
go vet ./...
go mod verify
```

Os testes de UUID e datas civis estão em `test/`; não há script intermediário nem cópia de manifesto. A primeira execução de `go test ./...` encontrou o bootstrap existente vazio; foi acrescentada a declaração `package main` para permitir a análise do módulo. Não há servidor implementado.

Na #8, implementar CPF/CNPJ em TDD com os exemplos identificados; na #9, encapsular data civil e regras de nascimento. Chi, GORM, PostgreSQL e SDKs de infraestrutura continuam fora das dependências atuais. Revisão/homologação permanecem pendentes. A renomeação local do RNF foi preservada; o plano ainda referencia seu nome versionado original.
