package main

import "fmt"

// 给定 n 个非负整数，用来表示柱状图中各个柱子的高度。每个柱子彼此相邻，且宽度为 1 。
// 求在该柱状图中，能够勾勒出来的矩形的最大面积。

func main() {
	heights := []int{2,1,5,6,2,3}

	fmt.Printf("%d", largestRectangleArea(heights))
}

// 单调栈：
func largestRectangleArea(heights []int) int {
	n := len(heights)
	left, right := make([]int, n), make([]int, n)
	mono_stack := []int{}
	for i := 0; i < n; i++ {
		for len(mono_stack) > 0 && heights[mono_stack[len(mono_stack)-1]] >= heights[i] {
			mono_stack = mono_stack[:len(mono_stack)-1]
		}
		if len(mono_stack) == 0 {
			left[i] = -1
		} else {
			left[i] = mono_stack[len(mono_stack)-1]
		}
		mono_stack = append(mono_stack, i)
	}
	mono_stack = []int{}
	for i := n - 1; i >= 0; i-- {
		for len(mono_stack) > 0 && heights[mono_stack[len(mono_stack)-1]] >= heights[i] {
			mono_stack = mono_stack[:len(mono_stack)-1]
		}
		if len(mono_stack) == 0 {
			right[i] = n
		} else {
			right[i] = mono_stack[len(mono_stack)-1]
		}
		mono_stack = append(mono_stack, i)
	}
	ans := 0
	for i := 0; i < n; i++ {
		ans = MAX(ans, (right[i] - left[i] - 1) * heights[i])
	}
	return ans
}

func largestRectangleAreaTimeout(heights []int) int {
	// 暴力解法
	n := len(heights)
	max := 0
	for i := 0; i < n; i++ {
		// 往左扩散
		lmin := 0
		rmin := n-1
		for j := i-1; j >= 0; j-- {
			// 找到比 i 位置小的坐标
			if heights[j] < heights[i] {
				lmin = j+1
				break
			}
		}
		// 往右扩散
		for k := i + 1; k < n; k++ {
			// 找到比 i 位置小的坐标
			if heights[k] < heights[i] {
				rmin = k-1
				break
			}
		}
		max = MAX(max, MAX(heights[i], (rmin-lmin+1) * heights[i]))
	}
	return max
}

func MAX(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
