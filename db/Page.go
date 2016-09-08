package db

type Page struct{
	Page int
	PageSize int
	TotalRows int
	List interface{}
}

func (p Page) PageCount() int{
	return p.TotalRows / p.PageSize + ((p.TotalRows % p.PageSize + 1) / (p.TotalRows % p.PageSize + 1))
}
