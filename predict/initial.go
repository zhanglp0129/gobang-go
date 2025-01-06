package predict

import (
	"reflect"
)

var (
	M *Model
)

func init() {
	model, err := loadModel(ModelName)
	if err != nil {
		panic(err)
	}
	M = model
}

func getShape(slice any) []int {
	var shape []int

	// 检查输入是否是切片
	v := reflect.ValueOf(slice)
	for v.Kind() == reflect.Slice {
		shape = append(shape, v.Len())
		if v.Len() > 0 {
			v = v.Index(0) // 进入下一层
		} else {
			break // 空切片，停止递归
		}
	}

	// 如果最终没有切片结构，返回错误
	if len(shape) == 0 {
		return nil
	}

	return shape
}
