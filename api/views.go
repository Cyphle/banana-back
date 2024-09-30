package api

type CreateAccountRequest struct {
	Name string `json:"name"`
}

type UpdateAccountCommandView struct {
	ID   int64  `param:"id" json:"id"`
	Name string `json:"name"`
}

type AccountIdPathParam struct {
	ID int64 `param:"id" query:"id"`
}

// TODO il faut ajouter les autres champs
type AccountView struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ArrayResponse[T any] struct {
	Data []AccountView `json:"data"`
}
