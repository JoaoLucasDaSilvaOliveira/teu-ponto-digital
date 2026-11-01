---
status: accepted
---
# Contas por CNPJ e novos cadastros na recontratação

Decidido em 01/11/2026 durante a especificação da issue #4: cada usuário pertence a um único CNPJ. A mesma pessoa pode possuir contas distintas em empresas diferentes, inclusive com o mesmo e-mail; alteração de dados, confirmação, bloqueio e encerramento de uma conta não alteram a outra. Esta decisão substitui a identidade de acesso compartilhada entre tenants prevista inicialmente em RF-01 e no ADR 0004, preservando todas as exigências de isolamento e autorização daquele ADR.

Escolhemos independência cadastral por empresa em vez da conveniência de uma conta global com múltiplos vínculos. Pessoa representa dados pessoais no CNPJ; usuário representa acesso; funcionário representa o trabalhador; vínculo representa autorização. Essas relações são expressas por composição e referências, sem hierarquia de herança nem duplicação dos dados pessoais por perfil. Responsável legal é uma relação com a empresa, não um quinto perfil de acesso. Pessoas externas podem ter conta e perfil sem cadastro de funcionário; funcionário pode existir sem conta.

Uma recontratação cria novo cadastro de funcionário, nova matrícula e, quando houver acesso, novo usuário, preservando os cadastros encerrados e suas referências históricas. CPF e e-mail podem se repetir no histórico, mas só pode haver um funcionário ativo por CPF e um usuário não encerrado por e-mail no mesmo CNPJ. Bloqueio não libera o e-mail. Matrículas antigas continuam reservadas no CNPJ. A pessoa cadastral pode ser reutilizada dentro desse CNPJ; não reassociar registros antigos ao novo usuário ou funcionário.

Consequência: login, recuperação de acesso, consultas e garantias de unicidade deverão incluir o CNPJ; e-mail ou CPF isolados não selecionam uma conta global. Os contratos detalhados e limites das validações locais estão em [Contratos de domínio](../ESPECIFICACAO_DOMINIO_IDENTITY_SERVICE.md).
