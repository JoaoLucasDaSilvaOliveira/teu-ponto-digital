# Revisão de coerência da documentação

Data: 26/10/2026. Status: diagnóstico inicial, com decisões posteriores registradas abaixo. As demais propostas seguem pendentes.

## Decisões confirmadas após a revisão

- REP-P próprio: [ADR 0001](adr/0001-rep-p-proprio.md).
- CLT como base inicial, sem CCT/ACT de referência nesta etapa.
- Priorizar sincronização em background onde possível; abrir/retomar como fallback: [ADR 0002](adr/0002-sincronizacao-offline.md).
- Diferenciar pendência local e registro durável no servidor foi confirmado; ver RF-10.1.
- Geofencing segue a política da empresa. Comportamento sem localização aprovado e registrado em RF-27; precisão aceitável e validação regulatória permanecem pendentes.
- Banco de horas carrega créditos e débitos durante o ciclo; saldo negativo sem regra válida de liquidação fica pendente, sem desconto automático (RF-17.1).
- Fluxo de marcação esquecida aceito provisoriamente: a marcação atual independe da aprovação do ajuste histórico; ver RF-12.1 a RF-13.1 e [detalhamento](TRATAMENTO_MARCACAO_ESQUECIDA.md).
- Flutter nativo é o foco; PWA é alternativa condicionada ao atendimento dos requisitos, podendo ser descartada: [ADR 0003](adr/0003-flutter-nativo-e-pwa-condicional.md).
- Cada CNPJ é um tenant com dados, permissões e configurações próprios: [ADR 0004](adr/0004-isolamento-por-cnpj.md).

A tabela a seguir descreve os achados da revisão original; os documentos foram atualizados apenas nos pontos confirmados acima.

- Modalidade de coleta definida pela empresa por funcionário: [ADR 0005](adr/0005-modalidade-coleta-por-funcionario.md). Pesquisa de offline em [IDENTIFICACAO_OFFLINE.md](IDENTIFICACAO_OFFLINE.md): avaliar execução nativa local, sem exigir Python. Captura com validação posterior foi descartada; coletor somente online é a última alternativa. SDK, hardware e metas ainda não aprovados.

- Prova de vida passiva escolhida em 27/10/2026: [ADR 0006](adr/0006-prova-de-vida-passiva.md). Fornecedores e modelos continuam sem seleção; pesquisas anteriores não alteram essa escolha.

## Resultado

Atualização em 30/10/2026: os perfis de acesso passam a ser Gestor, Supervisor, Colaborador e Contabilidade/DP, sem Gestor Master separado. O gestor define as permissões da Contabilidade/DP por ação e por CNPJ, incluindo possibilidades como exportar PDF, criar horários e criar jornadas; ver RF-06, RF-06.1 e RF-19. Esta decisão resolve a ausência desse perfil e a origem de suas permissões. A matriz completa de ações e as regras de aprovação do próprio ajuste e acesso a atestados continuam pendentes.

Na mesma data, foi aprovado o [plano inicial do identity-service](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md). RF-03 passa a exigir nome e CPF do responsável legal e delimita o certificado ao envio de A1, com A3 futuro; a custódia técnica será detalhada antes de implementar o upload. RF-05.1 define matrícula interna aleatória, reutilizável no coletor, e RF-20 preserva a matrícula da folha como referência externa. O plano separa usuário de funcionário e de vínculo com empresa. As decisões não resolvem signatários dos artefatos, todo o onboarding, a matriz completa de permissões nem os fluxos biométricos e de mudança de modalidade.

A visão do produto é consistente: controle de jornada multiempresa, registro, tratamento, apuração e integração com folha. Os documentos ainda não formam uma especificação inequívoca para implementação. Resolver primeiro enquadramento regulatório, experiência offline e políticas de apuração; depois consolidar as instruções no AGENTS.md.

O checkout contém README e documentação, sem aplicação disponível para confrontar os requisitos com comportamento implementado. Há exclusões preexistentes em identity_service; esta revisão não as altera.

## Conflitos e lacunas

| Prioridade | Local | Problema | Encaminhamento proposto |
|---|---|---|---|
| Crítica | LAYOUT_DOMINIO; RF-24 | AFD associado à Portaria 1.510/2009 é apresentado como REP-P. REP-P e REP-A aparecem como alternativas equivalentes. | Definir enquadramento. Substituir o esboço por referência e implementação do leiaute aplicável da Portaria 671; REP-A depende de instrumento coletivo autorizador. |
| Crítica | RF-23 a RF-25; REGRA_NEGOCIO | Certificado da empresa é tratado como solução geral de validade jurídica; não distingue responsável pela assinatura e tipo de artefato. AEJ/PTRP não consta no escopo. | Separar AFD, AEJ, espelho e comprovante, seus emissores, signatários e formatos. O MTE atribui a assinatura do AFD REP-P/REP-A ao fabricante/desenvolvedor. Não presumir que assinatura do empregador supre todas as obrigações. |
| Crítica | LOGICA-CALCULOS-CORE | Exceder 2h extras ou 10h trabalhadas vira automaticamente adicional de 100%, apresentado como regra legal universal. | Separar limites de jornada, alertas de irregularidade e remuneração. Essa conversão automática não decorre do art. 59 da CLT; percentuais precisam de fundamento e vigência. |
| Alta | LOGICA-CALCULOS-CORE | Exemplos aritméticos errados. | 9:50 − 7:20 = 2:30; 10:30 − 7:20 = 3:10; 9:50 − 8:48 = 1:02; 10:30 − 8:48 = 1:42. Corrigir também a distribuição das faixas após decidir a política. |
| Alta | Banco de horas | Compensação durante o ciclo contradiz desconto de todo débito no fechamento mensal. | Decidir se débito passa ao mês seguinte, limites, vencimento, rescisão e política aplicável. Distinguir fechamento mensal de liquidação do ciclo; não apresentar desconto mensal como regra universal. |
| Alta | RF-27; REGRA_NEGOCIO | Bloqueio por geofence e invalidação de marcação podem conflitar com registro fiel da jornada. | Avaliar restrição à luz do enquadramento e art. 74 da Portaria 671. Proposta: preservar ocorrência e evidências e sinalizar exceção para tratamento. Não descartar marcação silenciosamente. |
| Alta | RF-11.1 | WebAuthn é chamado de validação biométrica e confundido com reconhecimento facial do trabalhador. | Distinguir autenticação por credencial, verificação local do usuário e identificação facial. WebAuthn pode usar PIN; a aplicação não recebe a biometria. Definir ameaça que se pretende mitigar, inclusive aparelhos compartilhados. |
| Alta | RF-11.2 | Hash é chamado de assinatura; fuso e uptime são tratados como mecanismos suficientes de antifraude. | Hash sem chave não autentica origem. Definir chaves, vínculo do dispositivo, replay, relógio alterado e grau de confiança. Uptime no navegador não é prova confiável entre reinícios. Offline exige modelo próprio de evidências. |
| Alta | RNF-09.02; RNF-18; RF-25 | Batida assíncrona não define quando há aceitação durável e comprovante. | Proposta: cliente envia HTTPS, servidor confirma somente após persistência durável; eventos posteriores seguem por outbox e RabbitMQ. Distinguir pendência local, registro aceito e apuração. Cliente não acessa broker diretamente. |
| Alta | RNF-10 a RNF-12 | “Criptografia síncrona simples” não define mecanismo; JWT sem política temporal e de reentrega é incompleto. | Especificar TLS, credenciais por serviço, autorização nas filas e proteção em repouso conforme necessidade. Se JWT permanecer, definir emissor, audiência, expiração, rotação e comportamento em retries/DLQ. Assinatura não fornece confidencialidade. |
| Alta | RF-01, RF-04, RF-17 | Tenant é confundido com CNPJ/empresa/filial; configuração “global” do banco é ambígua. | Definir fronteira de isolamento e vínculo empregatício. Autorizar tenant no servidor e em jobs, exports e eventos; seleção no login não basta. Esclarecer se global significa por empresa. |
| Alta | Cálculos | Marcações são definidas como pares, sem separar evento individual de pareamento para apuração. Falta tratamento de marcação ímpar. | Preservar eventos originais; apuração sinaliza incompletude. Ajustes aprovados complementam o tratamento com autoria e justificativa. |
| Alta | Cálculos | Domingo, feriado, 12x36, “11x36” e trabalho noturno estão generalizados. | Definir regimes suportados e fundamento. Cobrir compensação de descanso, hora noturna reduzida, adicional, prorrogação e regras coletivas. 22h–05h sozinho não especifica o cálculo. |
| Alta | Cálculos | Faltam virada de dia, fuso da jornada, intervalos, tolerância com desvios opostos e vigência das regras. | Especificar casos concretos e resultados esperados. Jornada diária não deve ser cortada automaticamente à meia-noite. Preservar regra aplicada em reprocessamentos e fechamentos. |
| Média | RF-03, RF-05, RF-07 | E-mail obrigatório no primeiro acesso não aparece no cadastro; matrícula só aparece na exportação; A1/A3 não têm fluxo definido. | Documentar credenciais separadas do vínculo, onboarding e matrícula. Definir custódia/uso de certificado, especialmente A3 em tarefas de servidor. |
| Média | RF-06, RF-19 | Contador não tem papel/permissões explícitos; RH/gestor/supervisor variam entre textos. | Consolidar matriz de permissões por empresa e equipe, incluindo aprovação do próprio ajuste e acesso a atestados. |
| Média | RF-27 | NAO_EXIGIDO e QUALQUER têm comportamento aparentemente igual. | Definir diferença em captura de coordenadas, consentimento/permissão, indisponibilidade e auditoria; considerar unificação. |
| Média | RNFs | Microserviços são justificados como garantia de disponibilidade, mas há VPS única e dependências centrais. Sem metas mensuráveis. | Registrar custos operacionais e falhas únicas; definir recuperação, backup/restauração, latência, carga, retenção e observabilidade. Microserviços não garantem disponibilidade por si. |
| Média | LGPD/auditoria | Criptografia e logs inalteráveis não cobrem retenção, minimização e acesso; payload pode incluir atestados e segredos. | Definir classificação, controle de acesso, retenção, descarte, proteção de backups e redação de logs. Não copiar indiscriminadamente payloads sensíveis. |
| Média | LAYOUT_DOMINIO | MMAAAA tem 6 caracteres, mas campos têm tamanho 7; exemplos não demonstram posições e preenchimentos declarados. TIPO V sem escala monetária. | Obter manual oficial do importador Domínio e versão antes de chamar o leiaute de compatível. Definir encoding, terminador, escala, separadores e validar arquivo real de referência. |
| Baixa | RNFs; links; README | RNF-06 duplicado, links LOGIC inválidos, nome de arquivo com espaço acidental e README limitado ao backend. | Corrigir referências e IDs com rastreabilidade após consolidar requisitos. Atualizar README para o escopo real. |

## Flutter e PWA: avaliação técnica

Flutter é um framework. PWA descreve capacidades e distribuição de uma aplicação web. Flutter Web pode ser a base de uma PWA; isso não elimina restrições do navegador.

| Necessidade | Flutter nativo Android/iOS | Flutter Web/PWA |
|---|---|---|
| Distribuição | Build, assinatura e distribuição específicos por plataforma; lojas adicionam operação. | Acesso por URL e instalação conforme navegador, sem publicação obrigatória em loja. |
| Registro offline | Banco local e integrações com armazenamento seguro; requer projeto de sincronização. | IndexedDB para fila estruturada e service worker para recursos offline; requer implementação e gestão de cache. LocalStorage não deve ser a fila principal. |
| Envio com app fechado | APIs de background e restrições do SO; execução imediata não é garantida. | Background Sync não é universal e não é suportado no Safari/iOS conforme compatibilidade consultada. Garantir reenvio ao abrir/retomar, com ação manual disponível. |
| Localização na batida | APIs nativas, com permissões e limitações do aparelho. | Geolocation exige contexto seguro e permissão; precisão e disponibilidade variam. Distância medida no momento da batida é diferente de monitoramento contínuo por geofence. |
| Biometria | Autenticação local e chaves protegidas conforme plataforma; não equivale automaticamente a identificar o empregado. | WebAuthn autentica por credencial e pode verificar usuário por biometria ou PIN; não fornece imagem facial nem biometria bruta ao servidor. |
| Antifraude | Mais integrações de plataforma, sem garantia absoluta contra relógio/GPS adulterados. | Menos controle sobre ambiente, armazenamento e execução. Hash local isolado não prova autoria nem horário. |
| Reuso | Pode compartilhar Dart e interface com web, mantendo adaptações por plataforma. | Reusa Flutter, mas integrações web, cache e ciclo de atualização têm manutenção própria. |

A documentação Flutter atual informa que o framework deixou de gerar e administrar service worker de cache por padrão. Portanto “Flutter gera PWA” precisa ser desdobrado em requisitos implementáveis: manifesto/instalação, cache, persistência, migração, atualização e testes em aparelhos reais.

Recomendação condicional: se o MVP aceita sincronizar ao reabrir/retomar a aplicação, Flutter Web/PWA é uma opção a validar no Safari/iPhone. Se integrações locais e execução em background forem centrais, priorizar Flutter nativo e avaliar distribuição iOS. Nenhuma opção deve prometer sincronização imediata e infalível com o app fechado.

Antes de decidir, validar em iPhone/Safari e Android: carregamento offline após primeira visita, registro e reinício, reabertura após perda de rede, duplicação/replay, troca de usuário/empresa, falta de permissão de GPS, baixa precisão, armazenamento removido e atualização com batidas pendentes.

## Decisões iniciais para discussão

1. O produto será REP-P com tratamento de ponto próprio, ou integração com registrador externo? Evitar REP-P/REP-A indefinido.
2. No MVP, pendências offline podem ser enviadas ao reabrir a aplicação no iPhone, ou existe requisito de envio com app fechado?
3. Há CCT/ACT e categoria de referência para validar apuração, percentuais e banco de horas? Até existir, tratar regras controversas como propostas, não como legislação universal.

Depois dessas respostas: detalhar modelo de marcação/aceitação/tratamento, política do banco e permissões; atualizar os documentos de origem e registrar ADR apenas para decisões efetivamente tomadas. O AGENTS.md deve apontar para essas fontes por tipo de tarefa, sem duplicar regras nem inventar comandos de execução ausentes no repositório.

## Fontes consultadas

- [CLT, especialmente arts. 58, 59, 59-A e 73](https://www.planalto.gov.br/ccivil_03/decreto-lei/del5452.htm).
- [Portaria MTP 671 compilada disponibilizada pelo MTE](https://www.gov.br/trabalho-e-emprego/pt-br/assuntos/legislacao/portarias-1/portarias-vigentes-3/PDFPortarian671de8denovembrode2021compilada07.01.2026.pdf).
- [Perguntas e respostas do MTE: REP, assinaturas, AEJ, NSR e offline](https://www.gov.br/trabalho-e-emprego/pt-br/assuntos/inspecao-do-trabalho/fiscalizacao-do-trabalho/Perguntas%20e%20Respostas%20REP).
- [Flutter Web FAQ: service worker e limites da plataforma](https://docs.flutter.dev/platform-integration/web/faq).
- [W3C WebAuthn](https://www.w3.org/TR/webauthn-3/).
- [MDN Background Synchronization, incluindo compatibilidade](https://developer.mozilla.org/en-US/docs/Web/API/Background_Synchronization_API).
- [MDN Geolocation API](https://developer.mozilla.org/en-US/docs/Web/API/Geolocation_API).

Esta revisão identifica inconsistências e requisitos de validação; não certifica conformidade regulatória do produto.

## Contexto adicional: banco de horas

Lido o [chat compartilhado pelo usuário](https://chatgpt.com/share/6adfa930-49e0-83e9-929a-666c9262b008). A discussão distingue competência mensal, encerramento do ciclo e rescisão. Dentro do ciclo válido, o saldo negativo pode seguir para o mês seguinte. Desconto, abono e transporte para outro ciclo não são consequências automáticas; dependem de fundamento e condições aplicáveis. Também separa a origem dos débitos.

Encaminhamento para a especificação inicial baseada na CLT: retirar a conversão automática de saldo negativo mensal em horas falta, preservar saldo e origem durante o ciclo e deixar a liquidação negativa pendente quando não houver regra válida definida. Encaminhamento confirmado pelo usuário e aplicado ao texto de cálculos e ao RF-17.1. Não importar os exemplos de CCTs do chat como regras gerais do produto.
