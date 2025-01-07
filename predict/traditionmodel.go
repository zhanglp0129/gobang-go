package predict

import (
	"math/rand/v2"
	"slices"
	"strings"
)

type TraditionModel struct{}

const (
	p0 = 100000000
	p1 = 1000
	p2 = 100
	p3 = 10
	p4 = 1
	d4 = p1
	d3 = p2
	d2 = p3
	d1 = p4
)

var scores = map[string][2]float64{
	"11111":   {p0, p0},
	"011110":  {p1, d4},
	"0111010": {p1, d3 + d1},
	"0110110": {p1, d2 * 2},
	"211110":  {p1, 0},
	"2111010": {p1, d1},
	"2101112": {p1, 0},
	"2110112": {p1, 0},

	"01110":   {p2, d3},
	"010110":  {p2, d2 + d1},
	"0100110": {p2, d2 + d1},
	"2011100": {p2, 0},
	"211100":  {p2, 0},
	"210110":  {p2, 0},
	"2100110": {p2, d2},
	"2001112": {p2, 0},
	"2011102": {p2, 0},
	"2101102": {p2, 0},
	"2101012": {p2, 0},

	"0101010": {p2, d2 + d1},
	"2110100": {p2, 0},
	"2110102": {p2, 0},
	"001100":  {p3, d2},
	"01010":   {p3, 2 * d1},
	"211000":  {p3, 0},
	"210100":  {p3, 0},
	"210010":  {p3, 0},
	// 1
	"00100": {p4, d1},
}

var cache = make(map[string]float64)

func init() {
	keys := make([]string, 0)
	for k := range scores {
		keys = append(keys, k)
	}
	for _, k := range keys {
		key := []byte(k)
		slices.Reverse(key)
		scores[string(key)] = scores[k]
	}

	keys = make([]string, 0)
	for k := range scores {
		keys = append(keys, k)
	}
	for _, k := range keys {
		key := []byte(k)
		for i := 0; i < len(key); i++ {
			if key[i] == '1' {
				key[i] = '2'
			} else if key[i] == '2' {
				key[i] = '1'
			}
		}
		scores[string(key)] = [2]float64{-scores[k][0], -scores[k][1]}
	}
}

func (tm *TraditionModel) Evaluate(boards [][]int, cur int, noise float64, isMy bool) float64 {
	// 提取每一行、每一列、每条对角线，并组成字符串
	// 0空位 1我方 2对方
	// 计算总分
	var score float64
	// 行
	for i := 0; i < 15; i++ {
		var builder strings.Builder
		for j := 0; j < 15; j++ {
			write(&builder, boards[i][j], cur)
		}
		score += computeScore(&builder, isMy)
	}
	// 列
	for j := 0; j < 15; j++ {
		var builder strings.Builder
		for i := 0; i < 15; i++ {
			write(&builder, boards[i][j], cur)
		}
		score += computeScore(&builder, isMy)
	}
	// 上主对角线
	for j := 0; j < 15; j++ {
		var builder strings.Builder
		for d := 0; j+d < 15; d++ {
			write(&builder, boards[d][j+d], cur)
		}
		score += computeScore(&builder, isMy)
	}
	// 下主对角线
	for i := 1; i < 15; i++ {
		var builder strings.Builder
		for d := 0; i+d < 15; d++ {
			write(&builder, boards[i+d][d], cur)
		}
		score += computeScore(&builder, isMy)
	}
	// 上副对角线
	for j := 0; j < 15; j++ {
		var builder strings.Builder
		for d := 0; d < 15 && j-d >= 0; d++ {
			write(&builder, boards[d][j-d], cur)
		}
		score += computeScore(&builder, isMy)
	}
	// 下副对角线
	for i := 1; i < 15; i++ {
		var builder strings.Builder
		for d := 0; d+i < 15 && 14-d >= 0; d++ {
			write(&builder, boards[d+i][14-d], cur)
		}
		score += computeScore(&builder, isMy)
	}

	return score + min(noise, max(-noise, rand.NormFloat64()*noise/3))
}

func write(builder *strings.Builder, cur int, board int) {
	if board*cur == 0 {
		builder.WriteRune('0')
	} else if board*cur == 1 {
		builder.WriteRune('1')
	} else if board*cur == -1 {
		builder.WriteRune('2')
	}
}

func computeScore(builder *strings.Builder, isMy bool) float64 {
	line := builder.String()
	if r, ok := cache[line]; ok {
		return r
	}
	var res float64
	if len(line) >= 5 {
		for k, v := range scores {
			if strings.Contains(line, k) {
				if isMy {
					res += v[0]
				} else {
					res += v[1]
				}
			}
		}
	}
	cache[line] = res
	return res
}
