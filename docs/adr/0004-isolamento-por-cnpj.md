---
status: accepted
---
# Isolamento por CNPJ

Decidido em 26/10/2026: cada CNPJ constitui um tenant com dados, permissões e configurações próprios. A previsão inicial de uma identidade de gestor compartilhada entre tenants foi substituída em 01/11/2026 pelo [ADR 0007](0007-contas-por-cnpj-e-recontratacao.md): as contas são independentes por CNPJ. Selecionar um contexto não concede autorização sobre os demais.

Toda operação deve verificar no servidor o vínculo e a permissão no CNPJ correspondente, inclusive processamento assíncrono, exportações e acesso a arquivos. Referências entre registros devem pertencer ao mesmo contexto autorizado. A fronteira lógica não exige um banco físico por CNPJ; a estratégia de persistência continua a ser definida.
