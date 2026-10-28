# Marcação esquecida e continuidade da jornada

Status: aceito provisoriamente pelo usuário em 26/10/2026. Este fluxo orienta a especificação atual e pode ser revisto com a validação de uso. Ver RF-12.1 a RF-13.1.

## Problema

Uma saída ausente no dia anterior não permite concluir se houve esquecimento ou trabalho atravessando a noite. Alternar entrada/saída apenas pela contagem de marcações pode classificar a entrada seguinte como saída anterior e contaminar a sequência. Jornada noturna legítima também impede usar meia-noite como corte obrigatório.

## Fluxo aceito provisoriamente

Ao detectar ambiguidade, apresentar a pendência e permitir declarar a intenção atual: registrar entrada, intervalo, retorno ou saída. A sugestão do sistema é corrigível e não impede registrar o horário atual.

- Esquecimento: registrar entrada agora e, opcionalmente no mesmo fluxo, solicitar inclusão da saída anterior, informando data, horário declarado e justificativa. Oferecer concluir o pedido depois sem impedir a marcação atual.
- Trabalho contínuo: registrar saída agora e sinalizar duração incomum para conferência, sem cortar horas por ultrapassar limite previsto.
- Dúvida ou intenção indisponível: preservar a marcação atual e sinalizar classificação pendente, evitando pareamento definitivo por simples alternância.

A informação histórica declarada não é uma batida original retroativa. O pedido segue o fluxo de aprovação RF-12/RF-13; a marcação atual independe dessa aprovação. A rejeição exige motivo e mantém visíveis o pedido, a divergência e as marcações posteriores, sem transformá-las automaticamente em outro tipo. Permitir reapresentação ou revisão conforme política a definir.

## Apuração e concorrência

Manter ocorrência de marcação incompleta ou associação ambígua; os segmentos não afetados podem ser apurados. Não liquidar automaticamente o trecho incerto como falta, desconto ou horas extras. Exibir qualquer prévia como provisória até resolução.

Uma marcação offline recebida depois pode resolver a suposta ausência. Antes de aprovar o ajuste, conferir novamente os registros do período; se a marcação original já chegou, sinalizar conflito ou solicitação sem objeto, evitando dupla contagem. Duas aprovações simultâneas para a mesma ausência também exigem conferência da versão atual do tratamento.

Lembretes de saída podem reduzir esquecimentos, mas não geram batidas nem preenchem horário contratual automaticamente.

## Evidência de mercado

- [Pontotel: fluxo de marcação pelo aplicativo](https://www.pontotel.com.br/como-funciona-a-pontotel/): descreve escolha do tipo de marcação após identificação. Isso mostra que o tipo não precisa depender exclusivamente da alternância automática.
- [Pontotel: ajustes na folha](https://suporte.pontotel.com.br/pt-br/article/como-realizar-ajustes-e-correcoes-de-marcacoes-na-folha-de-ponto-1oboo1w/): documenta inclusão, ordenação e tratamento de registros, inclusive sinalização de dia seguinte; registros reais desconsiderados ficam tachados.
- [Pontomais/VR: ajustes, justificativas e abonos](https://materiais.vr.com.br/central-de-ajuda/gestao-de-controle-de-ponto-como-fazer-ajustes-de-horarios-justificativas-no-ponto-e-abono/): documenta solicitação pelo colaborador via app/web e ajuste pelo gestor.

Essas fontes sustentam componentes da proposta. Não foi encontrado nelas um fluxo idêntico de resolução imediata da saída esquecida sem impedir a entrada seguinte; essa combinação é recomendação para o Teu Ponto.

## Referência regulatória

A [Portaria 671, art. 74](https://www.gov.br/trabalho-e-emprego/pt-br/assuntos/legislacao/portarias-1/portarias-vigentes-3/PDFPortarian671de8denovembrode2021compilada07.01.2026.pdf) exige fidelidade das marcações e veda restrições de horário, marcação automática predeterminada e autorização prévia de sobrejornada. O tratamento deve preservar o registro original, conforme [orientações do MTE](https://www.gov.br/trabalho-e-emprego/pt-br/assuntos/inspecao-do-trabalho/fiscalizacao-do-trabalho/Perguntas%20e%20Respostas%20REP). Isso orienta a separação entre marcação atual, declaração histórica e ajuste aprovado.
