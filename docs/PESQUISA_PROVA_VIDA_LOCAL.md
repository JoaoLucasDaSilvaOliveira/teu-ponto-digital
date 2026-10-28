# Prova de vida local: fornecedores, operação offline e limites de evidência

> Atualização de 27/10/2026: o [ADR 0006](adr/0006-prova-de-vida-passiva.md) define prova de vida passiva. Opções ativas citadas nesta pesquisa não representam o fluxo escolhido. Modelo/SDK e desempenho seguem pendentes de validação.

Pesquisa documental em 26/10/2026, complementar à [pesquisa do pipeline facial móvel](PESQUISA_PIPELINE_FACIAL_MOBILE.md). Nenhum SDK foi instalado, executado ou contratado. Não há fornecedor selecionado. Preservam-se as decisões do usuário: validação de identidade posterior está descartada; coletor somente online é a última alternativa.

## Conclusão para o projeto

Há componentes comerciais e abertos capazes de executar detecção de apresentação facial no dispositivo. Contudo, **capacidade local, prazo de funcionamento offline, resistência a fraude e licença de produção são quatro verificações distintas**. A pesquisa encontrou uma candidata comercial com módulos móveis e licença offline documentados (3DiVi), uma candidata com limitações operacionais relevantes (FacePing) e uma base aberta para experimentação (MiniFASNet/Silent-Face). Nenhuma está homologada para um coletor fixo autônomo do Teu Ponto.

Essa classificação é uma recomendação de pesquisa, não um resultado comparativo de laboratório.

## Matriz de evidência e testes pendentes

| Opção | Evidência documentada | Limitação relevante | Verificação pendente |
|---|---|---|---|
| 3DiVi Face SDK | Componentes RGB passivos e ativos; variantes móveis; Android/iOS e Flutter; licença offline Android vinculada ao hardware. [Módulos](https://docs.3divi.ai/face_sdk/components/), [liveness](https://docs.3divi.ai/face_sdk/development/processing_block_api/dev_pb_liveness_estimation/), [APIs](https://docs.3divi.ai/face_sdk/language_apis/), [licença](https://docs.3divi.ai/face_sdk/licenses/). | Bibliotecas fechadas; licença online exige validação periódica; variante leve troca proteção/desempenho e não equivale automaticamente às demais. | Confirmar pacote e API Flutter exatos; cotação; modelo comercial; ativação/renovação offline; ensaio do binário contratado no aparelho real e evidência independente aplicável. |
| FacePing | Tutorial Android em modo avião; documentação Flutter; verificação local e distribuição protegida de cadastros. [Android](https://docs.faceping.ai/start/android/), [Flutter](https://docs.faceping.ai/start/flutter/). | Android em preview; produção depende de provisionamento do serviço; lease offline de 24 h. A documentação admite risco de vídeo contra o desafio de giro de cabeça. [Operação](https://docs.faceping.ai/devices/). | Avaliação de ataques, termos dos modelos e SDK, adequação a coletor sem supervisão, custo e contrato; aceitar ou rejeitar explicitamente a janela de 24 h. |
| Silent-Face-Anti-Spoofing / MiniFASNet | Código de treino/inferência, modelos e implantação Android publicados pelo mantenedor; repositórios sob Apache-2.0. [Projeto](https://github.com/minivision-ai/Silent-Face-Anti-Spoofing), [Android](https://github.com/minivision-ai/Silent-Face-Anti-Spoofing-APK). | O mantenedor avisa que RGB depende da câmera e do cenário. Não fornece neste material uma homologação independente para nosso caso. | Auditar artefatos e pesos exatos, pré-processamento, conversão e integração; testar ataques/desempenho; definir manutenção e atualização de ameaças. |

## 3DiVi: o que está documentado de maneira operacional

O SDK distingue liveness ativo, RGB passivo, profundidade e infravermelho. Não pressupor que um telefone comum possui os sensores das duas últimas modalidades. O módulo RGB pode operar com imagem de câmera convencional; o ativo analisa ações conforme um roteiro. [Componentes](https://docs.3divi.ai/face_sdk/components/).

A documentação específica oferece `2d_light` e `2d_ensemble_light` como variantes adequadas a dispositivos móveis. As tabelas publicadas nessa página, porém, medem latência em **Intel Xeon E5-2683 v4 e NVIDIA GTX 10xx**, não em celular. Os resultados de `2d` e `2d_light` explicitamente não incluem ataques com máscaras 2D/3D. Portanto, esses números não demonstram desempenho ou cobertura do futuro coletor. [Página técnica de liveness](https://docs.3divi.ai/face_sdk/development/processing_block_api/dev_pb_liveness_estimation/).

Há APIs nativas e plugin Flutter para Android/iOS. Isso documenta a rota de integração; ainda é preciso conferir se o wrapper expõe o módulo selecionado e se será necessária ponte nativa adicional. O exemplo Android do fornecedor demonstra liveness RGB em câmera do aparelho. [APIs](https://docs.3divi.ai/face_sdk/language_apis/), [demonstração Android](https://docs.3divi.ai/face_sdk/samples/java/java_demo/).

O licenciamento é decisivo: a opção online exige ativação e validação periódica, enquanto a offline é vinculada a hardware ou token; a página descreve obtenção de Hardware ID no Android e emissão do arquivo de licença. A página também apresenta vínculo por aplicativo móvel. Precisamos confirmar comercialmente a combinação válida de módulos, plataforma e bloqueio, sem presumir que todas as modalidades são intercambiáveis ou perpétuas. [Licenciamento oficial](https://docs.3divi.ai/face_sdk/licenses/).

O produto é comercial, entregue como bibliotecas fechadas; o fornecedor anuncia operação offline e modalidades anuais/perpétuas, sujeitas à proposta final. Nenhum preço foi assumido nesta pesquisa. [Página oficial do produto](https://3divi.ai/products/software/face-sdk).

## FacePing: offline documentado, mas limitado

O tutorial Android usa a versão `1.0.0-preview.3`, inclui ONNX Runtime e modelos na biblioteca e demonstra verificação com giro de cabeça em modo avião. É evidência de implementação documentada, não teste executado por nós. [Tutorial Android](https://docs.faceping.ai/start/android/).

Na operação de produção, o dispositivo sincroniza os cadastros autorizados e continua verificando sem conexão **por 24 horas**; quando o lease expira, para até sincronizar. A documentação diz que a verificação ao vivo combina avaliação passiva e direção aleatória da cabeça, mas admite que um vídeo reproduzido pode superar o giro e recomenda manter equipe à vista em dispositivos sem operador. Isso é incompatível com prometer um coletor autônomo resistente a replay sem avaliação adicional. Revogação só chega na próxima sincronização. [Dispositivos, anti-spoofing e retirada](https://docs.faceping.ai/devices/).

Não confundir a licença sandbox renovada pelo menos a cada 30 dias com o lease de produção de 24 horas. O fornecedor descreve núcleo nativo em Rust, templates criptografados, modelos ONNX e avisos de terceiros dentro do pacote. Não abrimos esse pacote; os direitos dos modelos específicos ainda precisam de conferência. [Segurança](https://docs.faceping.ai/security/).

Produção exige aprovação do caso e contrato; cadastros são feitos via serviço e sincronizados nos dispositivos. Assim, verificação local não significa independência operacional total do fornecedor. [Entrada em produção](https://docs.faceping.ai/going-live/).

## Alternativa aberta: Silent-Face/MiniFASNet

O projeto do mantenedor inclui arquitetura, pré-processamento, treino/teste e modelos, com implantação Android separada. O README distingue modelo aberto e variante de maior precisão privada, além de registrar que RGB é sensível ao modelo da câmera e ao cenário. Não extrapolar resultados da variante privada para os pesos abertos. Os números de latência publicados são resultados declarados pelo autor e não medição independente nem tempo total de uma batida. [README do mantenedor](https://github.com/minivision-ai/Silent-Face-Anti-Spoofing/blob/master/README_EN.md).

O repositório principal e o Android apresentam Apache-2.0. Não encontrei nas páginas consultadas uma exceção explícita impondo uso não comercial aos pesos originais, mas também não há ali uma ficha detalhada de proveniência/licença por peso e conjunto de treino. A próxima auditoria deve fixar hashes e origem de cada arquivo e conferir dependências; não copiar automaticamente a declaração de licença de um exportador terceiro de ONNX. [Licença do projeto](https://github.com/minivision-ai/Silent-Face-Anti-Spoofing/blob/master/LICENSE), [repositório Android](https://github.com/minivision-ai/Silent-Face-Anti-Spoofing-APK).

É uma base razoável para estudar integração e medir uma linha de base. Transformá-la em defesa mantida contra ataques exige testes, política de atualização e responsabilidade técnica do produto. Essa é uma avaliação de engenharia, não afirmação de que código aberto seja intrinsecamente mais ou menos seguro.

## Evidência independente: conferir o que foi testado

O NIST separa avaliações FRTE de reconhecimento das avaliações FATE de análise/PAD. Boa colocação de matching não comprova resistência a foto ou vídeo. [Escopo FRTE/FATE](https://pages.nist.gov/frvt/).

O relatório independente NISTIR 8491 avaliou 82 algoritmos de 45 desenvolvedores em imagens arquivadas. Mostra grande variação entre algoritmos e tipos de ataque; não avaliou interações em tempo real, profundidade ou infravermelho. Esse resultado apoia medir cobertura por ameaça, mas não certifica um aplicativo móvel completo. [Relatório NISTIR 8491](https://nvlpubs.nist.gov/nistpubs/ir/2023/NIST.IR.8491.pdf).

Não localizei, nesta pesquisa, relatório independente que estabeleça resistência PAD da combinação exata **3DiVi modelo selecionado + plugin Flutter + aparelho + operação offline**, nem relatório equivalente do FacePing ou dos pesos MiniFASNet selecionáveis. Isso descreve o limite da busca, não prova inexistência de avaliação.

Um exemplo do cuidado necessário: Luxand anuncia add-on móvel com avaliação iBeta; a carta pública citada foi realizada no aplicativo LivenessRecognition 1.3, Pixel 8/Android 14, **dependente do backend `idface-server-prod 1.41.0`**. Portanto, essa carta, isoladamente, não demonstra que a versão local offline herdará o resultado. [Documentação Luxand](https://www.luxand.com/facesdk/documentation/certifiedliveness.php), [carta primária iBeta, 05/01/2024](https://www.ibeta.com/wp-content/uploads/2024/01/240105-Luxand-PAD-Level-1-Confirmation-Letter.pdf).

## Recomendação de próxima etapa

1. **Primeira candidata comercial:** avaliar 3DiVi com módulo RGB móvel e licença offline específicos, condicionando a seleção a prova de conceito e evidência de ataques.
2. **Linha de base aberta:** avaliar MiniFASNet em paralelo técnico, sem classificá-lo como prova de vida homologada.
3. **FacePing como candidata secundária:** a janela de 24 h e a ressalva sobre replay/supervisão precisam ser aceitáveis; hoje não recomendar para prometer coletor autônomo robusto.
4. Solicitar de qualquer fornecedor versão/modelo exatos, fluxo de captura, ataques cobertos, métricas APCER/BPCER e relatório completo, aparelhos testados, licença, renovação, atualização de modelos e comportamento com relógio alterado/reinicialização offline.
5. Em aparelho empresarial real, medir fluxo completo, filas, memória/aquecimento, iluminação, alturas, óculos, rejeições legítimas e ataques com papel/tela/vídeo. Testes internos não substituem certificação, mas revelam inadequação operacional.

Critério de aceite proposto: identidade decidida no instante da coleta, sem enviar captura para validação posterior; sem queda silenciosa para método mais fraco; cadastros e licença previamente provisionados; limite offline e contingência explicitamente definidos. Nenhuma política foi aprovada ou fornecedor contratado como consequência deste documento.
