# Plano de implementação — identity-service

Status: aprovado em 30/10/2026. Esta aprovação autoriza o registro do plano e a criação das issues; a implementação será conduzida pelas issues, em rodadas posteriores.

## Objetivo e limite da rodada

Construir o domínio do primeiro serviço, responsável pelas identidades de acesso, cadastros de empresas e funcionários e vínculos de autorização, preparando o fornecimento de dados a outros serviços. Entregar entidades, objetos de valor, interfaces necessárias, biblioteca própria de matrícula e testes.

API, persistência, login completo, tokens, integrações, upload e armazenamento do A1, biometria, jornadas e cálculos são etapas posteriores. Por orientação do usuário em 02/11/2026, o módulo foi inicializado diretamente em `services/identity-service/`, com testes das dependências adotadas em `test/`; os comandos reais e os resultados da inicialização da #7 estão no [README do módulo](../services/identity-service/README.md).

## Fontes

- [Vocabulário](../CONTEXT.md).
- [Requisitos funcionais](REQUISITOS_FUNCIONAIS.md): RF-01 a RF-08, RF-11.1, RF-19 e RF-20.
- [Requisitos não funcionais](<REQUISITOS_NÃO _FUNCIONAIS.md>).
- [Revisão de coerência](REVISAO_DOCUMENTACAO.md).
- [Isolamento por CNPJ](adr/0004-isolamento-por-cnpj.md) e [modalidade de coleta](adr/0005-modalidade-coleta-por-funcionario.md).
- [Arquitetura aceita do backend](../services/docs/adr/0001-arquitetura-backend.md) e [instruções do backend](../services/AGENTS.md).

## Decisões aprovadas

- Aplicar encapsulamento, abstração, composição e polimorfismo conforme os comportamentos do domínio. Interfaces pequenas representam contratos; implementações específicas expressam comportamento variável. `Pessoa` foi um exemplo de generalização, não uma entidade obrigatória.
- Campos privados; construtores retornam erro e usam métodos próprios para validar e atribuir. Alterações inválidas preservam o estado anterior.
- Separar usuário, funcionário e vínculo com empresa. Funcionário pode existir sem conta; uma pessoa pode ser gestora e funcionária. Compartilhar dados e comportamentos comuns conforme a modelagem, sem duplicação entre perfis.
- Detalhamento confirmado em 01/11/2026: Pessoa concentra CPF/nome/nascimento no CNPJ; Usuário e Funcionário a referenciam. Contas são independentes por CNPJ e recontratação preserva os cadastros encerrados, conforme [ADR 0007](adr/0007-contas-por-cnpj-e-recontratacao.md). Responsável legal passa a ser uma relação com Pessoa, sem cópia cadastral por perfil.
- Quatro perfis: Gestor, Supervisor, Colaborador e Contabilidade/DP. Não existe Gestor Master separado.
- O gestor define permissões da Contabilidade/DP por ação e por CNPJ, como exportar PDF, criar horários e criar jornadas (RF-06.1). Contexto selecionado não concede autorização.
- Detalhamento confirmado na issue #5: Gestor exclusivo, Supervisor compondo ações próprias do Colaborador e concessões sobre equipe diretamente atribuída, sem autoaprovação. Configuração, catálogo inicial e contratos constam em [Autorização](AUTORIZACAO_IDENTITY_SERVICE.md); todos os acessos e serviços da aplicação serão inventariados na [issue #17](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/17), trabalho relacionado fora das 13 etapas/estimativas originais.
- Responsável legal tem nome e CPF obrigatórios.
- Certificado A1 é enviado à aplicação com senha quando necessária; validar o arquivo e proteger arquivo/senha na implementação futura. A3 fica para o futuro. A representação de referência no domínio não homologa signatários nem conformidade dos relatórios.
- Matrícula interna aleatória, não sequencial, estável e gerada por biblioteca própria. Pode ser apresentada no coletor; não substitui comparação facial e prova de vida. Matrícula da folha continua como referência externa separada.

## Modelo inicial a especificar

O detalhamento da issue #4 está em [Contratos de domínio](ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md): entidades, campos, cardinalidades, assinaturas, erros, estados, matrícula e coordenação de invariantes. A [especificação da issue #5](AUTORIZACAO_IDENTITY_SERVICE.md) completa configuração/política de autorização e referência direta de supervisão no Funcionário, substituindo setor como critério de subordinação.

Os limites de responsabilidade, os dados fornecidos aos demais módulos e a rastreabilidade da etapa 1 estão detalhados em [Limites e responsabilidades do identity-service](LIMITES_IDENTITY_SERVICE.md), entrega documental da issue #3. Esse detalhamento não antecipa a implementação nem encerra as pendências das etapas seguintes.

| Conceito | Dados principais | Comportamentos previstos |
|---|---|---|
| Empresa | ID/tenant, CNPJ, razão social, apelido, pessoa responsável e referência do A1 | Alterar identificação, definir responsável, associar referência de certificado |
| Pessoa | ID, empresa, CPF, nome e nascimento opcional para externos | Alterar nome, corrigir CPF, informar/corrigir nascimento |
| Usuário | ID, empresa, pessoa, e-mail, confirmação e situação de acesso | Alterar e-mail invalidando confirmação, confirmar e-mail, bloquear e encerrar acesso |
| Funcionário | ID, empresa, pessoa, matrícula, setor, usuário e vínculo Supervisor opcionais, modalidade e situação | Atribuir setor e supervisão, associar usuário, encerrar cadastro |
| Vínculo com empresa | ID, empresa, usuário, perfil, permissões, escopo e situação | Conceder/revogar acesso, alterar perfil, avaliar permissão |
| Setor | ID, empresa e nome | Renomear e representar agrupamento de supervisão |
| Responsável legal | Relação da empresa com Pessoa identificada por nome/CPF | Designar pessoa sem conceder acesso automaticamente |
| Referência do A1 | Empresa, identificação do arquivo protegido e metadados necessários | Representar associação sem expor segredos |

Os contratos da issue #4 definem entidades versus objetos de valor, relações, cardinalidades, campos, construtores, métodos, erros e estados. Unicidade persistente e carregamento autoritativo de vínculos continuam responsabilidades futuras de coordenação; validação local não os garante.

## Dependências

| Componente | Decisão |
|---|---|
| IDs | `github.com/google/uuid v1.6.0`, avaliado na #6 |
| Datas civis | Objeto próprio com `cloud.google.com/go/civil.Date`, fornecido por `cloud.google.com/go v0.123.0`; instantes usam `time.Time` |
| CPF/CNPJ | Tipos próprios com validação pela biblioteca padrão a implementar na #8; `github.com/cll0ud/cpfcnpj v0.1.0` avaliado e rejeitado por completar documentos e aceitar caracteres indevidos |
| E-mail e testes | Biblioteca padrão Go |
| Matrícula | Biblioteca própria com aleatoriedade segura |
| HTTP/persistência | Chi, GORM e driver PostgreSQL são escolhas para próximas camadas, sem adoção antecipada |

A [avaliação da issue #6](DEPENDENCIAS_DOMINIO_IDENTITY_SERVICE.md) registra versões, licenças, procedência e resultados; a [pesquisa de documentos](PESQUISA_VALIDACAO_CPF_CNPJ.md) distingue fontes oficiais e exemplos sintéticos. As dependências adotadas estão no `go.mod` real do serviço, com Go 1.24.0 como mínimo exigido pelo módulo de datas e hashes gerados por `go mod tidy`. Essa parte da #7 foi antecipada por solicitação do usuário.

## Biblioteca de matrícula

`internal/matriculagen/` gera valores aleatórios com fonte segura e propaga falhas. `internal/domain/value_object/matricula.go` valida e representa matrículas, inclusive existentes. O construtor do funcionário recebe uma matrícula já gerada.

Formato definido na issue #4: seis caracteres de `0-9A-Z`, normalizados para maiúsculas, únicos por CNPJ inclusive no histórico. Recontratação recebe nova matrícula; antigas permanecem reservadas. Os contratos detalham geração uniforme e tentativas em colisões. Não derivar de CPF/nome ou sequência. A persistência futura impõe unicidade; o gerador não promete unicidade absoluta. Testar fonte controlável e falhas, sem testes probabilísticos frágeis.

## Estrutura prevista

```text
services/identity-service/
├── go.mod
├── go.sum
├── README.md
└── internal/
    ├── domain/
    │   ├── entity/
    │   ├── value_object/
    │   └── policy/
    └── matriculagen/
```

Cada objeto de valor tem arquivo próprio, por exemplo `cnpj.go`, `cpf.go`, `birth_date.go` e `matricula.go`, com testes correspondentes. Criar pacotes/interfaces somente com consumidor real. Na #7, o bootstrap vazio, já movido para `services/identity-service/cmd/main/main.go`, foi removido por não haver executável nesta etapa. A estrutura futura de `cmd/` e das demais camadas segue o ADR; o bootstrap será criado quando houver comportamento executável a compor.

## Sequência e backlog

Cronograma por sequência e dependências, sem datas. Estimativas relativas em pontos (1, 2, 3, 5), sem conversão para dias; total inicial das subissues: 33 pontos.

| Ordem | Ação | Depende de | Pontos | Prioridade |
|---|---|---|---|---|
| 1 | Detalhar limites e responsabilidades do identity-service | — | 2 | HIGH |
| 2 | Especificar entidades, atributos, métodos e construtores | 1 | 3 | HIGH |
| 3 | Especificar perfis, permissões e escopos de autorização | 1, 2 | 3 | HIGH |
| 4 | Avaliar e definir dependências externas do domínio | 2 | 2 | HIGH |
| 5 | Inicializar módulo e comandos do identity-service | 1, 4 | 1 | HIGH |
| 6 | Implementar objetos de valor CPF e CNPJ | 2, 4, 5 | 3 | HIGH |
| 7 | Implementar e-mail, nascimento e demais objetos de valor | 2, 4, 5 | 2 | MEDIUM |
| 8 | Implementar biblioteca matriculagen e objeto de valor Matricula | 2, 5 | 3 | HIGH |
| 9 | Implementar pessoa, relação de responsável legal, empresa e referência do A1 | 2, 5, 6, 7 | 3 | HIGH |
| 10 | Implementar usuário e associação com os cadastros | 2, 5, 7, 9 | 2 | HIGH |
| 11 | Implementar setor, funcionário e referência de supervisão | 2, 3, 5, 6, 7, 8, 10 | 2 | HIGH |
| 12 | Implementar vínculo com empresa e políticas dos quatro perfis | 2, 3, 5, 7, 9, 10, 11 | 5 | HIGH |
| 13 | Revisar domínio, documentação e verificações da rodada | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 | 2 | MEDIUM |

Criar uma nova issue principal e estas 13 subissues no repositório atual, sem reutilizar a issue de teste #1. Cada uma terá objetivo, fontes, escopo, critérios verificáveis e dependências. Projeto: Teu Ponto — Desenvolvimento; view Backlog; Status Todo; Area Backend; responsável JoaoLucasDaSilvaOliveira; Priority e Estimate conforme tabela. A principal é um agrupador com Estimate 0, para não duplicar o esforço das subissues. Progresso de subissues é calculado pelo GitHub; vincular PRs quando existirem.

## Verificação e execução posteriores

Executar TDD em cada implementação: cenário observável, teste que falha pelo comportamento ausente, implementação mínima e refatoração. Cobrir invariantes, documentos inválidos, datas civis, falhas de aleatoriedade, preservação de estado e isolamento por tenant/setor/funcionário. Políticas devem testar ações não concedidas, vínculos revogados e permissões diferentes de contas da mesma pessoa em CNPJs distintos, sem banco ou HTTP.

Testes, análise estática e formatação terão comandos documentados no módulo. Não declarar o domínio como serviço seguro em produção antes de implementar a autorização autoritativa nas entradas, a persistência e a infraestrutura. Seguir branches `feat/identity-service/...`, commits convencionais e fluxo de revisão/homologação do repositório.

## Pendências preservadas

- Formato/escopo da matrícula, cardinalidades, estados e contratos comuns estão definidos na issue #4; configuração, política e supervisão direta estão especificadas na #5 antes da implementação do vínculo.
- Matriz completa de ações/serviços na #17, sem concessão implícita. Autoaprovação foi proibida; tratamento de atestado inclui acesso ao documento e abono da ocorrência, conforme a #5. Demais fluxos de documentos e lançamentos exigem especificação própria.
- Custódia técnica do A1, signatários de artefatos e regras regulatórias: resolver antes dos respectivos adaptadores e relatórios.
- Cadastro facial inicial, mudança de modalidade, fallbacks e detalhes do coletor/offline: permanecem nas fontes próprias.


## Issues publicadas

Issue principal: [#2 — Domínio de identidades e cadastros](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/2).

| Ordem | Issue |
|---|---|
| 1 | [#3 — Detalhar limites e responsabilidades do identity-service](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/3) |
| 2 | [#4 — Especificar entidades, atributos, métodos e construtores](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/4) |
| 3 | [#5 — Especificar perfis, permissões e escopos de autorização](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/5) |
| 4 | [#6 — Avaliar e definir dependências externas do domínio](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/6) |
| 5 | [#7 — Inicializar módulo e comandos do identity-service](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/7) |
| 6 | [#8 — Implementar objetos de valor CPF e CNPJ](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/8) |
| 7 | [#9 — Implementar e-mail, nascimento e demais objetos de valor](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/9) |
| 8 | [#10 — Implementar biblioteca matriculagen e objeto de valor Matricula](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/10) |
| 9 | [#11 — Implementar responsável legal, empresa e referência do A1](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/11) |
| 10 | [#12 — Implementar usuário e associação com os cadastros](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/12) |
| 11 | [#13 — Implementar setor e funcionário](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/13) |
| 12 | [#14 — Implementar vínculo com empresa e políticas dos quatro perfis](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/14) |
| 13 | [#15 — Revisar domínio, documentação e verificações da rodada](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/15) |
