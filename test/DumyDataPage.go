package test

type DummyDataPage struct {
	Total *int64  `json:"total"`
	Data  []Dummy `json:"data"`
}

func NewEmptyDummyDataPage() *DummyDataPage {
	return &DummyDataPage{}
}

func NewDummyDataPage(total *int64, data []Dummy) *DummyDataPage {
	return &DummyDataPage{Total: total, Data: data}
}
