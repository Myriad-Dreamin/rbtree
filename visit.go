package rbtree

import (
	"fmt"
	"reflect"
)

type Node interface {
	GetValue() interface{}
	Index(index int) interface{}
}

type NodeTrace interface {
	GetValueTrace() interface{}
	Index(index int) interface{}
}

func Inorder(treeNode Node) {
	fmt.Printf("%v ", treeNode.GetValue())
	for idx, node := 0, treeNode.Index(0).(Node); node!= nil; idx, node = idx+1, treeNode.Index(idx+1).(Node) {
		Inorder(node)
	}
}

func Preorder(treeNode Node) {
	if reflect.ValueOf(treeNode).IsNil() {
		return
	}
	Preorder(treeNode.Index(0).(Node))
	fmt.Printf("%v ", treeNode.GetValue())
	Preorder(treeNode.Index(1).(Node))
}

func printL(n Node, dep int) {
	if reflect.ValueOf(n).IsNil() {
		return
	}
	printL(n.Index(0).(Node), dep+1)
	fmt.Println()
	for i := 0; i < dep; i++ {
		fmt.Printf("    ")
	}
	fmt.Println(n.GetValue())
	printL(n.Index(1).(Node), dep+1)
}

func printLTrace(n NodeTrace, dep int) {
	if reflect.ValueOf(n).IsNil() {
		return
	}
	printLTrace(n.Index(0).(NodeTrace), dep+1)
	fmt.Println()
	for i := 0; i < dep; i++ {
		fmt.Printf("    ")
	}
	fmt.Println(n.GetValueTrace())
	printLTrace(n.Index(1).(NodeTrace), dep+1)
}

func PrintL(n interface{}) {
	if nt, ok := n.(NodeTrace); ok {
		printLTrace(nt, 0)
	} else if nt, ok := n.(Node); ok {
		printL(nt, 0)
	} else {
		panic("")
	}
	return
}

