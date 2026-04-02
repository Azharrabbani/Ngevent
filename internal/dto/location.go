package dto

type ReverseResponse struct {
	DisplayName string         `json:"display_name"`
	Address     ReverseAddress `json:"address"`
}

type ReverseAddress struct {
	City    string `json:"city"`
	Town    string `json:"town"`
	Village string `json:"village"`
	Country string `json:"country"`
}

type LocationResp struct {
	Coordinates *string
	Address     *string
	City        *string
	Country     *string
	Err         *error
}
