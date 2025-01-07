package train

import "math/rand"

// BoardData 棋盘数据
type BoardData struct {
	Feature [15][15]int
	Label   int
}

// GetData 获取AI自我对弈数据集。minData为需要的最小数据规模
func GetData(minData int, difficulty int, noise float64) []BoardData {
	res := make([]BoardData, 0)
	for {
		first, back := difficulty, difficulty
		data, win := Play(first, back, noise)
		// 数据增强
		for i := range data {
			label := win
			if i%2 == 1 {
				// 白棋
				label = -win
				for j := 0; j < 15; j++ {
					for k := 0; k < 15; k++ {
						data[i][j][k] = -data[i][j][k]
					}
				}
			}
			d := dataAugmentation(data[i])
			for j := range d {
				res = append(res, BoardData{
					Feature: d[j],
					Label:   label,
				})
			}
		}
		if len(res) > minData {
			break
		}
	}
	// 打乱顺序
	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return res
}

// 数据增强。会进行旋转，翻转，平移，再过滤
func dataAugmentation(data [15][15]int) [][15][15]int {
	// 获取落子数
	cnt := 0
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if data[i][j] != 0 {
				cnt++
			}
		}
	}
	res := make([][15][15]int, 0)
	if filter(cnt) {
		res = append(res, data)
	}
	// 翻转
	var flip [15][15]int
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			flip[i][j] = data[15-i-1][j]
		}
	}
	if filter(cnt) {
		res = append(res, flip)
	}
	// 旋转
	var t1, t2 [15][15]int
	// 旋转 90
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			t1[i][j] = data[15-j-1][i]
			t2[i][j] = flip[15-j-1][i]
		}
	}
	if filter(cnt) {
		res = append(res, t1)
	}
	if filter(cnt) {
		res = append(res, t2)
	}
	// 平移
	for _, d := range moveData(t1) {
		if filter(cnt) {
			res = append(res, d)
		}
	}
	for _, d := range moveData(t2) {
		if filter(cnt) {
			res = append(res, d)
		}
	}

	// 旋转 180
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			t1[i][j] = data[15-i-1][15-j-1]
			t2[i][j] = flip[15-i-1][15-j-1]
		}
	}
	if filter(cnt) {
		res = append(res, t1)
	}
	if filter(cnt) {
		res = append(res, t2)
	}
	// 平移
	for _, d := range moveData(t1) {
		if filter(cnt) {
			res = append(res, d)
		}
	}
	for _, d := range moveData(t2) {
		if filter(cnt) {
			res = append(res, d)
		}
	}

	// 旋转 270
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			t1[i][j] = data[j][15-i-1]
			t2[i][j] = flip[j][15-i-1]
		}
	}
	if filter(cnt) {
		res = append(res, t1)
	}
	if filter(cnt) {
		res = append(res, t2)
	}
	// 平移
	for _, d := range moveData(t1) {
		if filter(cnt) {
			res = append(res, d)
		}
	}
	for _, d := range moveData(t2) {
		if filter(cnt) {
			res = append(res, d)
		}
	}

	return res
}

// 平移数据
func moveData(data [15][15]int) [][15][15]int {
	// 获取最上，最左，最右，最下
	top, left, right, bottom := -1, -1, -1, -1
topLabel:
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if data[i][j] != 0 {
				top = i
				break topLabel
			}
		}
	}
leftLabel:
	for j := 0; j < 15; j++ {
		for i := 0; i < 15; i++ {
			if data[i][j] != 0 {
				left = i
				break leftLabel
			}
		}
	}
rightLabel:
	for j := 14; j >= 0; j-- {
		for i := 0; i < 15; i++ {
			if data[i][j] != 0 {
				right = i
				break rightLabel
			}
		}
	}
bottomLabel:
	for i := 14; i >= 0; i-- {
		for j := 0; j < 15; j++ {
			if data[i][j] != 0 {
				bottom = i
				break bottomLabel
			}
		}
	}

	res := make([][15][15]int, 0, 8)
	// 上右 下左 右 左 下 上 下右 上左
	ways := [8][2]int{{-1, 1}, {1, -1}, {0, 1}, {0, -1}, {1, 0}, {-1, 0}, {1, 1}, {-1, -1}}
	maxMove := [8]int{min(top, right), min(bottom, left), right, left, bottom, top, min(bottom, right), min(top, left)}
	for k, way := range ways {
		if maxMove[k] == 0 {
			continue
		}
		move := rand.Int()%maxMove[k] + 1
		dx, dy := way[0]*move, way[1]*move
		t := [15][15]int{}
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				if i+dx < 0 || i+dx >= 15 || j+dy < 0 || j+dy >= 15 {
					continue
				}
				t[i+dx][j+dy] = data[i][j]
			}
		}
		res = append(res, t)
	}
	return res
}

// 返回是否被过滤。被过滤的概率为1-落子数/225
func filter(cnt int) bool {
	return cnt > rand.Int()%225
}
