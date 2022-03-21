/**
 * @Title
 * @Description
 * @Author hyman
 * @Date 2022-03-20
 **/
package cmd

import (
	"fmt"
)

// 垂直竖线
type VerticalLines []int

func (vv VerticalLines) Show(index int) bool {
	for _, v := range vv {
		if v == index {
			return true
		}
	}
	return false
}
func (vv VerticalLines) Clone() []int {
	var _vv VerticalLines
	_vv = append(_vv, vv...)
	return _vv
}

/*
打印树形结构
vv：垂直的竖线
indent：缩进
t：文件
*/
func PrintFileTree(vv VerticalLines, indent int, t *TargetFile) {
	l := len(t.Children)
	for i, tf := range t.Children {
		for j := 0; j < indent; j++ {
			if vv.Show(j) { // 显示竖线
				fmt.Print("│ ")
			} else {
				fmt.Print("  ")
			}
		}
		nextVV := vv.Clone()
		if i == l-1 { // 最后一个
			fmt.Print("└")
		} else {
			fmt.Print("├")
			nextVV = append(nextVV, indent)
		}
		fmt.Print("──")
		// 蓝色表示正常替换
		// 黄色表示只有
		fmt.Print(tf.Label())
		fmt.Println()
		PrintFileTree(nextVV, indent+2, tf)
	}
}
