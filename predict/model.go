package predict

// Model 定义通用模型接口
type Model interface {
	Evaluate(boards [][]int, cur int, noise float64, isMy bool) float64
}
