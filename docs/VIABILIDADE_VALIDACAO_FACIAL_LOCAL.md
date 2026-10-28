# Viabilidade da validação facial no aparelho

> Atualização de 27/10/2026: o [ADR 0006](adr/0006-prova-de-vida-passiva.md) define prova de vida passiva. Opções ativas citadas nesta pesquisa não representam o fluxo escolhido. Modelo/SDK e desempenho seguem pendentes de validação.

Pesquisa em 26/10/2026. Resultado: tecnicamente viável em aplicativo nativo; aprovação para produção condicionada a modelo/licença, prova de vida, segurança e ensaio no hardware. Não há medição executada em celular neste projeto. Este relatório não altera as modalidades de coleta aprovadas.

## Parecer

| Questão | Conclusão | Evidência ou limite |
|---|---|---|
| Extrair vetor e comparar sem servidor Python | Viável. | ONNX Runtime executa inferência Android/iOS; o projeto já usa modelos ONNX. |
| Fazer comparação 1:1 offline | Viável com cadastro/modelo previamente provisionados. | A comparação independe de rede; atualização de permissões e revogações depende de sincronização. |
| Executar o buffalo_l atual com desempenho aceitável | Candidato a experimento, ainda não comprovado. | É ResNet50 com detector SCRFD-10GF; os tempos de MobileFaceNets não se aplicam a ele. |
| Usar modelo leve para melhorar desempenho | Viável, com revalidação de reconhecimento e cadastros. | Modelo menor muda o espaço vetorial; 512 dimensões por si só não garantem compatibilidade. |
| Reconhecimento e prova de vida locais | Existem componentes e SDKs móveis para isso. | Resistência a ataques e licença offline devem ser avaliadas para o módulo/binário exato. |
| Prometer operação em qualquer celular intermediário | Sem sustentação. | Homologar equipamento do coletor e condições de instalação. |
| Capturar offline e verificar identidade depois | Descartado pelo usuário. | Pendência de envio não é pendência de identificação. |

Fontes: [ONNX Runtime Mobile](https://onnxruntime.ai/docs/get-started/with-mobile.html), [model zoo InsightFace](https://github.com/deepinsight/insightface/blob/master/model_zoo/README.md), [pesquisa do pipeline e código privado](PESQUISA_PIPELINE_FACIAL_MOBILE.md), [pesquisa específica de prova de vida](PESQUISA_PROVA_VIDA_LOCAL.md).

## O projeto anterior já tem extração no navegador

Na referência privada `613bc0ddffa61a106c515bc5652b616b9ab4fe66`, [useFaceCapture.ts](https://github.com/JoaoLucasDaSilvaOliveira/TeuPontoDigital/blob/613bc0ddffa61a106c515bc5652b616b9ab4fe66/frontend/src/composables/useFaceCapture.ts) importa `@vladmandic/face-api`, carrega detector, landmarks e rede de reconhecimento e chama `withFaceDescriptors()`. O descritor é identificado como `FACE_API_128D_V1`.

No backend, [facial_usecase.go](https://github.com/JoaoLucasDaSilvaOliveira/TeuPontoDigital/blob/613bc0ddffa61a106c515bc5652b616b9ab4fe66/services/identity-service/internal/core/usecases/facial_usecase.go) extrai novamente o vetor pelo validador de imagem no servidor e compara esse vetor ao cadastro InsightFace. O descritor enviado pelo navegador não é a fonte usada na comparação final.

Portanto, já existe extração local no código, mas não uma validação offline completa. Não é possível comparar o vetor face-api 128D com o cadastro InsightFace 512D; preservar dois pipelines sem necessidade também aumenta custo e manutenção. Escolher uma implementação canônica para cadastro e verificação do coletor é recomendação de arquitetura, ainda não decisão do produto.

## Três caminhos de implementação

| Caminho | Adequação | Custos e consequência |
|---|---|---|
| Mesmos pesos de reconhecimento do buffalo_l + runtime nativo | Primeiro candidato se preservar cadastros for importante. | Portar detecção, alinhamento e normalização; medir consumo e latência; adicionar PAD local. Não basta copiar o limiar 0.84. |
| Reconhecimento mobile, como a família usada em buffalo_sc | Candidato se tamanho/latência do R50 forem inadequados. | Avaliar perda de desempenho biométrico e gerar cadastros compatíveis; verificar licença. |
| SDK comercial local com reconhecimento e PAD | Candidato se reduzir manutenção de ML e ter suporte forem prioritários. | Licença, dependência de fornecedor, formato próprio dos templates e prova de operação offline da combinação contratada. |

O catálogo oficial lista buffalo_l com ResNet50 e SCRFD-10GF, pacote de 326 MB; buffalo_sc com MBF e SCRFD-500MF, pacote de 16 MB sem módulos de atributos adicionais. Os números representam pacotes publicados, não RAM nem tamanho final do app. A tabela do fornecedor também mostra desempenho biométrico diferente entre famílias; a troca não é uma otimização sem custo. [Model zoo](https://github.com/deepinsight/insightface/blob/master/model_zoo/README.md).

O código FaceAnalysis da linha v0.7 permite selecionar `allowed_modules`; sem filtragem, prepara e executa os módulos carregados para cada rosto. O script privado não seleciona somente reconhecimento/detecção. Como proposta para o port, manter exclusivamente módulos necessários evita transportar cálculo de atributos que não contribui para identificação. As economias reais dependem de medir a implementação. [FaceAnalysis v0.7](https://github.com/deepinsight/insightface/blob/v0.7/python-package/insightface/app/face_analysis.py).

## Evidência adicional de PAD no celular

O [README primário do Silent-Face](https://github.com/minivision-ai/Silent-Face-Anti-Spoofing/blob/master/README_EN.md) publica implantação Android de PAD RGB e tempos declarados de 24 ms em Qualcomm 845 e 25 ms em Kirin 810. Esses dados mostram uma implementação móvel leve segundo o autor; não são benchmark independente, não descrevem a batida completa e não demonstram resistência a todos os ataques em nosso cenário. O próprio mantenedor registra dependência de câmera e ambiente. Essa evidência complementa, sem substituir, o ensaio do módulo escolhido.

## Integração e hardware

Arquitetura candidata: interface/câmera Flutter → ponte nativa → detecção e alinhamento → reconhecimento e PAD local → decisão de identidade → evento local protegido → sincronização. Modelos e templates precisam existir antes da falta de conexão. Um resultado PAD deve se referir à mesma sessão/captura da comparação; passar no PAD com uma pessoa e comparar outra captura não compõe uma validação confiável.

ONNX Runtime é o caminho inicial com menor mudança de formato para o modelo já usado. Iniciar avaliação em CPU/XNNPACK; só adotar aceleração adicional após medir suporte dos operadores e resultados. A documentação recomenda testar desempenho real, porque particionar modelos entre acelerador e CPU pode piorar o resultado. [Deploy móvel](https://onnxruntime.ai/docs/tutorials/mobile/), [verificador de usabilidade do modelo](https://onnxruntime.ai/docs/tutorials/mobile/helpers/model-usability-checker.html).

O usuário dispõe de Poco X7 Pro. A [especificação oficial](https://www.po.co/global/product/poco-x7-pro/specs/) informa MediaTek Dimensity 8400-Ultra e GPU Mali-G720. É um candidato para o primeiro ensaio; isso não determina a latência do modelo. Não planejar dependência do Qualcomm QNN para esse aparelho: o provider é específico do ecossistema Qualcomm. [QNN provider](https://onnxruntime.ai/docs/execution-providers/QNN-ExecutionProvider.html).

Em 26/10/2026, Flutter, Android SDK e adb estão disponíveis no ambiente de desenvolvimento, mas `adb devices -l` não apresentou aparelho conectado. Não foi possível medir desempenho em hardware real; não substituir essa medição por resultado de emulador ou CPU do computador.

## Condições de aprovação para o coletor

### Funcionalidade e compatibilidade

- Mesmas imagens processadas no servidor e no celular: comparar recortes alinhados, vetores e escores, incluindo decisões próximas do limiar. Definir tolerância numérica antes do ensaio; não exigir igualdade bit a bit de runtimes diferentes.
- Guardar hash/versão dos pesos e do pré-processamento junto da versão de cadastro. Se trocar modelo, planejar recadastro ou migração explícita.
- Permitir exatamente um rosto e validar qualidade da captura e associação da sessão com o funcionário selecionado.

### Desempenho

Medir separadamente inicialização fria, detecção, alinhamento, embedding, PAD e tempo total. Proposta inicial de experiência: p95 de até 5 segundos por sessão válida e fila de 50 sessões consecutivas sem falhas de memória ou degradação progressiva inaceitável. Esses números são metas propostas, não requisito aceito nem resultado observado. Registrar RAM, consumo, aquecimento e efeito da instalação fixa/iluminação. Repetir em aparelho de menor custo antes de declarar suporte a intermediários.

### Qualidade biométrica e ataques

Medir falso aceite e falsa rejeição com limiar explícito. Um piloto pequeno detecta problemas, mas não demonstra taxas de erro extremamente baixas: não anunciar FMR 1e-5 a partir de algumas dezenas de pessoas. Para PAD, exigir evidência referente à versão, dispositivo, tipo de câmera e ataques avaliados; testar foto impressa, tela e vídeo e considerar injeção de câmera separadamente. NIST FRTE de reconhecimento não certifica PAD. [NIST biometria](https://pages.nist.gov/800-63-4/sp800-63b.html), [guia de avaliação InsightFace](https://www.insightface.ai/guides/choose-face-recognition-model-and-evaluate).

### Offline operacional

Reiniciar em modo avião e validar: modelos, cadastro, licença, identidade, persistência e reenvio sem duplicação. Cobrir vencimento de autorização, relógio alterado, desligamento, atualização do app com pendências e revogação recebida depois. Janela máxima offline e comportamento após seu vencimento ainda precisam de decisão; a disponibilidade do reconhecimento local não resolve esses pontos automaticamente.

### Confiança e dados

Provisionar apenas templates autorizados no CNPJ/coletor e protegê-los em repouso. Assinatura do evento identifica a chave do dispositivo, mas não comprova sozinha execução íntegra de reconhecimento/PAD. Proteger aparelho empresarial, app e atualização de modelos; documentar quais ameaças são cobertas. Template facial continua sendo dado biométrico, mesmo sem guardar a foto. Ver [pesquisa de confiança local](PESQUISA_PIPELINE_FACIAL_MOBILE.md).

## Recomendação resultante

Avançar com validação nativa local como opção preferencial para o coletor, sem assumir Python no celular. Primeira trilha: mesmos pesos ONNX, pipeline mínimo e avaliação de paridade/desempenho. Segunda trilha: PAD local avaliado, com SDK comercial como candidato, não escolha implícita. Só mudar o modelo de reconhecimento se houver ganho necessário e plano para os cadastros.

A pesquisa responde positivamente à possibilidade técnica. A escolha de tecnologia e a declaração de prontidão para produção permanecem condicionadas aos ensaios acima. Coletor somente online continua como última alternativa; validação posterior continua descartada.
