---
status: accepted
---
# Prova de vida passiva na validação facial

Decidido pelo usuário em 27/10/2026: utilizar prova de vida passiva na modalidade de aparelho coletivo, junto da comparação facial com o cadastro do funcionário. A escolha busca reduzir atrito na marcação, sem exigir ações como piscar, sorrir ou virar a cabeça como fluxo padrão.

A verificação deve analisar a captura da mesma sessão usada na comparação facial. Detecção de rosto, comparação de embeddings e prova de vida são verificações distintas. Não substituir prova de vida passiva por desafio ativo ou por simples detecção de rosto sem revisar esta decisão.

Priorizar avaliação da execução local para preservar a validação no momento da coleta quando offline. Capturar para validar identidade posteriormente permanece descartado; coletor somente online com contingência é a última alternativa. Modelo, SDK, licença, limiares e hardware suportado ainda dependem de avaliação. Esta decisão de experiência não homologa a resistência a ataques de qualquer fornecedor.

As pesquisas de candidatos anteriores a esta decisão são evidências técnicas, não escolhas aprovadas. Candidatos que dependem de movimentos como etapa obrigatória não atendem ao fluxo passivo escolhido. Ver [viabilidade local](../VIABILIDADE_VALIDACAO_FACIAL_LOCAL.md) e [pesquisa de PAD](../PESQUISA_PROVA_VIDA_LOCAL.md).
