package main

import "fmt"

// 给定一个长度为 n 的整数数组 height
// 有 n 条垂线，第 i 条线的两个端点(i, 0)和(i, height[i])。
// 找出其中的两条线，使得它们与x轴共同构成的容器可以容纳最多的水。

func main() {
	arr := []int{1,8,6,2,5,4,8,3,7}

	fmt.Printf("max=%d", maxArea(arr))
}

// maxArea 双指针方式
func maxArea(height []int) int {
	n := len(height)
	i, j := 0, n-1
	res := 0
	for j > i {
		if height[j] > height[i] {
			res = MAX(res, height[i] * (j-i))
			i++
		} else {
			res = MAX(res, height[j] * (j-i))
			j--
		}
	}
	return res
}

// maxAreaTimeout 暴力求解，超出时间范围
func maxAreaTimeout(height []int) int {
	n := len(height)
	maxWater := 0
	for i := 0; i < n; i++ {
		for j := n - 1; j > i; j-- {
			maxWater = MAX(maxWater, MIN(height[i], height[j]) * (j - i))
		}
	}
	return maxWater
}

func MIN(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func MAX(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
