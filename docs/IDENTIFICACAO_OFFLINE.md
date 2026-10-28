# Identificação e coleta offline

> Atualização de 27/10/2026: o [ADR 0006](adr/0006-prova-de-vida-passiva.md) define prova de vida passiva. Opções ativas citadas nesta pesquisa não representam o fluxo escolhido. Modelo/SDK e desempenho seguem pendentes de validação.

Status: pesquisa de viabilidade local em andamento documental. Captura offline para validação posterior foi descartada pelo usuário; coletor somente online com contingência é a última alternativa. Nenhum SDK foi selecionado nem executado neste projeto.

## Decisão já tomada

A empresa atribui modalidade por funcionário: aparelho próprio com autenticação do dispositivo ou aparelho coletivo com identificador e comparação facial. O identificador e os fallbacks ainda serão definidos; a prova de vida será passiva (ADR 0006). Ver [ADR 0005](adr/0005-modalidade-coleta-por-funcionario.md).

## Alternativas para o aparelho coletivo

| Alternativa | Benefício | Custo ou limitação |
|---|---|---|
| Comparação facial e prova de vida locais | Permite verificar no momento da coleta sem conexão. | Exige motor compatível, avaliação de fraude e desempenho e proteção dos cadastros biométricos no aparelho. Modelos e cadastros devem estar provisionados antes da queda de rede. |
| Captura local e validação posterior | Alternativa descartada pelo usuário. | Não implementar como fallback. |
| Coletor online com contingência externa | Menor complexidade inicial de reconhecimento offline. | Não oferece a mesma continuidade dentro do app e exige procedimento de contingência definido com a empresa. |

Recomendação técnica: se o coletor coletivo deve continuar com o mesmo nível de verificação offline, avaliar um motor local de comparação e prova de vida antes de assumir esse compromisso. Não construir prova de vida apenas com instrução de piscar ou sorrir e tratá-la como equivalente a uma solução avaliada. Se o motor não passar nos testes, considerar coletor somente online com contingência externa como última alternativa; não reduzir o método para PIN em silêncio. Captura para validação posterior foi rejeitada.

## Aparelho próprio

Proposta: após cadastro e vínculo online do aparelho, permitir autenticação local e guardar eventos pendentes em armazenamento protegido. Enviar quando houver conexão, priorizando background e mantendo fallback ao retomar. Proteger a autoria do evento com chave do aparelho vinculada à conta; um hash isolado não autentica o remetente. Definir duração da autorização offline e comportamento após alteração/revogação das credenciais. Uma revogação no servidor não chega instantaneamente a um aparelho desconectado.

## Regras comuns propostas

- Separar envio pendente ou confirmado de tratamento de jornada resolvido ou pendente. No coletivo, pendência de envio não autoriza adiar a validação de identidade, opção rejeitada pelo usuário.
- Registrar horário capturado, horário de recebimento no servidor e evidências de sincronização do relógio; assinatura não comprova que o relógio local estava correto.
- Preservar eventos e reenviar sem duplicação; apresentar a quantidade de pendências. Limpar cache, sair da conta ou reinstalar não deve ser recomendado sem verificar pendências.
- No coletivo, armazenar apenas os cadastros autorizados para aquele coletor, criptografados e com política de atualização e expiração a definir. Considerar bloqueio físico, modo quiosque e proteção contra retirada de dados.
- Mudança de modalidade ou desligamento durante uma queda de rede deve gerar conferência na sincronização. Não apagar silenciosamente eventos capturados antes de conhecer a mudança.

## Fontes e limites

- [MTE: perguntas e respostas da Portaria 671](https://www.gov.br/trabalho-e-emprego/pt-br/assuntos/inspecao-do-trabalho/fiscalizacao-do-trabalho/Perguntas%20e%20Respostas%20REP): offline não é obrigatório em toda arquitetura de REP-P.
- [Portaria 671 compilada, Anexo IX](https://www.gov.br/trabalho-e-emprego/pt-br/assuntos/legislacao/portarias-1/portarias-vigentes-3/PDFPortarian671de8denovembrode2021compilada13.05.2025.pdf): admite offline excepcional e exige envio no primeiro momento de conexão, além de requisitos de segurança e relógio.
- [Google ML Kit Face Detection](https://developers.google.com/ml-kit/vision/face-detection): detecta rostos no dispositivo, mas não reconhece indivíduos. Não resolve sozinho a comparação facial.
- [AWS Face Liveness: requisitos](https://docs.aws.amazon.com/rekognition/latest/dg/face-liveness-requirements.html) e [funcionamento](https://docs.aws.amazon.com/rekognition/latest/dg/face-liveness.html): exemplo de serviço com sessão e transmissão de vídeo, exigindo conectividade. Não representa avaliação ou escolha desse fornecedor para o projeto.

Nenhum fornecedor de comparação facial offline foi selecionado ou validado nesta pesquisa. A recomendação define capacidades a avaliar, não uma promessa de disponibilidade em Flutter ou PWA.


## Viabilidade de execução no celular — pesquisa de 26/10/2026

### Python não é obrigatório

Python pode ser usado para treinar ou converter modelos no ambiente de desenvolvimento, mas não precisa fazer parte do aplicativo. O [LiteRT para Android](https://developers.google.com/edge/litert/android) oferece execução nativa com CPU e aceleradores. Isso é infraestrutura de inferência; não fornece sozinho uma solução de reconhecimento e prova de vida validada.

Flutter pode chamar um SDK nativo por plugin, canal de plataforma ou FFI, conforme a API fornecida. O processamento pesado deve ocorrer fora da execução responsável por manter a interface responsiva.

### Candidatos encontrados, sem seleção

| Candidato | Evidência publicada pelo fornecedor | O que ainda verificar |
|---|---|---|
| [3DiVi Face SDK](https://docs.3divi.ai/face_sdk/overview/) | Documenta verificação 1:1, liveness 2D ativo/passivo e interfaces móveis/Flutter. A [página do produto](https://3divi.ai/products/software/face-sdk) declara operação offline e licença comercial com bibliotecas fechadas. | Combinação exata de modelo, módulo de liveness e plugin; desempenho, licença por coletor, ativação/renovação sem rede e evidência independente de resistência a fraude. |
| [FacePing Android](https://docs.faceping.ai/start/android/) | Tutorial oficial descreve verificação em modo avião e liveness com movimento de cabeça em Android 64-bit a partir de Android 7. | Eficácia de liveness, maturidade, suporte Flutter, condições comerciais, limites do cadastro local e comportamento da licença sem rede. Demonstração não equivale a validação de segurança. |

Como indício histórico de viabilidade, a [documentação 3DiVi 3.12](https://docs.3divi.ai/face_sdk/3.12/system_req/) recomenda Snapdragon 821 e 4 GB de RAM para um fluxo VGA com um rosto e modelos de reconhecimento específicos. Esse dado antigo não mede a combinação atual de reconhecimento e prova de vida, não fornece latência para nosso fluxo e não deve virar requisito mínimo do Teu Ponto.

Conclusão: há evidência suficiente para considerar viável uma avaliação em celular intermediário; não há benchmark executado nem garantia para qualquer aparelho. O identificador permite comparar com um cadastro (1:1), evitando busca entre todos os funcionários. A prova de vida precisa de avaliação própria e pode exigir vários frames, interação e iluminação adequada.

### Plano de validação antes de comprometer o requisito

1. Selecionar SDK com licença de avaliação e confirmar comparação e liveness realmente locais, incluindo inicialização após reiniciar em modo avião.
2. Escolher um Android intermediário real como candidato a coletor suportado. Registrar modelo, sistema, memória e câmera; não usar somente emulador como evidência.
3. Medir tempo total da sessão, consumo de memória, aquecimento e falhas em sequência de marcações. Proposta de meta inicial: 95% das sessões válidas concluídas em até 5 segundos; meta ainda não aprovada e não apresentada como desempenho observado.
4. Testar posição fixa, pessoas de alturas diferentes, contraluz, óculos e condições reais de fila. Variar usuários; medir rejeições legítimas separadamente de ataques.
5. Avaliar ataques com foto, tela de outro aparelho e vídeo, com plano de testes e métricas do SDK. Testes caseiros isolados não certificam resistência a ataques.
6. Conferir proteção dos templates, provisionamento, alterações de cadastro, revogação, expiração de licença e sincronização sem duplicação.

O coletor é um equipamento empresarial: é possível definir hardware e condições de instalação suportados. Isso reduz variabilidade e pode ser mais realista que garantir o recurso em todo celular pessoal. Fabricante e aparelho ainda não foram escolhidos; não instalar SDK ou contratar licença como consequência desta pesquisa documental.

## Pesquisa detalhada do pipeline

Ver [execução local, compatibilidade dos vetores e fronteira de confiança](PESQUISA_PIPELINE_FACIAL_MOBILE.md). O repositório privado foi inspecionado pela autenticação local: o pipeline usa InsightFace 0.7.3 e ONNX Runtime, padrão buffalo_l e comparação por cosseno em Go. A pesquisa detalhada registra arquivos e referência inspecionada.

Parecer atualizado: [viabilidade da validação facial local](VIABILIDADE_VALIDACAO_FACIAL_LOCAL.md). Execução nativa é tecnicamente viável, com desempenho do modelo atual e resistência a ataques ainda condicionados à prova de conceito.
