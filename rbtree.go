package rbtree

import (
	"fmt"
)

const (
	ColorBlack = true
	ColorRed = false
)

type Lesser interface {
	Less(interface{}) bool
}

type Equaler interface {
	Equal(interface{}) bool
}

type LessEqualer interface {
	Lesser
	Equaler
}


type RBNode struct {
	Color bool
	V LessEqualer
	Prev, Next *RBNode
	Faz *RBNode
	Ch[2] *RBNode
}

func (r *RBNode) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf(`{color:%v, value: %v, child[]:{%v, %v}}`, r.Color, r.V, r.Ch[0], r.Ch[1])
}

func (r *RBNode) GetValue() interface{} {
	return r.V
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func (r *RBNode) Depth() int {
	if r == nil {
		return 0
	}
	return max(r.Ch[0].Depth(), r.Ch[1].Depth()) + 1
}

type trace struct {
	V interface{}
	Color string
}

func (r *RBNode) GetValueTrace() interface{} {
	if r.Red() {
		return trace{V:r.V, Color:"R"}
	} else if r != nil {
		return trace{V:r.V, Color:"B"}
	} else {
		return trace{Color:"B"}
	}
}

func (r *RBNode) Index(i int) interface{} {
	return r.Ch[i]
}

func (r *RBNode) SetFaz(faz *RBNode) *RBNode {
	if r == nil {
		return r
	}
	r.Faz = faz
	return r
}

func (r *RBNode) IsLeftChild() uint8 {
	if r.Faz.Ch[0] == r {
		return 1
	}
	return 0
}

func (r *RBNode) IsRightChild() uint8 {
	if r.Faz.Ch[1] == r {
		return 1
	}
	return 0
}

func (r *RBNode) Sibling() *RBNode {
	if r == nil {
		return nil
	}
	return r.Faz.Ch[1 - r.IsRightChild()]
}

func (r *RBNode) HasFaz() bool {
	return r.Faz != nil
}

func (r *RBNode) HasLeftChild() bool {
	return r.Ch[0] != nil
}

func (r *RBNode) HasRightChild() bool {
	return r.Ch[1] != nil
}

func (r *RBNode) Black() bool {
	if r == nil {
		return true
	}
	return r.Color
}

func (r *RBNode) Red() bool {
	if r == nil {
		return false
	}
	return !r.Color
}

func (r *RBNode) Rotate() *RBNode {
	if !r.HasFaz() {
		return r
	}
	faz, gra := r.Faz, r.Faz.Faz
	if r.IsRightChild() == 1 {
		if gra != nil {
			gra.Ch[faz.IsRightChild()] = r.SetFaz(gra)
		} else {
			r.Faz = nil
		}
		faz.Ch[1], r.Ch[0] = r.Ch[0].SetFaz(faz), faz.SetFaz(r)
	} else {
		if gra != nil {
			gra.Ch[faz.IsRightChild()] = r.SetFaz(gra)
		} else {
			r.Faz = nil
		}
		faz.Ch[0], r.Ch[1] = r.Ch[1].SetFaz(faz), faz.SetFaz(r)
	}
	return r
}

func find(r *RBNode, leer LessEqualer) *RBNode {
	if r == nil {
		return nil
	} else if r.V.Equal(leer) {
		return r
	} else if r.V.Less(leer) {
		return find(r.Ch[1], leer)
	} else {
		return find(r.Ch[0], leer)
	}
}


func (r *RBNode) Insert(leer LessEqualer) *RBNode {
	return insert(r, NewRBNode(leer))
}

func (r *RBNode) Find(leer LessEqualer) *RBNode {
	return find(r, leer)
}

func (r *RBNode) Prec() (res *RBNode) {
	if r.Ch[0] == nil {
		return nil
	}
	res = r.Ch[0]
	for res.HasRightChild() {
		res = res.Ch[1]
	}
	return
}

// fix
func fix(rt, r, sibling *RBNode) *RBNode {
	// must have sibling
	faz := sibling.Faz
	if faz == nil {
		sibling.Color = ColorBlack
		return sibling
	}

	if sibling.Red() {
		sibling.Color, faz.Color = ColorBlack, ColorRed
		sibling.Rotate()
		if sibling.Faz == nil {
			rt = sibling
		}
		if faz.IsLeftChild() == 1 {
			return fix(rt, r, faz.Ch[1])
		} else {
			return fix(rt, r, faz.Ch[0])
		}
	} else {
		if sibling.Ch[0].Black() && sibling.Ch[1].Black() {
			if faz.Red() {
				sibling.Color, faz.Color = ColorRed, ColorBlack
				return rt
			} else {
				sibling.Color = ColorRed
				if faz.Faz == nil {
					return faz
				}
				return fix(rt, faz, faz.Sibling())
			}
		} else {
			rch := sibling.IsRightChild()
			// sibling.Ch[0] must not be nil
			if sibling.Ch[rch ^ 1].Red() && sibling.Ch[rch].Black() {
				sibling.Ch[rch ^ 1].Color, sibling.Color = ColorBlack, ColorRed
				sibling = sibling.Ch[rch ^ 1].Rotate()
			}
			// sibling.Ch[1] must be red
			sibling.Color, faz.Color = faz.Color, sibling.Color
			sibling.Rotate()
			sibling.Ch[rch].Color = ColorBlack
			if sibling.Faz == nil {
				return sibling
			} else {
				return rt
			}
		}
	}
}

func deleteX(rt, r *RBNode) *RBNode {
	if r == nil {
		return rt
	}
	//    *
	//    |
	//    *
	//  /   \
	// *     *
	if r.HasLeftChild() && r.HasRightChild() {
		y := r.Prec()
		r.V, y.V = y.V, r.V
		return deleteX(rt, y)
	} else if r.HasLeftChild() || r.HasRightChild() || r.Color == ColorBlack {
		//   sA       sB       sC
		//    *        *        *
		//    |        |        |
		//    B        B        B
		//  /   \    /   \    /   \
		// B     *  *     B  B     B
		var ch *RBNode
		if r.HasLeftChild() {
			ch = r.Ch[0]
		} else {
			ch = r.Ch[1]
		}
		faz := r.Faz
		ch.SetFaz(faz)
		var sib *RBNode
		if faz != nil {
			sib = r.Sibling()
			faz.Ch[r.IsRightChild()] = ch
		} else {
			if ch != nil {
				ch.Color = ColorBlack
			}
			return ch
		}
		if ch.Red() {
			ch.Color = ColorBlack
			return rt
		} else {
			return fix(rt, ch, sib)
		}
	} else {
		//    *
		//    |
		//    R
		//  /   \
		// n     n
		faz := r.Faz
		if faz != nil {
			faz.Ch[r.IsRightChild()] = nil
		}
		return rt
	}
}

func (r *RBNode) Delete(leer LessEqualer) *RBNode {
	return deleteX(r, find(r, leer))
}


func insertBinary(rt, n *RBNode) *RBNode {
	if rt == nil {
		return n
	}
	if rt.V.Less(n.V) {
		rt.Ch[1] = insertBinary(rt.Ch[1], n).SetFaz(rt)
	} else {
		rt.Ch[0] = insertBinary(rt.Ch[0], n).SetFaz(rt)
	}
	return rt
}

func proc(n *RBNode) *RBNode {
	if n.Faz == nil {
		n.Color = ColorBlack
		return n
	}
	// where faz must be black
	if n.Faz.Faz == nil {
		return n.Faz
	}
	faz, gra := n.Faz, n.Faz.Faz
	unc := faz.Sibling()

	if unc.Red() {
		unc.Color, faz.Color, gra.Color = ColorBlack, ColorBlack, ColorRed
		if gra.Faz == nil {
			gra.Color = ColorBlack
			return gra
		}
		if gra.Faz.Black() {
			return nil
		}
		return proc(gra)
	} else {
		if n.IsRightChild() == faz.IsRightChild() {
			faz.Color, gra.Color = gra.Color, faz.Color
			faz.Rotate()
			if faz.HasFaz() {
				return nil
			} else {
				return faz
			}
		} else {
			n.Color, gra.Color = gra.Color, n.Color
			n.Rotate().Rotate()
			if n.HasFaz() {
				return nil
			} else {
				return n
			}
		}
	}
}

func insert(rt, n *RBNode) *RBNode {
	if rt == nil {
		n.Color = ColorBlack
		return n
	}
	rt = insertBinary(rt, n)
	if n.Faz.Black() {
		return rt
	}
	n = proc(n)
	if n == nil {
		return rt
	} else {
		return n
	}
}

func NewRBNode(leer LessEqualer) *RBNode {
	return &RBNode{
		V:     leer,
	}
}

