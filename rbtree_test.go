package rbtree

import (
	"fmt"
	"math/rand"
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

func checkBST(rt *RBNode) bool {
	if rt == nil {
		return true
	}
	if rt.Ch[0] != nil {
		if rt.V.Less(rt.Ch[0].V) {
			return false
		}
		if !checkBST(rt.Ch[0]) {
			return false
		}
	}
	if rt.Ch[1] != nil {
		if rt.Ch[1].V.Less(rt.V) {
			return false
		}
		if !checkBST(rt.Ch[1]) {
			return false
		}
	}
	return true
}

type rbCheckContext struct {
	cnt int
	ccnt int
}

func (ctx *rbCheckContext) _checkRBProperty(rt *RBNode) {
	if rt == nil {
		if ctx.cnt == -1 {
			ctx.cnt = ctx.ccnt
		} else if ctx.ccnt != ctx.cnt {
			panic("bad count property")
		}
		return
	}
	if rt.Color == ColorBlack {
		ctx.ccnt++
	}

	if rt.HasLeftChild() {
		if rt.Color == ColorRed && rt.Ch[0].Color == ColorRed {
			panic(fmt.Errorf("bad color property %v", rt))
		}
	}
	ctx._checkRBProperty(rt.Ch[0])

	if rt.HasRightChild() {
		if rt.Color == ColorRed && rt.Ch[1].Color == ColorRed {
			panic(fmt.Errorf("bad color property %v", rt))
		}
	}
	ctx._checkRBProperty(rt.Ch[1])

	if rt.Color == ColorBlack {
		ctx.ccnt--
	}
}

func (ctx *rbCheckContext) checkRBProperty(rt *RBNode) (good bool) {
	ctx.cnt = -1
	ctx.ccnt = 0
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			good = false
		} else {
			good = true
		}
	}()
	ctx._checkRBProperty(rt)
	return
}

func assertRBTree(tree *RBNode) {
	if !checkRB(tree) {
		PrintL(tree)
		panic("bad")
	}
}


func checkRB(rt *RBNode) bool {
	if !checkFather(rt) {
		fmt.Println("err father")
	}
	if !checkBST(rt) {
		fmt.Println("err bst")
	}
	return (&rbCheckContext{}).checkRBProperty(rt)
}


func TestRBTree(t *testing.T) {
	var tree *RBNode

	for i := 0; i < 1000; i++ {
		is := rand.Perm(100)
		for j := range is {
			tree = tree.Insert(Int(is[j]))
			assertRBTree(tree)
		}
		for j := range is {
			tree = tree.Delete(Int(is[j]))
			assertRBTree(tree)
		}
		tree = tree.Delete(Int(233))
		assertRBTree(tree)
		if tree != nil {
			panic("bad delete")
		}
	}
}



func TestDepth(t *testing.T) {
	var tree *RBNode
	for j := 0; j < 10000000;j++ {
		tree = tree.Insert(Int(j))
	}
	fmt.Println(tree.Depth())
}

func BenchmarkRBNode_Insert1e6(b *testing.B) {
	for i := 0; i < b.N;i++ {
		var tree *RBNode
		for j := 0; j < 1000000;j++ {
			tree = tree.Insert(Int(i))
		}
	}
}

func BenchmarkRBNode_Insert1e7(b *testing.B) {
	for i := 0; i < b.N;i++ {
		var tree *RBNode
		for j := 0; j < 10000000;j++ {
			tree = tree.Insert(Int(i))
		}
	}
}


func BenchmarkRBNode_InsertDelete1e6(b *testing.B) {
	for i := 0; i < b.N;i++ {
		var tree *RBNode
		for j := 0; j < 1000000;j++ {
			tree = tree.Insert(Int(i))
		}
		for j := 0; j < 1000000;j++ {
			tree = tree.Delete(Int(i))
		}
	}
}


func BenchmarkRBNode_RandomBehavior1e6(b *testing.B) {
	var rnd = 233333
	var values []Int
	for i := 0; i < b.N;i++ {
		var tree *RBNode
		for j := 0; j < 1000000;j++ {
			rnd = rnd * 2323129414 % 1000000007
			if (rnd & 3) == 0 || (rnd & 3) == 1 {
				tree = tree.Insert(Int(rnd))
				values = append(values, Int(rnd))
			} else if (rnd & 3) == 2 {
				tree = tree.Delete(values[rnd % len(values)])
			} else {
				_ = tree.Find(values[rnd % len(values)])
			}
		}
	}
}

func BenchmarkRBNode_RandomBehaviorPure1e6(b *testing.B) {
	var rnd = 233333
	var values []Int
	for i := 0; i < b.N;i++ {
		for j := 0; j < 1000000;j++ {
			rnd = rnd * 2323129414 % 1000000007
			if (rnd & 3) == 0 || (rnd & 3) == 1 {
				values = append(values, Int(rnd))
			} else if (rnd & 3) == 2 {
				_ = values[rnd % len(values)]
			} else {
				_ = values[rnd % len(values)]
			}
		}
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