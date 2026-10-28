---
status: accepted
---
# Flutter nativo prioritário; PWA condicionada às capacidades exigidas

O foco do Teu Ponto é o aplicativo nativo em Flutter. Em 26/10/2026, foi considerada distribuição web/PWA como alternativa temporária de acesso por iPhone, pois o responsável pelo projeto ainda não dispõe de Mac para gerar o aplicativo iOS. Essa conveniência de distribuição não autoriza reduzir requisitos do produto.

A PWA só será adotada se atender aos requisitos obrigatórios, incluindo localização para geofencing, captura e validação de identidade facial conforme o método que vier a ser definido, persistência offline e sincronização nos limites já aceitos no ADR 0002. Se a validação de capacidades falhar, manter somente a distribuição nativa e resolver o release iOS posteriormente. Não há compromisso de lançar a PWA.

A existência de uma API não comprova atendimento: validar em aparelhos reais, especialmente iPhone/Safari. Definir precisão e idade máximas da localização, comportamento sem permissão e método de reconhecimento facial antes de aprovar a alternativa web. Geolocalização de alta precisão é uma solicitação, não garantia de coordenadas exatas, inclusive no aplicativo nativo. Captura de câmera e reconhecimento facial são capacidades distintas; WebAuthn não substitui reconhecimento facial do empregado.

Referências: [precisão de localização](https://developer.mozilla.org/en-US/docs/Web/API/GeolocationCoordinates/accuracy), [solicitação de localização](https://developer.mozilla.org/en-US/docs/Web/API/Geolocation/getCurrentPosition), [câmera no navegador](https://developer.mozilla.org/en-US/docs/Web/API/MediaDevices/getUserMedia) e [WebAuthn](https://www.w3.org/TR/webauthn-3/).
