package main

import "fmt"

// 你是一个专业的小偷，计划偷窃沿街的房屋。
// 每间房内都藏有一定的现金，影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，
// 如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。

func main() {
	arr := []int{1,2,3,1}

	fmt.Printf("%d", rob(arr))
}

// dp[i]: 表示前 i 个房间偷的最大现金
// dp[i] = MAX(dp[i-1], dp[i-2] + f[i])
func rob(nums []int) int {
	// 动态规划
	n := len(nums)

	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}
	if n == 2 {
		return MAX(nums[0], nums[1])
	}

	max := 0
	dp := make([]int, n)
	dp[0] = nums[0]
	dp[1] = MAX(nums[0], nums[1])
	for i := 2; i < n; i++ {
		dp[i] = MAX(dp[i-1], dp[i-2] + nums[i])
	}
	for i := 0; i < n; i++ {
		max = MAX(dp[i], max)
	}
	return max
}

func MAX(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
