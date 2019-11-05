package rbtree

import (
	"fmt"
	"reflect"
)

type Node interface {
	GetValue() interface{}
	Index(index int) Node
}

func Inorder(treeNode Node) {
	fmt.Printf("%v ", treeNode.GetValue())
	for idx, node := 0, treeNode.Index(0); node!= nil; idx, node = idx+1, treeNode.Index(idx+1) {
		Inorder(node)
	}
}

func Preorder(treeNode Node) {
	if reflect.ValueOf(treeNode).IsNil() {
		return
	}
	Preorder(treeNode.Index(0))
	fmt.Printf("%v ", treeNode.GetValue())
	Preorder(treeNode.Index(1))
}





