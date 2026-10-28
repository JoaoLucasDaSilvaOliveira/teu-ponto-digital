---
status: accepted
---
# Priorizar sincronização em background com fallback ao retomar

Decidido em 26/10/2026: tentar sincronizar marcações offline em background até onde o navegador ou sistema operacional permitir. Sincronizar ao abrir ou retomar a aplicação é um fallback obrigatório, não o único mecanismo previsto. A decisão busca reduzir pendências sem depender de garantias de execução que a plataforma não oferece.

Flutter permanece a tecnologia de interface. O [ADR 0003](0003-flutter-nativo-e-pwa-condicional.md) define a distribuição nativa como prioridade e a PWA como alternativa condicionada à validação de capacidades. A ausência de Background Sync em uma plataforma exige o fallback; um worker comum não substitui uma API de execução em background. A escolha de distribuição deve considerar testes em aparelhos reais e os limites descritos na [revisão técnica](../REVISAO_DOCUMENTACAO.md).
