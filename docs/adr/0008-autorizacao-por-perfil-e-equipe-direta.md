---
status: accepted
---
# Autorização por perfil, concessões e equipe direta

Decidido no detalhamento da issue #5, consolidado em 02/11/2026: a configuração de vínculo possui um perfil. Supervisor compõe as capacidades próprias do Colaborador e recebe do Gestor concessões por ação sobre subordinados diretamente atribuídos; setor é agrupamento organizacional, não critério de autorização da equipe. Gestor é exclusivo, sem ações próprias de ponto, e Contabilidade/DP depende de concessões por ação no CNPJ. Não modelar Supervisor como subtipo de entidade Funcionário nem duplicar Pessoa por perfil: políticas específicas implementam o mesmo contrato Go, com validações comuns de conta, vínculo, tenant e recurso.

Escolhemos atribuição direta em vez de supervisão inferida por setor, permitindo equipes que atravessam setores sem ampliar acesso a colegas do mesmo setor. Funcionário referencia o vínculo específico de supervisão; substituir vínculo não transporta equipe automaticamente. A referência cadastral não comprova autoridade atual: a coordenação futura verifica estado, perfil, funcionário associado e autorização de atribuição.

Autoaprovação e lançamentos administrativos em benefício próprio são proibidos. Pedidos de Supervisor exigem Gestor; concessões não ultrapassam essa hierarquia. Tratar ocorrência com atestado reúne acesso ao documento e decisão/lançamento do abono correspondente: poder abonar com base no documento exige poder visualizá-lo. Relatório/dashboard não concede acesso ao documento. Perfil Gestor é escolha de acesso do produto, não enquadramento jurídico de dispensa de jornada.

Os [contratos e cenários](../AUTORIZACAO_IDENTITY_SERVICE.md) definem o catálogo inicial com rejeição de ações indefinidas. A [issue #17](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/17) definirá a matriz completa dos acessos/serviços da aplicação e suas concessões. Essas decisões substituem as pendências de composição, escopo por setor e autoaprovação mencionadas nos documentos anteriores; login, persistência, propagação de revogações e fluxos de documentos continuam sujeitos às respectivas implementações.
