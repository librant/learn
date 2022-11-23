package main

import (
	"fmt"
	"sort"
)

// 给你一个整数数组 nums ，判断是否存在三元组 [nums[i], nums[j], nums[k]] 满足 i != j、i != k 且 j != k ，
// 同时还满足 nums[i] + nums[j] + nums[k] == 0

func main() {
	nums := []int{-1, 0, 1, 2, -1, -4}

	ret := threeSumTimeout(nums)
	fmt.Printf("ret: %v", ret)
}

func threeSumTimeout(nums []int) [][]int {
	n := len(nums)
	// 先对数组进行升序排序，主要是剔除重复数组
	sort.Ints(nums)
	ans := make([][]int, 0)

	// 从第一个数开始遍历
	for i := 0; i < n; i++ {
		// 如果第一个数重复，则是重复数组
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		// 第三个数从最右边进行遍历
		k := n - 1
		target := -1 * nums[i]

		for j := i + 1; j < n; j++ {
			// 第二个数也需要不相同，相同则是重复数组
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			// 需要保证 b 的指针在 c 的指针的左侧, 这里是因为是正序，也就是说，后面的数大了，需要往左移动 c 的指针
			for j < k && nums[j]+nums[k] > target {
				k--
			}
			// 如果指针重合，随着 b 后续的增加
			// 就不会有满足 a+b+c=0 并且 b<c 的 c 了，可以退出循环
			if j == k {
				break
			}
			if nums[j]+nums[k] == target {
				ans = append(ans, []int{nums[i], nums[j], nums[k]})
			}
		}
	}
	return ans
}
