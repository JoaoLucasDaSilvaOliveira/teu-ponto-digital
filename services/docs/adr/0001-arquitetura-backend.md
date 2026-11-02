---
status: accepted
---
# Serviços separados com dependências orientadas ao domínio

Decidido em 28/10/2026: cada serviço terá sua estrutura em `services/<nome-do-serviço>/`, com `cmd/` e `internal/` próprios, em vez de um backend único diretamente em `services/`. Isso mantém explícita a organização por serviço prevista para o backend; os nomes, limites de negócio e a divisão dos módulos Go serão definidos conforme a implementação.

Dentro de cada serviço, `domain` concentra entidades, value objects, regras e contratos de domínio; `application` coordena casos de uso; `interfaces` adapta HTTP e eventos de entrada; `infra` implementa persistência e integrações. Contratos pequenos ficam na camada que os consome; os bootstraps em `cmd/` compõem as implementações por injeção explícita de dependências. A separação mantém regras de negócio independentes de transporte e persistência e permite testá-las antes de implementar os adaptadores.

As instruções operacionais estão em [services/AGENTS.md](../../AGENTS.md). O bootstrap existente em `services/cmd/main/main.go` ainda precisa ser alocado a um serviço quando seu nome e responsabilidade forem definidos; este registro não move nem implementa código.

## Complemento aprovado em 30/10/2026

O primeiro serviço será `identity-service`, em `services/identity-service/`, responsável pelas identidades de acesso, cadastros de empresas e funcionários e vínculos de autorização. A primeira rodada entrega o modelo de domínio e testes, conforme o [plano aprovado](../../../docs/PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md); API, persistência e autenticação completa são etapas posteriores.

Aplicar encapsulamento, abstração, composição e polimorfismo conforme as necessidades do domínio. Dados e comportamentos comuns devem ser reutilizados por composição; quando o comportamento variar por tipo, usar interfaces pequenas e implementações específicas, com seleção inicial no bootstrap. O exemplo de uma entidade generalista `Pessoa` não obriga a criação desse tipo nem de uma hierarquia de herança. Go não oferece herança de classes; composição e incorporação de tipos não tornam a estrutura externa um subtipo da estrutura incorporada.

Os campos das entidades devem ser privados. Construtores retornam erro para dados inválidos e utilizam os próprios métodos de validação e atribuição da entidade. Métodos de alteração preservam invariantes e mantêm o estado anterior em caso de erro. Dependências e regras entre entidades não devem ser simuladas por setters locais: unicidade e verificação de referências exigem coordenação com os contratos apropriados.

Para datas civis, como nascimento, usar um objeto de valor próprio que encapsule `cloud.google.com/go/civil.Date`, mantendo instantes em `time.Time`. Essa dependência de valores é compatível com o domínio independente de infraestrutura. CPF e CNPJ também serão tipos próprios; `github.com/cll0ud/cpfcnpj` é candidato à implementação interna, sujeito à avaliação de licença, versão e comportamento antes da adoção.

A geração de matrícula interna será uma biblioteca própria em `internal/matriculagen/`, separada do objeto de valor em `internal/domain/value_object/matricula.go`. O gerador utiliza aleatoriedade segura; o objeto de valor valida e representa a matrícula. O construtor do funcionário recebe uma matrícula já gerada. A garantia persistente de unicidade e o tratamento de colisões pertencem à coordenação do cadastro com a persistência.

## Detalhamento do domínio em 01/11/2026

Os [contratos da issue #4](../../../docs/ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md) adotam Pessoa como entidade do CNPJ para dados pessoais, referenciada por Usuário, Funcionário e pela relação de responsável legal. Substituem a previsão inicial de responsável legal como objeto independente com dados duplicados, sem criar herança. O [ADR 0007](../../../docs/adr/0007-contas-por-cnpj-e-recontratacao.md) estabelece contas independentes por CNPJ e preservação de cadastros encerrados na recontratação. Pessoa será implementada antes de seus consumidores na etapa de empresa/responsável, mantendo as demais dependências do plano.

## Avaliação de dependências em 02/11/2026

A [issue #6](../../../docs/DEPENDENCIAS_DOMINIO_IDENTITY_SERVICE.md) fixa `github.com/google/uuid v1.6.0` e o pacote `cloud.google.com/go/civil` do módulo `cloud.google.com/go v0.123.0`. A candidata `github.com/cll0ud/cpfcnpj v0.1.0` foi rejeitada após testes por completar entradas incompletas e aceitar caracteres arbitrários/vazio; CPF/CNPJ serão implementados com biblioteca padrão nos tipos próprios da #8. Essa conclusão substitui sua condição anterior de candidata nesta rodada. Por solicitação do usuário, módulo e hashes foram criados diretamente em `services/identity-service/`, com testes em `test/`, antecipando essa parte da #7. O domínio permanece independente de transporte/persistência.

## Destino do bootstrap na issue #7

Em 02/11/2026, o bootstrap vazio, movido anteriormente de `services/cmd/main/main.go` para `services/identity-service/cmd/main/main.go`, foi removido. Sua única declaração, `package main`, permitia teste/análise, mas falhava no build por não haver função `main`. A etapa de domínio não precisa de executável; o bootstrap será criado no caminho do serviço quando houver comportamento executável e dependências concretas a compor. Isso mantém a organização aceita sem criar servidor fictício. Os comandos e resultados estão no [README do módulo](../../identity-service/README.md).
