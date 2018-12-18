package http

type EcsGetAgentStatusRequest struct {
	EcsId int `json:"ecs_id"`
}

type Status struct {
	Status int `json:"status"`
}

