# Perfis, permissões e escopos — identity-service

Especificação da [issue #5](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/5), parte da issue pai #2. Em elaboração: equipe por atribuição direta e exclusividade do Gestor foram confirmadas em 02/11/2026; a separação das concessões de atestados continua em esclarecimento. Não representa implementação Go nem autorização de produção.

Fontes: [RF-01, RF-04, RF-06, RF-06.1, RF-13 e RF-19](REQUISITOS_FUNCIONAIS.md), [contratos da issue #4](ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md), [limites](LIMITES_IDENTITY_SERVICE.md), [plano](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md), [vocabulário](../CONTEXT.md), ADRs [0004](adr/0004-isolamento-por-cnpj.md) e [0007](adr/0007-contas-por-cnpj-e-recontratacao.md) e [arquitetura do backend](../services/docs/adr/0001-arquitetura-backend.md). A [revisão de coerência](REVISAO_DOCUMENTACAO.md) mantém as lacunas ainda não resolvidas.

## Invariantes comuns

- Existem exatamente os perfis Gestor, Supervisor, Colaborador e Contabilidade/DP. Responsável legal é relação cadastral, não perfil; não existe Gestor Master.
- Conta pertence a um CNPJ. Conta ativa com e-mail confirmado e vínculo ativo são condições distintas, todas necessárias para uma ação de usuário autorizada por este módulo.
- O CNPJ do contexto, da conta, do vínculo e do recurso deve coincidir. O vínculo precisa pertencer à conta avaliada. UUID válido e empresa selecionada no cliente não comprovam autorização.
- E-mail, CPF ou igualdade de dados pessoais entre empresas não relacionam contas para fins de autorização. Bloqueio, encerramento e revogação não podem ser contornados selecionando outro perfil ou usando referência histórica.
- Ausência de concessão, ação desconhecida, recurso incompatível, referência não resolvida ou escopo insuficiente resultam em rejeição. Não interpretar ação inexistente como ação genérica permitida.
- Contabilidade/DP só executa ações explicitamente concedidas pelo gestor daquele CNPJ. Exportar PDF, criar horários e criar jornadas são concessões distintas; uma não implica a outra.
- Preservar o vínculo entre ação e recurso concreto, inclusive em arquivos, exportações, jobs e filas. Consulta/listagem precisa restringir resultados; autorizar a tela não autoriza todos os registros retornados.
- Autoaprovação de ajuste é proibida, inclusive ao acumular perfis. Nenhum perfil administrativo ou concessão permite ignorar essa proibição. Acesso ao conteúdo de atestados continua sem autorização até fechar as concessões específicas, distintas das ações genéricas de ajuste/relatório.

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

Criar funcionário sob o perfil Supervisor exige concessão específica e já o vincula à sua equipe. Essa capacidade não permite capturar um funcionário existente sem responsável ou de outra equipe. Transferência/atribuição de subordinados existentes exige ação do Gestor; o contrato final identificará a referência de supervisão e suas validações, preservando histórico e isolamento por CNPJ.

CRUD de subordinados não permite apagar registros originais ou histórico. Encerramento de cadastro deve seguir os contratos da #4. Cadastro pessoal e de funcionário não equivale a administrar conta de acesso, vínculo, perfil ou concessão; estas últimas ações permanecem do Gestor. Operações que alterem recursos compartilhados precisam verificar todo o alcance da alteração, sem atingir outra equipe por referência indireta.

### Perfil do produto e enquadramento trabalhista

O perfil Gestor é autorização administrativa, não uma classificação jurídica. O [art. 62, II e parágrafo único da CLT](https://www2.camara.leg.br/legin/fed/declei/1940-1949/decreto-lei-5452-1-maio-1943-415500-normaatualizada-pe.html) trata do exercício de cargo de gestão e da condição remuneratória ali prevista; não estabelece dispensa universal de controle de jornada para qualquer pessoa chamada gerente. Não derivar isenção de jornada do nome do perfil. A especificação do enquadramento individual e suas consequências para apuração continuam na área de jornada, fora desta issue.

## Autoridade dos dados e limites

O futuro caso de uso carrega conta/vínculo atuais e o recurso no CNPJ autorizado, resolve suas relações e solicita a decisão ao domínio. Dados vindos do cliente ou do cabeçalho de uma mensagem não são prova de vínculo, perfil ou escopo. A política pura não abre conexão e não consulta HTTP, ORM ou broker.

O consumidor verifica as relações que conhece: ajuste e seu funcionário, arquivo e seu proprietário, funcionário e setor. Se um arquivo corresponde a vários funcionários, autorização sobre um deles não autoriza o arquivo inteiro. A política nunca infere propriedade apenas do ID do solicitante informado no transporte.

Permissão para uma ação não elimina suas precondições de negócio: autorização de exportação não valida A1, leiaute, signatário ou fechamento; autorização de ajuste não resolve concorrência e sincronização tardia. Também não permite editar uma marcação original.

Esta política cobre ações de contas de usuário. Funcionário sem conta no coletor coletivo não recebe um usuário/perfil fictício: autenticação do coletor, identificação facial, prova de vida e autorização do dispositivo exigem o fluxo próprio. Ter matrícula ou saber o CPF não autoriza uma ação de usuário sobre aquele funcionário.

## Verificação da rodada

Entrega documental parcial. Permanecem por fechar a matriz do catálogo inicial, os valores de configuração, a interface comum, os adaptadores e os cenários específicos dos perfis. Não há módulo ou código Go alterado nesta rodada.
