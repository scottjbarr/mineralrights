package mineralrights

type Validator func(int64) bool

func buildOreSellValidator(storage int64) Validator {
	return func(i int64) bool {
		return i >= 0 && i <= storage
	}
}

func buildMinesSellValidator(mines int64) Validator {
	return func(i int64) bool {
		return i >= 0 && i <= mines
	}
}
