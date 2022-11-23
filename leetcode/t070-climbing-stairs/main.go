package main

import "fmt"

// 每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？

func main() {
	fmt.Println(climbStairs(3))
}

// 动态规划：
// dp[i]: 表示到达第 i 台阶的方法
// dp[i] = dp[i-1] + dp[i-2]
func climbStairs(n int) int {
	if n <= 2 {
		return n
	}
	dp := make([]int, n)
	dp[0] = 1
	dp[1] = 2
	for i := 2; i < n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n-1]
}
