# PROTOCOLO TÉCNICO DE VALIDAÇÃO DE SOFTWARE METROLÓGICO
**Identificação:** VAL-SW-001
**Assunto:** Homologação do Algoritmo de Validação Autônoma de Planilhas de Cálculo
**Elaboração:** Auditoria Técnica ISO/IEC 17025

---

## 1. Objetivo e Escopo
Estabelecer os requisitos arquiteturais, lógicos e matemáticos implementados no software validador de planilhas eletrônicas (`.xlsx`). O escopo engloba a garantia da precisão de ponto flutuante, a integridade da árvore de dependências computacionais (AST) e a mitigação de erros de propagação em cálculos analíticos laboratoriais.

## 2. Referências Normativas
* **ABNT NBR ISO/IEC 17025:2017** - Seção 7.11 (Controle de dados e gestão da informação).
* **ABNT NBR ISO/IEC 17025:2017** - Seção 7.7 (Garantia da validade dos resultados).

## 3. Arquitetura de Processamento
O algoritmo opera através de concorrência computacional assíncrona, isolando variáveis por aba para garantir a segurança da memória no processamento de matrizes densas. A extração de dados obedece à arquitetura de **Compilação em 3 Etapas (*Naked Sheet & AST Pruning*)**:

1. **Nudez Total (Naked Sheet):** O sistema varre a matriz e remove temporariamente a formatação visual (estilos, máscaras e casas decimais truncadas) de todas as células, alterando-as para a base estrita do formato "Geral". Isso força o motor de cálculo a trabalhar diretamente com os dados numéricos brutos em precisão máxima de 64-bits salvos no código-fonte do documento original.
2. **Poda de Condicionais (Descascamento de IFs):** Máscaras visuais estruturadas como condicionais de bloqueio (ex: `IF(A1="","", Matemática)`) são interceptadas por um analisador sintático recursivo. O invólucro condicional é removido e apenas a expressão matemática pura é injetada na memória RAM, prevenindo o colapso do motor de leitura ao lidar com dependências referencialmente vazias ou recursivas.
3. **Restauração de Segurança:** Imediatamente após a consolidação do cálculo estrito, o algoritmo restaura as matrizes de estilo e as fórmulas literais originais em suas respectivas coordenadas temporais, preservando a integridade passiva do arquivo primário para auditorias subsequentes.

## 4. Critérios de Aceitação e Análise Metrológica
A comparação entre o valor consolidado estático ("Valor Salvo") e a simulação de recálculo em vácuo ("Recalculado") rejeita tolerâncias absolutas engessadas, adotando o princípio de **Tolerância Dinâmica** para conter o viés em sequências com Propagação de Erro.

O critério de conformidade da célula obedece à seguinte expressão de limite dinâmico avaliada de forma independente:
$Limite_{Dinamico} = \max(Tol_{Absoluta}, Tol_{Relativa} \times \max(|V_{Orig}|, |V_{Calc}|))$

**Parâmetros Adotados pela Rotina Matemática:**
*   **Tolerância Absoluta ($Tol_{Absoluta}$):** `0.0005` (Absorve flutuações e ruídos microscópicos de ponto flutuante próximos à nulidade).
*   **Tolerância Relativa ($Tol_{Relativa}$):** `0.001` (Garante um teto de desvio admissível de proporção geométrica de `0.1%`, acompanhando o crescimento ou colapso logarítmico dos resultados provenientes de constantes multiplicativas).

Apenas resultados cuja diferença absoluta obedeçam à restrição máxima estipulada pelo teto dinâmico recebem a homologação sistemática de exatidão.

## 5. Tratamento de Exceções e Interrupções
Anomalias oriundas de arquitetura de planilha ou corrupção de variáveis são mitigadas e rastreadas pelas seguintes regras condicionais:
*   **Erros Nativos e Incompatibilidade Sintática:** Avisos provenientes da estrutura computacional das matrizes (`#DIV/0!`, `#VALUE!`, `#N/A`, `#NAME?`) ou limitações estritas de complexidade do motor subjacente (`formula not supported`) não determinam falha imediata na validação das variáveis estáticas, sendo realocados sob a classificação moderada de **"REVISÃO MANUAL"**, induzindo intervenção humana.
*   **Equivalência de Vazio Metrológico:** Cadeias representacionais literais e ocultas, tais como `""`, `"0"`, `"0.0000"`, `"-"` (Traço Contábil) e `"-0.000001"`, são reduzidas, higienizadas e interpretadas de maneira coesa como zero absoluto, garantindo a aprovação se espelhadas entre o estado nativo e a simulação matemática.

## 6. Emissão de Relatório e Governança
A saída do algoritmo culmina na transcrição dos eventos validados em documentação imutável com layout otimizado, contemplando as seguintes regras de auditoria:
*   Registro compulsório do conjunto de metadados de rastreabilidade inseridos em tempo de operação (Identificador, Código, Revisão).
*   Injeção de assinatura de tempo global (Timestamp) em sincronização com a compilação do relatório em formato `PDF`.
*   Priorização de inconformidades: Os quadros matriciais contendo reprovações categóricas ("FALHA") ou avisos de intervenção estrutural ("REVISÃO MANUAL") antecedem de forma compulsória a transcrição das avaliações confirmadas, otimizando o isolamento imediato da não conformidade perante investigações corretivas e auditorias setoriais.