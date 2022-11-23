package main

import (
	"fmt"
)

// 给你两个单词word1和word2， 请返回将word1转换成word2 所使用的最少操作数。
// 你可以对一个单词进行如下三种操作：
// 插入一个字符
// 删除一个字符
// 替换一个字符

func main() {
	word1 := "horse"
	word2 := "ros"
	fmt.Println(minDistance(word1, word2))
}

// dp[i][j]: 表示 word1 的前 i 个字符到 word2 的前 j 个字符的编辑距离
// word1 的第 i 个字符和 word2 的第 j 个字符相同：dp[i][j] = 1 + MIN(dp[i-1][j], dp[i][j-1], dp[i-1][j-1]-1)
// word1 的第 i 个字符和 word2 的第 j 个字符不相同： dp[i][j] = 1 + MIN(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
// dp[i-1][j] 到 dp[i][j]: 删除操作，将 word1 的第 i 个字符删除
// dp[i][j-1] 到 dp[i][j]: 插入操作，将 word2 的第 j 个字符插入
// dp[i-1][j-1] 到 dp[i][j]: 替换操作， 将 word1 的第 i 个字符 替换成 word2 第 j 个字符
func minDistance(word1 string, word2 string) int {
	// 动态规划
	n := len(word1)
	m := len(word2)
	// 只要其中有一个是空串
	if n * m == 0 {
		return n+m
	}
	dp := make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dp[i] = make([]int, m+1)
	}
	// 边界条件初始化
	for i := 0; i < n+1; i++ {
		dp[i][0] = i
	}
	for j := 0; j < m+1; j++ {
		dp[0][j] = j
	}
	for i := 1; i < n+1; i++ {
		for j := 1; j < m+1; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = 1 + MIN(dp[i-1][j], dp[i][j-1], dp[i-1][j-1]-1)
			} else {
				dp[i][j] = 1 + MIN(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return dp[n][m]
}

func MIN(a, b, c int) int {
	if a <= b && a <= c{
		return a
	}
	if b <= a && b <= c {
		return b
	}
	if c <= a && c <= b {
		return c
	}
	return 0
}


