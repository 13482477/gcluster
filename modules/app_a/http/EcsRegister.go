package http

type EcsRegisterRequest struct {
	CheckUrl     string `json:"url"`
	CheckUrlSign string `json:"sign"`
}
