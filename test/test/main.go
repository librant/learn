package main

import "fmt"

type DLinkNode struct {
	Key, Value int
	Prev, Next *DLinkNode
}

type HaspMap struct {
	Head, Tail *DLinkNode
}

func main() {

}

func InitDLinkNode(key, value int) *DLinkNode {
	return &DLinkNode{
		Key:   key,
		Value: value,
		Prev:  &DLinkNode{},
		Next:  &DLinkNode{},
	}
}

func New() *HaspMap {
	head := InitDLinkNode(0, 0)
	tail := InitDLinkNode(0, 0)
	head.Next = tail
	tail.Prev = head
	return &HaspMap{
		Head: head,
		Tail: tail,
	}
}

func (h *HaspMap) Get(key int) (int, error) {
	temp := h.Head.Next
	for temp != nil {
		if temp.Key == key {
			return temp.Value, nil
		}
	}
	return -1, fmt.Errorf("not found")
}

func (h *HaspMap) Insert(key, value int) {
	_, err := h.Get(key)
	if err != nil {
		// 当前的 Key 不在链表中, 增加对应的节点
		temp := InitDLinkNode(key, value)
		// 获取 head 之后的节点指针
		next := h.Head.Next
		h.Head.Next = temp
		temp.Prev = h.Head
		temp.Next = next
		next.Prev = temp
		return
	}
	// 当前的 key 在链表中
	temp := h.Head.Next
	for temp != nil {
		if temp.Key == key {
			temp.Value = value
			return
		}
	}
}
