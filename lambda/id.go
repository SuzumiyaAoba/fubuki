package lambda

var uniqueId int

func SetID(x *Var) *Var {
	uniqueId++
	x.ID = uniqueId
	return x
}
