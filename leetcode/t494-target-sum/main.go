package main

import "fmt"

func main() {
	nums := []int{1,1,1,1,1}
	target := 3

	fmt.Printf("%d\n", findTargetSumWays(nums, target))
}

// 回溯过程中维护一个计数器 count，当遇到一种表达式的结果等于目标数 target 时，
// 将 count 的值加 1。遍历完所有的表达式之后，即可得到结果等于目标数 target 的表达式的数目。

func findTargetSumWays(nums []int, target int) (count int) {
	var backtrack func(int, int)
	backtrack = func(index, sum int) {
		// 当所有的数据都回溯完成后，计算 sum 和 target 的值
		if index == len(nums) {
			if sum == target {
				count++
			}
			return
		}
		backtrack(index+1, sum+nums[index])
		backtrack(index+1, sum-nums[index])
	}
	// 从 index = 0 的位置开始回溯
	backtrack(0, 0)
	return
}
