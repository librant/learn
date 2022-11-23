package main

import "fmt"

// 给定 n 个非负整数表示每个宽度为 1 的柱子的高度图，计算按此排列的柱子，下雨之后能接多少雨水。

func main() {
	height := []int{4,2,0,3,2,5}

	fmt.Printf("%d\n", trap(height))
}

// 动态规划
// dp[i]: 标识 第 i 位置能接的最大高度 等于 左右两边各自最大值中的最小值 减去 f[i] 位置的值
// dp[i] = MIN(MAX(f[0 ~ i-1])，MAX(f[i+1, n-1]))) - f[i]

func trap(height []int) int {
	n := len(height)
	if n <= 2 {
		return 0
	}

	sum := 0
	for i := 1; i < n-1; i++ {
		lmax := 0
		rmax := 0
		for j := 0; j < i; j++ {
			lmax = MAX(height[j], lmax)
		}
		for k := i+1; k < n; k++ {
			rmax = MAX(height[k], rmax)
		}
		sum += MAX(0, MIN(lmax, rmax) - height[i])
	}
	return sum
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
