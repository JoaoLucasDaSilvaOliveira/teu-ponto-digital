# Backend — arquitetura

Estas instruções se aplicam a `services/` e complementam o [AGENTS.md da raiz](../AGENTS.md). Requisitos de produto e backlog permanecem nas suas fontes; este arquivo define a arquitetura do backend. A decisão de organização está registrada no [ADR de arquitetura](docs/adr/0001-arquitetura-backend.md).

As regras de branches, commits e promoção entre ambientes são as de [Versionamento de código no AGENTS.md da raiz](../AGENTS.md#versionamento-de-código). Aplique-as também ao trabalho em cada serviço.

## Dependências e domínio

Use injeção explícita de dependências. Os pontos de entrada em `cmd/` carregam configuração, inicializam implementações concretas e distribuem dependências por construtores. Casos de uso recebem suas dependências; abertura de conexões, escolha de implementações e composição da aplicação pertencem ao bootstrap.

O domínio concentra entidades, value objects e regras de negócio. As dependências de negócio apontam para o domínio, direta ou indiretamente. Entidades preservam suas invariantes; regras de negócio não devem depender de HTTP, ORM ou broker.

Entidades podem usar biblioteca padrão, bibliotecas de valores como `github.com/google/uuid` e bibliotecas próprias de validação/value objects que respeitem a mesma independência. Elas não importam camadas externas da aplicação nem frameworks de transporte ou persistência. As exceções para UUID/validação não autorizam dependências arbitrárias de infraestrutura.

## Polimorfismo

Quando o comportamento variar conforme o tipo do objeto, expresse a variação por interfaces Go e implementações específicas. O chamador trabalha com o contrato, sem inspecionar o tipo concreto para decidir como o objeto deve responder.

Use interfaces pequenas, correspondentes ao comportamento exigido pelo consumidor. Validações e decisões próprias de uma regra continuam sendo condições normais; a orientação de polimorfismo se aplica à seleção de comportamento por tipo. A seleção inicial de implementações concretas pertence ao bootstrap. Go não requer classes abstratas para aplicar esse padrão.

## Estrutura

Cada serviço fica em `services/<nome-do-serviço>/` e possui seus próprios `cmd/` e `internal/`. Os caminhos abaixo são relativos à raiz de cada serviço. A divisão dos módulos Go e os nomes dos primeiros serviços serão definidos ao iniciar sua implementação.

| Caminho | Responsabilidade |
|---|---|
| `/` | README e configuração da aplicação, incluindo suporte a `.env`. Segredos reais ficam fora do controle de versão. |
| `cmd/` | Pontos de entrada e composição de dependências. |
| `cmd/main/main.go` | Bootstrap principal, por exemplo do servidor HTTP. |
| `cmd/workers/` | Bootstraps dos processos assíncronos. Cada executável possui seu próprio diretório e `main.go`. |
| `internal/` | Código interno de produção: domínio e demais camadas. |

Um producer é uma capacidade de publicação, não necessariamente um processo independente. Se participar de um caso de uso, sua implementação pertence a `internal/`; `cmd/workers/` é reservado à composição e execução de processos.

### Responsabilidades das camadas

| Camada | Responsabilidade |
|---|---|
| `internal/domain/` | Entidades, value objects, regras e contratos de domínio. |
| `internal/application/` | Casos de uso que coordenam o domínio e dependências por contratos. |
| `internal/interfaces/` | Entradas da aplicação: handlers HTTP e consumidores de eventos; adaptação de dados de transporte para casos de uso. |
| `internal/infra/` | Persistência, clientes de serviços e publicação de mensagens; implementações concretas dos contratos consumidos internamente. |

Direção de imports: `application` → `domain`; `interfaces` → `application`/`domain`; `infra` → contratos e tipos das camadas internas. `cmd` reúne as implementações. Contratos pequenos de repositórios/gateways ficam na camada que os consome, preservando o domínio independente. Crie subpacotes e interfaces à medida que houver comportamento que os utilize.

## Pacotes escolhidos

- UUID: `github.com/google/uuid`.
- Roteamento HTTP: `github.com/go-chi/chi/v5`.
- Persistência: `gorm.io/gorm` e o driver do banco adotado. Os requisitos atuais definem PostgreSQL; nesse caso, `gorm.io/driver/postgres`.

Chi pertence à adaptação HTTP e sua composição; GORM pertence à persistência e sua inicialização. Preserve o domínio independente do ORM, incluindo tipos como `gorm.Model` e `*gorm.DB`. A forma de mapeamento entre entidades e registros de persistência deve ser definida ao detalhar as camadas. As versões serão fixadas no `go.mod` ao iniciar a implementação.

## TDD

Antes de implementar uma feature, defina os cenários de negócio e escreva o próximo teste unitário ou de aplicação para o comportamento esperado. Execute-o e confirme que falha pelo comportamento ausente; implemente o mínimo para passar e refatore mantendo os testes verdes. Avance em ciclos pequenos.

Testes de domínio verificam invariantes e cálculos. Testes de aplicação exercitam casos de uso pelos seus contratos, com dependências controláveis. Regras de negócio devem ser verificáveis sem servidor HTTP ou banco real. Para persistência, transporte e demais integrações, acrescente testes de integração quando necessários para verificar comportamento que testes isolados não cobrem.

Os comandos de verificação e integração com CI/CD serão documentados quando o módulo e a implementação existirem; este arquivo não afirma que já há uma suíte executável.
