# Plano de implementação — identity-service

Status: aprovado em 30/10/2026. Esta aprovação autoriza o registro do plano e a criação das issues; a implementação será conduzida pelas issues, em rodadas posteriores.

## Objetivo e limite da rodada

Construir o domínio do primeiro serviço, responsável pelas identidades de acesso, cadastros de empresas e funcionários e vínculos de autorização, preparando o fornecimento de dados a outros serviços. Entregar entidades, objetos de valor, interfaces necessárias, biblioteca própria de matrícula e testes.

API, persistência, login completo, tokens, integrações, upload e armazenamento do A1, biometria, jornadas e cálculos são etapas posteriores. Não há módulo Go ou suíte executável nesta etapa de planejamento; comandos reais serão registrados quando o módulo existir.

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
- Quatro perfis: Gestor, Supervisor, Colaborador e Contabilidade/DP. Não existe Gestor Master separado.
- O gestor define permissões da Contabilidade/DP por ação e por CNPJ, como exportar PDF, criar horários e criar jornadas (RF-06.1). Contexto selecionado não concede autorização.
- Responsável legal tem nome e CPF obrigatórios.
- Certificado A1 é enviado à aplicação com senha quando necessária; validar o arquivo e proteger arquivo/senha na implementação futura. A3 fica para o futuro. A representação de referência no domínio não homologa signatários nem conformidade dos relatórios.
- Matrícula interna aleatória, não sequencial, estável e gerada por biblioteca própria. Pode ser apresentada no coletor; não substitui comparação facial e prova de vida. Matrícula da folha continua como referência externa separada.

## Modelo inicial a especificar

| Conceito | Dados principais | Comportamentos previstos |
|---|---|---|
| Empresa | ID/tenant, CNPJ, razão social, apelido, responsável e referência do A1 | Alterar cadastro, definir responsável, associar referência de certificado |
| Usuário | ID, e-mail, confirmação e situação de acesso; associação aos cadastros pessoais definida na modelagem | Alterar e-mail invalidando confirmação, confirmar e-mail, bloquear acesso |
| Funcionário | ID, empresa, matrícula, CPF, nome, nascimento, setor e usuário opcionais, modalidade | Atualizar dados pessoais, atribuir setor, associar usuário |
| Vínculo com empresa | ID, empresa, usuário, perfil, permissões, escopo e situação | Conceder/revogar acesso, alterar perfil, avaliar permissão |
| Setor | ID, empresa e nome | Renomear e representar agrupamento de supervisão |
| Responsável legal | Nome e CPF | Validar e atualizar identificação |
| Referência do A1 | Empresa, identificação do arquivo protegido e metadados necessários | Representar associação sem expor segredos |

A especificação determinará entidades versus objetos de valor, relações/cardinalidades, campos obrigatórios/opcionais/imutáveis, tipos, assinaturas dos construtores, métodos, erros e estados. Definirá estados de cadastro incompleto e invariantes entre referências. Unicidade persistente e carregamento autoritativo de vínculos são responsabilidades futuras de coordenação; validação local não os garante.

## Dependências

| Componente | Decisão |
|---|---|
| IDs | `github.com/google/uuid`, já escolhido |
| Datas civis | Objeto próprio com `cloud.google.com/go/civil.Date`; instantes usam `time.Time` |
| CPF/CNPJ | Tipos próprios; `github.com/cll0ud/cpfcnpj` candidato, condicionado à avaliação de licença, versão e comportamento |
| E-mail e testes | Biblioteca padrão Go |
| Matrícula | Biblioteca própria com aleatoriedade segura |
| HTTP/persistência | Chi, GORM e driver PostgreSQL são escolhas para próximas camadas, sem adoção antecipada |

Fixar versões reproduzíveis no módulo após avaliação. Conferir exemplos oficiais de CNPJ numérico/alfanumérico e rejeição de dados incompletos, sem completar zeros silenciosamente. Referências de avaliação: [civil](https://pkg.go.dev/cloud.google.com/go/civil), [biblioteca candidata](https://github.com/cll0ud/cpfcnpj) e [CNPJ alfanumérico — Receita Federal](https://www.gov.br/receitafederal/pt-br/acesso-a-informacao/acoes-e-programas/programas-e-atividades/cnpj-alfanumerico).

## Biblioteca de matrícula

`internal/matriculagen/` gera valores aleatórios com fonte segura e propaga falhas. `internal/domain/value_object/matricula.go` valida e representa matrículas, inclusive existentes. O construtor do funcionário recebe uma matrícula já gerada.

Antes do código, definir tamanho, alfabeto, formato e escopo de unicidade. Não derivar de CPF/nome ou sequência. A persistência futura impõe unicidade e a coordenação do cadastro permite nova tentativa em colisões; o gerador não promete unicidade absoluta. Testar fonte controlável e falhas, sem testes probabilísticos frágeis.

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

Cada objeto de valor tem arquivo próprio, por exemplo `cnpj.go`, `cpf.go`, `birth_date.go` e `matricula.go`, com testes correspondentes. Criar pacotes/interfaces somente com consumidor real. A estrutura futura de `cmd/` e das demais camadas segue o ADR; definir o destino do bootstrap preexistente sem inventar um servidor nesta rodada.

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
| 9 | Implementar responsável legal, empresa e referência do A1 | 2, 5, 6, 7 | 3 | HIGH |
| 10 | Implementar usuário e associação com os cadastros | 2, 5, 7 | 2 | HIGH |
| 11 | Implementar setor e funcionário | 2, 5, 6, 7, 8, 10 | 2 | HIGH |
| 12 | Implementar vínculo com empresa e políticas dos quatro perfis | 2, 3, 5, 7, 9, 10, 11 | 5 | HIGH |
| 13 | Revisar domínio, documentação e verificações da rodada | 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 | 2 | MEDIUM |

Criar uma nova issue principal e estas 13 subissues no repositório atual, sem reutilizar a issue de teste #1. Cada uma terá objetivo, fontes, escopo, critérios verificáveis e dependências. Projeto: Teu Ponto — Desenvolvimento; view Backlog; Status Todo; Area Backend; responsável JoaoLucasDaSilvaOliveira; Priority e Estimate conforme tabela. A principal é um agrupador com Estimate 0, para não duplicar o esforço das subissues. Progresso de subissues é calculado pelo GitHub; vincular PRs quando existirem.

## Verificação e execução posteriores

Executar TDD em cada implementação: cenário observável, teste que falha pelo comportamento ausente, implementação mínima e refatoração. Cobrir invariantes, documentos inválidos, datas civis, falhas de aleatoriedade, preservação de estado e isolamento por tenant/setor/funcionário. Políticas devem testar ações não concedidas, vínculos revogados e permissões diferentes do mesmo usuário em CNPJs distintos, sem banco ou HTTP.

Testes, análise estática e formatação terão comandos documentados no módulo. Não declarar o domínio como serviço seguro em produção antes de implementar a autorização autoritativa nas entradas, a persistência e a infraestrutura. Seguir branches `feat/identity-service/...`, commits convencionais e fluxo de revisão/homologação do repositório.

## Pendências preservadas

- Formato e escopo de unicidade da matrícula, contratos finais, cardinalidades e estados de cadastro incompleto: resolver nas especificações antes do código.
- Matriz completa de ações, autoaprovação e acesso a atestados: especificar por área, sem conceder implicitamente.
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
