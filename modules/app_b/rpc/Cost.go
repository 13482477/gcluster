package rpc

type Cost struct {
	Cost float64 `json:"cost"`
}

type EcsCostResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Cost   `json:"data"`
}
