package predict

import (
	"encoding/json"
	"os"
)

type Model struct {
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

func loadModel(filename string) (*Model, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var model Model
	err = json.Unmarshal(file, &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
