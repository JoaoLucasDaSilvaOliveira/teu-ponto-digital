# Contratos de domínio — identity-service

Especificação da [issue #4](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/4), parte da issue pai #2. Decisões de produto confirmadas em 01/11/2026; entrega documental, sem implementar entidades Go. As interfaces descritas aqui orientam as próximas issues; autorização detalhada continua na #5, dependências na #6 e inicialização do módulo na #7.

Fontes: [plano](PLANO_IMPLEMENTACAO_IDENTITY_SERVICE.md), [limites](LIMITES_IDENTITY_SERVICE.md), [vocabulário](../CONTEXT.md), [requisitos funcionais](REQUISITOS_FUNCIONAIS.md), [arquitetura do backend](../services/docs/adr/0001-arquitetura-backend.md), [isolamento por CNPJ](adr/0004-isolamento-por-cnpj.md), [contas e recontratação](adr/0007-contas-por-cnpj-e-recontratacao.md) e [revisão de coerência](REVISAO_DOCUMENTACAO.md).

## Entidades, valores e composição

| Conceito | Natureza | Identidade e composição |
|---|---|---|
| Empresa | Entidade | Seu ID também é o tenant ID; CNPJ imutável. Referencia a pessoa responsável legal e compõe a referência opcional do A1. |
| Pessoa | Entidade do CNPJ | CPF, nome e nascimento ficam somente aqui. Referenciada por usuários, funcionários e empresa. Sem conta global nem sincronização automática entre CNPJs. |
| Usuário | Entidade do CNPJ | Referencia uma pessoa e compõe e-mail e estados de acesso/confirmação. Não contém cópia de CPF/nome/nascimento. |
| Funcionário | Entidade do CNPJ | Cadastro do trabalhador, referencia pessoa, setor e usuário opcional; compõe matrícula e modalidade de coleta. Não é um perfil de acesso. |
| Setor | Entidade do CNPJ | Agrupamento identificado, renomeável; não contém lista mutável de funcionários. A relação fica no funcionário. |
| Vínculo com empresa | Entidade do CNPJ | Referencia usuário e compõe configuração de autorização, completada pela [issue #5](AUTORIZACAO_IDENTITY_SERVICE.md). Um perfil por configuração; Supervisor inclui capacidades próprias de Colaborador, Gestor é exclusivo. |
| Responsável legal | Relação | Empresa referencia Pessoa identificada por nome/CPF. Não é entidade duplicada nem perfil de acesso. Não exige nascimento ou conta. |
| CPF, CNPJ, Email, Nome, DataNascimento, Matricula, ModalidadeColeta, ReferenciaA1 | Objetos de valor | Imutáveis, comparáveis por conteúdo e substituídos por inteiro. Não possuem ciclo de identidade próprio. |

Dados pessoais não são repetidos em estruturas para Gestor, Supervisor, Colaborador ou Contabilidade/DP. Não há herança `Pessoa → Usuário → Funcionário`: uma pessoa pode existir sem conta, um usuário externo pode existir sem funcionário e um gestor pode também ser funcionário. Não criar interfaces genéricas de pessoa, CRUD ou repositório sem consumidor real.

Os comportamentos variáveis dos quatro perfis são especificados como adaptadores de uma interface pequena de política na [issue #5](AUTORIZACAO_IDENTITY_SERVICE.md). A entidade vínculo recebe essa política para avaliar ações, sem selecionar comportamento pelo tipo concreto. A #4 especifica armazenamento/invariantes comuns; a #5 completa matriz inicial e assinatura da política, sem implementar código nesta rodada documental.

## Relações e cardinalidades

- Empresa tem zero ou muitas pessoas, contas, funcionários, setores e vínculos durante o cadastro. Todas essas entidades pertencem a exatamente um CNPJ, imutável.
- Empresa referencia zero ou uma pessoa como responsável legal e zero ou uma referência corrente do A1. Sem responsável ou A1 o cadastro permanece incompleto; não apagar versões históricas ao substituí-los na persistência futura.
- Pessoa pertence a exatamente uma empresa. Cada CPF corresponde a no máximo uma Pessoa naquele CNPJ; a mesma pessoa natural em outro CNPJ recebe outro cadastro e ID. Não há propagação entre eles.
- Pessoa pode ter vários funcionários históricos e várias contas encerradas; no máximo um funcionário ativo e um usuário não encerrado naquele CNPJ. Uma pessoa externa pode ter zero funcionários. Nascimento é opcional para externos e obrigatório para criar funcionário.
- Funcionário referencia exatamente uma Pessoa; opcionalmente um Usuário da mesma pessoa e CNPJ e um Setor daquele CNPJ. Um usuário pode estar associado a no máximo um funcionário ativo. Associações históricas ficam preservadas.
- Usuário possui zero ou um vínculo corrente com sua empresa; vínculos revogados podem permanecer no histórico. No máximo um vínculo ativo por usuário/CNPJ, cuja configuração representará os perfis/escopos definidos na #5.
- A pessoa responsável legal pode também possuir conta, vínculo e funcionário. Designá-la não concede permissões. Contador externo pode ter conta e perfil sem funcionário; Supervisor não precisa ser um subtipo de Funcionário.

## Convenções dos contratos Go

`entity` contém as entidades e `vo` representa o alias de importação do pacote de objetos de valor. IDs são `uuid.UUID` (`github.com/google/uuid`) e rejeitam `uuid.Nil`. O bootstrap/caso de uso gera IDs; construtores não consultam banco nem geram matrículas. Campos das entidades e valores são privados; somente parâmetros auxiliares de criação podem ser públicos, sem representar entidades mutáveis.

Construtores retornam `(*Entidade, error)` e valores retornam `(Valor, error)`. Toda entidade fornece `ID() uuid.UUID`; as entidades de tenant fornecem `EmpresaID() uuid.UUID`. Leituras retornam valores ou cópias; opcionais retornam `(valor, bool)`. Coleções de permissões/escopos retornam cópias, inclusive dos elementos mutáveis. Nenhum getter expõe ponteiro, slice ou mapa interno para alteração externa.

Getters sem tipo de retorno repetido nas seções abaixo devolvem exatamente o tipo do campo correspondente na tabela. `SituacaoUsuario`, `SituacaoFuncionario` e `SituacaoVinculo` são tipos nomeados de string com os valores fechados listados em cada entidade; valores desconhecidos são inválidos. `PendenciaEmpresa` é um tipo nomeado de string com `RESPONSAVEL_AUSENTE` e `REFERENCIA_A1_AUSENTE`, nessa ordem quando ambas estiverem presentes. Ausência de pendência local não significa autorização operacional.

Construtores usam os próprios métodos de validação/atribuição, inclusive métodos privados para IDs e relações imutáveis. Validam em uma instância temporária e só a retornam em caso de sucesso. Alterações com vários campos validam primeiro os novos valores/referências; só então substituem o estado. Em erro, todos os campos anteriores permanecem iguais. O valor zero de entidade/objeto de valor não é um cadastro válido: operações públicas que os recebem devem validá-los, mesmo que tenham sido construídos sem o construtor.

Referências recebidas como entidades permitem verificar localmente empresa, pessoa e estado; não comprovam carregamento autoritativo. Casos de uso futuros carregam dados atuais, autorizam, controlam concorrência e persistem atomicamente. Não colocar `context.Context`, ORM, HTTP, credenciais ou loggers nas entidades para simular essa coordenação.

## Objetos de valor

| Valor e construtor | Regras e operações |
|---|---|
| `NovoCPF(entrada string) (CPF, error)` | Exatamente 11 dígitos ASCII ou máscara exata `DDD.DDD.DDD-DD`; rejeitar espaços, outros caracteres, sequências de um dígito repetido e DV inválido, sem completar zeros. Representação canônica sem máscara. `String() string`, `Validar() error`. Implementação própria com biblioteca padrão na #8, conforme [avaliação da #6](DEPENDENCIAS_DOMINIO_IDENTITY_SERVICE.md); não comprova identidade, titularidade ou cadastro oficial. |
| `NovoCNPJ(entrada string) (CNPJ, error)` | Exatamente 12 posições ASCII alfanuméricas e dois DVs numéricos, sem máscara ou com máscara exata `XX.XXX.XXX/XXXX-DD`; converter `a-z` para maiúsculas explicitamente, rejeitar caracteres extras, zerado e DV inválido. Guardar string canônica, sem inteiro/padding. `String() string`, `Validar() error`. Implementação própria com biblioteca padrão na #8; não comprova situação cadastral. |
| `NovoEmail(entrada string) (Email, error)` | Endereço simples, sem nome de exibição/lista/CRLF; aparar espaços externos; domínio em minúsculas. Para unicidade nesta aplicação, chave de comparação insensível a caixa em todo endereço; preservar parte local na apresentação. Não remover pontos nem sufixos `+`. `String() string`, `Chave() string`, `Validar() error`. Confirmação continua separada. |
| `NovoNome(entrada string) (Nome, error)` | Aparar espaços externos e rejeitar vazio, UTF-8 inválido e caracteres de controle. Preservar grafia e acentos; não exigir sobrenome nem impor limite inventado por coluna de banco. Usado em nome pessoal, setor, razão social e apelido. `String() string`, `Validar() error`. |
| `NovaDataNascimento(ano int, mes time.Month, dia int, hoje civil.Date) (DataNascimento, error)` | Encapsular `civil.Date`; data real, não zero, ano de 1 a 9999 e não futura em relação a `hoje` válido recebido explicitamente. Não consultar relógio global, nem impor idade mínima não especificada. `Data() civil.Date`, `ValidarEm(hoje civil.Date) error`. |
| `NovaMatricula(entrada string) (Matricula, error)` | Exatamente seis caracteres ASCII em `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`; converter `a-z` para maiúsculas. Rejeitar espaços, hífens, outros símbolos e tamanho diferente, sem aparar/completar. `String() string`, `Validar() error`. |
| `NovaModalidadeColeta(codigo string) (ModalidadeColeta, error)` | Somente `APARELHO_PROPRIO` ou `APARELHO_COLETIVO`; sem modalidade padrão silenciosa. `Codigo() string`, `Validar() error`. Não executa biometria, autenticação ou fallback. |
| `NovaReferenciaA1(empresaID, arquivoID uuid.UUID, naoAntes, naoDepois time.Time) (ReferenciaA1, error)` | IDs não nulos; instantes não zero e `naoAntes < naoDepois`, normalizados a UTC. Contém somente tenant, identificador opaco do arquivo protegido e intervalo informado pelo futuro validador. `EmpresaID()`, `ArquivoID()`, `NaoAntes()`, `NaoDepois()` e `VigenteEm(instante time.Time) bool`, com intervalo `[naoAntes, naoDepois)`. Referência/data não comprova validação criptográfica, cadeia ou revogação. Sem arquivo, caminho de armazenamento, senha ou chave privada. |

Operações retornam erro para entradas inválidas; nenhuma delas escreve logs com CPF, e-mail ou segredos. Instantes de transição são recebidos explicitamente e precisam ser não zero; não são prova de horário confiável ou auditoria por si.

## Empresa

| Campo privado | Tipo | Regra |
|---|---|---|
| id | `uuid.UUID` | Obrigatório e imutável; também identifica o tenant. |
| cnpj | `vo.CNPJ` | Obrigatório e imutável; único entre empresas, inclusive históricas. |
| razaoSocial, apelido | `vo.Nome` | Obrigatórios e alteráveis. |
| responsavelID | `uuid.UUID` opcional | Ausência explícita; referência a Pessoa do mesmo tenant. |
| referenciaA1 | `vo.ReferenciaA1` opcional | Ausência explícita; referência do mesmo tenant. |

Construtor: `NovaEmpresa(id uuid.UUID, cnpj vo.CNPJ, razaoSocial, apelido vo.Nome) (*Empresa, error)`. Cria cadastro inicial com identificação válida, sem responsável/A1. Não permite documentos vazios como rascunho.

Métodos: `AlterarIdentificacao(razaoSocial, apelido vo.Nome) error`, `DefinirResponsavel(pessoa *Pessoa) error`, `AssociarA1(referencia vo.ReferenciaA1) error`. Getters: `CNPJ()`, `RazaoSocial()`, `Apelido()`, `ResponsavelID() (uuid.UUID, bool)`, `ReferenciaA1() (vo.ReferenciaA1, bool)` e `PendenciasLocais() []PendenciaEmpresa`. Esta última informa ausência de responsável e referência do A1; não recebe booleano externo para fingir verificar vínculo gestor ou validade criptográfica.

Estado cadastral completo é resultado coordenado de identificação, responsável, referência validada/vigente e pelo menos um vínculo gestor ativo. Não armazenar um booleano local que fique obsoleto quando certificado expirar ou vínculo for revogado. Rascunho permite completar o cadastro, sem autorizar produção de artefatos. Ausência ou vencimento do A1 não descarta nem bloqueia a marcação do trabalhador: preservar e sinalizar pendência. Exportação que dependa de assinatura válida fica pendente até regularização; não prometer conformidade dos signatários. O fluxo de entrada em operação e a validação regulatória continuam para a área de coleta.

Erros locais: ID/documento/nome inválido, referência nula/inválida ou de outro tenant. CNPJ duplicado, ausência de gestor e falha criptográfica são erros da coordenação futura, não do construtor.

## Pessoa

| Campo privado | Tipo | Regra |
|---|---|---|
| id, empresaID | `uuid.UUID` | Obrigatórios e imutáveis. |
| cpf | `vo.CPF` | Obrigatório; correção autorizada preserva ID e demais relações. |
| nome | `vo.Nome` | Obrigatório e alterável. |
| nascimento | `vo.DataNascimento` opcional | Necessário antes de criar funcionário. Depois de informado, pode ser corrigido, mas não apagado nesta rodada. |

Construtor: `NovaPessoa(id, empresaID uuid.UUID, cpf vo.CPF, nome vo.Nome) (*Pessoa, error)`. Sem nascimento inicial para permitir responsável legal e participante externo sem coleta desnecessária dessa informação.

Métodos: `AlterarNome(nome vo.Nome) error`, `CorrigirCPF(cpf vo.CPF) error`, `DefinirNascimento(data vo.DataNascimento, hoje civil.Date) error`. Getters: `CPF()`, `Nome()`, `Nascimento() (vo.DataNascimento, bool)`.

Pessoa não tem perfil nem situação de login/emprego. `CorrigirCPF` apenas valida/atribui localmente; o caso de uso futuro exige autorização, unicidade por CNPJ e auditoria. A correção não muda contas, IDs ou matrículas. Mudanças ficam dentro do CNPJ. A pessoa pode ser reutilizada na recontratação, mas as evidências históricas da época de cada operação devem continuar preservadas na auditoria futura; não usar dados pessoais atuais como prova de conteúdo histórico.

## Usuário

| Campo privado | Tipo | Regra |
|---|---|---|
| id, empresaID, pessoaID | `uuid.UUID` | Obrigatórios e imutáveis. |
| email | `vo.Email` | Obrigatório e alterável enquanto não encerrado. |
| revisaoEmail | `uint64` | Inicia em 1; incrementa ao alterar a chave do e-mail. Revisão esgotada rejeita alteração. |
| confirmadoEm | `time.Time` opcional | Ausente inicialmente; confirmação vinculada ao e-mail/revisão atuais. |
| situacao | `SituacaoUsuario` | `ATIVO`, `BLOQUEADO` ou `ENCERRADO`. |
| encerradoEm | `time.Time` opcional | Somente no estado encerrado, terminal. |

Construtor: `NovoUsuario(id uuid.UUID, pessoa *Pessoa, email vo.Email) (*Usuario, error)`. Extrai empresa/pessoa de cadastro válido; inicia ATIVO e e-mail não confirmado. ATIVO significa conta não bloqueada, não autenticação concluída.

Métodos: `AlterarEmail(email vo.Email) error`, `ConfirmarEmail(email vo.Email, revisao uint64, quando time.Time) error`, `Bloquear() error`, `Desbloquear() error`, `Encerrar(quando time.Time) error`. Getters: `PessoaID()`, `Email()`, `RevisaoEmail()`, `ConfirmadoEm() (time.Time, bool)`, `Situacao()`, `EncerradoEm() (time.Time, bool)`.

Alterar a chave do e-mail invalida confirmação e incrementa revisão atomicamente. Mudar somente a apresentação mantendo a chave preserva confirmação. Confirmar exige e-mail/revisão atuais; isso evita aplicar confirmação atrasada após trocas, inclusive ida e volta ao mesmo endereço. O domínio não verifica OTP/Magic Link: o futuro caso de uso comprova o desafio antes de chamar o método. Bloqueado não recebe confirmação nem pode ser usado para acesso; pode ter e-mail corrigido, continuando bloqueado. Encerrado rejeita alterações/reativação. Bloquear ATIVO/BLOQUEADO e desbloquear ATIVO/BLOQUEADO são idempotentes. Encerrar novamente preserva o primeiro instante, sem recriar acesso.

E-mail confirmado, conta ATIVA e vínculo autorizado são condições distintas. Bloqueio não libera e-mail para outro usuário. Conta encerrada não é selecionada para login só porque possui o mesmo e-mail da nova conta. Login futuro requer CNPJ e seleção inequívoca do cadastro corrente.

## Funcionário

| Campo privado | Tipo | Regra |
|---|---|---|
| id, empresaID, pessoaID | `uuid.UUID` | Obrigatórios e imutáveis. |
| matricula | `vo.Matricula` | Obrigatória e imutável por cadastro, reservada após encerramento. |
| modalidade | `vo.ModalidadeColeta` | Obrigatória, atribuída inicialmente pela empresa. |
| setorID, usuarioID | `uuid.UUID` opcionais | Mesmo tenant; usuário da mesma Pessoa e não encerrado ao associar. |
| supervisorVinculoID | `uuid.UUID` opcional | Vínculo Supervisor específico do mesmo CNPJ, conforme os contratos da #5. Não é inferido pelo setor. |
| situacao | `SituacaoFuncionario` | `ATIVO` ou `ENCERRADO`. |
| encerradoEm | `time.Time` opcional | Preenchido no encerramento terminal. |

Construtor: `NovoFuncionario(id uuid.UUID, pessoa *Pessoa, matricula vo.Matricula, modalidade vo.ModalidadeColeta, hoje civil.Date) (*Funcionario, error)`. Exige Pessoa válida com nascimento válido em `hoje`. Inicia ATIVO, sem setor nem conta. Não gera matrícula nem concede perfil Colaborador.

Métodos: `AtribuirSetor(setor *Setor) error`, `RemoverSetor() error`, `AssociarUsuario(usuario *Usuario) error`, `Encerrar(quando time.Time) error`. Getters: `PessoaID()`, `Matricula()`, `Modalidade()`, `SetorID() (uuid.UUID, bool)`, `UsuarioID() (uuid.UUID, bool)`, `Situacao()`, `EncerradoEm() (time.Time, bool)`.

Complemento da #5: `AtribuirSupervisor(supervisor vo.ReferenciaSupervisor) error`, `RemoverSupervisor() error`, `SupervisorVinculoID() (uuid.UUID, bool)`. Inicialmente sem supervisão. O valor de referência contém os IDs de empresa, vínculo, usuário e Pessoa do Supervisor, sem afirmar validade atual do acesso. A equipe é derivada dessas referências diretas; existência, validade e autorização da atribuição/transferência seguem [Autorização](AUTORIZACAO_IDENTITY_SERVICE.md), com preservação de estado em erro e histórico no encerramento. Funcionário na #13 pode compilar sem a implementação de Vínculo/políticas da #14; a coordenação valida os registros autoritativos antes de atribuir.

Associar novamente o mesmo usuário é idempotente; associar outro a funcionário já associado retorna conflito, sem substituir a conta. Encerrado rejeita alterações e preserva referências; repetir encerramento preserva o primeiro instante. Não há reativação, troca de CNPJ, regeneração de matrícula ou mudança de modalidade nesta rodada. Migração de modalidade e substituição/desassociação de conta exigem fluxo próprio antes de expor métodos; essas pendências não impedem definir o cadastro inicial. CPF/nome/nascimento são acessados pela Pessoa, sem setters duplicados no Funcionário.

Encerrar funcionário não encerra usuário automaticamente: ele pode continuar gestor ou contador. Se for necessária nova conta para recontratação, a coordenação deve encerrar a anterior e revogar seu vínculo corrente, com autorização e histórico, antes de criar a nova. Não transferir o vínculo anterior automaticamente. Novos eventos referenciam novo funcionário/usuário; eventos antigos continuam referenciando os cadastros originais.

## Setor

Campos privados: `id`, `empresaID` (`uuid.UUID`, obrigatórios/imutáveis) e `nome` (`vo.Nome`, obrigatório/alterável). Construtor: `NovoSetor(id, empresaID uuid.UUID, nome vo.Nome) (*Setor, error)`. Método: `Renomear(nome vo.Nome) error`; getter `Nome()`.

Nomes iguais não significam o mesmo setor; somente ID define identidade. Não impor unicidade de nome não prevista nos requisitos. Exclusão, transferência em massa e histórico de lotação ficam fora desta rodada. Atribuir setor a funcionário não concede supervisão ao usuário.

## Vínculo com empresa

Campos privados: `id`, `empresaID`, `usuarioID` (`uuid.UUID`, obrigatórios/imutáveis), `autorizacao` (`vo.ConfiguracaoAutorizacao`, obrigatório), `situacao` (`SituacaoVinculo`, `ATIVO` ou `REVOGADO`) e `revogadoEm` (`time.Time` opcional). Configuração é valor imutável de um perfil e concessões, com alcance definido por perfil e recurso conforme a #5; não publicar um tipo vazio/`any` como solução provisória.

Contratos previstos para implementação na #14, após a #5: `NovoVinculo(id uuid.UUID, usuario *Usuario, autorizacao ConfiguracaoAutorizacao) (*Vinculo, error)`, `AlterarAutorizacao(autorizacao ConfiguracaoAutorizacao) error`, `Revogar(quando time.Time) error`. Getters: `UsuarioID()`, `Autorizacao()`, `Situacao()`, `RevogadoEm() (time.Time, bool)`.

Criação exige usuário válido e não encerrado; perfil não é inferido de CPF, empresa responsável ou matrícula. O vínculo inicia ATIVO, mas usuário bloqueado/não confirmado continua sem acesso. Alteração troca a configuração inteira atomicamente, inclusive ao mudar perfil, sem preservar concessões incompatíveis. Revogação é terminal; repetir preserva o instante original. Conceder novamente cria novo vínculo, condicionado à ausência de outro ativo. Verificar funcionário do Colaborador/Supervisor, subordinados diretos e concessões requer dados autoritativos e as regras da #5.

Não há método local que garanta autorizações entre tenants só recebendo IDs do cliente. A avaliação precisa da conta/vínculo atuais, ação, empresa e recurso/escopo. O método `Autorizar`, a interface e os quatro adaptadores estão definidos na #5, que proíbe autoaprovação e vincula acesso ao atestado à capacidade de tratar sua ocorrência. Essas decisões substituem as pendências iniciais, sem autorizar ações extras ainda indefinidas.

## Unicidade, coordenação e concorrência

| Garantia | Escopo e tratamento |
|---|---|
| Empresa | CNPJ único e ID imutável; não reutilizar um tenant para outro CNPJ. |
| Pessoa | CPF único no CNPJ, inclusive quando seus funcionários/usuários estão encerrados; reutilizar Pessoa na recontratação. Corrigir CPF para o de outra Pessoa conflita, sem fusão silenciosa. |
| Funcionário ativo | No máximo um por Pessoa/CPF no CNPJ. Histórico encerrado não impede recontratação; bloqueio de conta não encerra emprego. |
| Usuário corrente | E-mail pela `Chave()` e Pessoa únicos entre usuários não encerrados no CNPJ. E-mail pode se repetir entre CNPJs e entre contas encerradas. |
| Associação usuário/funcionário | Mesmo tenant e Pessoa; no máximo um funcionário ativo para o usuário. Nenhuma associação implícita por e-mail/CPF isolado. |
| Vínculo | No máximo um ativo por usuário/CNPJ; histórico revogado preservado. |
| Matrícula | Única no CNPJ entre todos os funcionários, inclusive encerrados. Pode repetir em outro CNPJ; coletor precisa conhecer o CNPJ. |

Essas garantias exigem consulta autoritativa e restrições/transações na futura persistência. Fazer `Exists` e depois gravar sem proteção de concorrência não basta. Dois cadastros concorrentes não podem produzir dois funcionários ativos, duas contas correntes ou matrículas iguais; conflito deve retornar resultado controlado. Alterar dados, associar referências e autorizar também exige controle contra atualização concorrente. A forma técnica fica para os respectivos casos de uso/adaptadores.

Encerramento da conta, revogação do vínculo e criação da nova conta não podem produzir reutilização ambígua de credenciais; até a coordenação concluir, o novo acesso não é concedido. Preservar as referências históricas. Projeções/cache futuros não podem transformar estado antigo em autorização atual; a estratégia de propagação/revogação continua pendente da integração.

## Biblioteca de matrícula

`internal/matriculagen` fornece `Gerar() (vo.Matricula, error)` ao consumidor, com seis escolhas uniformes no alfabeto de 36 caracteres. A implementação padrão usa aleatoriedade criptograficamente segura; leitura/fonte controlável fica em seam interna para testes determinísticos. Rejection sampling ou operação limitada uniforme deve evitar viés de módulo. Não usar contador, timestamp, CPF, nome, `math/rand` ou prefixo do tenant como entropia.

O espaço contém `36^6 = 2.176.782.336` valores por CNPJ; colisões são possíveis muito antes de esgotá-lo. O gerador faz uma geração e não promete unicidade. O cadastro futuro tenta persistir com restrição de unicidade e, somente em conflito de matrícula, gera novamente, até dez tentativas por operação. Após isso retorna erro controlado sem cadastro parcial; não repetir indefinidamente nem repetir outros erros de persistência. Falha da fonte retorna imediatamente. A reserva histórica nunca é liberada por encerramento. Matrícula não é segredo nem fator de autenticação.

Implementação da [issue #10](https://github.com/JoaoLucasDaSilvaOliveira/teu-ponto-digital/issues/10): `matriculagen.Gerar()` usa `crypto/rand.Reader`, com descarte dos bytes de 252 a 255 antes da redução ao alfabeto. A fonte é recebida pela função interna `gerar(io.Reader)`, sem configuração pública de fonte alternativa. Falhas retornam matrícula zero e erro com contexto e causa preservada. `value_object.NovaMatricula`, `String` e `Validar` representam valores independentemente do gerador; erros de formato usam `value_object.ErrValorInvalido`, inclusive para o valor zero. Persistência, novas tentativas em colisão e o construtor do funcionário permanecem nas etapas próprias.

## Erros e reconstrução

Erros locais de domínio serão distinguíveis por `errors.Is`, com contexto por `%w`, sem logar e retornar o mesmo erro. Mensagens em minúsculas, sem valores pessoais. Categorias mínimas: `ErrIDInvalido`, `ErrValorInvalido`, `ErrReferenciaInvalida`, `ErrTenantDivergente`, `ErrPessoaDivergente`, `ErrNascimentoObrigatorio`, `ErrEstadoIncompativel`, `ErrAssociacaoExistente`, `ErrConfirmacaoObsoleta` e `ErrRevisaoEsgotada`. Validação pode acrescentar tipo/campo ao erro sem expor o conteúdo sensível.

Conflitos persistentes, não encontrado, não autorizado, concorrência, cadastro incompleto, certificado indisponível e esgotamento de tentativas são resultados da coordenação, sem fingir que um setter verifica banco/assinatura. A biblioteca propaga a falha de aleatoriedade com contexto, preservando sua causa. Erros não usam `panic` para entrada inválida.

Construtores de criação não devem ser usados para reconstruir usuários confirmados, bloqueados ou encerrados, nem vínculos revogados: isso apagaria estado. A persistência futura deverá oferecer reconstrução validada por tipos de estado explícitos, internos ao módulo, reutilizando validadores e verificando coerência de campos opcionais/estados. Não expor setters públicos de ID, confirmação ou situação para facilitar ORM. A assinatura do mapeamento/reconstrução será definida junto ao adaptador, sem criar um repositório nesta issue.

## Cenários observáveis para as próximas implementações

1. Empresa com identificação válida pode iniciar incompleta; responsável ou A1 de outro CNPJ é rejeitado sem alteração. Troca de CNPJ não existe.
2. Pessoa externa sem nascimento pode ser responsável/ter acesso; criar funcionário exige nascimento. CPF/nome não são duplicados por perfil. Correção de CPF preserva IDs/matrícula e conflita com CPF já cadastrado no tenant.
3. Mesmo CPF/e-mail em A e B produz registros independentes; alteração de nome/e-mail/bloqueio em A não altera B. Referências cruzadas são rejeitadas.
4. Funcionário sem conta pode ser cadastrado; matrícula sozinha não autentica no coletivo. Associar conta da mesma empresa mas de outra Pessoa falha; conta encerrada falha.
5. Gestor também funcionário usa a mesma Pessoa/conta no CNPJ, sem cadastro por perfil. Conta bloqueada ou vínculo revogado não autoriza ações, mesmo com e-mail confirmado.
6. Alterar e-mail invalida confirmação; resposta antiga falha após nova alteração, inclusive retorno ao endereço anterior. Dados inválidos preservam o estado completo.
7. Recontratação preserva funcionário/conta anteriores, cria novos IDs/matrícula e não reutiliza permissões automaticamente. Histórico mantém IDs antigos; e-mail fica disponível só após encerrar conta anterior.
8. Matrícula minúscula válida normaliza; hífen/espaço/tamanho errado falha. Fonte controlável exercita extremos do alfabeto e falha de leitura. Conflitos controlados exercitam nova tentativa e limite, sem teste probabilístico frágil.
9. Datas inválidas/futuras e valores zero falham. A1 com intervalo inválido falha; referência cronologicamente vigente não equivale a certificado validado. Vencimento não descarta marcações e deixa exportação dependente de assinatura pendente.
10. Encerramento/revogação preserva referências e instante na repetição; métodos posteriores incompatíveis falham. Getters não permitem modificar estado por referência. Construtores e alterações usam as mesmas invariantes.

TDD nas issues de código: próximo cenário, teste falhando pelo comportamento ausente, implementação mínima e refatoração. Testes de domínio pela interface pública, sem HTTP/banco, e testes de integração pertinentes para concorrência/unicidade quando houver adaptadores. A #5 acrescentará os cenários da matriz de autorização.

## Rastreabilidade e pendências preservadas

| Fonte | Contrato especificado | Continuidade |
|---|---|---|
| RF-01/RF-04; ADRs 0004/0007 | Contas e dados por CNPJ; contexto não concede acesso | Seleção do CNPJ/login futuros, sem conta global. |
| RF-02/RF-06/RF-06.1 | Vínculo separado de Funcionário e quatro perfis | Configuração, política, acúmulo e escopos na #5; implementação na #14. |
| RF-03 | Empresa, Pessoa responsável e referência A1 sem segredos; incompletude explícita | Implementação na #11 agora depende de Pessoa; upload/custódia e signatários futuros. |
| RF-05 | Pessoa com nascimento e Funcionário separados; recontratação preserva histórico | Pessoa deve ser implementada antes dos consumidores das issues #11/#12/#13. |
| RF-05.1 | Seis caracteres, alfabeto, normalização, unicidade por CNPJ e reserva histórica | Biblioteca/valor na #10; unicidade e tentativas em coordenação futura. |
| RF-07/RF-08 | E-mail, confirmação e estados por conta/CNPJ | Desafio, sessões, tokens, recuperação e login completos posteriores. |
| RF-11.1; ADR 0005 | Modalidade inicial explícita e matrícula como identificador | Cadastro facial, mudança de modalidade e fallback ainda nas fontes próprias. |
| RF-19/RF-20 | Conta Contabilidade/DP externa possível, autorização por empresa, matrícula interna própria | Mapeamento da folha/relatórios permanece fora do módulo. |
| ADR backend 0001 | Campos privados, composição, construtores com validação, erros sem mutação parcial | #6 avalia dependências; #7 inicializa módulo; #15 verifica rodada. |

Assinatura de eventos auditáveis solicitada pelo usuário fica registrada como trabalho futuro: conteúdo canônico, ator, contexto, instante/evidência temporal, identificação do certificado, armazenamento da assinatura, proteção da chave e verificação ainda precisam ser especificados. Não colocar hash como garantia de decisão empresarial; certificado da empresa não identifica sozinho o usuário que operou. Isso não escolhe signatários de AFD/AEJ/comprovantes nem implementa assinatura/auditoria nesta issue.

Nascimento de externos continua opcional; limites de tamanho de transporte/persistência, exclusão/retenção, onboarding, segurança das mensagens e evidências offline exigem especificação nas respectivas áreas. A inicialização e a validade do certificado não são comprovadas por tipos/UUIDs. Não há bloqueio de marcação por expiração do A1 nesta decisão; a validação regulatória das capacidades de coleta permanece pendente.

## Verificação desta entrega

Verificações executadas: `git diff --check`, checagem com Python de referências locais e IDs RF nos nove documentos desta rodada, presença dos 14 construtores e conferência de `36^6`. Resultado: 52 referências locais existentes, 37 IDs RF válidos e nenhuma falha nas verificações de construtores ou no espaço de matrículas. Revisão documental confrontou estados/transições, preservação de histórico, dependências e referências de tenant com as decisões confirmadas.

Foi detectada uma referência preexistente no plano ao caminho `REQUISITOS_NÃO _FUNCIONAIS.md`, ausente no checkout por renomeação local para `REQUISITOS_NÃO_FUNCIONAIS.md`. O conteúdo renomeado foi comparado ao original versionado e é idêntico. A renomeação e a referência associada foram preservadas fora desta entrega para não incorporar alterações alheias; a referência continua válida no estado versionado que conserva o nome original.

Não há código Go alterado ou módulo/suíte executável criado nesta issue; testes Go, análise estática e formatação Go não foram executados. Revisão e homologação continuam pelo fluxo do repositório. A política de autorização da #5 e as áreas futuras explicitadas acima não são declaradas implementadas ou homologadas.
