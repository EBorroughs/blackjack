package api

type upsertAction string

const (
	upsertCreate upsertAction = "create"
	upsertHit    upsertAction = "hit"
	upsertStand  upsertAction = "stand"
)

type upsertGameRequestBody struct {
	Action upsertAction `json:"action"`
}
