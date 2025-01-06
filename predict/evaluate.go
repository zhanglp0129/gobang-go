package predict

import (
	"math"
	"sync"
)

// 评估棋盘的价值。boards为当前棋盘，node为下一步落子位置，cur为当前是黑棋还是白棋在落子
// 其中boards[node[0]][node[1]]应该为0，即此时还没有落子
// 返回值为[-1, 1]范围内的浮点数，表示当前盘面的得分
func evaluate(boards [][]int, node [2]int, cur int) float64 {
	// 将boards转为[][][]float64
	feature := make([][][]float64, 1)
	feature[0] = make([][]float64, 15)
	for i := 0; i < 15; i++ {
		feature[0][i] = make([]float64, 15)
		for j := 0; j < 15; j++ {
			feature[0][i][j] = float64(cur) * float64(boards[i][j])
		}
	}
	feature[0][node[0]][node[1]] = float64(cur)
	// 计算前向传播
	return float64(cur) * forward(feature)[0]
}

// 前向传播
func forward(input [][][]float64) []float64 {
	// conv1
	input = conv2d(input, M.Conv1Weight, M.Conv1Bias, 0, 1)
	relu3d(input)
	// conv2
	input = conv2d(input, M.Conv2Weight, M.Conv2Bias, 1, 1)
	relu3d(input)
	// pool
	input = maxPool2d(input, 2, 2)
	// conv3
	input = conv2d(input, M.Conv3Weight, M.Conv3Bias, 1, 1)
	relu3d(input)
	// flatten
	output := flatten(input)
	// fc1
	output = linear(output, M.Fc1Weight, M.Fc1Bias)
	relu1d(output)
	// fc2
	output = linear(output, M.Fc2Weight, M.Fc2Bias)
	tanh1d(output)
	return output
}

// 卷积。input为输入，shape为in h w
// weight为权重，shape为out in kernel_size kernel_size
// bias为偏置参数，shape为out
func conv2d(input [][][]float64, weight [][][][]float64, bias []float64, padding int, stride int) [][][]float64 {
	in, out := len(input), len(weight)
	kernelSize := len(weight[0][0])
	res := make([][][]float64, out)
	var wait sync.WaitGroup
	wait.Add(out)
	for k := 0; k < out; k++ {
		m, n := len(input[0])+2*padding, len(input[0][0])+2*padding
		get := func(l, i, j int) float64 {
			if i < padding || j < padding || i >= m-padding || j >= n-padding {
				return 0.0
			}
			return input[l][i-padding][j-padding]
		}
		go func() {
			t := make([][]float64, (m-kernelSize+stride)/stride)
			for i := range t {
				t[i] = make([]float64, (n-kernelSize+stride)/stride)
			}
			for l := 0; l < in; l++ {
				for i := 0; i < m && i+kernelSize <= m; i += stride {
					for j := 0; j < n && j+kernelSize <= n; j += stride {
						var sum = 0.0
						for x := 0; x < kernelSize; x++ {
							for y := 0; y < kernelSize; y++ {
								sum += get(l, i+x, j+y) * weight[k][l][x][y]
							}
						}
						t[i][j] += sum
					}
				}
			}
			for i := range t {
				for j := range t[i] {
					t[i][j] += bias[k]
				}
			}
			res[k] = t
			wait.Done()
		}()
	}
	wait.Wait()
	return res
}

// 二维最大池化
func maxPool2d(input [][][]float64, kernelSize, stride int) [][][]float64 {
	in := len(input)
	var wait sync.WaitGroup
	wait.Add(in)
	res := make([][][]float64, in)
	for k := 0; k < in; k++ {
		m, n := len(input[0]), len(input[0][0])
		go func() {
			t := make([][]float64, (m-kernelSize+stride)/stride)
			for i := range t {
				t[i] = make([]float64, (n-kernelSize+stride)/stride)
			}
			for i := 0; i < m && i+kernelSize <= m; i += stride {
				for j := 0; j < n && j+kernelSize <= n; j += stride {
					var mx = -1e9
					for x := 0; x < kernelSize; x++ {
						for y := 0; y < kernelSize; y++ {
							mx = max(mx, input[k][i+x][j+y])
						}
					}
					t[i/stride][j/stride] = mx
				}
			}
			res[k] = t
			wait.Done()
		}()
	}
	wait.Wait()
	return res
}

// 展平
func flatten(input [][][]float64) []float64 {
	res := make([]float64, 0, len(input)*len(input[0])*len(input[0][0]))
	for i := range input {
		for j := range input[i] {
			for k := range input[j] {
				res = append(res, input[i][j][k])
			}
		}
	}
	return res
}

// 线性全连接。weight为权重，shape为out in。bias为偏置参数，shape为out
func linear(input []float64, weight [][]float64, bias []float64) []float64 {
	out, in := len(weight), len(weight[0])
	var wait sync.WaitGroup
	wait.Add(out)
	res := make([]float64, out)
	for i := 0; i < out; i++ {
		go func() {
			// 计算矩阵乘法
			var t float64
			for j := 0; j < in; j++ {
				t += input[j] * weight[i][j]
			}
			// 计算偏置参数
			t += bias[i]
			res[i] = t
			wait.Done()
		}()
	}
	wait.Wait()
	return res
}

func relu1d(input []float64) {
	for i := range input {
		input[i] = max(0, input[i])
	}
}

func relu3d(input [][][]float64) {
	for i := range input {
		for j := range input[i] {
			for k := range input[j] {
				input[i][j][k] = max(0, input[i][j][k])
			}
		}
	}
}

func tanh1d(input []float64) {
	for i := range input {
		input[i] = math.Tanh(input[i])
	}
}
