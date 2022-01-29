package serverDomain

type ZentaoResp struct {
	Status string
	Data   string
}

type ZentaoLang struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ZentaoProduct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ZentaoModule struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ZentaoSuite struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ZentaoTask struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ZentaoResult struct {
	Seq       string `json:"seq"`
	ProductId string `json:"productId"`
	TaskId    string `json:"taskId"`
}
