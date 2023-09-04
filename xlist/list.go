package xlist

/*
author:liujinyin
date:2022-3-8
desc:双休链表list
*/

import "fmt"

type ListIterCallback func(idx int, node *Node) bool

type List struct {
	root   Node // 头节点
	length int  // list长度
}

type Node struct {
	list  *List
	Value any
	next  *Node
	prev  *Node
}

func (n *Node) Remove() {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.next = nil
	n.prev = nil

	n.list.length--
	n.list = nil
}

func (n *Node) String() string {
	return fmt.Sprint(n.Value)
}

func NewList() *List {
	l := &List{} // 获取List{}的地址
	l.length = 0 // list初始长度为0
	l.root.next = &l.root
	l.root.prev = &l.root
	return l
}

func (n *List) String() string {
	return fmt.Sprint(n.length)
}

func (l *List) IsEmpty() bool {
	return l.root.next == &l.root
}

func (l *List) Length() int {
	return l.length
}

func (l *List) Append(items ...any) *List {

	for _, item := range items {
		var n *Node
		switch v := item.(type) {
		case Node:
			n = &v
		case *Node:
			n = v
		default:
			n = &Node{Value: item}
		}
		n.list = l
		n.next = &l.root
		n.prev = l.root.prev
		l.root.prev.next = n
		l.root.prev = n
		l.length++
	}
	return l
}

func (l *List) Remove(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.next = nil
	n.prev = nil
	l.length--
}

func (l *List) Iter(callback ListIterCallback) {
	node := l.root.next
	if l.IsEmpty() {
		return
	}
	idx := 0
	for {
		next := node.next
		if !callback(idx, node) {
			return
		}
		if next == &l.root {
			return
		}
		idx++
		node = next
	}
}
