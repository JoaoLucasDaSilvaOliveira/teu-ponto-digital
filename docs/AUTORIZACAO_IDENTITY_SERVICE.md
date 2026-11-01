# Perfis, permissões e escopos — identity-service

Especificação da [issue #5](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/5), parte da issue pai #2. Em elaboração: as invariantes comuns abaixo já decorrem das fontes aceitas; composição de perfis, escopo do Supervisor e delegação administrativa estão em esclarecimento. Não representa implementação Go nem autorização de produção.

Fontes: [RF-01, RF-04, RF-06, RF-06.1, RF-13 e RF-19](REQUISITOS_FUNCIONAIS.md), [contratos da issue #4](ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md), [limites](LIMITES_IDENTITY_SERVICE.md), [plano](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md), [vocabulário](../CONTEXT.md), ADRs [0004](adr/0004-isolamento-por-cnpj.md) e [0007](adr/0007-contas-por-cnpj-e-recontratacao.md) e [arquitetura do backend](../services/docs/adr/0001-arquitetura-backend.md). A [revisão de coerência](REVISAO_DOCUMENTACAO.md) mantém as lacunas ainda não resolvidas.

## Invariantes comuns

- Existem exatamente os perfis Gestor, Supervisor, Colaborador e Contabilidade/DP. Responsável legal é relação cadastral, não perfil; não existe Gestor Master.
- Conta pertence a um CNPJ. Conta ativa com e-mail confirmado e vínculo ativo são condições distintas, todas necessárias para uma ação de usuário autorizada por este módulo.
- O CNPJ do contexto, da conta, do vínculo e do recurso deve coincidir. O vínculo precisa pertencer à conta avaliada. UUID válido e empresa selecionada no cliente não comprovam autorização.
- E-mail, CPF ou igualdade de dados pessoais entre empresas não relacionam contas para fins de autorização. Bloqueio, encerramento e revogação não podem ser contornados selecionando outro perfil ou usando referência histórica.
- Ausência de concessão, ação desconhecida, recurso incompatível, referência não resolvida ou escopo insuficiente resultam em rejeição. Não interpretar ação inexistente como ação genérica permitida.
- Contabilidade/DP só executa ações explicitamente concedidas pelo gestor daquele CNPJ. Exportar PDF, criar horários e criar jornadas são concessões distintas; uma não implica a outra.
- Preservar o vínculo entre ação e recurso concreto, inclusive em arquivos, exportações, jobs e filas. Consulta/listagem precisa restringir resultados; autorizar a tela não autoriza todos os registros retornados.
- Autoaprovação de ajuste e acesso ao conteúdo de atestados continuam pendentes e são rejeitados nesta rodada. Um perfil administrativo não resolve essas pendências implicitamente.

## Autoridade dos dados e limites

O futuro caso de uso carrega conta/vínculo atuais e o recurso no CNPJ autorizado, resolve suas relações e solicita a decisão ao domínio. Dados vindos do cliente ou do cabeçalho de uma mensagem não são prova de vínculo, perfil ou escopo. A política pura não abre conexão e não consulta HTTP, ORM ou broker.

O consumidor verifica as relações que conhece: ajuste e seu funcionário, arquivo e seu proprietário, funcionário e setor. Se um arquivo corresponde a vários funcionários, autorização sobre um deles não autoriza o arquivo inteiro. A política nunca infere propriedade apenas do ID do solicitante informado no transporte.

Permissão para uma ação não elimina suas precondições de negócio: autorização de exportação não valida A1, leiaute, signatário ou fechamento; autorização de ajuste não resolve concorrência e sincronização tardia. Também não permite editar uma marcação original.

Esta política cobre ações de contas de usuário. Funcionário sem conta no coletor coletivo não recebe um usuário/perfil fictício: autenticação do coletor, identificação facial, prova de vida e autorização do dispositivo exigem o fluxo próprio. Ter matrícula ou saber o CPF não autoriza uma ação de usuário sobre aquele funcionário.

## Verificação da rodada

Entrega documental parcial. Permanecem por fechar a matriz do catálogo inicial, os valores de configuração, a interface comum, os adaptadores e os cenários específicos dos perfis. Não há módulo ou código Go alterado nesta rodada.
