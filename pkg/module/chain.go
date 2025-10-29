package module

type Chain struct {
	Name string `json:"name"`
	Next *Chain `json:"next"`
}
