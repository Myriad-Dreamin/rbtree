package rbtree

import "fmt"

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
	//fmt.Println(&r, &r.Ch[0], &r.Ch[1])
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf(`{color:%v, child[]:{%v, %v}}`, r.Color, r.Ch[0], r.Ch[1])
}

func (r *RBNode) GetValue() interface{} {
	return r.V
}

func (r *RBNode) Index(i int) Node {
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

func (r *RBNode) Uncle() *RBNode {
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

func deleteX(rt, r *RBNode) *RBNode {
	if r.HasLeftChild() && r.HasRightChild() {
		y := r.Prec()
		r.V, y.V = y.V, r.V
		return deleteX(rt, y)
	} else if r.HasLeftChild() || r.HasRightChild() {
		var ch *RBNode
		if r.HasLeftChild() {
			ch = r.Ch[0]
			r.Ch[0] = nil
		} else {
			ch = r.Ch[1]
			r.Ch[1] = nil
		}
		if r.Black() {
			ch.Color = ColorBlack
		}
		faz := r.Faz
		if r.Faz != nil {
			faz.Ch[r.IsRightChild()] = ch.SetFaz(faz)
			r.Faz = nil
			return rt
		} else {
			ch.Faz = nil
			return ch
		}
	} else {
		if r.Faz == nil {
			return nil
		}
		unc, faz := r.Uncle(), r.Faz
		faz.Ch[r.IsRightChild()] = nil
		r.Faz = nil

		if unc == nil {
			if rt == r {
				return nil
			}
			return rt
		} else {
			if unc.HasLeftChild() && unc.HasRightChild() {
				unc.Ch[0].Color, unc.Ch[1].Color, unc.Color, faz.Color = ColorRed, ColorBlack, ColorBlack, ColorBlack
				unc.Rotate()
				if unc.Faz == nil {
					return unc
				}
				return rt
			} else if unc.HasLeftChild() || unc.HasRightChild() {
				var nz *RBNode
				if unc.HasLeftChild() {
					nz = unc.Ch[0]
				} else {
					nz = unc.Ch[1]
				}
				if nz.IsRightChild() == unc.IsRightChild() {
					faz.Color = ColorRed
					unc.Rotate()
					if unc.Faz == nil {
						return unc
					}
					return rt
				} else {
					faz.Color, unc.Color, nz.Color = ColorRed, ColorRed, ColorBlack
					nz.Rotate().Rotate()
					if nz.Faz == nil {
						return nz
					}
					return rt
				}
			} else {
				faz.Color, unc.Color = ColorBlack, ColorRed
				return rt
			}
		}
	}
}

func deleteN(rt *RBNode, leer LessEqualer) *RBNode {
	r := find(rt, leer)
	if r == nil {
		return rt
	}
	return deleteX(rt, r)
}

func (r *RBNode) Delete(leer LessEqualer) *RBNode {
	return deleteN(r, leer)
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
	unc := faz.Uncle()

	if unc.Red() {
		unc.Color, faz.Color, gra.Color = ColorBlack, ColorBlack, ColorRed
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


func main() {
}

