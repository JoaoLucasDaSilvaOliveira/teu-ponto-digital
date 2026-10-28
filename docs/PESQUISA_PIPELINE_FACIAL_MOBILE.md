# Pipeline facial no celular: execução, compatibilidade e confiança

Pesquisa documental em 26/10/2026. Fontes primárias consultadas: documentação oficial, código dos mantenedores e artigos dos autores. Nenhum motor foi instalado ou medido no projeto. Este documento informa uma decisão futura; não seleciona fornecedor e não altera as modalidades já aceitas no ADR 0005.

## Complemento de viabilidade

A pesquisa foi aprofundada em [parecer de viabilidade do fluxo completo](VIABILIDADE_VALIDACAO_FACIAL_LOCAL.md) e [prova de vida local](PESQUISA_PROVA_VIDA_LOCAL.md), incluindo análise do frontend existente, opções de modelos e limites de licença/offline.

## Conclusão

**Sim: captura, extração do vetor facial e comparação podem ocorrer no próprio celular, sem servidor de inferência e sem embarcar um interpretador Python.** O ONNX Runtime oferece inferência em Android e iOS. Flutter pode integrar bibliotecas nativas pelos canais de plataforma. A disponibilidade do runtime não garante que qualquer modelo caiba, seja rápido ou tenha licença apropriada. [ONNX Runtime Mobile](https://onnxruntime.ai/docs/get-started/with-mobile.html), [integração nativa Flutter](https://docs.flutter.dev/platform-integration/platform-channels).

Há evidência acadêmica concreta de extração facial rápida em telefone, descrita abaixo. Ainda assim, desempenho de um coletor completo, compatibilidade com os vetores existentes e prova de vida precisam ser avaliados separadamente.

## O pipeline atual pode continuar conceitualmente igual

O exemplo oficial da dlib faz em C++ o que costuma ser chamado de pipeline Python: detecta o rosto, alinha um recorte de 150×150 com landmarks, executa a rede e produz um vetor de 128 dimensões. Seu exemplo Python expõe o mesmo conceito. Python é uma interface possível; não é uma propriedade obrigatória do reconhecimento. [Exemplo C++ dlib](https://www.dlib.net/dnn_face_recognition_ex.cpp.html), [exemplo Python dlib](https://github.com/davisking/dlib/blob/master/python_examples/face_recognition.py).

Fluxo proposto para avaliação no Teu Ponto:

1. Provisionar no coletor os modelos e os cadastros autorizados para seu CNPJ.
2. Receber o identificador do funcionário e selecionar seu cadastro.
3. Capturar imagem ou sequência da câmera.
4. Detectar rosto e landmarks; verificar condições de captura e alinhar o recorte.
5. Executar o modelo de embeddings localmente.
6. Comparar o vetor recém-extraído com o vetor daquele cadastro, usando a métrica e o limiar validados para o modelo.
7. Executar a prova de vida exigida, também localmente, antes de confirmar a identidade.
8. Persistir o evento e sincronizar depois, caso esteja offline.

Essa sequência é uma proposta de arquitetura, não um SDK já implementado. O modelo de prova de vida pode organizar os passos de captura de outra maneira. A validação de identidade posterior continua descartada pelo usuário.

**Detectar não é reconhecer.** O ML Kit detecta faces, landmarks, expressões e rastreia faces em frames, mas declara expressamente que não reconhece indivíduos. Um identificador de rastreamento não é a matrícula da pessoa. [ML Kit Face Detection](https://developers.google.com/ml-kit/vision/face-detection).

## Reutilizar os vetores atuais exige identificar o modelo exato

O repositório privado indicado pelo usuário foi acessado pela autenticação local do GitHub após a conexão do aplicativo retornar 404. Inspeção da branch main, referência observada `613bc0ddffa61a106c515bc5652b616b9ab4fe66`:

- [requirements-face.txt](https://github.com/JoaoLucasDaSilvaOliveira/TeuPontoDigital/blob/613bc0ddffa61a106c515bc5652b616b9ab4fe66/services/identity-service/requirements-face.txt): InsightFace 0.7.3, ONNX Runtime 1.24.4, NumPy e OpenCV.
- [face_embedder.py](https://github.com/JoaoLucasDaSilvaOliveira/TeuPontoDigital/blob/613bc0ddffa61a106c515bc5652b616b9ab4fe66/services/identity-service/scripts/face_embedder.py): decodifica imagem, instancia FaceAnalysis com CPUExecutionProvider, extrai normed_embedding; o modelo vem de argumento de configuração.
- [config.go](https://github.com/JoaoLucasDaSilvaOliveira/TeuPontoDigital/blob/613bc0ddffa61a106c515bc5652b616b9ab4fe66/services/identity-service/internal/config/config.go): padrão buffalo_l e detecção 640x640, substituíveis por configuração. Isso não confirma valores do ambiente de produção nem os hashes dos pesos baixados.
- [face_embedding.go](https://github.com/JoaoLucasDaSilvaOliveira/TeuPontoDigital/blob/613bc0ddffa61a106c515bc5652b616b9ab4fe66/services/identity-service/internal/core/domain/face_embedding.go): identifica templates InsightFace como INSIGHTFACE_BUFFALO_L_512D e distingue legado FACE_API_128D_V1.
- [facial_usecase.go](https://github.com/JoaoLucasDaSilvaOliveira/TeuPontoDigital/blob/613bc0ddffa61a106c515bc5652b616b9ab4fe66/services/identity-service/internal/core/usecases/facial_usecase.go): normaliza embeddings e compara cosseno em Go; limiar padrão 0.84. Valida imagem pelo servidor, não usa simplesmente o descritor enviado pelo cliente como fonte de verdade.

Consequência: estudar primeiro a execução dos mesmos pesos ONNX e pipeline no aplicativo nativo, sem portar Python inteiro. Reutilização dos cadastros ainda exige obter os arquivos efetivamente usados, conferir versões/hash e medir paridade. Não aplicar os benchmarks MobileFaceNets ao buffalo_l. O script inspecionado extrai embedding de imagem estática e não demonstra prova de vida; isso não é uma auditoria completa do frontend e de todos os fluxos de antifraude.

O exemplo dlib abaixo é uma alternativa explicativa, não a biblioteca identificada no projeto.

No código oficial ArcFaceONNX do InsightFace, a inferência depende de alinhamento `norm_crop`, dimensões de entrada, troca de canais e normalização. A comparação usa similaridade de cosseno. Transportar apenas o arquivo ONNX sem reproduzir o restante não reproduz automaticamente os resultados do servidor. [Implementação ArcFaceONNX](https://github.com/deepinsight/insightface/blob/master/python-package/insightface/model_zoo/arcface_onnx.py).

Como inferência técnica dessas implementações, a migração tem três situações:

| Situação | Consequência |
|---|---|
| Mesmo modelo/pesos e pipeline equivalente | Pode permitir reutilizar os vetores; comprovar paridade com imagens e pares de referência. Diferenças numéricas de runtime ainda precisam ser medidas. |
| Troca por outro modelo menor | Não presumir compatibilidade com os vetores anteriores, mesmo que tenham o mesmo número de dimensões. Regerar cadastros com o novo modelo ou manter explicitamente a versão anterior. |
| Quantização ou mudança de precisão | Reavaliar acurácia, diferenças de escores e decisões próximas do limiar antes de publicar. |

Compatibilidade entre modelos é um problema de pesquisa próprio: o trabalho de treinamento retrocompatível do CVPR 2020 propõe um método específico para comparar representações novas com antigas sem recalcular toda a base. Portanto, não é uma propriedade automática de redes faciais. [Towards Backward-Compatible Representation Learning](https://openaccess.thecvf.com/content_CVPR_2020/html/Shen_Towards_Backward-Compatible_Representation_Learning_CVPR_2020_paper.html).

A documentação do ONNX Runtime registra que quantização pode reduzir a acurácia e recomenda avaliar e depurar essa perda. Não tratar conversão para INT8 como uma otimização transparente. [Quantização ONNX Runtime](https://onnxruntime.ai/docs/performance/model-optimizations/quantization.html).

Proposta de metadados para cada cadastro: identificador da pessoa e CNPJ, identificador/hash dos pesos, versão do pré-processamento, dimensão, normalização, métrica, política de limiar e versão do cadastro. Essa é uma recomendação do projeto para impedir mistura acidental de espaços vetoriais.

Se só houver vetores antigos e o novo modelo não for compatível, eles não substituem imagens de cadastro para recalcular embeddings. As opções práticas são manter o modelo anterior ou realizar novo cadastro; não há conversão universal comprovada nesta pesquisa. A política de retenção de imagens deve ser decidida por necessidade e proteção de dados, não presumida como consequência da migração.

## Um celular intermediário consegue?

O artigo MobileFaceNets apresenta um experimento concreto: na tabela 2, CPU Snapdragon 820, quatro threads e biblioteca NCNN, a configuração 112×112 levou 24 ms por inferência; a variante 96×96 levou 18 ms. O modelo MobileFaceNet reportado tem arquivo de 4 MB. **Esses números são de extração da representação nas condições do artigo: não são tempo total da batida, não incluem prova de vida e não representam RAM total do aplicativo.** [MobileFaceNets, especialmente tabela 2](https://arxiv.org/pdf/1804.07573).

Isso evidencia que a extração facial móvel não depende necessariamente de hardware topo de linha. Como inferência de engenharia, um intermediário atual é um candidato razoável para avaliação com modelo apropriado. Não há evidência neste projeto para garantir qualquer aparelho intermediário, câmera, modelo pesado já utilizado no servidor ou SDK escolhido.

A comparação 1:1 de dois vetores é uma operação numérica pequena em relação à execução da rede; reduzir a busca a um funcionário ajuda, mas não elimina detecção, extração e prova de vida. No coletor compartilhado, homologar alguns aparelhos e condições de instalação é mais controlável do que garantir a mesma experiência em qualquer celular pessoal.

## Prova de vida é outro requisito

A extração de embeddings e a comparação de identidades não demonstram, sozinhas, que a captura veio de uma pessoa presente. O NIST diferencia comparação biométrica, detecção de apresentação e ataques por injeção; admite execução local de comparação/PAD e recomenda avaliar a integridade do sensor e do dispositivo. A referência é orientação técnica de identidade digital dos EUA, não certificação REP-P brasileira. [NIST SP 800-63B-4, biometria](https://pages.nist.gov/800-63-4/sp800-63b.html).

Logo, portar o pipeline Python existente não acrescenta automaticamente resistência a foto impressa, vídeo em outra tela ou câmera adulterada. Piscar/sorrir detectado não deve ser rotulado como prova de vida robusta sem avaliação específica.

Como evidência de que há produtos móveis com esses componentes, o Face SDK da 3DiVi documenta verificação 1:1, qualidade de imagem, liveness 2D ativo/passivo e interfaces Android/iOS/Flutter. A documentação de capacidades não comprova desempenho, resistência a ataques ou operação offline da combinação contratada; esses itens exigem confirmação e ensaio. [Introdução oficial Face SDK](https://docs.3divi.ai/face_sdk/overview/).

## Confiar no resultado local exige proteger o coletor

Mover inferência para o aparelho muda a fronteira de confiança: o servidor passa a receber uma afirmação produzida naquele aparelho. Assinar um evento autentica a chave que o assinou; não demonstra por si só que a inferência e a câmera não foram adulteradas. Essa é uma inferência de arquitetura de segurança.

O Android Keystore pode manter chaves não exportáveis e protegidas por hardware, mas a própria documentação avisa que um processo comprometido pode conseguir usar a chave sem extraí-la. Assim, não prometer que assinatura local torna impossível fabricar um resultado facial. [Android Keystore](https://developer.android.com/privacy-and-security/keystore).

Play Integrity oferece sinais sobre app/dispositivo e o Google recomenda combiná-los com outras medidas; não é atestado de que um modelo facial específico foi executado corretamente. Não presumir verificação remota nova durante uma queda de rede. [Play Integrity](https://developer.android.com/google/play/integrity/overview).

Propostas para o protótipo: coletores gerenciados em modo quiosque, cadastros locais mínimos criptografados, provisionamento autenticado, modelos e políticas versionados, chave por coletor, eventos identificáveis para reenvio sem duplicação e auditoria. Persistência deve distinguir validação local concluída de envio ainda pendente. Perda de rede impede conhecer imediatamente revogações novas; definir janela de autorização offline e tratamento na reconexão.

## Licença pode mudar a escolha

Código disponível e modelo liberado para uso comercial são verificações diferentes. O InsightFace informa licença MIT para seu código, mas seus modelos pré-treinados disponibilizados seguem política de pesquisa não comercial; o repositório indica contato para licença de reconhecimento. Não incorporar pesos ao aplicativo comercial presumindo que a licença do código os cobre. [Repositório oficial e política de licença](https://github.com/deepinsight/insightface).

## Próximo experimento proposto, sem compromisso de implementação

1. Identificar a biblioteca atual, modelo e pesos, versões, limiar e método de cadastro; verificar licença.
2. Reproduzir em Android uma imagem de referência com o mesmo pipeline, sem Python no aplicativo, e comparar vetores/escores com o servidor.
3. Medir captura até decisão em aparelho real: inicialização fria, p50/p95, memória, filas de marcações e aquecimento.
4. Avaliar reconhecimento e prova de vida separadamente, incluindo falsos aceites/rejeições e ataques; não usar somente acurácia em benchmark público.
5. Reiniciar o coletor em modo avião e verificar modelos, cadastro, licença, validação imediata e persistência. Nenhuma identidade pode depender de validação posterior.
6. Só então decidir entre portar modelo, contratar SDK local ou adotar online como última alternativa.

O escopo aceito continua: aparelho pessoal com autenticação do dispositivo; coletivo com identificador e comparação facial. Reconhecimento facial local também pode existir em celular pessoal, mas essa possibilidade técnica não muda automaticamente a política de produto escolhida.
