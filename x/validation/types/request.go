package types

type GetValidationsRequest struct {
	ExternalAddress string `json:"external_address,omitempty"`
}

func NewGetValidationsRequest(addr string) *GetValidationsRequest {
	return &GetValidationsRequest{
		ExternalAddress: addr,
	}
}
