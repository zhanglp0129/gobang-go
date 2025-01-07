package predict

import (
	"math"
	"math/rand/v2"
	"sync"
)

type CNNModel struct {
	Conv1Weight [][][][]float64 `json:"conv1.weight"`
	Conv1Bias   []float64       `json:"conv1.bias"`
	Conv2Weight [][][][]float64 `json:"conv2.weight"`
	Conv2Bias   []float64       `json:"conv2.bias"`
	Conv3Weight [][][][]float64 `json:"conv3.weight"`
	Conv3Bias   []float64       `json:"conv3.bias"`
	Fc1Weight   [][]float64     `json:"fc1.weight"`
	Fc1Bias     []float64       `json:"fc1.bias"`
	Fc2Weight   [][]float64     `json:"fc2.weight"`
	Fc2Bias     []float64       `json:"fc2.bias"`
}

func (cm *CNNModel) Evaluate(boards [][]int, cur int, noise float64, isMy bool) float64 {
	// 将boards转为[][][]float64
	feature := make([][][]float64, 1)
	feature[0] = make([][]float64, 15)
	for i := 0; i < 15; i++ {
		feature[0][i] = make([]float64, 15)
		for j := 0; j < 15; j++ {
			feature[0][i][j] = float64(cur) * float64(boards[i][j])
		}
	}
	// 计算前向传播
	return float64(cur)*cm.forward(feature)[0] + min(noise, max(-noise, rand.NormFloat64()*noise/3))
}

// 前向传播
func (cm *CNNModel) forward(input [][][]float64) []float64 {
	// conv1
	input = cm.conv2d(input, cm.Conv1Weight, cm.Conv1Bias, 0, 1)
	cm.relu3d(input)
	// conv2
	input = cm.conv2d(input, cm.Conv2Weight, cm.Conv2Bias, 1, 1)
	cm.relu3d(input)
	// pool
	input = cm.maxPool2d(input, 2, 2)
	// conv3
	input = cm.conv2d(input, cm.Conv3Weight, cm.Conv3Bias, 1, 1)
	cm.relu3d(input)
	// flatten
	output := cm.flatten(input)
	// fc1
	output = cm.linear(output, cm.Fc1Weight, cm.Fc1Bias)
	cm.relu1d(output)
	// fc2
	output = cm.linear(output, cm.Fc2Weight, cm.Fc2Bias)
	cm.tanh1d(output)
	return output
}

// 卷积。input为输入，shape为in h w
// weight为权重，shape为out in kernel_size kernel_size
// bias为偏置参数，shape为out
func (cm *CNNModel) conv2d(input [][][]float64, weight [][][][]float64, bias []float64, padding int, stride int) [][][]float64 {
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
func (cm *CNNModel) maxPool2d(input [][][]float64, kernelSize, stride int) [][][]float64 {
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
func (cm *CNNModel) flatten(input [][][]float64) []float64 {
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
func (cm *CNNModel) linear(input []float64, weight [][]float64, bias []float64) []float64 {
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

func (cm *CNNModel) relu1d(input []float64) {
	for i := range input {
		input[i] = max(0, input[i])
	}
}

func (cm *CNNModel) relu3d(input [][][]float64) {
	for i := range input {
		for j := range input[i] {
			for k := range input[j] {
				input[i][j][k] = max(0, input[i][j][k])
			}
		}
	}
}

func (cm *CNNModel) tanh1d(input []float64) {
	for i := range input {
		input[i] = math.Tanh(input[i])
	}
}
