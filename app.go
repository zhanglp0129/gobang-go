package main

import (
	"context"
	"gobang-go/predict"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Predict 预测下一步应该下在哪里。boards表示棋盘，cur表示当前轮到谁下了 1黑棋 -1白棋，difficulty表示难度
// model为模型，1传统模型 2卷积模型
// 返回值为预测位置的x y坐标
func (a *App) Predict(boards [][]int, cur int, difficulty int, model int) [2]int {
	// 校验数据合法性
	if (cur != -1 && cur != 1) || difficulty <= 0 || difficulty > 3 {
		return [2]int{-1, -1}
	}
	// 校验棋盘切片长度
	if len(boards) != 15 {
		return [2]int{-1, -1}
	}
	for i := range boards {
		if len(boards[i]) != 15 {
			return [2]int{-1, -1}
		}
	}
	// 预测结果
	var m predict.Model
	switch model {
	case 1:
		m = &predict.TraditionModel{}
	case 2:
		m = predict.M
	default:
		return [2]int{-1, -1}
	}
	prediction := predict.Prediction{
		Boards:     boards,
		Cur:        cur,
		Difficulty: difficulty,
		Noise:      0,
		Model:      m,
	}
	return prediction.GetMaxScoreNode()
}
