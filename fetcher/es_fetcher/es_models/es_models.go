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
	Expand      Expand `json:"expand"`
}

type Expand struct {
	CurrentSeason []Season `json:"current_season"`
}

type Pagination struct {
	Page        int  `json:"page"`
	HasNextPage bool `json:"hasNextPage"`
}

type Season struct {
	ID int `json:"id"`
}

type Teams struct {
	Data       []Team     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Team struct {
	FullName string `json:"fullName"`
	Country  string `json:"country"`
}
