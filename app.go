package main

import (
	"context"
	"time"
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
// 返回值为预测位置的x y坐标
func (a *App) Predict(boards [15][15]int, cur int, difficulty int) [2]int {
	time.Sleep(5 * time.Second)
	for i := range boards {
		for j := range boards[i] {
			if boards[i][j] == 0 {
				return [2]int{i, j}
			}
		}
	}
	return [2]int{-1, -1}
}
