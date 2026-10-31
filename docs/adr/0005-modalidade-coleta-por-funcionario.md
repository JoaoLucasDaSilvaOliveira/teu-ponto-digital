---
status: accepted
---
# Modalidade de coleta definida pela empresa por funcionário

A empresa define, por funcionário no contexto do CNPJ, a modalidade de coleta: aparelho próprio ou aparelho coletivo. No aparelho próprio, utilizar autenticação do dispositivo para facilitar o uso diário. No aparelho coletivo, o funcionário apresenta um identificador e passa por comparação facial com seu cadastro e prova de vida passiva, conforme [ADR 0006](0006-prova-de-vida-passiva.md). Atualização aprovada em 30/10/2026: a matrícula interna aleatória do funcionário pode ser usada como esse identificador, conforme RF-05.1 e RF-11.1. A forma de apresentação na interface do coletor será detalhada antes de sua implementação.

A modalidade é uma configuração autorizada pela empresa; o funcionário não escolhe livremente outro método durante a marcação. Identificador no coletor é uma declaração de identidade, não um substituto da comparação facial. Prova de vida passiva foi escolhida no ADR 0006. Políticas de fallback, cadastro inicial, mudança de modalidade e implementação da validação offline ainda precisam ser especificadas. Não presumir que PIN do aparelho identifica o funcionário no coletor coletivo.
