package main

import "fmt"

// 给定一个整数数组temperatures，表示每天的温度，返回一个数组answer，其中answer[i]是指对于第 i 天，下一个更高温度出现在几天后。
// 如果气温在这之后都不会升高，请在该位置用0 来代替。

func main() {
	temperatures := []int{73,74,75,71,69,72,76,73}

	fmt.Printf("%v", dailyTemperatures(temperatures))
}

// 单调栈解法：维护一个不增的单调栈，栈中存入对应的下标值
func dailyTemperatures(temperatures []int) []int {
	n := len(temperatures)
	answer := make([]int, n)
	var stack []int

	for i := 0; i < n; i++ {
		temperature := temperatures[i]
		// 如果栈不为空，且当前温度比栈顶的温度高，需要将温度较低元素的下标出栈
		for len(stack) > 0 && temperature > temperatures[stack[len(stack)-1]] {
			// 首先，栈顶的元素找到第一个高温下标
			prevIndex := stack[len(stack)-1]
			// 将栈顶的元素出栈
			stack = stack[:len(stack)-1]
			answer[prevIndex] = i - prevIndex
		}
		stack = append(stack, i)
	}
	return answer
}

// 暴力超时
func dailyTemperaturesTimeout(temperatures []int) []int {
	n := len(temperatures)
	answer := make([]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j ++ {
			if temperatures[j] > temperatures[i] {
				answer[i] = j - i
				break
			}
		}
	}
	return answer
}
