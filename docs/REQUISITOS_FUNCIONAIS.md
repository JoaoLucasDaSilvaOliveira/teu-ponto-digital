# Requisitos Funcionais - Projeto Teu Ponto (Core)

## Resumo
Este documento descreve as funcionalidades diretas do sistema de ponto digital "Teu Ponto". Os requisitos estão agrupados por domínios de negócio (Épicos) e possuem identificadores únicos (RF-XX) para rastreabilidade no desenvolvimento.

## Épico 1: Gestão Corporativa, Autenticação e Multi-Tenancy
- **RF-01:** O sistema deve suportar uma arquitetura multilocatário (Multi-Tenant), com cada CNPJ constituindo um tenant com dados, permissões e configurações próprios. Um único usuário gestor pode possuir vínculos e gerenciar múltiplos CNPJs (empresas/filiais) com a mesma credencial, respeitando as permissões de cada contexto.
- **RF-02:** O sistema deve permitir o cadastro e gerenciamento de perfis de Gestores.
- **RF-03:** O sistema deve permitir o cadastro de empresas (Empregador), vinculando-as a um ou mais perfis de Gestores, coletando obrigatoriamente: CNPJ, Razão Social, Apelido, Responsável Legal (nome e CPF) e Certificado Digital A1. O A1 deve ser enviado à aplicação com sua senha quando necessária à abertura; validar o arquivo e proteger o armazenamento do certificado e da senha. A3 fica para uma etapa futura. A representação do certificado no domínio não define seus signatários nem resolve as pendências regulatórias de RF-23 a RF-25.
- **RF-04:** O sistema deve exigir, durante o fluxo de login de usuários com múltiplos vínculos, a seleção do contexto da empresa que desejam operar na sessão ativa, garantindo o isolamento dos dados.
- **RF-05:** O sistema deve permitir o cadastro de funcionários (dentro do contexto da empresa selecionada), coletando: CPF, Nome Completo e Data de Nascimento.
  - **RF-05.1:** A aplicação deve atribuir uma matrícula interna aleatória, não sequencial e estável ao cadastro do funcionário, gerada por biblioteca própria e independente de seus dados pessoais. A matrícula poderá ser usada como identificador apresentado no coletor coletivo. Formato, tamanho e escopo de unicidade serão especificados antes da implementação; a persistência deverá garantir a unicidade nesse escopo e permitir nova tentativa de geração em caso de colisão.
- **RF-06:** O sistema deve suportar os perfis de acesso Gestor (administração da empresa autorizada), Supervisor (acesso ao seu setor/time), Colaborador (acesso aos próprios dados) e Contabilidade/DP (acesso conforme permissões concedidas pelo gestor da empresa). Um gestor pode administrar múltiplos CNPJs mediante vínculos distintos, conforme RF-01; não há um perfil separado de Gestor Master.
  - **RF-06.1:** O gestor de cada empresa deve definir as permissões concedidas à Contabilidade/DP naquele CNPJ, por ação, incluindo, por exemplo, exportar PDF, criar horários e criar jornadas. O perfil de Contabilidade/DP não concede essas capacidades automaticamente. As permissões concedidas em uma empresa não se estendem às demais; o servidor deve verificar o vínculo e a permissão da ação no CNPJ correspondente, conforme o ADR 0004.
- **RF-07:** O sistema deve exigir uma validação de identidade via confirmação de e-mail (Magic Link ou OTP) no primeiro acesso do usuário.
- **RF-08:** O sistema deve implementar gerenciamento de sessão de longa duração (ex: via Refresh Tokens) para manter o aparelho logado após a primeira autenticação bem-sucedida, minimizando atritos diários de login.

## Épico 2: Jornadas e Registro de Ponto
- **RF-09:** O sistema deve permitir o cadastro de horários e jornadas de trabalho, validando os limites básicos estabelecidos pela CLT (ex: interjornada de 11h).
- **RF-10:** O sistema deve permitir que o funcionário registre suas marcações de ponto (entrada, saída para intervalo, retorno do intervalo e saída).
  - **RF-10.1:** A interface deve distinguir marcação salva somente no aparelho, aguardando envio, de marcação registrada no servidor. Confirmar o registro no servidor somente após armazenamento durável; apuração e demais processamentos posteriores podem ocorrer em background. Reenvios da mesma marcação não devem gerar registros duplicados.
- **RF-11:** O sistema deve realizar a captura e validação da identidade do funcionário no momento da marcação.
  - **RF-11.1:** A empresa deve definir a modalidade de coleta por funcionário no contexto do CNPJ. Aparelho próprio utiliza autenticação do dispositivo. Aparelho coletivo utiliza identificador do funcionário e comparação facial com seu cadastro; a matrícula interna pode ser usada como identificador, conforme RF-05.1. A forma de apresentação na interface do coletor será detalhada antes de sua implementação. Ver [ADR 0005](adr/0005-modalidade-coleta-por-funcionario.md). WebAuthn é uma opção de autenticação para a plataforma web, não uma API de comparação facial do funcionário.
  - **RF-11.2:** As batidas offline devem preservar evidências de horário e origem e permitir reenvio sem duplicação. O mecanismo de proteção de autoria e detecção de adulteração deve ser especificado antes da implementação; hash isolado não equivale a assinatura e uptime não comprova o horário real da marcação.
  - **RF-11.3:** A comparação facial no aparelho coletivo deve incluir prova de vida passiva na mesma sessão, conforme [ADR 0006](adr/0006-prova-de-vida-passiva.md). Fallbacks de identificação, cadastro inicial e implementação da validação facial offline permanecem pendentes de especificação. A ausência de conectividade não autoriza substituir silenciosamente a comparação facial por identificador ou PIN. Captura offline para validação facial posterior foi descartada. Avaliar reconhecimento e prova de vida locais; coletor somente online com contingência é a última alternativa, conforme [pesquisa de viabilidade](IDENTIFICACAO_OFFLINE.md).

## Épico 3: Tratamento e Gestão de Ponto
- **RF-12:** O sistema deve permitir que o colaborador solicite ajustes manuais para batidas incorretas, faltantes ou esquecidas, exigindo uma justificativa.
  - **RF-12.1:** Diante de possível marcação esquecida, permitir registrar a marcação atual com intenção explícita (entrada, intervalo, retorno ou saída), independentemente de solicitar ou aprovar um ajuste histórico. Oferecer “Registrar entrada agora e informar saída esquecida”, permitindo concluir a solicitação depois. Fluxo aceito provisoriamente; detalhes em [Tratamento de marcação esquecida](TRATAMENTO_MARCACAO_ESQUECIDA.md).
  - **RF-12.2:** A informação de horário histórico deve gerar solicitação de ajuste com data, horário declarado e justificativa, preservando os registros originais. Não criar uma batida original retroativa nem preencher automaticamente o horário contratual.
- **RF-13:** O sistema deve permitir que o superior hierárquico (Supervisor/Gestor) avalie (aprove ou reprove) as solicitações de ajuste de ponto feitas pelos colaboradores.
  - **RF-13.1:** A aprovação do ajuste histórico não deve condicionar novas marcações. A rejeição deve conter motivo, preservar a divergência e não reclassificar automaticamente as marcações posteriores. Antes da aprovação, verificar marcações recebidas por sincronização tardia e ajustes concorrentes, evitando dupla contagem. Trechos ambíguos permanecem pendentes de tratamento, sem liquidação automática como falta, desconto ou horas extras.
- **RF-14:** O sistema deve permitir que o superior hierárquico lance abonos de faltas ou atrasos, anexando justificativas (ex: atestados médicos).
- **RF-15:** O sistema deve gerenciar a regra de Aviso Prévio, aplicando a redução de jornada escolhida: redução de 7 dias no período final ou redução de 2 horas na jornada diária, ajustando os cálculos de horas esperadas para não gerar faltas ou horas extras indevidas.

## Épico 4: Lógica de Apuração e Banco de Horas
- **RF-16:** O sistema deve permitir configurar, por funcionário ou por grupo, a modalidade de apuração de horas excedentes: Apuração Direta (pagamento em folha) ou Compensação (Banco de Horas).
- **RF-17:** O sistema deve possuir uma configuração de Banco de Horas por CNPJ (ex: validade do banco, limite de horas positivas/negativas) que ditará a regra de negócio do cômputo.
  - **RF-17.1:** Créditos e débitos devem transitar entre competências dentro do ciclo válido. Não converter saldo negativo mensal automaticamente em desconto. Ao encerrar o ciclo ou na rescisão, saldo negativo sem regra válida de liquidação deve ficar pendente de tratamento, preservando sua origem e histórico.
- **RF-18:** O sistema deve calcular automaticamente o saldo diário do colaborador, confrontando as marcações realizadas com a jornada cadastrada.

## Épico 5: Integração Contábil e Relatórios
- **RF-19:** O sistema deve possuir uma interface dedicada à Contabilidade/DP, incluindo Contador/Escritório de Contabilidade, para o mapeamento (De-Para) de dados. O acesso ao mapeamento e às demais funcionalidades deve respeitar as permissões concedidas pelo gestor de cada empresa, conforme RF-06.1.
- **RF-20:** O sistema deve permitir a vinculação da matrícula interna do funcionário no Teu Ponto (RF-05.1) com o código de matrícula do sistema de folha de pagamento. A matrícula da folha é uma referência externa separada e não substitui a matrícula interna.
- **RF-21:** O sistema deve permitir o cadastro e vinculação de Rubricas (ex: "Hora Extra 50%", "Atraso", "Adicional Noturno") para compatibilidade com o sistema contábil.
- **RF-22:** O sistema deve exportar os dados apurados em um formato de integração compatível com hardwares/softwares de folha de pagamento de mercado (foco inicial: Domínio Sistemas, utilizando importação .TXT).
- **RF-23:** O sistema deve gerar relatórios de espelho de ponto sob demanda, assinados digitalmente com o certificado da empresa para garantir validade jurídica.
- **RF-24:** O sistema deve ser um REP-P próprio, gerando obrigatoriamente o Arquivo Fonte de Dados (AFD) no layout exigido pelo MTP.
- **RF-25:** O sistema deve gerar um Comprovante de Registro de Ponto do Trabalhador para cada batida, disponibilizando-o para download (PDF) ou envio automático por e-mail.

## Épico 6: Geolocalização e Cercas Virtuais (Geofencing)
- **RF-26:** O sistema deve permitir o cadastro das coordenadas geográficas (Latitude e Longitude) da base da empresa ou de suas filiais, juntamente com a definição de um raio de abrangência aceitável (em metros) para a marcação de ponto.
- **RF-27:** O sistema deve permitir configurar a política de restrição de localidade com quatro opções:
  - **1. Não exigido (`NAO_EXIGIDO`)**: a batida é registrada sem exigência de geofence.
  - **2. Restrita (`DENTRO_RAIO`)**: a batida é bloqueada se estiver fora do raio.
  - **3. Flexível com Tag (`QUALQUER_COM_TAG`)**: a batida é permitida e marca "Batida Externa" quando fora do raio.
  - **4. Livre (`QUALQUER`)**: a batida é permitida sem validação de raio.

A política de geofencing é configurada pela empresa e permanece independente da apuração de limites de jornada. Quando a localização estiver indisponível:

| Política | Comportamento aprovado |
|---|---|
| `NAO_EXIGIDO` | Permitir a marcação sem localização. |
| `DENTRO_RAIO` | Exigir localização válida para verificar o raio. |
| `QUALQUER_COM_TAG` | Permitir e sinalizar “Localização indisponível”. |
| `QUALQUER` | Permitir a marcação sem localização. |

Localização indisponível não significa estar fora do raio. O limiar de precisão aceitável ainda precisa ser especificado. A restrição da política `DENTRO_RAIO` mantém pendente a validação regulatória indicada na revisão documental.

## Restrições de Negócio (Base Legal)
- **RN-01:** O sistema baseia-se na Consolidação das Leis do Trabalho (CLT) vigente.
- **RN-02:** A base inicial de especificação será a CLT, sem CCT/ACT de referência nesta etapa. Extensões para instrumentos coletivos exigem regras e vigência explícitas; não presumir percentuais ou regimes específicos.
