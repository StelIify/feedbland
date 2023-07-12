package validator

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	// Check that the sort parameter matches a value in the safelist.
	v.Check(v.In(f.Sort, f.SortSafelist), "sort", "invalid sort value")
}

func (f *Filters) Limit() int {
	return f.PageSize
}
func (f *Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}
