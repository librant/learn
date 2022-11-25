package main

import "fmt"

// 给定一个未排序的整数数组 nums ，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。
// 请你设计并实现时间复杂度为 O(n) 的算法解决此问题。

func main() {
	nums := []int{100,4,200,1,3,2}

	fmt.Printf("%d\n", longestConsecutive(nums))
}

func longestConsecutive(nums []int) int {
	n := len(nums)
	arrMap := make(map[int]bool)
	for i := 0; i < n; i++ {
		arrMap[nums[i]] = true
	}
	maxCnt := 0
	for k := range arrMap {
		// 左边的值存在则，则直接进入下一个循环
		if !arrMap[k-1] {
			currentNum := k
			currentStreak := 1
			for arrMap[currentNum+1] {
				currentNum++
				currentStreak++
			}
			if maxCnt < currentStreak {
				maxCnt = currentStreak
			}
		}
	}
	return maxCnt
}
