package privatebin

type Response struct {
	Id    string `json:"id"`
	Adata []any  `json:"adata"`
	Ct    string `json:"ct"`
}
