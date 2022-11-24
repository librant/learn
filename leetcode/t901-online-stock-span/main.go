package main

import (
	"fmt"
	"math"
)

func main() {
	s := Constructor()
	fmt.Printf("%v\n", s)

	fmt.Printf("%d\n", s.Next(100))
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.Next(80))
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.Next(60))
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.Next(70))
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.Next(60))
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.Next(75))
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.Next(85))
	fmt.Printf("%v\n", s)
	fmt.Printf("%d\n", s.Next(110))
	fmt.Printf("%v\n", s)
}

// 在调 Next 的函数时，将数据存入，并给出之前下降的天数

type StockSpanner struct {
	stack [][2]int // 二维数组主要是存放当前的天数和股票的金额
	idx   int      // 当前最后操作的索引
}

func Constructor() StockSpanner {
	return StockSpanner{[][2]int{{-1, math.MaxInt32}}, -1}
}

// 使用单调不增的栈的方式，当输入的金额超过栈顶元素，需要将栈中的小的元素弹出

func (s *StockSpanner) Next(price int) int {
	s.idx++
	// 当前的元素超过栈顶的元素，弹出所有小于当前存入的金额
	for price >= s.stack[len(s.stack)-1][1] {
		s.stack = s.stack[:len(s.stack)-1]
	}
	// 再将当前的栈元素加入
	s.stack = append(s.stack, [2]int{s.idx, price})
	return s.idx - s.stack[len(s.stack)-2][0]
}

/**
 * Your StockSpanner object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Next(price);
 */
