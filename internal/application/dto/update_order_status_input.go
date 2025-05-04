package dto

type UpdateOrderStatusInput struct {
	ExternalRef string `json:"external_reference"`
	Status      string `json:"status"`
}
