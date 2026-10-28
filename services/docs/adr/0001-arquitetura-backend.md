---
status: accepted
---
# Serviços separados com dependências orientadas ao domínio

Decidido em 28/10/2026: cada serviço terá sua estrutura em `services/<nome-do-serviço>/`, com `cmd/` e `internal/` próprios, em vez de um backend único diretamente em `services/`. Isso mantém explícita a organização por serviço prevista para o backend; os nomes, limites de negócio e a divisão dos módulos Go serão definidos conforme a implementação.

Dentro de cada serviço, `domain` concentra entidades, value objects, regras e contratos de domínio; `application` coordena casos de uso; `interfaces` adapta HTTP e eventos de entrada; `infra` implementa persistência e integrações. Contratos pequenos ficam na camada que os consome; os bootstraps em `cmd/` compõem as implementações por injeção explícita de dependências. A separação mantém regras de negócio independentes de transporte e persistência e permite testá-las antes de implementar os adaptadores.

As instruções operacionais estão em [services/AGENTS.md](../../AGENTS.md). O bootstrap existente em `services/cmd/main/main.go` ainda precisa ser alocado a um serviço quando seu nome e responsabilidade forem definidos; este registro não move nem implementa código.
