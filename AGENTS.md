# Teu Ponto — instruções para agentes

## Escopo e fontes

O Teu Ponto é um REP-P próprio com tratamento de jornada. Trabalhe neste checkout; o repositório privado anterior `JoaoLucasDaSilvaOliveira/TeuPontoDigital` é uma referência histórica, não uma implementação automaticamente válida para este projeto.

Ao iniciar uma tarefa de produto, consulte [visão geral](docs/REGRA_NEGOCIO.md) e [vocabulário](CONTEXT.md). Carregue os documentos específicos conforme a tarefa:

| Tarefa | Fontes a consultar |
|---|---|
| Funcionalidades, permissões, geofencing ou ajustes | [Requisitos funcionais](docs/REQUISITOS_FUNCIONAIS.md) |
| Arquitetura, infraestrutura, segurança ou sincronização | [Requisitos não funcionais](<docs/REQUISITOS_NÃO _FUNCIONAIS.md>) |
| Cálculos, banco de horas ou enquadramento normativo | [Revisão de coerência](docs/REVISAO_DOCUMENTACAO.md), depois [lógica de cálculos](docs/LOGICA-CALCULOS-CORE.md) |
| Marcação esquecida, virada de dia ou correção concorrente | [Fluxo de tratamento](docs/TRATAMENTO_MARCACAO_ESQUECIDA.md) |
| AFD, comprovantes, espelho, AEJ ou integração Domínio | [Revisão de coerência](docs/REVISAO_DOCUMENTACAO.md), depois [esboços de leiautes](docs/LAYOUT_DOMINIO.md) |
| Identificação facial, prova de vida ou coleta offline | ADRs 0005/0006 abaixo, depois [parecer de viabilidade](docs/VIABILIDADE_VALIDACAO_FACIAL_LOCAL.md) e as pesquisas nele referenciadas |

Os [ADRs aceitos](docs/adr/) registram as decisões tomadas e prevalecem sobre propostas anteriores de arquitetura/produto. Os requisitos definem o comportamento; pesquisas, recomendações de fornecedores e metas propostas não equivalem a decisões aprovadas. A revisão registra conflitos ainda abertos: ao encontrar um deles, resolva o ponto necessário à tarefa antes de implementá-lo e atualize sua fonte. Decisões do produto não comprovam conformidade legal; use a norma oficial aplicável para validar exigências regulatórias.

## Decisões que orientam o trabalho

- [0001 — REP-P próprio](docs/adr/0001-rep-p-proprio.md): CLT como base inicial, sem CCT/ACT de referência. Regras dependentes desses instrumentos exigem especificação própria.
- [0002 — Sincronização](docs/adr/0002-sincronizacao-offline.md): priorizar background onde permitido, com reenvio ao abrir/retomar como fallback. Distinguir pendência local de registro durável no servidor e tolerar reenvios sem duplicação.
- [0003 — Flutter nativo](docs/adr/0003-flutter-nativo-e-pwa-condicional.md): nativo é o foco. PWA é opcional e depende de atender às capacidades obrigatórias; pode ser descartada.
- [0004 — Isolamento por CNPJ](docs/adr/0004-isolamento-por-cnpj.md): cada CNPJ é um tenant. Verificar vínculo e autorização no servidor, também em filas, arquivos e exportações.
- [0005 — Coleta por funcionário](docs/adr/0005-modalidade-coleta-por-funcionario.md): a empresa atribui aparelho próprio com autenticação do dispositivo ou coletivo com identificador e comparação facial.
- [0006 — Prova de vida passiva](docs/adr/0006-prova-de-vida-passiva.md): aplicar no coletivo, na mesma sessão da comparação. Validação de identidade posterior foi descartada; execução local é a opção a avaliar e online com contingência é a última alternativa.

## Cuidados de implementação

Preserve marcações originais. Ajustes históricos seguem tratamento auditável e aprovação, sem impedir novas marcações nem deslocar automaticamente sua classificação quando um pedido é rejeitado. O fluxo de marcação esquecida foi aceito provisoriamente; consulte sua especificação para conflitos e sincronização tardia.

Consulte RF-17.1 para banco de horas: fechamento mensal e liquidação do ciclo são eventos distintos. Mantenha saldo negativo sem regra válida de liquidação pendente de tratamento, sem desconto automático.

Trate as políticas de localização conforme RF-27. Localização indisponível, baixa precisão e estar fora do raio são situações distintas; preserve a pendência de validação regulatória da política restritiva.

Em biometria, mantenha separadas autenticação do aparelho, comparação facial e prova de vida. Preserve versão dos modelos e compatibilidade dos cadastros; um embedding não é intercambiável entre modelos por ter a mesma dimensão. Verifique licença dos pesos além da licença do código. Declare benchmark como resultado apenas após medição, identificando aparelho e etapas medidas.

## Pendências conhecidas

Antes de implementar cada área afetada, consulte os achados da revisão. Merecem atenção especial: leiaute REP-P/AEJ e signatários; esboço Domínio sem validação; regras de adicionais e regimes de jornada; proteção/autorização das mensagens; precisão mínima de localização; cadastro facial inicial, SDK/modelos, hardware e janela offline. A documentação de cálculos ainda contém exemplos e generalizações apontados na revisão: não os transforme diretamente em regras de produção.

## Versionamento de código

As regras desta seção se aplicam ao repositório inteiro, incluindo `services/`.

### Branches

| Branch | Finalidade |
|---|---|
| `main` | Código revisado, testado e integrado; destino da promoção para produção. |
| `development` | Integração e homologação, com ambiente e banco de dados próprios para testes. |
| `feat/<nome-do-serviço>/<funcionalidade>` | Desenvolvimento de uma funcionalidade específica. |
| `bugfix/<nome-do-serviço>/<correcao>` | Correção de erros mais graves. |
| `hotfix/<nome-do-serviço>/<correcao>` | Correção de erros de simples resolução. |

`feat`, `bugfix` e `hotfix` são prefixos de branches de trabalho, não branches permanentes. As definições de `bugfix` e `hotfix` acima são as adotadas por este projeto; não trocar sua classificação por convenções externas.

### Fluxo

1. Parta da `main` atualizada e crie uma branch para a funcionalidade ou correção, vinculando o trabalho à issue correspondente quando existir. Exemplo: `feat/identity-service/entidades-dominio`.
2. Desenvolva e execute as verificações pertinentes na branch de trabalho. Mantenha `main` e `development` como destinos de integração, sem desenvolver diretamente nelas.
3. Integre o trabalho em `development` para homologação e testes no ambiente próprio. Homologação não equivale a liberação para produção.
4. Após homologação e revisão, abra um pull request da própria branch de trabalho para `main`. Promova somente a funcionalidade ou correção homologada, sem incluir outras mudanças presentes em `development`.
5. Faça correções de homologação na branch de trabalho e volte a integrá-las em `development` antes de promover a versão corrigida. Evite incorporar `development` inteira na branch de trabalho, pois isso pode levar funcionalidades não relacionadas para `main`.

### Commits

Crie commits ao final de cada rodada de trabalho ou ao concluir uma implementação relevante. Revise o diff e inclua somente alterações pertencentes ao trabalho; preserve mudanças alheias. Os commits devem representar unidades coerentes e informar o que foi alterado. Registre as verificações executadas e eventuais falhas na entrega, sem declarar uma implementação validada quando não foi possível verificá-la.

Siga [Conventional Commits](https://www.conventionalcommits.org/pt-br/v1.0.0/): `<tipo>[escopo opcional]: <descrição>`. Use escopo quando ajudar a identificar o serviço ou área alterada. Exemplos: `feat(identity-service): adiciona entidades de domínio`, `fix(clock-service): corrige pareamento de marcações` e `docs(projeto): define fluxo de versionamento`. Use `!` ou o rodapé `BREAKING CHANGE:` para mudanças incompatíveis. Os prefixos de branches não determinam diretamente o tipo do commit: correções em `bugfix/` ou `hotfix/` usam `fix`.

## Execução e conclusão

Inspecione os arquivos e manifestos existentes para determinar estrutura, versões e comandos. Este projeto não define ainda um roteiro verificado de build/teste; registre os comandos reais quando a implementação existir, sem copiar rotinas do repositório anterior como se já funcionassem aqui.

Ao alterar comportamento, atualize a fonte correspondente e preserve os IDs dos requisitos; acrescente ADR quando houver uma nova decisão arquitetural relevante. Mantenha pesquisas separadas de requisitos aceitos.

Para código, execute as verificações disponíveis e pertinentes. Em cálculos, autorização entre tenants, sincronização e concorrência de ajustes, cubra casos observáveis e falhas relevantes. Para alterações só documentais, confira referências, IDs e coerência com decisões aceitas. Na entrega, informe o que mudou, o que foi verificado e o que continua pendente.
