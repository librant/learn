package main

import "fmt"

func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

func main() {
	nums := []int{2, 7, 11, 15, 9, 2, 5, 3, 4}
	target := 9

	ts := twoSum(nums, target)
	fmt.Printf("%v\n", ts)
}
