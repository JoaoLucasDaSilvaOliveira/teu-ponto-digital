# Teu Ponto

Vocabulário do controle de jornada e da organização das empresas no produto.

## Linguagem

**Tenant**: Contexto de um CNPJ, com dados, contas, permissões e configurações próprios. Uma pessoa pode atuar em vários tenants por cadastros e contas independentes.

**Empresa**: Cadastro do empregador identificado por um CNPJ, correspondente a um tenant no produto.

**Pessoa**: Cadastro de dados pessoais no contexto de um CNPJ, com CPF e nome e, quando necessário, nascimento. Pode estar relacionado a acesso, ao cadastro de trabalhador e à responsabilidade legal, sem duplicação dos dados por perfil.

**Responsável legal**: Pessoa identificada por nome e CPF, designada no cadastro da empresa. Essa condição é distinta de possuir identidade de acesso ou perfil de Gestor.

**Gestor**: Perfil exclusivo de administração da empresa no CNPJ da conta, incluindo concessão de ações aos demais perfis, sem ações próprias de ponto. Não determina o enquadramento trabalhista da pessoa.

**Usuário**: Conta de acesso de uma pessoa em um único CNPJ, com e-mail, confirmação e situação de acesso próprios. Contas com o mesmo e-mail em empresas diferentes são independentes.

**E-mail confirmado**: Condição da conta cujo endereço atual passou pela confirmação de titularidade. É distinta da validade sintática do endereço, da situação de acesso da conta e da autorização no CNPJ.

**Funcionário**: Cadastro do trabalhador no contexto de uma empresa, associado a uma pessoa e a uma matrícula interna. Pode existir sem conta; uma recontratação cria outro cadastro, preservando o anterior encerrado.

**Vínculo com empresa**: Relação que autoriza um usuário a atuar em um CNPJ, com perfil, permissões e escopo próprios daquele contexto. É distinto do vínculo empregatício do trabalhador.

**Supervisor**: Perfil que inclui as capacidades próprias do Colaborador e ações sobre seus subordinados conforme concessões do Gestor no CNPJ. Seus próprios pedidos de ajuste dependem de avaliação por Gestor.

**Colaborador**: Perfil de acesso aos próprios dados e às operações do trabalhador, como registrar ponto e solicitar ajustes.

**Contabilidade/DP**: Perfil cujo acesso às funcionalidades de cada empresa depende das permissões concedidas por seu gestor.

**Setor**: Agrupamento organizacional de funcionários de uma mesma empresa. Não determina automaticamente a subordinação a um Supervisor.

**Equipe**: Conjunto de funcionários diretamente atribuídos a um Supervisor no mesmo CNPJ, podendo pertencer a setores distintos.

**Subordinado**: Funcionário atribuído diretamente a um Supervisor no contexto da empresa. Cada funcionário possui no máximo um Supervisor responsável.

**Matrícula interna**: Código aleatório e não sequencial de seis caracteres, com dígitos de 0 a 9 e letras de A a Z em maiúsculas, atribuído pela aplicação ao cadastro do funcionário. Pode ser apresentado no coletor coletivo para declarar a identidade.

**Matrícula da folha**: Código do funcionário no sistema externo de folha, relacionado à matrícula interna por um mapeamento de integração.

**Solicitação de ajuste**: Pedido de correção do tratamento de ponto, acompanhado de justificativa e sujeito à avaliação de um superior. Não equivale a uma marcação original retroativa.

**Modalidade de coleta**: Atribuição da empresa ao funcionário que define o uso de aparelho próprio ou aparelho coletivo para registrar ponto.

**Aparelho próprio**: Modalidade de coleta atribuída ao funcionário para uso de seu dispositivo pessoal com autenticação do dispositivo.

**Aparelho coletivo**: Modalidade de coleta em dispositivo compartilhado da empresa, com identificação do funcionário, comparação facial com seu cadastro e prova de vida passiva.

**Comparação facial**: Verificação de correspondência entre o rosto capturado e o cadastro da identidade declarada. Não determina por si só se a captura veio de uma pessoa presente.

**Prova de vida passiva**: Verificação de sinais de tentativa de apresentação por foto, tela ou outro artefato, sem solicitar ações ao funcionário como parte do fluxo padrão.
