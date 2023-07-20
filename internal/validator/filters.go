package validator

func ValidateFilters(v *Validator, offset, limit int, sort string, sortSafeList []string) {
	v.Check(offset >= 0, "offset", "must be zero or greater")
	v.Check(offset <= 10_000_000, "offset", "must be a maximum of 10 million")
	v.Check(limit > 0, "limit", "must be greater than zero")
	v.Check(limit <= 100, "limit", "must be a maximum of 100")
	// Check that the sort parameter matches a value in the safelist.
	v.Check(v.In(sort, sortSafeList), "sort", "invalid sort value")
}
