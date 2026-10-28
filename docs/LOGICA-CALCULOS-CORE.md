# Este é um arquivo contendo a lógica do clock-service: O micro serviço reponsável por batidas de ponto e banco de horas.

## Objective
Garantir que a lógica da aplicação obedeça as leis, em especial a Consolidação das Leis do Trabalho (CLT) brasileira.

## Entidades que compõe um serviço de ponto
1. Entradas: São espaços do tempo (timestampz) onde o empregado inicia as atividades laborais diária.
1. Saidas: São espaços do tempo (timestampz) onde o empregado finaliza as atividade laborais diária.
1. Intervalo: Trata-se de um conjunto de 1 saída seguida de 1 entrada entre a entrada e saída diária.
1. Horário: Menor unidade de uma jornada, possui um conjunto específico de entradas, intervalo e saídas.
1. Jornada: Conjunto de horários ao longo do período de 7 dias. Jornadas devem possuir horários definidos de segunda à domingo.
1. Apuração: Entidade fictícia que representa um cálculo baseado na quantidade de horas trabalhadas em um dia ou um período específico.
1. Saldo de horas: Entidade abstrata resultado de uma apuração.

## Regras gerais
- Um registro ponto sempre será composto de um par de entrada e saída, exemplo: Entrada: 08:00, Saída: 12:00. Apuração: 4 horas trabalhadas.
- Nunca deve ser realizada uma simples comparação entre, por exemplo, entrada esperada (cadastro) e entrada obtida e avaliar a diferença. 
- Podem ser trabalhadas no máximo 2 horas extras por dia, no máximo 10 horas trabalhadas no dia. O que ultrapassar 2 horas extra ou 10 horas trabalhadas no dia, verifique [Regras apuração direta](<LOGIC#Regras apuração direta>) para ver caso de uso.
O cálculo é mais complexo que isso:
  - Como definido, um registro sempre será composto por um par de entrada e saída. Normalmente um registro ponto possui 2 pares de entradas e saídas.
  - A entidade Apuração deve ser capaz de avaliar o período total trabalhado, avaliando cada par de entrada e saída, e produzir a quantidade de horas trabalhadas.
  - As horas trabalhadas devem ser subtraídas da quantidade de horas esperadas (Carga Horária) no dia apurado. Obs: A quantidade de horas esperadas (Carga Horária) deve ser extraída do registro de horário do empregado. 
  - Exemplo: Entrada: 08:00, Saída 12:00, Entrada: 14:00, Saída: 18:00. Nesse exemplo, o primeiro par possui um total de 4 horas trabalhadas e o segundo par um total de 4 horas trabalhadas, totalizando 8 horas no dia. Se, nesse exemplo, a quantidade de horas esperada era 07:20, há um saldo de 40 minutos positivo.
  - **IMPORTANTE**: Por mais que a apuração não realize uma simples comparação entre, por exemplo, entrada esperada (cadastro) e entrada obtida, existe um caso de excessão: Caso o saldo apurado seja diferente de 0 é preciso avaliar algumas situações:
    - Lei brasileira define: 
      - "A Tolerância de Ponto refere-se ao período máximo de tempo em que o registro da jornada de trabalho  na entrada, na saída ou nos intervalos – pode variar em relação ao horário contratual, sem que essa diferença seja considerada atraso, hora extra ou falta.
      Não serão descontadas nem computadas como jornada extraordinária as variações de horário no registro de ponto que não excederem o limite de cinco minutos na entrada ou na saída.
      No entanto, a lei impõe um limite máximo para o total dessas variações diárias: o somatório dessas variações não pode ultrapassar dez minutos diários."
      <a id="regras-apuracao-diaria"></a>
      - ##### Exemplos práticos da regra:
        - Entrada com Atraso: Se o colaborador marca o ponto de entrada às 08h05min (horário contratual 08h00min), a diferença de 5 minutos não é considerada atraso e não pode ser descontada.
        - Saída Antecipada: Se o colaborador marca o ponto de saída às 17h55min (horário contratual 18h00min), a antecipação de 5 minutos não é considerada falta de horário e o tempo não pode ser descontado.
        - Variações Múltiplas: Se o colaborador marca 3 minutos antes na entrada e 4 minutos depois na saída, o total de variação é de 7 minutos. Como não ultrapassa o limite diário de 10 minutos, essa variação não é computada.
        - Excesso na Variação: Se o colaborador marca 6 minutos depois na entrada, a totalidade do tempo de variação será contabilizada. Ou seja, os 6 minutos serão considerados atraso ou, dependendo da situação, hora extra, pois o limite de 5 minutos para um único registro foi excedido.
      É fundamental que o RH e a gestão de pessoal estejam atentos a essa regra: a tolerância é de até 5 minutos por marcação, respeitando o limite global de 10 minutos por dia.

## Regras legais (CLT)
1. As empresas podem optar, individual ou coletivamente, por diferentes formas de apuração dos períodos trabalhados:
- Banco de horas/Compensação de horas
  - Apuração periódica.
  - Cálculo baseia-se em saldo de horas do período
    - Saldo de horas positivo chama-se: Crédito.
    - Saldo de horas negativo chama-se: Débito.
  - Apuração realizada por abatimentos
    - Exemplo: Se no dia 01 houverem 30 minutos positivos e no dia 02 houverem 30 minutos negativos, a apuração deve apontar um saldo igual a 0.
-  Apuração direta/pagamento de horas extras

## Regras apuração diretapor dia
- Saldo deve ser apurado diariamente.
- Cálculo baseia-se da diferença entre o total de horas trabalhadas no dia e a quantidade de horas esperadas (Carga Horária) para o mesmo dia.
- Apuração realizada por dia
  - Exemplo: Se no dia 01 houverem 30 minutos positivos, deve ser somado em horas extras. Estudo de caso mais complexo em [Exemplos práticos da regra](<LOGIC#Regras gerais>).
- Resultados da apuração direta:
  - Horas trabalhadas: Somatório bruto do intervalo entre pares de entrada e saída.
  - Horas extras: Devem ser aplicadas quando o a diferença entre as horas trabalhadas e as horas esperadas (Carga Horária) é positivo. OBS.: É considerada hora extra toda a apuração onde o saldo é positivo e  que não esteja abatida pela tolerância legal.
    - Horas extras 50%: Aplicada em dias trabalhados onde o saldo é positivo. **IMPORTANTE:** Horas extras 50% devem ser aplicadas até a soma de 2 horas diárias ou 10 horas de jornada. O sistema não deve impedir saldos positivos superiores a 2 horas diárias ou 10 horas de jornada e sim remunerar o excedente em 100%.
      - Exemplo: Jornada de 7:20 por dia 
        - Funcionário trabalhou 9:20 -> 2 horas extras 50%
        - Funcionário trabalhou 9:50 -> 2 horas extras 50% e 30 minutos de horas extras
        - Funcionário trabalhou 10:30 -> 2 horas extras 50% e 1 hora e 30 minutos extra 100%.
      - Exemplo: Jornada de 8:48 por dia
        - Funcionário trabalhou 9:20 -> 32 minutos extras 50%
        - Funcionário trabalhou 9:50 -> 1 horas e 2 minutos extras 50%
        - Funcionário trabalhou 10:30 -> 1 horas e 12 minutos extras 50% e 30 minutos extra 100%. (limite de 10h diárias).
      - **IMPORTANTE:** Para a escala 12x36 ou 11x36 ou afins a regra de 10 horas de jornada é eliminada, permanecendo todas as regras anteriores. 
      - Dias da semana para aplicar horas extras 50% por padrão:
        - Segunda-feira.
        - Terça-feira.
        - Quarta-feira.
        - Quinta-feira.
        - Sexta-feira.
        - Sábado.
      - Dias da semana para aplicar horas extras 100% por padrão:
        - Domingos (somente o excedente da diferença entre as horas trabalhadas e a carga horária).
        - Feriados obrigatórios (qualquer hora trabalhada).
      - Observação: a porcentagem sobre a horas (50% ou 100%) pode variar de acordo com leis específicas, como Acordos Coletivos de Trabalho (CCT) e ou Dissídios Coletivos. O sistema deve possuir o default como 50% e 100% mas permitir a adequação a outras porcentagens.   
  - Horas falta: Devem ser aplicadas quando o a diferença entre as horas trabalhadas e as horas esperadas (Carga Horária) é negativo. OBS.: É considerada hora falta toda a apuração onde o saldo é negativo e que não esteja abatida pela tolerância legal.
    - Observação: Deve ser considerado como horas falta o dia que, além do saldo negativo, possuir um ou mais registros. Dias que possuem carga horária e não possuem nenhum registro serão considerados como **Dia de falta**.
  - Dia de falta: Dias que possuem carga horária e não possuem nenhum registro.
  - Faltas justificadas: Dias que possuem nenhum ou vários registros (completos(par de dois) ou não) onde há alguma justificativa. Justificativas permitidas: Atestado, Aviso prévio dispensa, Liberação, Férias.
  - Feriado: 
    - Quando não há nenhum registro no dia, porém o dia é um feriado obrigatório. Saldo igual a 0.
    - Qualquer hora trabalhada em feriados obrigatórios é 100% (ou a porcentagem definida).
  - Horas noturnas: Período trabalhado entre 22h e 05h.

## Regras compensação de horas (Banco de horas)
- Saldo deve ser apurado de ciclicamente.
- É considerado um ciclo o período de 6 meses ou 1 ano. Existem CCTs que mudam essa data, o sistema deve permitir alteração de período do início/encerramento do ciclo.
- Cálculo baseia-se da diferença entre o total de horas trabalhadas no dia e a quantidade de horas esperadas (Carga Horária) para o mesmo dia.
- Durante o período do ciclo, os Créditos e os Débitos se compensam. 
- Se, ao final do ciclo, houver saldo positivo, este deve ser destacado e pago como horas extras 50% (ou porcentagem customizada).
- Deve ser feito um documento chamado "Acordo Banco de Horas" onde, resumidamente, um funcionário, um departamento ou todos os empregados da empresa concordam e permitem a implementação do sistema de Compensação de Horas denominado Banco de Horas. Nesse documento é definido o período de ciclo, a porcentagem incidente sobre as horas extras, caso seja necessário o pagamento, além de outras informações. 
- Demonstrativo mensal e liquidação por ciclo, preservando o histórico
  - Exemplo: Se no dia 01 houverem 30 minutos positivos e no dia 02 houverem 30 minutos negativos, a apuração deve apontar um saldo igual a 0.
  - Ao final do mês, créditos e débitos devem ser carregados para o mês seguinte dentro do ciclo válido. Saldo negativo mensal não se transforma automaticamente em horas falta ou desconto.
  - No encerramento do ciclo ou na rescisão, saldo negativo sem regra válida de liquidação permanece pendente de tratamento. Não descontar, abonar ou transferir para novo ciclo automaticamente. Preservar origem, período e histórico dos lançamentos.
  - Exemplo: Se, ao final do ciclo, o valor salvo (que foi carregado entre os meses até o fim do ciclo) for positivo, esse salvo deve ser tranformado em **Horas extras 50% (mínimo, ou de acordo com a CCT)**.
- Resultados da compensação de horas:
  - Horas trabalhadas: Somatório bruto do intervalo entre pares de entrada e saída.
  - Crédito: Deve ser aplicado quando o a diferença entre as horas trabalhadas e as horas esperadas (Carga Horária) é positivo. OBS.: É considerada crédito toda a apuração onde o saldo é positivo e que não esteja abatida pela tolerância legal.
    - Exemplo: Jornada de 7:20 por dia 
      - Funcionário trabalhou 9:20 -> 2 horas de crédito.
      - Funcionário trabalhou 9:50 -> 2 horas e 30 minutos de crédito.
      - Funcionário trabalhou 10:30 -> 3 horas e 30 minutos de crédito.
    - Exemplo: Jornada de 8:48 por dia
      - Funcionário trabalhou 9:20 -> 32 minutos de crédito.
      - Funcionário trabalhou 9:50 -> 1 horas e 2 minutos de crédito.
      - Funcionário trabalhou 10:30 -> 1 horas e 42 minutos de crédito.
    - **IMPORTANTE:** Para a escala 12x36 ou 11x36 ou afins a regra de 10 horas de jornada é eliminada, permanecendo todas as regras anteriores. 
    - Dias da semana para aplicar crédito por padrão:
      - Segunda-feira.
      - Terça-feira.
      - Quarta-feira.
      - Quinta-feira.
      - Sexta-feira.
      - Sábado.
      - Domingos (somente o excedente da diferença entre as horas trabalhadas e a carga horária).
      - Feriados obrigatórios (qualquer hora trabalhada).
      - Observação: a porcentagem sobre a horas (50% ou 100%) pagas ao final do ciclo, se houver, podem variar de acordo com leis específicas, como Acordos Coletivos de Trabalho (CCT) e ou Dissídios Coletivos. O sistema deve possuir o default como 50% e 100% mas permitir a adequação a outras porcentagens.   
  - Débito: Devem ser aplicadas quando a diferença entre as horas trabalhadas e as horas esperadas (Carga Horária) é negativo. OBS.: É considerada débito toda a apuração onde o saldo é negativo e que não esteja abatida pela tolerância legal.
  - Saldo negativo: Débitos superiores aos créditos permanecem no banco durante o ciclo válido. A competência mensal não autoriza conversão automática em horas falta.
  - Faltas justificadas: Dias que possuem nenhum ou vários registros (completos(par de dois) ou não) onde há alguma justificativa. Justificativas permitidas: Atestado, Aviso prévio dispensa, Liberação, Férias, Compensação Banco de horas.
  - Feriado: 
    - Quando não há nenhum registro no dia, porém o dia é um feriado obrigatório. Saldo igual a 0.
    - Qualquer hora trabalhada em feriados obrigatórios é crédito.
  - - Horas noturnas: Período trabalhado entre 22h e 05h.