package train

import "gobang-go/predict"

// Play AI自己下棋。first和back分别为先手难度和后手难度，
// noise为噪声，返回值为每一步的棋盘和游戏结果
func Play(first, back int, noise float64) ([][15][15]int, int) {
	res := make([][15][15]int, 0)
	boards := make([][]int, 15)
	for i := 0; i < 15; i++ {
		boards[i] = make([]int, 15)
	}
	for i := 1; i <= 225; i++ {
		cur, difficulty := 1, first
		if i%2 == 0 {
			cur = -1
			difficulty = back
		}

		node := predict.GetMaxScoreNode(boards, cur, difficulty, noise)
		if node[0] == -1 || node[1] == -1 {
			return res, 0
		}
		over := gameOver(boards, node, cur)
		boards[node[0]][node[1]] = cur
		res = append(res, boardsSliceToArray(boards))
		if over {
			return res, cur
		}
	}
	return res, 0
}

// 判断是否游戏结束
func gameOver(boards [][]int, node [2]int, cur int) bool {
	ways := [8][2]int{{1, -1}, {-1, 0}, {0, 1}, {-1, -1}}
	for _, way := range ways {
		dx, dy := way[0], way[1]
		cnt := 1
		for i := 1; true; i++ {
			nx, ny := node[0]+dx*i, node[1]+dy*i
			if nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || boards[nx][ny] != cur {
				break
			}
			cnt++
		}
		dx *= -1
		dy *= -1
		for i := 1; true; i++ {
			nx, ny := node[0]+dx*i, node[1]+dy*i
			if nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || boards[nx][ny] != cur {
				break
			}
			cnt++
		}
		if cnt >= 5 {
			return true
		}
	}
	return false
}

// 将切片类型的棋盘转为数组类型
func boardsSliceToArray(boards [][]int) [15][15]int {
	var res [15][15]int
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			res[i][j] = boards[i][j]
		}
	}
	return res
}
