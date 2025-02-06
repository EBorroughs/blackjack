package api

type UpsertAction string

const (
	UpsertCreate UpsertAction = "create"
	UpsertHit    UpsertAction = "hit"
	UpsertStand  UpsertAction = "stand"
)

type UpsertGameRequestBody struct {
	Action UpsertAction `json:"action"`
}
