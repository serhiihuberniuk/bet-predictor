package esmodels

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type Leagues struct {
	Data       []League   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type League struct {
	CountryName string `json:"countryName"`
	Name        string `json:"name"`
}

type Pagination struct {
	Page        int  `json:"page"`
	HasNextPage bool `json:"hasNextPage"`
}
