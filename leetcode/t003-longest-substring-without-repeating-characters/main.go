package main

import "fmt"

// 给定一个字符串 s ，请你找出其中不含有重复字符的 最长子串 的长度。

func main() {
	s := "pwwkew"

	max := lengthOfLongestSubstring(s)
	fmt.Printf("max = %d\n", max)
}

func lengthOfLongestSubstring(s string) int {
	// 哈希集合，记录每个字符是否出现过
	m := map[byte]int{}
	n := len(s)

	// 右指针，初始值为 -1，相当于我们在字符串的左边界的左侧，还没有开始移动
	rk, ans := -1, 0
	for i := 0; i < n; i++ {
		if i != 0 {
			// 左指针向右移动一格，移除一个字符
			delete(m, s[i-1])
		}
		for rk + 1 < n && m[s[rk+1]] == 0 {
			// 不断地移动右指针
			m[s[rk+1]]++
			rk++
		}
		// 第 i 到 rk 个字符是一个极长无重复字符子串
		ans = MAX(ans, rk - i + 1)
	}
	return ans
}

func MAX(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
