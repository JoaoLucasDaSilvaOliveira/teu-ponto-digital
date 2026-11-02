# Perfis, permissões e escopos — identity-service

Especificação documental da [issue #5](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/5), parte da issue pai #2, consolidada com as decisões confirmadas até 02/11/2026. Define os contratos e cenários dos quatro perfis, sem implementar código Go ou autorização de produção. A [issue #17](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/17) inventaria todos os acessos e serviços da aplicação, padrões e concessões; esta especificação cobre o catálogo inicial explicitado abaixo, rejeitando ações indefinidas.

Fontes: [RF-01, RF-04, RF-06, RF-06.1, RF-13 e RF-19](REQUISITOS_FUNCIONAIS.md), [contratos da issue #4](ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md), [limites](LIMITES_IDENTITY_SERVICE.md), [plano](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md), [vocabulário](../CONTEXT.md), ADRs [0004](adr/0004-isolamento-por-cnpj.md) e [0007](adr/0007-contas-por-cnpj-e-recontratacao.md) e [arquitetura do backend](../services/docs/adr/0001-arquitetura-backend.md). A [revisão de coerência](REVISAO_DOCUMENTACAO.md) mantém as lacunas ainda não resolvidas.

## Invariantes comuns

- Existem exatamente os perfis Gestor, Supervisor, Colaborador e Contabilidade/DP. Responsável legal é relação cadastral, não perfil; não existe Gestor Master.
- Conta pertence a um CNPJ. Conta ativa com e-mail confirmado e vínculo ativo são condições distintas, todas necessárias para uma ação de usuário autorizada por este módulo.
- O CNPJ do contexto, da conta, do vínculo e do recurso deve coincidir. O vínculo precisa pertencer à conta avaliada. UUID válido e empresa selecionada no cliente não comprovam autorização.
- E-mail, CPF ou igualdade de dados pessoais entre empresas não relacionam contas para fins de autorização. Bloqueio, encerramento e revogação não podem ser contornados selecionando outro perfil ou usando referência histórica.
- Ausência de concessão, ação desconhecida, recurso incompatível, referência não resolvida ou escopo insuficiente resultam em rejeição. Não interpretar ação inexistente como ação genérica permitida.
- Contabilidade/DP só executa ações explicitamente concedidas pelo gestor daquele CNPJ. Exportar PDF, criar horários e criar jornadas são concessões distintas; uma não implica a outra.
- Preservar o vínculo entre ação e recurso concreto, inclusive em arquivos, exportações, jobs e filas. Consulta/listagem precisa restringir resultados; autorizar a tela não autoriza todos os registros retornados.
- Autoaprovação de ajuste é proibida, inclusive por troca de conta/perfil. Nenhum perfil administrativo ou concessão permite ignorar essa proibição. Tratar ocorrência com atestado inclui visualizar o documento e decidir/lançar o abono correspondente; relatório/dashboard e avaliação genérica de ajuste não concedem esse acesso.

## Decisões esclarecidas sobre os perfis

- **Colaborador:** ações sobre os próprios registros, incluindo marcação, consulta/dashboard próprio, solicitação de relatório e de ajuste. As precondições da coleta, o cadastro ativo e a associação da conta ao funcionário continuam necessários.
- **Supervisor:** inclui automaticamente as capacidades próprias do Colaborador, por composição de comportamento, sem exigir duas identidades ou duplicar dados pessoais. Deve ter cadastro de funcionário para essas ações próprias; isso não lhe permite agir como colaborador de outro funcionário.
- **Gestão da equipe:** o Gestor concede ao Supervisor as ações permitidas sobre subordinados. O perfil sozinho não concede administração irrestrita. O catálogo deve contemplar administração de cadastros de subordinados, horários/jornadas da equipe, consulta/dashboard, relatórios e avaliação de ajustes; cada ação depende da concessão e do recurso dentro da equipe.
- **Hierarquia de decisões:** Supervisor pode avaliar ajustes de subordinados se autorizado; seu próprio pedido deve ser avaliado por Gestor. Outro perfil ou outra conta da mesma Pessoa não pode ser usado para decidir o próprio pedido. A identidade cadastral da Pessoa e as referências históricas precisam ser consideradas, não apenas comparar IDs de login.
- **Lançamentos sobre si:** Supervisor não pode lançar/decidir abonos ou outros lançamentos administrativos sobre si, mesmo possuindo concessão sobre subordinados. Solicitar ajuste próprio continua uma ação distinta de avaliá-lo.
- **Delegação:** Gestor administra acessos e concede ações a Supervisor e Contabilidade/DP no CNPJ. Esses perfis não concedem acesso, não se promovem a Gestor e não delegam suas concessões a terceiros. Responsabilização/auditoria não substitui a verificação prévia do escopo.
- **Gestor:** perfil exclusivo, incompatível com Colaborador, Supervisor e Contabilidade/DP na mesma conta. Não concede ações próprias de ponto, solicitação de ajuste ou lançamento de atestado próprio, mesmo que exista cadastro de funcionário associado. Essa decisão de acesso não determina dispensa legal de controle de jornada nem apaga cadastros/histórico. Substitui a possibilidade de combinação de Gestor com Colaborador discutida anteriormente.

### Equipe por atribuição direta

Cada funcionário tem zero ou um Supervisor responsável no CNPJ. Equipe é o conjunto de funcionários atribuídos diretamente a esse Supervisor, podendo atravessar setores. Integrar o mesmo setor não concede acesso automaticamente; mudar o setor não transfere subordinação.

Criar funcionário sob o perfil Supervisor exige concessão específica e já o vincula à sua equipe. Essa capacidade não permite capturar um funcionário existente sem responsável ou de outra equipe. Transferência/atribuição de subordinados existentes exige ação do Gestor; a referência de supervisão e suas validações são definidas abaixo, preservando histórico e isolamento por CNPJ.

CRUD de subordinados não permite apagar registros originais ou histórico. Encerramento de cadastro deve seguir os contratos da #4. Cadastro pessoal e de funcionário não equivale a administrar conta de acesso, vínculo, perfil ou concessão; estas últimas ações permanecem do Gestor. Operações que alterem recursos compartilhados precisam verificar todo o alcance da alteração, sem atingir outra equipe por referência indireta.

### Perfil do produto e enquadramento trabalhista

O perfil Gestor é autorização administrativa, não uma classificação jurídica. O [art. 62, II e parágrafo único da CLT](https://www2.camara.leg.br/legin/fed/declei/1940-1949/decreto-lei-5452-1-maio-1943-415500-normaatualizada-pe.html) trata do exercício de cargo de gestão e da condição remuneratória ali prevista; não estabelece dispensa universal de controle de jornada para qualquer pessoa chamada gerente. Não derivar isenção de jornada do nome do perfil. A especificação do enquadramento individual e suas consequências para apuração continuam na área de jornada, fora desta issue.

## Autoridade dos dados e limites

O futuro caso de uso carrega conta/vínculo atuais e o recurso no CNPJ autorizado, resolve suas relações e solicita a decisão ao domínio. Dados vindos do cliente ou do cabeçalho de uma mensagem não são prova de vínculo, perfil ou escopo. A política pura não abre conexão e não consulta HTTP, ORM ou broker.

O consumidor verifica as relações que conhece: ajuste e seu funcionário, arquivo e seu proprietário, funcionário e setor. Se um arquivo corresponde a vários funcionários, autorização sobre um deles não autoriza o arquivo inteiro. A política nunca infere propriedade apenas do ID do solicitante informado no transporte.

Permissão para uma ação não elimina suas precondições de negócio: autorização de exportação não valida A1, leiaute, signatário ou fechamento; autorização de ajuste não resolve concorrência e sincronização tardia. Também não permite editar uma marcação original.

Esta política cobre ações de contas de usuário. Funcionário sem conta no coletor coletivo não recebe um usuário/perfil fictício: autenticação do coletor, identificação facial, prova de vida e autorização do dispositivo exigem o fluxo próprio. Ter matrícula ou saber o CPF não autoriza uma ação de usuário sobre aquele funcionário.

## Catálogo inicial e matriz

Cada código é uma ação fechada; não há prefixo/curinga que autorize uma ação nova. As listas na mesma célula enumeram ações distintas. Legenda: **próprio** exige funcionário ativo associado à conta/Pessoa; **equipe C** exige concessão expressa do Gestor e subordinado direto; **empresa C** exige concessão expressa no CNPJ; **empresa** é padrão do Gestor para aquela ação conhecida; **—** é rejeitado. Sem concessões, Supervisor conserva apenas ações próprias e Contabilidade/DP não ganha ações de negócio.

| Ação | Gestor | Supervisor | Colaborador | Contabilidade/DP |
|---|---|---|---|---|
| `PONTO_REGISTRAR` | — | próprio | próprio | — |
| `PONTO_CONSULTAR`, `DASHBOARD_CONSULTAR` | empresa | próprio; equipe C | próprio | — |
| `RELATORIO_PDF_EXPORTAR` | empresa | próprio; equipe C | próprio | empresa C |
| `AJUSTE_SOLICITAR` | — | próprio | próprio | — |
| `AJUSTE_AVALIAR` | empresa, exceto próprio | equipe C, exceto próprio e pedido de Supervisor | — | — |
| `ATESTADO_ENVIAR` | — | próprio | próprio | — |
| `ATESTADO_TRATAR` | empresa, exceto próprio | equipe C, exceto próprio e pedido de Supervisor | — | — |
| `FUNCIONARIO_CRIAR`, `FUNCIONARIO_CONSULTAR`, `FUNCIONARIO_ALTERAR`, `FUNCIONARIO_ENCERRAR` | empresa | equipe C, com as restrições de criação/alteração abaixo | — | — |
| `HORARIO_CRIAR`, `JORNADA_CRIAR` | empresa | equipe C | — | empresa C |
| `HORARIO_CONSULTAR`, `HORARIO_ALTERAR`, `HORARIO_ENCERRAR` | empresa | equipe C | — | — |
| `JORNADA_CONSULTAR`, `JORNADA_ALTERAR`, `JORNADA_ENCERRAR` | empresa | equipe C | — | — |
| `ACESSO_ADMINISTRAR`, `CONCESSAO_DEFINIR`, `EQUIPE_ATRIBUIR` | empresa | — | — | — |

A tabela possui 23 códigos de ação. CRUD é representado por códigos independentes, e encerramento não significa exclusão física. A configuração rejeita concessões para ações fora das células C do perfil: conceder exportação não concede administração de contas, consulta de documentos médicos, criação de horários ou jornadas. Concessão não é executada pelo mero fato de constar no cadastro; precisa do ator Gestor autorizado no futuro caso de uso.

`ACESSO_ADMINISTRAR` agrupa nesta rodada operações administrativas sobre contas/vínculos; a #17 deverá decompor os acessos concretos antes de suas entradas serem implementadas. Não é uma permissão genérica para negócio. Mapeamento de folha, demais formatos de exportação, configurações corporativas não enumeradas e outras ações ficam pendentes na #17 e rejeitados. A criação de horários/jornadas pela Contabilidade/DP é concedível por RF-06.1; não presume o restante do CRUD nem acesso geral a dashboards. Os consumidores deverão fornecer os dados estritamente necessários à operação já autorizada, sem expor uma consulta administrativa adicional por implicação.

Gestor exclusivo não pode usar alcance de empresa para executar uma operação pessoal de ponto/ajuste/atestado. Administração de acesso não apaga funcionário nem determina regime de jornada. O conjunto de perfis simultâneos não foi autorizado nesta rodada: cada configuração possui um perfil, e Supervisor reutiliza internamente as capacidades próprias do Colaborador. Combinações adicionais são rejeitadas até definição explícita na #17; isso não cria outra identidade nem duplica Pessoa.

### Atestado e abono como capacidade conjunta

`ATESTADO_TRATAR` permite abrir o documento da ocorrência autorizada e decidir/lançar seu abono correspondente. Não existe concessão independente para abonar com base em atestado sem poder visualizá-lo. O documento precisa pertencer à mesma ocorrência/funcionário e ao CNPJ autorizado; a capacidade não permite navegar em todos os anexos da pessoa ou exportá-los por meio de relatório genérico.

O fluxo de negócio deve apresentar o documento antes da decisão e vincular abono, justificativa, ocorrência e ator na auditoria futura. Abrir o arquivo não comprova leitura humana, autenticidade médica ou validade jurídica. A política de acesso não implementa essas validações. Lançamentos sem atestado, revisão/cancelamento de abonos e regras clínicas/retenção continuam para suas especificações e para o catálogo da #17, sem concessão automática. O envio pelo próprio trabalhador é uma ação separada de tratar/decidir a ocorrência.

### Escopos, coleções e recursos compartilhados

- **Próprio:** conta, Pessoa e funcionário associados no CNPJ, com cadastro de funcionário ativo. Um papel administrativo não permite se passar por outro trabalhador. Consulta a histórico pelo usuário após encerramento requer fluxo próprio, fora das ações próprias desta rodada.
- **Equipe:** funcionário aponta para o vínculo Supervisor do ator; ser do mesmo setor não basta. A permissão vale para a operação sobre cada subordinado, sem incluir o Supervisor como beneficiário de operações administrativas.
- **Empresa:** recurso no CNPJ da conta/vínculo. Não se estende a filiais com outro CNPJ, ainda que usem os mesmos dados pessoais ou e-mail.
- **Criação:** não exigir ID de recurso inexistente. O alvo de criação inclui empresa e alcance pretendido. Funcionário criado por Supervisor recebe sua referência de supervisão na mesma operação; não aceitar alvo de outra equipe ou mudança indireta de acesso.
- **Lotes, dashboards e arquivos:** todos os registros materializados precisam estar no escopo. Aplicar filtro autorizado antes de agregar, gerar ou exportar; downloads conferem novamente o alcance do arquivo. Se houver recurso fora do escopo, rejeitar o lote explícito inteiro; consulta exploratória pode retornar somente o conjunto autorizado, sem vazar totais externos.
- **Horários/jornadas compartilhados:** Supervisor não altera/encerra definição que afete pessoa de outra equipe ou seu próprio tratamento. Caso de uso verifica referências e impacto completo, inclusive associações futuras; em dúvida rejeita. Definição global da empresa não se torna da equipe só porque um subordinado a utiliza. Criação restrita à equipe deve preservar esse alcance; não atribuir um conjunto vazio a alcance irrestrito.
- **Dados pessoais compartilhados:** alterar cadastro do subordinado não autoriza modificar conta, vínculo, responsável legal ou dados que afetem cadastros fora da equipe. Resolver referências antes de modificar; se a alteração de Pessoa tiver alcance maior, exigir Gestor.

## Referência de supervisão no Funcionário

Acrescentar à entidade da #4 o campo privado opcional `supervisorVinculoID uuid.UUID`. Ele referencia um vínculo específico do perfil Supervisor, no mesmo CNPJ, pertencente a usuário ativo e confirmado. Esse usuário deve estar associado a funcionário ativo para representar Supervisor nesta rodada. Não usar setor como substituto dessa referência, nem tornar Equipe uma nova entidade com lista duplicada.

Contratos adicionais: `AtribuirSupervisor(supervisor vo.ReferenciaSupervisor) error`, `RemoverSupervisor() error` e `SupervisorVinculoID() (uuid.UUID, bool)`. O valor é criado por `NovaReferenciaSupervisor(empresaID, vinculoID, usuarioID, pessoaID uuid.UUID) (ReferenciaSupervisor, error)`, com campos privados, IDs não nulos e getters por valor. Não carrega dados pessoais duplicados nem prova que a conta/vínculo está atualmente habilitada.

Métodos locais validam o valor, tenant, funcionário não encerrado e ausência de autorreferência pela Pessoa; só atribuem depois de validar, sem substituir o estado em erro. O caso de uso futuro carrega usuário e vínculo correspondentes e exige conta ativa/confirmada, perfil Supervisor, vínculo ativo, funcionário ativo associado ao Supervisor, autorização do Gestor e referências coerentes. Também verifica ciclos e concorrência. Não aceitar IDs do cliente como prova dessas condições. Usar o valor de referência evita exigir implementação de Vínculo/políticas da #14 para compilar Funcionário na #13; a validação autoritativa continua na coordenação.

A referência é ao vínculo para que revogação, troca de perfil ou nova concessão não herdem uma equipe silenciosamente. Vínculo revogado ou que deixou de ser Supervisor não permite atuar sobre referências antigas. Restabelecer acesso exige nova atribuição explícita pelo Gestor quando o vínculo for substituído. Troca de setor não muda essa referência. Encerramento preserva referência/histórico; recontratação cria novo funcionário e não transporta supervisor automaticamente.

Esses métodos não são entradas públicas autorizadas por si: Supervisor só participa da atribuição inicial durante criação concedida; transferência ou apropriação de existente é do Gestor. A hierarquia não concede supervisão transitiva. Mesmo se um Supervisor estiver cadastralmente sob outro, seus próprios pedidos de ajuste/atestado exigem Gestor.

## Valores da configuração e consulta

Todos os valores têm campos privados, valor zero inválido, construtores com erro, coleções defensivas e normalização sem inventar concessões. IDs são `uuid.UUID`; enumerações desconhecidas são rejeitadas. Os nomes abaixo completam os contratos da #4; implementá-los apenas nas issues de código correspondentes.

| Valor | Construtor e conteúdo |
|---|---|
| `Perfil` | `NovoPerfil(codigo string) (Perfil, error)`; somente `GESTOR`, `SUPERVISOR`, `COLABORADOR`, `CONTABILIDADE_DP`. |
| `Acao` | `NovaAcao(codigo string) (Acao, error)`; somente os códigos da matriz. `Codigo() string`. Código desconhecido não entra na política como string livre. |
| `ConfiguracaoAutorizacao` | `NovaConfiguracaoAutorizacao(perfil Perfil, concessoes []Acao) (ConfiguracaoAutorizacao, error)`; um perfil e conjunto sem duplicatas de ações concedíveis naquele perfil. `Perfil() Perfil`, `Concessoes() []Acao`, `Concede(acao Acao) bool`, `Validar() error`. |
| `AlvoAutorizacao` | Referência resolvida do recurso e alcance: empresa, tipo, ID existente ou intenção de criação, funcionário/Pessoa beneficiários quando aplicáveis, referência de supervisão e alcance coletivo. Sem documento/segredo no valor. Construção valida os campos exigidos para o tipo; ausência de dados necessários é inválida, não equivale a empresa inteira. |
| `PedidoAutorizacao` | `NovoPedidoAutorizacao(empresaID uuid.UUID, acao Acao, alvo AlvoAutorizacao) (PedidoAutorizacao, error)`; valida correspondência entre ação, tipo e alcance. Getters retornam os valores imutáveis. Não recebe perfil/concessões do transporte. |

Gestor/Colaborador usam conjunto vazio de concessões: seus padrões são fechados pela matriz; não se cadastra curinga `TODAS`. Supervisor recebe apenas ações C de equipe, além das ações próprias fixas. Contabilidade/DP recebe apenas as ações C de empresa. Construtores rejeitam tentativa de conceder ação incompatível, múltiplos perfis ou enumeração desconhecida. Alterar configuração substitui o valor inteiro atomicamente e não preserva capacidades incompatíveis do perfil anterior.

O construtor específico do alvo e sua conversão de DTO/evento pertencem ao módulo que conhece o recurso, não a um DTO global que o cliente possa preencher como prova de autoridade. Antes de cada decisão, ele deve fornecer relações verificadas, todos os beneficiários/alcance e se o pedido exige Gestor por ser de Supervisor. Esse nível de aprovação não pode ser reduzido apenas mudando o perfil após solicitar o ajuste. O contrato exige que essas evidências venham do estado autoritativo/histórico da operação; mecanismo de armazenamento será definido nos casos de uso.

## Interface comum e adaptadores

Em `domain/entity`, o vínculo consome o contrato pequeno abaixo. Tipos de valor são do pacote de valores; esse pacote não importa entidades ou políticas, evitando ciclo de imports.

```go
type PoliticaAutorizacao interface {
    Permite(config vo.ConfiguracaoAutorizacao, pedido vo.PedidoAutorizacao) bool
}
```

Método comum: `(*Vinculo).Autorizar(usuario *Usuario, funcionarioAtor *Funcionario, pedido vo.PedidoAutorizacao, politica PoliticaAutorizacao) error`. Recebe funcionário do ator quando a operação/perfil exigir; Gestor e Contabilidade/DP podem não ter funcionário. Antes de delegar, valida conta ATIVA/confirmada, vínculo ATIVO daquela conta, dados válidos e todos os tenants iguais. Quando exigido, verifica associação do funcionário ativo à mesma conta/Pessoa. Rejeita política ausente, autoaprovação/lançamento em benefício próprio e tentativa de decidir pedido de Supervisor sem Gestor. A política não recebe autoridade para sobrepor essas rejeições. Em seguida verifica se o adaptador permite ação e alcance; negação retorna erro, sem alterar estado.

Adaptadores de `internal/domain/policy`: `PoliticaGestor`, `PoliticaSupervisor`, `PoliticaColaborador` e `PoliticaContabilidadeDP`, todos com `Permite`. Cada um rejeita configuração cujo perfil não corresponda ao seu. Gestor permite somente padrões administrativos definidos; Colaborador somente próprios; Supervisor compõe internamente a mesma regra de ações próprias do Colaborador e acrescenta equipe/concessão; Contabilidade/DP exige concessão em alcance de empresa. Não duplicar dados pessoais ou introduzir herança de entidades.

O bootstrap seleciona os quatro adaptadores em registro por `Perfil`, validado antes de servir entradas, e injeta o contrato no consumidor. Seleção central por código de perfil é composição; consumidor/caso de uso não usa type assertion, reflexão ou `switch` de tipo para decidir permissões. Configuração inválida/sem adaptador falha fechada. A identidade do recurso, do tenant, da conta e do vínculo é sempre verificada pelo método comum, inclusive se um adaptador incorreto for injetado.

Categorias de erro distinguíveis por `errors.Is`: `ErrContextoInvalido`, `ErrTenantDivergente`, `ErrVinculoIncompativel`, `ErrContaSemAcesso`, `ErrAcaoNaoPermitida`, `ErrAutoaprovacao` e `ErrHierarquiaInsuficiente`; reaproveitar as categorias da #4 quando equivalentes, sem duplicar sentinelas para o mesmo significado. Configuração/valor inválido usa erro de validação com contexto, sem dados pessoais. Não logar e retornar o mesmo erro; tradução para HTTP/evento é futura.

## Revogação e pontos de verificação

Retirar concessão, trocar perfil/equipe, bloquear conta ou revogar vínculo deve afetar a próxima decisão com dados atuais. Não manter acesso com base apenas no token de login ou no fato de o job ter sido agendado quando permitido. Revalidar antes de executar/materializar/exportar e antes de baixar arquivo; mutações exigem controle contra mudança concorrente entre verificar e gravar. Lote não concluído não ganha imunidade à revogação. Estratégia técnica de sessão, cache, filas e atualização atômica continua futura.

Auditoria registra ator e alteração; não concede poder. A concessão do próprio Supervisor/Contabilidade não é autoadministrável, mesmo quando a pessoa possui responsabilidade cadastral na empresa. Primeiro Gestor é bootstrap de onboarding ainda a especificar, não promoção pública/autoconcessão. A #17 deve fechar proteção do último Gestor, recuperação de acesso e demais operações administrativas concretas antes de implementá-las.

## Cenários para TDD na implementação da #14

1. Cada adaptador aceita apenas seu perfil; Gestor Master, ação desconhecida, valor zero, configuração de vários perfis e concessão incompatível falham.
2. Gestor administra as ações empresariais enumeradas e concede ações aos destinatários compatíveis; não registra ponto ou solicita ajuste/atestado próprio.
3. Colaborador consulta/registra seus dados, dashboard, PDF, ajuste e envio de atestado; não opera sobre outro funcionário, concede acesso ou avalia o próprio pedido.
4. Supervisor sem concessões tem as mesmas ações próprias do Colaborador, mas não lê/decide/exporta dados da equipe só por possuir o perfil.
5. Supervisor com uma concessão atua somente na ação concedida e em subordinado direto. Mesmo setor sem referência de supervisão falha; subordinados de setores diferentes podem pertencer à equipe.
6. Supervisor pode criar funcionário já em sua equipe se autorizado; não captura existente, transfere subordinado, altera acessos nem promove perfis. Alterar recurso compartilhado fora de alcance falha.
7. Ajuste/abono do próprio Supervisor é rejeitado por ele e por outro Supervisor, mesmo se houver relação cadastral de equipe; Gestor diferente da Pessoa beneficiária pode avaliar. Conta histórica/troca de perfil não contorna a regra.
8. Concessão de tratar atestado permite consultar o documento daquela ocorrência e lançar o abono correspondente; PDF/dashboard/ajuste não libera o anexo. Sem concessão, outro funcionário ou outra ocorrência falha; conteúdo não vai no pedido da política.
9. Contabilidade/DP sem concessões não possui ações de negócio. Conceder PDF em A não permite criar horário/jornada em A nem exportar em B; cada ação/CNPJ exige configuração própria.
10. Conta bloqueada, não confirmada, encerrada, vínculo revogado, vínculo de outro usuário e qualquer referência de outro CNPJ falham antes da política específica. Reutilizar CPF/e-mail em outro CNPJ não altera o resultado.
11. Revogar concessão ou vínculo e transferir equipe muda a próxima decisão. Novo vínculo Supervisor não herda referências da equipe anterior; recontratação não herda autorização nem supervisão.
12. Getters retornam cópias; alterações inválidas mantêm configuração/referências anteriores. Autorizar não muda entidades e testes não precisam de HTTP/banco.
13. Lote explícito com um alvo fora do escopo é rejeitado; dashboard não agrega totais externos; permissão de gerar arquivo não equivale a permissão eterna de download. Consultas precisam filtrar antes de materializar.

Aplicar testes de interface e matriz nas issues de código; os cenários de dados autoritativos, concorrência, revogação e arquivos precisam de testes dos futuros casos de uso/adaptadores, além dos testes puros de domínio. Não declarar segurança ponta a ponta com os testes de política isolados.

## Rastreabilidade e limites da entrega

| Fonte | Definição nesta issue |
|---|---|
| RF-01/RF-04; ADRs 0004/0007 | Conta, contexto, vínculo e recurso no mesmo CNPJ, sem conta/permissão global. |
| RF-06/RF-06.1 | Quatro perfis, padrões fechados, concessões explícitas de Gestor, Supervisor compondo ações próprias e equipe direta. |
| RF-13 | Proibição de autoaprovação; pedidos de Supervisor exigem Gestor; tratamento auditável continua no módulo de ponto. |
| RF-14 | Concessão de tratar ocorrência com atestado inclui documento e abono; não concede lançamentos próprios ou navegação irrestrita em anexos. |
| RF-19/RF-20 | Contabilidade/DP sem concessões implícitas; catálogo de mapeamento/integração será fechado na #17. |
| ADR backend 0001; issue #4 | Valores imutáveis, interface pequena, adaptadores específicos e validações comuns independentes de infraestrutura. |

As decisões confirmadas substituem as pendências de autoaprovação e da dependência entre atestado/abono: autoaprovação passa a proibida; acesso administrativo ao documento fica vinculado à ocorrência e à capacidade de tratá-la. Permissões adicionais, serviços e padrões não enumerados são trabalho da [#17](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/17), sem liberação provisória. Cadastro facial, autoria do coletor, política offline, signatários, retenção e login continuam em suas fontes próprias.

## Verificação da entrega

Verificações executadas: `git diff --check`; checagem com Python nos sete documentos desta rodada, com 65 referências locais existentes e 37 IDs RF válidos; 23 códigos distintos na matriz; concessões explícitas de PDF, criar horário e criar jornada para Contabilidade/DP; presença dos quatro adaptadores. Revisão documental confrontou matriz, hierarquia, supervisão direta, alvos compartilhados, rejeições e contratos da #4. A checagem do estado versionado também deve considerar o caminho original do RNF, conforme a pendência abaixo.

Há uma referência preexistente no plano ao arquivo RNF com espaço no nome, ausente no checkout devido à renomeação local alheia. A alteração foi preservada fora desta entrega; o caminho original continua no estado versionado. A criação da issue #17 foi recusada pelo conector GitHub com HTTP 403 e concluída via `gh issue create --body-file`; `gh issue view 17` confirmou título, estado aberto, URL e responsável.

Nenhum módulo/código Go criado; testes, formatação e análise estática Go não foram executados nesta entrega documental. Implementação na #14, catálogo completo na #17 e revisão/homologação permanecem pendentes. Não declarar autorização ponta a ponta ou fluxos de atestado implementados.
