package lambda

type Bind map[string]int

func find(env Bind, x string) int {
	if id, ok := env[x]; ok {
		return id
	}
	return 0
}

func AlphaIn(env Bind, term Term) Term {
	switch t := term.(type) {
	case *Var:
		return &Var{t.Symbol, find(env, t.Symbol)}
	case *Abs:
		ID := t.Var.ID
		SetID(t.Var)
		env[t.Var.Symbol] = t.Var.ID

		ret := &Abs{&Var{t.Var.Symbol, t.Var.ID}, AlphaIn(env, t.Body)}

		t.Var.ID = ID
		env[t.Var.Symbol] = ID
		return ret
	case *App:
		return &App{AlphaIn(env, t.Lterm), AlphaIn(env, t.Rterm)}
	case *Def:
		return &Def{t.Symbol, AlphaIn(env, t.Bound)}
	}
	panic("unknown term")
}

func AlphaTerm(term Term) Term {
	return AlphaIn(make(Bind), term)
}

func Alpha(terms []Term) []Term {
	alpha := make([]Term, len(terms))
	for idx, term := range terms {
		alpha[idx] = AlphaTerm(term)
	}

	return alpha
}
