package lambda

type Env = map[string]Term

func subst(term Term, x *Var, s Term) Term {
	switch t := term.(type) {
	case *Var:
		if *t == *x {
			return s
		}
		return term
	case *Abs:
		return &Abs{t.Var, subst(t.Body, x, s)}
	case *App:
		return &App{subst(t.Lterm, x, s), subst(t.Rterm, x, s)}
	default:
		panic("unreachable")
	}
}

func BetaIn(env Env, term Term) Term {
	switch t := term.(type) {
	case *Var:
		if f, ok := env[t.Symbol]; ok && t.ID == 0 {
			return BetaIn(env, AlphaTerm(f))
		}
		return &Var{t.Symbol, t.ID}
	case *Abs:
		return &Abs{t.Var, BetaIn(env, t.Body)}
	case *App:
		rterm := BetaIn(env, t.Rterm)
		lterm := BetaIn(env, t.Lterm)
		switch lt := lterm.(type) {
		case *Var, *App:
			return &App{lterm, rterm}
		case *Abs:
			return BetaIn(env, subst(lt.Body, lt.Var, rterm))
		default:
			panic("unreachable")
		}
	case *Def:
		env[t.Symbol] = t.Bound
	}
	return term
}

func Beta(terms []Term) []Term {
	beta := make([]Term, len(terms))
	env := make(Env)
	for idx, term := range terms {
		beta[idx] = BetaIn(env, term)
	}
	return beta
}
