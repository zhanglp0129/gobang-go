package predict

import (
	"github.com/emirpasic/gods/v2/queues/linkedlistqueue"
	"math"
	"math/rand"
)

func GetMaxScoreNode(boards [][]int, cur int, difficulty int, noise float64) [2]int {
	nextNodes := getNextNodes(boards)
	// 判断是否为初次落子
	if len(nextNodes) == 0 && cur == BlackChess && boards[0][0] == Empty {
		return [2]int{7, 7}
	}
	maxNode, maxScore := [2]int{-1, -1}, MinScore
	noise = math.Abs(noise)
	for _, next := range nextNodes {
		score := alphaBeta(boards, next, cur, MinScore, MaxScore, true, difficulty-1, noise)
		if score > maxScore {
			maxScore = score
			maxNode = next
		}
	}
	return maxNode
}

// 运行alpha beta剪枝算法。depth为最大搜索层数，即最大递归深度
// maxPlayer是否为极大化玩家。node为当前结点，即当前下棋位置
func alphaBeta(boards [][]int, node [2]int, cur int, alpha, beta float64, maxPlayer bool, depth int, noise float64) float64 {
	if depth == 0 {
		// 到达最大深度
		return evaluate(boards, node, cur) + min(noise, max(-noise, rand.NormFloat64()*noise/3))
	}

	boards[node[0]][node[1]] = cur
	nextNodes := getNextNodes(boards)
	if len(nextNodes) == 0 {
		boards[node[0]][node[1]] = 0
		return evaluate(boards, node, cur) + min(noise, max(-noise, rand.NormFloat64()*noise/3))
	}
	var value float64
	if maxPlayer {
		maxEval := MinScore
		for _, next := range nextNodes {
			eval := alphaBeta(boards, next, -1*cur, alpha, beta, false, depth-1, noise)
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		value = maxEval
	} else {
		minEval := MaxScore
		for _, next := range nextNodes {
			eval := alphaBeta(boards, next, -1*cur, alpha, beta, true, depth-1, noise)
			minEval = min(minEval, eval)
			beta = min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		value = minEval
	}
	boards[node[0]][node[1]] = 0
	return value
}

// 获取所有下一个结点。采用广度优先搜索，获取所有4个格子以内存在棋子的格子，
// 且不包含有棋子的格子
func getNextNodes(boards [][]int) [][2]int {
	queue := linkedlistqueue.New[[3]int]()
	mp := make(map[[2]int]struct{}) // 去重
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if boards[i][j] != Empty {
				queue.Enqueue([3]int{i, j, 4})
				mp[[2]int{i, j}] = struct{}{}
			}
		}
	}
	res := make([][2]int, 0)
	// 8个方向
	ways := [8][2]int{{-1, 1}, {1, -1}, {1, 0}, {-1, 0}, {0, -1}, {0, 1}, {1, 1}, {-1, -1}}
	// 广搜
	for !queue.Empty() {
		c, _ := queue.Dequeue()
		for _, way := range ways {
			dx, dy := way[0], way[1]
			nx, ny := dx+c[0], dy+c[1]
			if nx < 0 || nx >= 15 || ny < 0 || ny >= 15 {
				continue
			}
			if _, ok := mp[[2]int{nx, ny}]; ok {
				continue
			}
			res = append(res, [2]int{nx, ny})
			mp[[2]int{nx, ny}] = struct{}{}
			if c[2] > 1 {
				queue.Enqueue([3]int{nx, ny, c[2] - 1})
			}
		}
	}
	return res
}
