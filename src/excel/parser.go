package excel

import "strings"

// descascarIF agora diferencia "Proteção Visual" de "Lógica de Negócio"
func descascarIF(formula string) string {
	f := strings.TrimSpace(formula)
	f = strings.TrimPrefix(f, "=")
	f = strings.ReplaceAll(f, "$", "")

	if !strings.HasPrefix(strings.ToUpper(f), "IF(") {
		return formula
	}

	args := capturarArgumentosIF(f)

	if len(args) >= 2 {
		valTrue := strings.TrimSpace(args[1])
		valFalse := ""
		if len(args) >= 3 {
			valFalse = strings.TrimSpace(args[2])
		}

		// Definimos o que é considerado "apenas visual"
		isTruePlaceholder := (valTrue == `""` || valTrue == `"-"` || valTrue == `" "` || valTrue == "0")
		isFalsePlaceholder := (valFalse == `""` || valFalse == `"-"` || valFalse == `" "` || valFalse == "0" || valFalse == "")

		// Caso 1: IF(condicao; ""; FORMULA_REAL) -> Retorna FORMULA_REAL
		if isTruePlaceholder && !isFalsePlaceholder {
			return descascarIF(valFalse)
		}

		// Caso 2: IF(condicao; FORMULA_REAL; "") -> Retorna FORMULA_REAL
		if isFalsePlaceholder && !isTruePlaceholder {
			return descascarIF(valTrue)
		}

		// Caso 3: LÓGICA DE NEGÓCIO REAL (ex: IF(A1>10; B1; C1))
		// Se ambos os lados são "úteis", NÃO podemos descascar.
		// Retornamos a fórmula original para que o excelize tente calcular o IF completo.
		return formula 
	}

	return formula
}

// capturarArgumentosIF isola os argumentos tratando parênteses aninhados
func capturarArgumentosIF(f string) []string {
	parenCount := 0
	inQuotes := false
	var args []string
	var currentArg strings.Builder

	// Começa após o "IF("
	for i := 3; i < len(f); i++ {
		ch := f[i]
		if ch == '"' {
			inQuotes = !inQuotes
		}
		if !inQuotes {
			if ch == '(' {
				parenCount++
			} else if ch == ')' {
				if parenCount == 0 {
					args = append(args, strings.TrimSpace(currentArg.String()))
					break
				}
				parenCount--
			} else if (ch == ',' || ch == ';') && parenCount == 0 {
				args = append(args, strings.TrimSpace(currentArg.String()))
				currentArg.Reset()
				continue
			}
		}
		currentArg.WriteByte(ch)
	}
	return args
}
