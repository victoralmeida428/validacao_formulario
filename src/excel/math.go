package excel

import "math"

func calcularSomaQuadrados(valores []float64) float64 {
	soma := 0.0
	for _, v := range valores {
		soma += v * v
	}
	return soma
}

func calcularMedia(valores []float64) float64 {
	if len(valores) == 0 {
		return 0
	}
	soma := 0.0
	for _, v := range valores {
		soma += v
	}
	return soma / float64(len(valores))
}

func calcularDesvioPadrao(valores []float64) float64 {
	if len(valores) < 2 {
		return 0
	}
	media := calcularMedia(valores)
	var variancia float64
	for _, v := range valores {
		variancia += (v - media) * (v - media)
	}

	return math.Sqrt(variancia / float64(len(valores)-1))
}
