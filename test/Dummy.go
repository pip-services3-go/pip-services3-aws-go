package test

type Dummy struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

func NewDummy(id string, key string, content string) *Dummy {
	return &Dummy{
		Id:      id,
		Key:     key,
		Content: content,
	}
}
