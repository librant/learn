package main

import "fmt"

// 给你一个字符串 s，找到 s 中最长回文子串。

// 动态规划:
// dp[i][j]: 标识 i， j 位置为 回文子串
// dp[i][j] = (s[i] == s[j]) && dp[i+1][j-1]

func main() {
	str := "abcdedcba"

	fmt.Printf(longestPalindrome(str))
}

func longestPalindrome(s string) string {
	n := len(s)
	if n < 2 {
		return s
	}
	dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]bool, n)
		// 初始状态, 对角线默认为 true
		dp[i][i] = true
	}

	result := s[0:1]  //初始化结果(最小的回文就是单个字符)
	// 状态转移方程： dp[i][j] == (s[i] == s[j]) && dp[i+1][j-1]  (j - i < 3)
	for i := 2; i <= n; i++ {
		for start := 0; start < n - i + 1; start++ {
			end := start + i - 1
			if s[start] != s[end] {
				dp[start][end] = false
			} else if i < 3 {
				dp[start][end] = true
			} else {
				dp[start][end] = dp[start+1][end-1]
			}
			if dp[start][end] && (end-start+1) > len(result) {
				result = s[start:end+1]
			}
		}
	}
	return result
}
