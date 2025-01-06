package main

import (
	"encoding/json"
	"fmt"
	"gobang-go/train"
	"os"
	"strconv"
)

func main() {
	// 获取命令行参数，格式为：gobang-train.exe <minData> difficulty noise
	minData, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	difficulty, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	noise, err := strconv.ParseFloat(os.Args[3], 64)
	res, err := json.Marshal(train.GetData(minData, difficulty, noise))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}
