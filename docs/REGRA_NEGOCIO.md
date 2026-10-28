# Visão Geral - Projeto Teu Ponto

## 1. O que é o sistema
O "Teu Ponto" é um sistema completo de controle de jornada de trabalho e marcação de ponto digital. O produto é um REP-P próprio e engloba processamento das regras trabalhistas para o fechamento da folha no RH/DP.

## 2. Arquitetura e Tecnologias Base
- **Backend:** Go (responsável pelas regras de negócio, cálculo de horas e API).
- **Frontend:** Flutter nativo como foco; PWA opcional, condicionada ao atendimento dos requisitos obrigatórios. Ver [ADR 0003](adr/0003-flutter-nativo-e-pwa-condicional.md).
- **Isolamento:** Cada CNPJ é um tenant com dados, permissões e configurações próprios. Ver [ADR 0004](adr/0004-isolamento-por-cnpj.md).
## 3. Regras de Negócio Globais (Core Domain)
- **Legislação Base:** A CLT brasileira é a base normativa inicial. Não há CCT/ACT de referência definida nesta etapa; políticas dependentes de instrumento coletivo permanecem pendentes de especificação. Essa delimitação de escopo não afasta instrumentos aplicáveis aos futuros empregadores.
- **Tolerância de Ponto:** Deve considerar as regras de tolerância de atrasos e horas extras (ex: 5 a 10 minutos diários).
- **Identificação na coleta:** A empresa define a modalidade por funcionário: autenticação do dispositivo no aparelho próprio; identificador, comparação facial e prova de vida passiva no aparelho coletivo. Fallbacks ainda precisam ser especificados; ver [ADR 0005](adr/0005-modalidade-coleta-por-funcionario.md).
- **Saídas:** O sistema atua como a fonte da verdade para a geração do espelho de ponto mensal e relatórios de absenteísmo.
- **Relatórios:** O sistema gera relatórios das estradas e saídas. Esses relatórios sempre devem ser assinados usando o certificado digital da empresa.
## 4. Atores Principais
- **Colaborador:** Realiza as marcações de entrada e saída e visualiza seu próprio espelho de ponto.
- **Gestor:** Aprova justificativas e monitora a equipe.
- **RH / Departamento Pessoal:** Configura regras de jornada, gerencia exceções e exporta os relatórios mensais para fechamento da folha.
