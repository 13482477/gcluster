package utils

type Page struct {
	Page     int32
	PageSize int32
	Total    int32
}

func NewPage() *Page {
	return &Page{1, 20, 0}
}

func (p *Page) SetFromRequestMap(req *map[string]interface{}) {
	r := *req
	if v, ok := r["page"]; ok {
		if v.(int32) > 0 {
			p.Page = v.(int32)
		}
		delete(r, "page")
	}
	if v, ok := r["page_size"]; ok {
		if v.(int32) > 0 {
			p.PageSize = v.(int32)
		}
		delete(r, "page_size")
	}
}

func (p *Page) SetTotal(total int32) {
	p.Total = total
}

func (p *Page) OffSet() int32 {
	return (p.Page - 1) * p.PageSize
}
