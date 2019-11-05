package rbtree

import (
	"fmt"
	"reflect"
	"testing"
)

type Int int

func (t Int) Less(i interface{}) bool {
	return t < i.(Int)
}

func (t Int) Equal(i interface{}) bool {
	return t == i.(Int)
}

func Equivalent(left, right *RBNode) bool {
	if (left == nil) != (right == nil) {
		return false
	}
	if left == nil {
		return true
	}
	if left.Color != right.Color || !reflect.DeepEqual(left.V, right.V) {
		return false
	}
	if left.HasLeftChild() != right.HasLeftChild() {
		return false
	}
	if left.HasRightChild() != right.HasRightChild() {
		return false
	}

	if left.HasLeftChild() && !Equivalent(left.Ch[0], right.Ch[0]) {
		return false
	}

	if left.HasRightChild() && !Equivalent(left.Ch[1], right.Ch[1]) {
		return false
	}
	return true
}

func Test_insert(t *testing.T) {
	type args struct {
		rt *RBNode
		n  *RBNode
	}

	var A, B, C, D, E, F, G Int = 1, 2, 3, 4, 5, 6, 7
	_, _, _, _, _, _, _ = A, B, C, D, E, F, G
	tests := []struct {
		name string
		args args
		want *RBNode
	}{
		{
			name: "test_no_root",
			args: args{
				rt: nil,
				n: &RBNode{
					Color: ColorRed, V: A,
					Ch: [2]*RBNode{},
				},
			},
			want: &RBNode{
				Color: ColorBlack, V: A,
				Ch: [2]*RBNode{},
			},
		},
		{
			name: "test_father_black",
			args: args{
				rt: &RBNode{
					Color: ColorBlack, V: A,
					Ch: [2]*RBNode{},
				},
				n: &RBNode{
					Color: ColorRed, V: B,
					Ch: [2]*RBNode{},
				},
			},
			want: &RBNode{
				Color: ColorBlack, V: A,
				Ch: [2]*RBNode{
					nil,
					{
						Color: ColorRed, V: B,
						Ch: [2]*RBNode{},
					},
				},
			},
		},
		{
			name: "test_father_red_uncle_red",
			args: args{
				rt: &RBNode{
					Color: ColorBlack,
					V:     B,
					Ch: [2]*RBNode{
						{
							Color: ColorRed, V: A,
							Ch: [2]*RBNode{},
						},
						{
							Color: ColorRed, V: C,
							Ch: [2]*RBNode{},
						},
					},
				},
				n: &RBNode{
					Color: ColorRed,
					V:     D,
					Ch:    [2]*RBNode{},
				},
			},
			want: &RBNode{
				Color: ColorBlack,
				V:     B,
				Ch: [2]*RBNode{
					{
						Color: ColorBlack, V: A,
						Ch: [2]*RBNode{},
					},
					{
						Color: ColorBlack, V: C,
						Ch: [2]*RBNode{
							nil,
							{
								Color: ColorRed,
								V:     D,
								Ch:    [2]*RBNode{},
							},
						},
					},
				},
			},
		},
		{
			name: "test_father_red_uncle_black",
			args: args{
				rt: &RBNode{
					Color: ColorBlack,
					V:     C,
					Ch: [2]*RBNode{
						{
							Color: ColorBlack, V: B,
							Ch: [2]*RBNode{},
						},
						{
							Color: ColorRed, V: D,
							Ch: [2]*RBNode{},
						},
					},
				},
				n: &RBNode{
					Color: ColorRed,
					V:     E,
					Ch:    [2]*RBNode{},
				},
			},

			//    C,B
			// B,B   D,R
			//          E,R

			//       D,B
			//    C,R   E,R
			// B,B
			want: &RBNode{
				Color: ColorBlack,
				V:     D,
				Ch: [2]*RBNode{
					{
						Color: ColorRed, V: C,
						Ch: [2]*RBNode{
							{
								Color: ColorBlack,
								V:     B,
								Ch:    [2]*RBNode{},
							},
						},
					},
					{
						Color: ColorRed, V: E,
						Ch: [2]*RBNode{
						},
					},
				},
			},
		},
		{
			name: "test_father_red_uncle_black",
			args: args{
				rt: &RBNode{
					Color: ColorBlack,
					V:     D,
					Ch: [2]*RBNode{
						{
							Color: ColorRed, V: B,
							Ch: [2]*RBNode{},
						},
						{
							Color: ColorBlack, V: E,
							Ch: [2]*RBNode{},
						},
					},
				},
				n: &RBNode{
					Color: ColorRed,
					V:     C,
					Ch:    [2]*RBNode{},
				},
			},

			//       D,B
			// B,R      E,B
			//    C,R
			//       D,R
			//    C,B   E,B
			// B,R
			//    C,B
			// B,R   D,R
			//          E,B
			want: &RBNode{
				Color: ColorBlack,
				V:     C,
				Ch: [2]*RBNode{
					{
						Color: ColorRed, V: B,
						Ch: [2]*RBNode{},
					},
					{
						Color: ColorRed, V: D,
						Ch: [2]*RBNode{
							nil,
							{
								Color: ColorBlack,
								V:     E,
								Ch:    [2]*RBNode{},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insert(tt.args.rt, tt.args.n); !Equivalent(got, tt.want) {
				t.Errorf("insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func checkFather(t *RBNode) bool {
	if t == nil {
		return true
	}
	if t.Ch[0] != nil {
		if t.Ch[0].Faz != t {
			return false
		}
		if !checkFather(t.Ch[0]) {
			return false
		}
	}
	if t.Ch[1] != nil {
		if t.Ch[1].Faz != t {
			return false
		}
		if !checkFather(t.Ch[1]) {
			return false
		}
	}
	return true
}


func TestRBTree(t *testing.T) {

	var values []Int

	var tree *RBNode
	for i := Int(10); i < 20; i++ {
		tree = tree.Insert(i)
		checkFather(tree)
		values = append(values, i)
	}
	for i := Int(1); i < 10; i++ {
		tree = tree.Insert(i)
		checkFather(tree)
		values = append(values, i)
	}
	for i := Int(20); i < 30; i++ {
		tree = tree.Insert(i)
		checkFather(tree)
		values = append(values, i)
	}
	for i := Int(10); i < 20; i++ {
		tree = tree.Insert(i)
		checkFather(tree)
		values = append(values, i)
	}
	for i := Int(1); i < 10; i++ {
		tree = tree.Insert(i)
		checkFather(tree)
		values = append(values, i)
	}
	for i := Int(20); i < 30; i++ {
		tree = tree.Insert(i)
		checkFather(tree)
		values = append(values, i)
	}
	for i := Int(10); i < 20; i++ {
		tree = tree.Insert(i)
		if !checkFather(tree) {
			panic("bad father")
		}
		values = append(values, i)
	}
	for i := Int(1); i < 10; i++ {
		tree = tree.Insert(i)
		if !checkFather(tree) {
			panic("bad father")
		}
		values = append(values, i)
	}
	for i := Int(20); i < 30; i++ {
		tree = tree.Insert(i)
		if !checkFather(tree) {
			panic("bad father")
		}
		values = append(values, i)
	}
	Preorder(tree)
	for i := range values {
		tree = tree.Delete(values[i])
		if !checkFather(tree) {
			panic("bad father")
		}
		fmt.Println()
		Preorder(tree)
	}
}

func Test_insertBinary(t *testing.T) {
	type args struct {
		rt *RBNode
		n  *RBNode
	}

	var A, B, C, D, E, F, G Int = 1, 2, 3, 4, 5, 6, 7
	_, _, _, _, _, _, _ = A, B, C, D, E, F, G
	tests := []struct {
		name string
		args args
		want *RBNode
	}{
		{
			name: "test_father_red_uncle_black",
			args: args{
				rt: &RBNode{
					Color: ColorBlack,
					V:     D,
					Ch: [2]*RBNode{
						{
							Color: ColorRed, V: B,
							Ch: [2]*RBNode{},
						},
						{
							Color: ColorBlack, V: E,
							Ch: [2]*RBNode{},
						},
					},
				},
				n: &RBNode{
					Color: ColorRed,
					V:     C,
					Ch:    [2]*RBNode{},
				},
			},

			//       D,B
			// B,R      E,B
			//    C,R
			want: &RBNode{
				Color: ColorBlack,
				V:     D,
				Ch: [2]*RBNode{
					{
						Color: ColorRed, V: B,
						Ch: [2]*RBNode{
							nil,
							{
								Color: ColorRed,
								V:     C,
								Ch:    [2]*RBNode{},
							},
						},
					},
					{
						Color: ColorBlack, V: E,
						Ch: [2]*RBNode{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insertBinary(tt.args.rt, tt.args.n); !Equivalent(got, tt.want) {
				t.Errorf("insertBinary() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}