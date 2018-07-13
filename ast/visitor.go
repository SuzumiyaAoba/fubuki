package ast

type Visitor interface {
	VisitTopdown(e Expr) Visitor
	VisitBottomup(e Expr)
}

func Visit(vis Visitor, e Expr) {
	v := vis.VisitTopdown(e)
	if v == nil {
		return
	}

	switch n := e.(type) {
	case *Abs:
		Visit(v, n.Body)
	case *App:
		Visit(v, n.Lexp)
		Visit(v, n.Rexp)
	case *Def:
		Visit(v, n.Bound)
	}
	vis.VisitBottomup(e)
}
