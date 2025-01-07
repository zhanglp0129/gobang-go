package predict

import (
	"github.com/emirpasic/gods/v2/queues/linkedlistqueue"
	"math"
	"math/rand"
)

type Prediction struct {
	Boards     [][]int
	Cur        int
	Difficulty int
	Noise      float64
	Model      Model
}

func (p *Prediction) GetMaxScoreNode() [2]int {
	nextNodes := getNextNodes(p.Boards)
	// 判断是否为初次落子
	if len(nextNodes) == 0 && p.Cur == BlackChess && p.Boards[7][7] == Empty {
		return [2]int{7, 7}
	}
	maxNodes, maxScore := make([][2]int, 0), MinScore
	p.Noise = math.Abs(p.Noise)
	for _, next := range nextNodes {
		p.Boards[next[0]][next[1]] = p.Cur
		score := p.alphaBeta(MinScore, MaxScore, false, p.Difficulty)
		if score > maxScore {
			maxScore = score
			maxNodes = [][2]int{next}
		} else if score == maxScore {
			maxNodes = append(maxNodes, next)
		}
		p.Boards[next[0]][next[1]] = Empty
	}
	if len(maxNodes) == 0 {
		return [2]int{-1, -1}
	}
	return maxNodes[rand.Int()%len(maxNodes)]
}

// 运行alpha beta剪枝算法。depth为最大搜索层数，即最大递归深度
// maxPlayer是否为极大化玩家。node为当前结点，即当前下棋位置
func (p *Prediction) alphaBeta(alpha, beta float64, maxPlayer bool, depth int) float64 {
	if depth == 0 {
		// 到达最大深度
		return p.Model.Evaluate(p.Boards, p.Cur, p.Noise, !maxPlayer)
	}

	nextNodes := getNextNodes(p.Boards)
	if len(nextNodes) == 0 {
		return p.Model.Evaluate(p.Boards, p.Cur, p.Noise, !maxPlayer)
	}
	if maxPlayer {
		maxEval := MinScore
		for _, next := range nextNodes {
			p.Boards[next[0]][next[1]] = p.Cur
			eval := p.alphaBeta(alpha, beta, false, depth-1)
			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)
			p.Boards[next[0]][next[1]] = Empty
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := MaxScore
		for _, next := range nextNodes {
			p.Boards[next[0]][next[1]] = p.Cur
			eval := p.alphaBeta(alpha, beta, true, depth-1)
			minEval = min(minEval, eval)
			beta = min(beta, eval)
			p.Boards[next[0]][next[1]] = Empty
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

// 获取所有下一个结点。采用广度优先搜索，获取所有4个格子以内存在棋子的格子，
// 且不包含有棋子的格子
func getNextNodes(boards [][]int) [][2]int {
	queue := linkedlistqueue.New[[3]int]()
	mp := make(map[[2]int]struct{}) // 去重
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if boards[i][j] != Empty {
				queue.Enqueue([3]int{i, j, 2})
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
