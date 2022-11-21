package main

import "fmt"

// ListNode 链表节点定义
type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	l1 := new(ListNode)
	l2 := new(ListNode)
	head1 := l1
	head2 := l2
	for i := 0; i < 3; i++ {
		if i != 2 {
			l1.Next = new(ListNode)
			l2.Next = new(ListNode)
		}
		l1.Val = i + 2
		l2.Val = i + 4
		l1 = l1.Next
		l2 = l2.Next
	}
	l3 := addTwoNumbers(head1, head2)
	for l3 != nil {
		fmt.Printf("%d ", l3.Val)
		l3 = l3.Next
	}
	fmt.Println()
}

// addTwoNumbers 两数相加
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	// 定义一个尾结点，或者可以理解为临时节点
	var tail *ListNode
	var head *ListNode
	// 判断是否进位
	carry := 0
	for l1 != nil || l2 != nil {
		n1, n2 := 0, 0
		if l1 != nil {
			n1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			n2 = l2.Val
			l2 = l2.Next
		}
		sum := (n2 + n1 + carry) % 10
		carry = (n2 + n1 + carry) / 10

		if head == nil {
			// 这里是保存首节点地址，用于返回
			head = &ListNode{Val: sum}
			tail = head
		} else {
			// 这里对节点指向的下一个 ListNode 进行赋值，并移动链表位置
			tail.Next = &ListNode{Val:sum}
			tail = tail.Next
		}
	}
	if carry > 0 {
		tail.Next = &ListNode{Val:carry}
	}
	return head
}

// 链表的遍历，抓住头节点
// 链接结束的判断

