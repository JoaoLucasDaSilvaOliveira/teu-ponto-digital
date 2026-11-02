# Validação estrutural de CPF e CNPJ

Pesquisa documental em 02/11/2026 para a issue #6. Escopo: formato, preservação de zeros e cálculo de dígitos verificadores (DV). Esta nota é uma pesquisa técnica; não altera requisitos aceitos, não seleciona biblioteca e não comprova existência cadastral, situação fiscal, titularidade ou conformidade jurídica. Nenhum CPF foi extraído de cadastro de pessoas.

## CPF

O manual e-Financeira da Receita Federal descreve `REGRA_VALIDA_CPF`: nove posições de base e duas de DV, calculadas por módulo 11. O exemplo didático publicado é `28001238938`. No algoritmo apresentado, a base recebe pesos `1,2,3,4,5,6,7,8,9`; o primeiro DV é o resto por 11, convertendo resto 10 em zero. Para o segundo DV, a base acrescida do primeiro recebe pesos `0,1,2,3,4,5,6,7,8,9`, com a mesma conversão. As somas do exemplo são 201 e 195, produzindo 3 e 8. [Manual e-Financeira, Anexo II, REGRA_VALIDA_CPF, p. 34](https://sped.rfb.gov.br/arquivo/download/7961).

A fórmula usual com pesos decrescentes `10..2` e `11..2`, usando zero para resto menor que 2 e `11-resto` nos demais casos, é uma formulação equivalente por aritmética modular; esta equivalência é uma dedução técnica, não uma transcrição do manual. O exemplo também consta da [versão 1.04, seção 3.2.1.35](https://sped.rfb.gov.br/estatico/15/D77411CA9890CFE6A86656BAB12B5FB59146AD/Versao%20104%20-%2020160525%20-%20e-Financeira%20-%20Manual%20de%20Preenchimento.pdf).

Limitação de acesso: os trechos acima foram encontrados no índice de busca dos próprios documentos oficiais; a abertura direta dos PDFs do Sped apresentou timeout nesta sessão. Não foi usada reprodução de terceiros como fonte.

O manual do ReceitanetBX, seção 6.1.1, define CPF como sequência de 11 dígitos, sem formatação, preservando zeros à esquerda. Para CNPJ numérico, define 14 dígitos com o mesmo cuidado. Esse manual de 2020 documenta o formato numérico e não deve fundamentar a exclusão de letras no CNPJ atual. [Manual oficial ReceitanetBX, p. 10](https://servicos.receita.fazenda.gov.br/publico/programas/ReceitanetBX/Servicos/manual_servico_receitanetbx.pdf).

Sequências repetidas precisam de verificação própria, pois só o checksum não basta. Há precedente oficial do Inep que rejeita todos os CPFs com um único algarismo repetido e também `00000000191`; essa lista é regra daquele sistema e não demonstra uma lista universal da Receita. [Manual de migração do Censo Superior 2012, p. 108](https://download.inep.gov.br/educacao_superior/censo_superior/migracao/2012/Manual_Migracao.pdf). A eventual rejeição adicional de `00000000191` no Teu Ponto exige decisão explícita, sem inferir que todo CPF com zeros iniciais é inválido.

## CNPJ numérico e alfanumérico

O FAQ da Receita confirma 14 posições, com `0–9` e letras maiúsculas `A–Z` nas 12 primeiras; as duas últimas são numéricas. Confirma quaisquer das 26 letras: não há fundamento nessa fonte para proibir `F`, `I`, `O`, `Q` ou `U`. Os formatos numérico e alfanumérico coexistem. [Perguntas e respostas oficiais, perguntas 1, 2, 6 e 7](https://www.gov.br/receitafederal/pt-br/centrais-de-conteudo/publicacoes/perguntas-e-respostas/cnpj/cnpj-alfanumerico.pdf).

O manual técnico define valor de cada caractere como `ASCII-48`: números mantêm seus valores e letras vão de `A=17` a `Z=42`. Para o primeiro DV, os pesos da esquerda para a direita são `5,4,3,2,9,8,7,6,5,4,3,2`. Para o segundo, acrescenta-se o primeiro DV e usam-se `6,5,4,3,2,9,8,7,6,5,4,3,2`. Em cada etapa, resto 0 ou 1 gera zero; demais restos geram `11-resto`. O exemplo `12.ABC.345/01DE-35` tem somas 459 e 424. Esse cálculo preserva a interpretação dos caracteres numéricos. [Manual oficial de cálculo, pp. 3–4](https://www.gov.br/receitafederal/pt-br/centrais-de-conteudo/publicacoes/documentos-tecnicos/cnpj/manual-dv-cnpj.pdf).

A estrutura canônica pode ser expressa por `[0-9]{11}` para CPF e `[A-Z0-9]{12}[0-9]{2}` para CNPJ, sempre ancorada ao texto inteiro; isso é uma tradução técnica das descrições, não validação completa. Aceitar máscaras ou converter minúsculas é política de entrada da aplicação. Remover caracteres arbitrários ou completar automaticamente entradas curtas não decorre das regras de checksum e pode ocultar erros de digitação.

## Massa de teste e origem

“DV correto” na tabela significa somente compatibilidade matemática. Mesmo valores construídos para testes não têm garantia de ausência em cadastros reais; não devem ser enviados a serviços de produção como identidade de alguém.

| Valor canônico | Origem | Uso |
|---|---|---|
| `28001238938` | Exemplo didático oficial e-Financeira acima | CPF: referência de DV `38` |
| `12ABC34501DE35` | Exemplo oficial do manual técnico acima | CNPJ com letras na raiz e na ordem; DV `35` |
| `12345678000195` | Exemplo ilustrativo do FAQ oficial, pergunta 17 | CNPJ numérico; DV `95` |
| `00000000000191` | Teste positivo de `src/typescript/test.ts`, ZIP oficial RFB abaixo | CNPJ numérico com vários zeros iniciais; são 14 posições |
| `ABCDEFGHIJKL80` | Teste positivo do mesmo arquivo oficial | Todas as 12 posições da base com letras |
| `01234567890` | Sintético: base escolhida `012345678`, DV calculado localmente | CPF: preservar um zero inicial |
| `00012345601` | Sintético: base escolhida `000123456`, DV calculado localmente | CPF: preservar vários zeros iniciais |
| `01234567000195` | Sintético: base escolhida `012345670001`, DV calculado localmente | CNPJ numérico com zero inicial |
| `00ABC34501DE56` | Sintético: base escolhida `00ABC34501DE`, DV calculado localmente | CNPJ alfanumérico com zeros iniciais |
| `FIOQU000000104` | Sintético: base escolhida `FIOQU0000001`, DV calculado localmente | Exercitar letras que não devem ser excluídas sem fonte |

Os DVs sintéticos foram calculados com um script temporário, sem biblioteca candidata. Não são fixtures publicadas pela Receita. Os documentos consultados não forneceram nesta pesquisa um exemplo didático de CPF com zero inicial e DV correto identificado como massa fictícia oficial.

A Receita também disponibiliza [Simulador Nacional de CNPJ](https://servicos.receitafederal.gov.br/servico/cnpj-alfa/simula), que identifica os números gerados como fictícios para testes, e [arquivos de referência para cálculo](https://www.gov.br/receitafederal/pt-br/centrais-de-conteudo/publicacoes/documentos-tecnicos/cnpj/codigos-cnpj.zip/view). Nesta pesquisa não foi congelada exportação do simulador. O [ZIP oficial](https://www.gov.br/receitafederal/pt-br/centrais-de-conteudo/publicacoes/documentos-tecnicos/cnpj/codigos-cnpj.zip) foi baixado e inspecionado, sem execução: `src/typescript/test.ts` também espera rejeição de vazio, CNPJ zerado com e sem máscara, letras nas posições de DV, caracteres arbitrários e minúsculas. São comportamentos do código de referência; a aplicação pode definir conversão explícita de caixa antes da validação. Não se atribuem ao pacote os valores marcados como sintéticos na tabela.

## Limites e verificações seguintes

Rastreabilidade do pacote inspecionado: SHA-256 `c86335cb9251b65aee2b3582df1b77bfc7728b1c5d0f22999a38b964ca6d7136`, baixado em 02/11/2026. O ZIP não foi incluído no repositório.

A seleção de biblioteca deve conferir separadamente tamanho e alfabeto, checksum, sequências inválidas e normalização; uma função que calcula DV não necessariamente valida entrada completa. Casos relevantes incluem vazio, somente máscara, caracteres extras, Unicode fora de ASCII, letras no DV, tamanho incorreto, alteração de cada DV e zeros iniciais preservados. Regras de situação cadastral e titularidade pertencem a outra verificação e não são atendidas por este cálculo.
