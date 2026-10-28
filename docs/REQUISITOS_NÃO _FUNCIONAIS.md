# Requisitos Não Funcionais - Projeto Teu Ponto (Core)

## Resumo
Este documento descreve as restrições arquiteturais, de segurança, performance, infraestrutura e usabilidade do sistema "Teu Ponto". Os requisitos possuem identificadores únicos (RNF-XX) para guiar as decisões técnicas do desenvolvimento.

## Épico 1: Arquitetura e Stack Tecnológica
- **RNF-01:** O backend deve ser arquitetado em Micro-serviços, dividindo os domínios da aplicação (ex: Gestão de Usuários, Processamento de Ponto, Relatórios) para garantir alta disponibilidade, escalabilidade isolada e recuperabilidade.
- **RNF-02:** O sistema deve utilizar a linguagem Go nos micro-serviços, priorizando concorrência (Goroutines) para baixa latência.
- **RNF-03:** O frontend deve utilizar Flutter, priorizando a distribuição nativa. A alternativa web/PWA depende da validação das capacidades obrigatórias e pode ser descartada; ver [ADR 0003](adr/0003-flutter-nativo-e-pwa-condicional.md).
- **RNF-04:** O sistema deve utilizar um API Gateway como ponto único de entrada (Single Point of Entry) para receber as requisições externas e roteá-las para os micro-serviços internos.
- **RNF-05:** O API Gateway escolhido será Nginx self-hosted.
- **RNF-06:** O sistema deve utilizar banco de dados relacional PostgreSQL, com isolamento lógico dos dados por CNPJ e autorização no servidor em operações síncronas e assíncronas; ver [ADR 0004](adr/0004-isolamento-por-cnpj.md).

## Épico 2: Infraestrutura e Deploy
- **RNF-06:** Toda a aplicação (API Gateway, Micro-serviços, Frontend, Banco de Dados e Broker) deve ser conteinerizada utilizando Docker.
- **RNF-07:** Os containers devem ser orquestrados dentro de uma rede interna (Docker Network) privada, isolando a comunicação entre os serviços e expondo para o ambiente externo apenas o API Gateway e o Frontend.
- **RNF-08:** O projeto será provisionado e publicado em uma VPS (Virtual Private Server).

## Épico 3: Mensageria e Comunicação Interna
- **RNF-09:** A comunicação assíncrona entre os micro-serviços deve ser realizada através do message broker RabbitMQ ou síncrona através do protocolo HTTP.
- **RNF-09.01:** Chamadas que precisem de resposta imediata devem usar chamadas diretas via HTTP, como, por exemplo: `Consulta de Funcionários existentes` e outros serviços que exigem feedback imediato ao usuário.
- **RNF-09.02:** Chamadas de processamento tardio que não exijam resposta imediata ao usuário devem usar filas assíncronas para a requisição. Exemplo: `batida de ponto`.
- **RNF-10:** O payload das mensagens publicadas no RabbitMQ deve trafegar em formato JSON nativo, com criptografia síncrona simples.
- **RNF-11:** Para garantir a autenticidade e o contexto das ações assíncronas, o API Gateway (ou serviço produtor) deve injetar um token JWT assinado nos *Headers* (propriedades) da mensagem do RabbitMQ, contendo informações essenciais (ex: `tenant_id` e `user_id`).
- **RNF-12:** O micro-serviço consumidor deve validar a assinatura do JWT lido do cabeçalho utilizando uma chave pública ou *secret* previamente carregada em memória na inicialização do serviço (via variável de ambiente), rejeitando sumariamente mensagens não assinadas sem realizar chamadas HTTP síncronas.

## Épico 4: Segurança e Autenticação
- **RNF-13:** O sistema deve suportar Login Social (Google/Microsoft) implementado via protocolo OAuth2/OpenID Connect.
- **RNF-14:** Senhas de usuários (fallback) devem ser hasheadas utilizando algoritmos fortes (ex: bcrypt ou Argon2) antes da persistência no PostgreSQL.
- **RNF-15:** A comunicação externa entre o client e o API Gateway deve ocorrer obrigatoriamente sob protocolo HTTPS.
- **RNF-16:** Todos os dados sensíveis dos funcionários (CPF, biometria) devem seguir as diretrizes da LGPD, garantindo criptografia em repouso no banco de dados.

## Épico 5: Disponibilidade e Resiliência
- **RNF-17:** A interface deve utilizar Flutter. Priorizar a distribuição nativa e avaliar web/PWA somente como alternativa que atenda aos requisitos obrigatórios, conforme o ADR 0003. Na web, implementar explicitamente cache offline e fila persistente em IndexedDB; Flutter Web não implica suporte offline pronto. Ver [decisão de sincronização](adr/0002-sincronizacao-offline.md).
- **RNF-18:** A sincronização deve priorizar execução em background quando houver API compatível e permissão da plataforma. Na web, utilizar service workers e Background Sync onde suportados; um worker comum não garante execução com a aplicação fechada. Nas versões nativas, utilizar os mecanismos de background do sistema operacional, respeitando suas restrições. Abrir ou retomar a aplicação deve disparar nova tentativa como fallback, com opção manual e indicação de pendências. Não prometer envio imediato com o aplicativo fechado. A sincronização deve tolerar interrupções e reenvios sem duplicar marcações. A fronteira de aceitação durável no servidor e a publicação no RabbitMQ ainda precisam ser detalhadas; este requisito não autoriza conexão direta do cliente ao broker.

## Épico 6: Auditoria e Performance
- **RNF-19:** Operações críticas (ex: edição de ponto manual) devem gerar logs inalteráveis de auditoria contendo: ID do usuário, data/hora, IP e payload da alteração.
- **RNF-20:** O cálculo do raio de geolocalização (Fórmula de Haversine) deve ser estritamente otimizado para não impactar a latência da API.
- **RNF-21:** A geração do Arquivo Fonte de Dados (AFD) deve ser processada em *streams* na memória do micro-serviço responsável, prevenindo *Out of Memory* (OOM) ao lidar com grandes volumes de batidas.