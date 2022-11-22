package main

// 给定一个长度为 n 的整数数组 height
// 有 n 条垂线，第 i 条线的两个端点(i, 0)和(i, height[i])。
// 找出其中的两条线，使得它们与x轴共同构成的容器可以容纳最多的水。

func main() {

}

func maxArea(height []int) int {
	n := len(height)

	// 采用双指针的方式，当哪边的水少，就移动哪边的指针
	maxWater := MIN(height[0], height[n-1]) * (n-1)
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
