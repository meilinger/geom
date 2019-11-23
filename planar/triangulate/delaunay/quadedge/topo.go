package quadedge

import (
	"log"
	"math"

	"github.com/go-spatial/geom/winding"

	"github.com/go-spatial/geom/planar"

	"github.com/go-spatial/geom"
)

// Splice operator affects the two edge rings around the origin of a and b,
// and, independently, the two edge rings around the left faces of a and b.
// In each case, (i) if the two rings are distinct, Splace will combine
// them into one; (ii) if the two are the same ring, Splice will break it
// into two separate pieces.
// Thus, Splice can be used both to attach the two edges together, and
// to break them apart. See Guibas and Stolfi (1985) p.96 for more details
// and illustrations.
func Splice(a, b *Edge) {
	if a == nil || b == nil {
		return
	}
	alpha := a.ONext().Rot()
	beta := b.ONext().Rot()

	t1 := b.ONext()
	t2 := a.ONext()
	t3 := beta.ONext()
	t4 := alpha.ONext()

	a.next = t1
	b.next = t2
	alpha.next = t3
	beta.next = t4
}

// Connect Adds a new edge (e) connecting the destination of a to the
// origin of b, in such a way that all three have the same
// left face after the connection is complete.
// Additionally, the data pointers of the new edge are set.
func Connect(a, b *Edge, order winding.Order) *Edge {
	if b == nil || a == nil {
		return nil
	}
	if debug {
		log.Printf("\n\n\tvalidate a:\n%v\n", a.DumpAllEdges())
		if err := Validate(a, order); err != nil {
			if err1, ok := err.(ErrInvalid); ok {
				for i, estr := range err1 {
					log.Printf("err: %03v : %v", i, estr)
				}
			}
		}
		log.Printf("\n\n\tvalidate b:\n%v\n", b.DumpAllEdges())
		if err := Validate(b, order); err != nil {
			if err1, ok := err.(ErrInvalid); ok {
				for i, estr := range err1 {
					log.Printf("err: %03v : %v", i, estr)
				}
			}
		}
		log.Printf("-------------------------\n")
	}
	if debug {
		log.Printf("\n\n\tConnect\n\n")
		log.Printf("Connecting %v to %v:", wkt.MustEncode(*a.Dest()), wkt.MustEncode(*b.Orig()))
	}
	bb, err := ResolveEdge(order, b, *a.Dest())
	if debug {
		if err != nil {
			panic(err)
		}
		log.Printf("splice e.Sym, bb: bb: %v", wkt.MustEncode(bb.AsLine()))
	}
	e := NewWithEndPoints(a.Dest(), bb.Orig())
	if debug {
		log.Printf("a: %v", wkt.MustEncode(a.AsLine()))
		log.Printf("a:LNext(): %v", wkt.MustEncode(a.LNext().AsLine()))
		log.Printf("a:LPrev(): %v", wkt.MustEncode(a.LPrev().AsLine()))
		log.Printf("splice e, a:LNext(): e: %v", wkt.MustEncode(e.AsLine()))
		log.Printf("splice e.Sym, b: b: %v", wkt.MustEncode(b.AsLine()))
	}

	Splice(e, a.LNext())
	Splice(e.Sym(), bb)
	if debug {
		log.Printf("\n\n\tvalidate e:\n%v\n", e.DumpAllEdges())
		if err := Validate(e, order); err != nil {
			if err1, ok := err.(ErrInvalid); ok {
				for i, estr := range err1 {
					log.Printf("err: %03v : %v", i, estr)
				}
			}
			log.Printf("Vertex Edges: %v", e.DumpAllEdges())
		}
		log.Printf("\n\n\tvalidate a:\n%v\n", a.DumpAllEdges())
		if err := Validate(a, order); err != nil {
			if err1, ok := err.(ErrInvalid); ok {
				for i, estr := range err1 {
					log.Printf("err: %03v : %v", i, estr)
				}
			}
			log.Printf("Vertex Edges: %v", e.DumpAllEdges())
		}
		log.Printf("\n\n\tvalidate b:\n%v\n", b.DumpAllEdges())
		if err := Validate(b, order); err != nil {
			if err1, ok := err.(ErrInvalid); ok {
				for i, estr := range err1 {
					log.Printf("err: %03v : %v", i, estr)
				}
			}
			log.Printf("Vertex Edges: %v", e.DumpAllEdges())
			panic("invalid edge b")
		}
		log.Printf("-------------------------\n")
	}
	return e
}

// Swap Essentially turns edge e counterclockwase inside its enclosing
// quadrilateral. The data pointers are modified accordingly.
func Swap(e *Edge) {
	a := e.OPrev()
	b := e.Sym().OPrev()
	Splice(e, a)
	Splice(e.Sym(), b)
	Splice(e, a.LNext())
	Splice(e.Sym(), b.LNext())
	e.EndPoints(a.Dest(), b.Dest())
}

// Delete will remove the edge from the ring
func Delete(e *Edge) {
	if e == nil {
		return
	}
	if debug {
		log.Printf("Deleting edge %p", e)
	}
	sym := e.Sym()

	Splice(e, e.OPrev())
	Splice(sym, sym.OPrev())
}

// OnEdge determines if the point x is on the edge e.
func OnEdge(pt geom.Point, e *Edge) bool {
	org := e.Orig()
	if org == nil {
		return false
	}
	dst := e.Dest()
	if dst == nil {
		return false
	}
	l := geom.Line{*org, *dst}
	return planar.IsPointOnLineSegment(pt, l)
}

// RightOf indicates if the point is right of the Edge
// If a point is below the line it is to it's right
// If a point is above the line it is to it's left
func RightOf(yflip bool, x geom.Point, e *Edge) bool {
	mul := 1.0
	if yflip {
		mul = -1.0
	}
	org := e.Orig()
	if org == nil {
		return false
	}
	dst := e.Dest()
	if dst == nil {
		return false
	}
	return sign(x.Subtract(*org).CrossProduct(dst.Subtract(*org))) == mul
}

func sign(f float64) float64 {
	if cmp.Float(f, 0.0) {
		return 0.0
	}
	if math.Signbit(f) {
		return -1.0
	}
	return 1.0
}
