# Layouts de Exportação

## Layout AFD — Portaria MTE 1.510/2009

Arquivo Fonte de Dados para conformidade REP-P.

### Tipo 1 — Cabeçalho

| Campo          | Tam | Conteúdo                                        |
|----------------|-----|-------------------------------------------------|
| NSR            | 9   | Sequencial numérico, zeros à esquerda           |
| Tipo           | 1   | "1"                                             |
| CNPJ           | 14  | Somente dígitos, zeros à esquerda               |
| CEI            | 12  | CEI ou zeros                                    |
| Razão Social   | 150 | Padded espaços à direita                        |
| Data início    | 8   | DDMMAAAA                                        |
| Data fim       | 8   | DDMMAAAA                                        |
| Data geração   | 8   | DDMMAAAA                                        |
| Hora geração   | 6   | HHMMSS                                          |
| Versão layout  | 9   | "001" padded zeros à esquerda                   |
| CRLF           | —   | `\r\n`                                          |

### Tipo 3 — Marcação

| Campo    | Tam | Conteúdo                                                         |
|----------|-----|------------------------------------------------------------------|
| NSR      | 9   | Sequencial                                                       |
| Tipo     | 1   | "3"                                                              |
| Data/Hora| 12  | DDMMAAAAHHmm                                                     |
| PIS/NIS  | 12  | PIS do funcionário, zeros à esquerda; fallback: CPF; fallback: zeros |
| CRLF     | —   | `\r\n`                                                           |

### Tipo 9 — Trailer

| Campo         | Tam | Conteúdo                                      |
|---------------|-----|-----------------------------------------------|
| NSR           | 9   | Último sequencial                             |
| Tipo          | 1   | "9"                                           |
| Total tipo 3  | 9   | Contagem de linhas tipo 3, zeros à esquerda   |
| CRLF          | —   | `\r\n`                                        |

---

# Layout de Exportação — Domínio Sistemas

Arquivo texto de largura fixa para importação no sistema **Domínio Contábil / Folha**.

## Linha de Cabeçalho

| Campo         | Posição | Tamanho | Descrição                        |
|---------------|---------|---------|----------------------------------|
| CNPJ          | 1–14    | 14      | CNPJ da empresa (somente dígitos)|
| RAZAO_SOCIAL  | 16–75   | 60      | Razão social (right-padded)      |
| REFERENCIA    | 77–83   | 7       | Período no formato `MMAAAA`      |

## Linhas de Detalhe (uma por rubrica/colaborador)

| Campo           | Posição | Tamanho | Descrição                                  |
|-----------------|---------|---------|---------------------------------------------|
| MATRICULA       | 1–6     | 6       | Matrícula do colaborador, left-padded com `0`|
| CODIGO_RUBRICA  | 7–11    | 5       | Código da rubrica configurada pelo tenant   |
| VALOR           | 12–21   | 10      | Valor em minutos (horas) ou dias, left-padded com `0`|
| TIPO            | 22      | 1       | `H` = horas, `D` = dias, `V` = valor       |
| REFERENCIA      | 23–29   | 7       | Período no formato `MMAAAA`                 |

## Tipos de Rubrica Suportados

| Tipo               | Descrição                    | Unidade |
|--------------------|------------------------------|---------|
| `HORA_EXTRA_50`    | Hora extra 50%               | H       |
| `HORA_EXTRA_100`   | Hora extra 100%              | H       |
| `ATRASO`           | Atrasos                      | H       |
| `FALTA`            | Faltas                       | D       |
| `ADICIONAL_NOTURNO`| Adicional noturno            | H       |
| `DSR`              | Descanso semanal remunerado  | H       |

## Exemplo

```
12345678901234 Empresa Exemplo Ltda                                  032026
000123100  0000000120H032026
000123200  0000000001D032026
000456100  0000000060H032026
```

Neste exemplo:
- Colaborador matrícula `000123`, rubrica `100` (HE50%), 120 min (2h), referência março/2026
- Colaborador matrícula `000123`, rubrica `200` (FALTA), 1 dia, referência março/2026
- Colaborador matrícula `000456`, rubrica `100` (HE50%), 60 min (1h), referência março/2026
