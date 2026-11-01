# Contratos de domínio — identity-service

Status: em elaboração na [issue #4](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/4), parte da issue pai #2. Este documento ainda não constitui a especificação completa para implementação.

Fontes: [plano](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md), [limites](LIMITES_IDENTITY_SERVICE.md), [vocabulário](../CONTEXT.md), [requisitos funcionais](REQUISITOS_FUNCIONAIS.md), [arquitetura do backend](../services/docs/adr/0001-arquitetura-backend.md) e [isolamento por CNPJ](adr/0004-isolamento-por-cnpj.md).

## Decisões esclarecidas em 01/11/2026

As decisões abaixo foram informadas pelo usuário durante o detalhamento. Os pontos em discussão na seção seguinte impedem fechar entidades, cardinalidades e assinaturas sem pressupor respostas.

### Matrícula interna — RF-05.1

- Exatamente seis caracteres no conjunto `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`.
- Entrada em letras minúsculas é normalizada para maiúsculas; símbolos fora do conjunto e comprimentos diferentes são rejeitados, sem completar posições.
- Geração aleatória por biblioteca própria com fonte criptograficamente segura, independente de CPF, nome e sequências. O gerador propaga falhas da fonte.
- A validação e representação pertencem ao objeto de valor; a geração permanece separada. O construtor do funcionário recebe o valor gerado.
- Aleatoriedade não garante unicidade. Escopo, reserva histórica e coordenação de colisões precisam ser fechados antes de implementar.

### Correção de documentos

- CNPJ imutável depois do cadastro, pois identifica o tenant.
- CPF do funcionário corrigível por operação autorizada e auditável, preservando ID e matrícula e verificando unicidade no contexto aplicável.
- A correção do CPF não altera automaticamente a associação ao usuário.
- Autorização, auditoria persistente e garantia de unicidade exigem coordenação futura; um método local não comprova essas propriedades.

### Empresa incompleta e certificado — RF-03

- Permitir iniciar o cadastro da empresa antes de receber o A1 e vincular o primeiro gestor.
- O usuário solicitou responsável legal e certificado válido como condições para registrar ponto e exportar relatórios. A consequência de perda de validade durante a operação, inclusive offline, ainda precisa ser conciliada com a preservação das marcações e a validação regulatória antes de implementar qualquer bloqueio.
- Uma referência cadastrada não comprova validade atual do certificado. Upload, validação e custódia continuam fora desta issue.
- O usuário solicitou estudar futuramente evidências assinadas nas ações auditáveis. O mecanismo, o conteúdo assinado, a atribuição de autoria e o momento da assinatura ainda não estão especificados; hash isolado não é assinatura digital.

## Mudanças de modelo em discussão

### Pessoa, usuário, funcionário e perfil

O usuário propôs `Pessoa` com CPF, nome e nascimento, relacionada a `Usuário` com e-mail e estado; `Funcionário` como pessoa com tipo/perfil, incluindo Responsável legal, Gestor, Contabilidade/DP, Supervisor e Colaborador.

É necessário esclarecer se “é um” representa associação/composição de dados ou obrigatoriedade de uma conta, e se “Funcionário” mantém o significado de trabalhador ou passa a incluir participantes externos. O vocabulário atual distingue cadastro do trabalhador, identidade de acesso, vínculo de autorização e responsável legal. Responsável legal não consta como quinto perfil de acesso nas decisões atuais.

Não converter essa proposta em herança de classes: a arquitetura Go determina composição e interfaces pequenas quando há comportamento variável. Entidades e contratos finais dependem de resolver o significado dos conceitos.

### Contas isoladas por CNPJ e recontratação

O usuário solicitou contas por empresa, sem interferência dos cadastros de um CNPJ em outro, inclusive no e-mail e no estado de acesso. Isso revisa a identidade compartilhada entre empresas descrita em RF-01 e no ADR 0004; as fontes deverão ser conciliadas ao fechar o modelo.

Também solicitou que uma recontratação possa criar novo usuário no mesmo CNPJ com o mesmo CPF e e-mail, mantendo unicidade apenas entre cadastros ativos. O tratamento das referências e contas históricas, do cadastro de pessoa e da matrícula entre contratações precisa ser fechado. Não sobrescrever nem reassociar automaticamente registros históricos ao novo cadastro.

## Verificação da rodada

Entrega parcial documental, sem entidades Go, módulo ou testes executáveis novos. Critérios da issue #4 ainda pendentes: entidades versus objetos de valor, cardinalidades finais, campos, assinaturas, métodos, erros, estados e invariantes dependentes das decisões acima.
