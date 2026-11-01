# Limites e responsabilidades — identity-service

Especificação documental da [issue #3](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/3), etapa 1 do [plano aprovado](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md), agrupada na [issue #2](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/2). Não representa implementação executável nem aprovação de políticas ainda pendentes.

## Fontes e alcance

Este detalhamento aplica as decisões do plano, o [vocabulário](../CONTEXT.md), [RF-01 a RF-08, RF-11.1, RF-19 e RF-20](REQUISITOS_FUNCIONAIS.md), os ADRs [0004](adr/0004-isolamento-por-cnpj.md) e [0005](adr/0005-modalidade-coleta-por-funcionario.md), a [arquitetura do backend](../services/docs/adr/0001-arquitetura-backend.md) e suas [instruções](../services/AGENTS.md). As lacunas da [revisão de coerência](REVISAO_DOCUMENTACAO.md) permanecem abertas onde explicitadas abaixo.

O `identity-service` é a fonte dos cadastros de identidades de acesso, empresas, funcionários, setores e vínculos de autorização. Seu módulo de domínio concentra as invariantes desses conceitos; os futuros casos de uso coordenam referências, unicidade e autorização com dados autoritativos. A existência de um cadastro não comprova autenticação nem concede acesso.

Esta issue entrega somente documentação. Até mesmo as entidades Go e seus testes pertencem às próximas etapas. Não inicializa módulo, não move o bootstrap existente e não cria interfaces Go sem consumidor. Tipos, cardinalidades, métodos, construtores e estados estão nos [contratos da issue #4](ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md); o catálogo de ações e as políticas continuam na #5. O [ADR 0007](adr/0007-contas-por-cnpj-e-recontratacao.md) atualiza o modelo para contas independentes por CNPJ e Pessoa como fonte de dados pessoais.

## Cadastros sob responsabilidade do módulo

| Conceito | Responsabilidade do identity-service | Limite |
|---|---|---|
| Usuário | Identidade de acesso, e-mail, confirmação e situação de acesso; associação aos cadastros conforme especificação futura | Não confundir confirmação de e-mail com identificação facial ou autenticação de aparelho. Login e sessões serão implementados depois. |
| Empresa | Cadastro do empregador por CNPJ: identificação, razão social, apelido e responsável legal; associação à referência do A1 | Cada CNPJ constitui seu próprio tenant. Não concentra toda configuração de jornada ou coleta só por pertencer à empresa. |
| Pessoa | Dados pessoais no CNPJ, referenciados por funcionário, usuário e empresa | Não é conta global nem perfil. Não sincronizar dados pessoais automaticamente entre CNPJs. |
| Responsável legal | Relação da empresa com Pessoa identificada por nome e CPF | Não equivale automaticamente a usuário, gestor ou signatário de todos os artefatos. |
| Referência do A1 | Associação à empresa e identificação do arquivo protegido, com metadados a especificar | Não contém senha nem conteúdo privado exposto aos consumidores. Upload, validação e custódia técnica são futuros. |
| Funcionário | Cadastro do trabalhador no CNPJ, referenciando Pessoa com CPF/nome/nascimento, com matrícula interna, setor e usuário opcionais e modalidade de coleta | Cadastro pode existir sem conta. Não contém marcações, resultados de apuração ou cadastro biométrico por implicação. |
| Setor | Agrupamento de funcionários da empresa usado para delimitar supervisão | Cadastro e escopo não substituem a verificação da ação sobre o recurso concreto. |
| Vínculo com empresa | Relação de autorização entre usuário e CNPJ, com perfil, permissões, escopo e situação | Não é vínculo empregatício, sessão ou simples empresa selecionada. Revogar acesso não equivale a excluir o cadastro do funcionário. |

A matrícula interna é estável, aleatória e não sequencial, com seis caracteres `0-9A-Z`, única por CNPJ inclusive no histórico. Sua geração cabe à biblioteca própria prevista no plano. O cadastro recebe a matrícula gerada; a futura coordenação com persistência garante unicidade e trata colisões. Recontratação cria novo funcionário/matrícula e, quando houver acesso, novo usuário, preservando os anteriores encerrados. A matrícula da folha é uma referência externa separada: o mapeamento para integração não torna o `identity-service` responsável por exportar a folha.

## Identidade, cadastro e autorização são relações distintas

- **Funcionário sem conta:** o cadastro do trabalhador pode existir sem usuário e, portanto, sem vínculo de acesso derivado desse usuário. O fluxo de criação/associação da conta permanece a especificar. No coletor coletivo, apresentar a matrícula declara uma identidade; não substitui comparação facial e prova de vida.
- **Gestor que também é funcionário:** a identidade de acesso é a mesma; o cadastro de funcionário representa o trabalhador no CNPJ e o vínculo representa suas autorizações. Não criar outra identidade apenas para separar esses papéis. Como acumular perfis e escopos será definido na issue #5, sem presumir autoaprovação de ajustes.
- **Gestor de dois CNPJs:** a mesma pessoa pode administrar A e B por contas e vínculos independentes. Alterar dados/bloquear a conta em A não altera B. Selecionar B não transporta para B permissões, sessão, funcionários ou setores de A.
- **Contabilidade/DP em duas empresas:** uma concessão para exportar PDF em A não permite exportar em B nem criar jornadas em A. Cada ação depende da concessão do gestor do CNPJ correspondente.
- **Supervisor e funcionário de outro setor:** o perfil Supervisor não basta para acessar o cadastro ou operação fora de seu escopo. Referências ao setor e ao funcionário precisam corresponder ao contexto autorizado.
- **Associação a usuário:** associar um usuário a um funcionário não concede automaticamente um perfil; ter vínculo de gestor não cria automaticamente cadastro de funcionário. Regras de unicidade da associação e sua cardinalidade final pertencem à issue #4.

Esses cenários detalham as separações já aprovadas. Os contratos da #4 definem encerramento/recontratação e associações, sem fixar exclusão, transferência entre empresas ou onboarding. Perfis simultâneos e escopos continuam na #5.

## Autorizações que o módulo representa

Os quatro perfis são Gestor, Supervisor, Colaborador e Contabilidade/DP. Não existe Gestor Master separado. O módulo representa a situação do acesso do usuário e do vínculo, o CNPJ, o perfil, as permissões por ação e o escopo necessário para avaliar uma operação. A matriz completa ainda será especificada; os exemplos de RF-06.1 não são uma concessão automática nem um catálogo exaustivo.

O `identity-service` administra os vínculos e concessões de seu domínio. O servidor que executa uma operação, inclusive em outro módulo, deve verificar vínculo e permissão no CNPJ e escopo do recurso, conforme ADR 0004. O `identity-service` fornece os dados de autorização; o consumidor conhece a ação e o recurso que pretende operar. Uma eventual política compartilhada não substitui essa verificação nas entradas, filas, arquivos e exportações.

Referências entre cadastros de funcionário, setor e vínculo devem ser coerentes com o tenant autorizado. Validação isolada de UUID ou presença de `tenant_id` não comprova essa coerência. O domínio não abre conexões nem consulta HTTP: coordenação e carregamento autoritativo pertencem aos futuros casos de uso e adaptadores.

## Dados fornecidos aos demais módulos

A tabela descreve necessidades semânticas, não endpoints, DTOs, eventos ou acesso direto a tabelas. Fornecer apenas os dados necessários à finalidade e ao escopo autorizado; possuir uma referência de usuário ou tenant não autoriza a consulta. Transporte, versionamento, consistência, propagação de revogações e eventual uso de cópias locais serão especificados antes das integrações.

| Consumidor por responsabilidade, sem fixar nomes de novos serviços | Dados necessários do identity-service | Responsabilidade que continua no consumidor |
|---|---|---|
| Registro e coleta de ponto | Empresa/tenant, referência e identificação do funcionário, matrícula interna e modalidade atribuída | Receber e preservar marcações; autenticar aparelho ou executar comparação facial e prova de vida conforme a modalidade; sincronizar e confirmar aceitação durável. |
| Tratamento e gestão de ponto | Referências de usuário, funcionário e setor; vínculo, perfil, permissões, escopo e situações de acesso | Verificar autorização sobre o ajuste concreto e aplicar o fluxo de aprovação, concorrência e preservação das marcações. |
| Jornadas e apuração | Referências de empresa e funcionário e dados cadastrais necessários; contexto de autorização para administração | Cadastrar horários/jornadas, associar regras de trabalho e calcular horas e saldos. Permitir criar jornada não significa possuir seus dados no identity-service. |
| Relatórios e integração com folha | Identificação do empregador e trabalhador, matrícula interna e contexto autorizado; referência do A1 quando aplicável e autorizado | Manter o mapeamento da matrícula da folha, gerar exportações e artefatos, resolver leiautes, signatários e processo de assinatura. Referência do A1 não dá acesso irrestrito ao arquivo ou à senha. |
| Autenticação e sessão futuras | Identidade de acesso, e-mail/confirmação, situação de acesso e vínculos elegíveis ao contexto | Executar o fluxo de autenticação e sessão e verificar autorização efetiva. A distribuição técnica dessa capacidade será detalhada depois; esta linha não cria outro serviço. |

CPF, nascimento, e-mail e outros dados pessoais não compõem um payload universal. A necessidade de cada consumidor será detalhada com controles de acesso e proteção conforme RNF-16; segredos do A1 e material biométrico não são dados cadastrais de livre distribuição.

## Fora desta issue e da primeira implementação do domínio

| Área | Delimitação |
|---|---|
| API e persistência | Sem handlers, rotas, banco, migrations ou repositórios concretos nesta issue. Chi, GORM e PostgreSQL continuam escolhas para camadas posteriores. |
| Login completo | RF-07/RF-08 continuam requisitos futuros: Magic Link versus OTP, login social, senhas, tokens, expiração, recuperação e revogação de sessões exigem especificação. Representar confirmação e situação de acesso no domínio não entrega esses fluxos. |
| A1 | A exigência cadastral de RF-03 permanece. Upload, validação do arquivo, armazenamento e proteção da senha não são implementados; estados de cadastro incompleto serão definidos na issue #4. A3 continua futuro. |
| Biometria e aparelhos | A modalidade pertence ao cadastro do funcionário. Cadastro facial, modelos, embeddings, prova de vida, credenciais do dispositivo, fallback e execução offline não são implementados nem recebem proprietário técnico definitivo aqui. |
| Ponto e cálculos | Sem marcações, ajustes, jornadas, cálculo, banco de horas ou geofencing. Identidade e autorização fornecem contexto; não determinam resultados dessas regras. |
| Relatórios e folha | Sem AFD, AEJ, espelho, comprovante, assinatura, exportação Domínio ou mapeamento de rubricas. A referência do A1 não resolve conformidade ou os conflitos de RF-23 a RF-25. |
| Integrações | Sem HTTP entre serviços, filas, eventos ou sincronização de cadastros. RNF-09 a RNF-12 não definem por si um contrato seguro e atualizado de autorização. |

## Rastreabilidade e próximas especificações

| Fonte | Responsabilidade coberta | Continuidade / pendência |
|---|---|---|
| RF-01; ADR 0004 | Empresa por CNPJ e vínculos independentes por tenant | Verificação autoritativa também nos consumidores; estratégia de persistência e propagação de revogação futuras. |
| RF-02 | Gestão de acesso do perfil Gestor | Entidades e métodos na #4; políticas na #5, sem Gestor Master. |
| RF-03 | Cadastro da empresa, responsável com nome/CPF e referência do A1 | Estados incompletos na #4; custódia e upload em etapa posterior; signatários seguem na revisão. |
| RF-04 | Empresa selecionada distinta de autorização | Fluxo de seleção/sessão na implementação futura do login. |
| RF-05 e RF-05.1 | Funcionário independente de usuário e matrícula interna | Cardinalidades na #4; formato e unicidade da matrícula antes da #10; cadastro na #13. |
| RF-06 e RF-06.1 | Quatro perfis, escopo e concessão por ação/CNPJ | Matriz, acúmulo de perfis, autoaprovação e acesso a atestados na #5 ou especificação da área afetada; implementação na #14. |
| RF-07 | Confirmação de e-mail da identidade de acesso | Estado no domínio na #12; mecanismo e onboarding futuros. |
| RF-08 | Identidade e vínculos como contexto da sessão | Sessão de longa duração, tokens e segurança fora da primeira rodada de domínio. |
| RF-11.1; ADR 0005 | Modalidade atribuída por empresa e funcionário | Sem implementar coleta; mudança de modalidade, cadastro inicial, fallbacks e offline continuam pendentes nas fontes próprias. |
| RF-19 | Acesso da Contabilidade/DP ao mapeamento De-Para e demais funcionalidades conforme concessões | Dados de autorização no identity-service; interface e operações do mapeamento na integração futura, sujeitas às políticas da #5. |
| RF-20 | Separação entre matrícula interna e matrícula da folha | Cadastro fornece matrícula interna; mapeamento externo pertence à integração futura, cujo contrato e leiaute ainda exigem validação. |
| ADR backend 0001 | Domínio independente e dependências explícitas | #4 define o modelo; #6 avalia dependências; #7 inicializa módulo/comandos; #8 a #14 implementam; #15 revisa a rodada. |

As pendências acima não impedem delimitar responsabilidades, mas impedem tratá-las como comportamento implementável sem sua especificação. Em particular, não se conclui aqui que um cadastro incompleto está apto a operar, que uma sessão antiga mantém permissões revogadas ou que um relatório está juridicamente válido.

## Verificação desta entrega

Verificações executadas nesta entrega:

- `git diff --check`, sem erros de whitespace.
- Checagem automatizada com Python das referências locais nos três documentos alterados e dos IDs RF citados neste detalhamento: 17 referências locais existentes, 15 IDs RF válidos e presença de todos os requisitos exigidos pela issue na matriz de rastreabilidade.
- Revisão documental dos cenários de funcionário sem conta, gestor também funcionário, múltiplos CNPJs e Contabilidade/DP sem concessão, confrontados com plano, vocabulário e ADRs aceitos.

Não há alteração Go nem módulo/suíte executável para esta entrega documental; testes Go e análise estática não foram executados. As skills `golang-testing` e `golang-error-handling` foram consultadas; sua aplicação ao código fica para as respectivas issues, com TDD e comandos reais registrados quando o módulo existir. Revisão e homologação pelo fluxo do repositório permanecem pendentes.
