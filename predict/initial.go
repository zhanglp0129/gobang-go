package predict

import (
	"encoding/json"
	"os"
)

var (
	M *CNNModel
)

func init() {
	// 加载CNN模型
	initCNNModel(ModelName)
}

func initCNNModel(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	var model CNNModel
	err = json.Unmarshal(file, &model)
	if err != nil {
		return
	}
	M = &model
}
