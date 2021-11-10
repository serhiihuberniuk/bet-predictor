package esmodels

//"access_token": "eyJraWQiOiJtaGdzWDNGZ1wvd0kyRzR3Z1JQZ2FJTmJmbjllQTFiWkszcnlnOUgzWFZoMD0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI0NWxrYjRwZ3JzcWQ4aTIwdXJlNDl0dnFpYiIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiZm9vdGJhbGwtZGVidWcuZWxlbmFzcG9ydC5pb1wvcmVhZDoqIiwiYXV0aF90aW1lIjoxNTk0Mzg4Mzc4LCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuZXUtd2VzdC0xLmFtYXpvbmF3cy5jb21cL2V1LXdlc3QtMV9sdFdMNzVoYnEiLCJleHAiOjE1OTQzOTE5NzgsImlhdCI6MTU5NDM4ODM3OCwidmVyc2lvbiI6MiwianRpIjoiMmIxNjczYzQtY2NkOS00MGMzLTliY2UtZDQxOGJlYmJhYWRmIiwiY2xpZW50X2lkIjoiNDVsa2I0cGdyc3FkOGkyMHVyZTQ5dHZxaWIifQ.L-QerQB_VcseWnNj9yQaQTgt_PpJxK0Qye1ufyTzEHcmA_vmM5aQkDkCTfYPHYvWvryLRvhf0iDpN3Pg2jdtmC1JzD8zIChpcIv02qY1O2ONU1oso97rQOipmE3Bcz_91TOtkz71dbrkQUv-LS1cEO6CdsKPM_KACtex7ysSVYbpuPQADrGtA9-4qLhhczOz5X0clUszH1AEvuSEdwgSt7hRYODFniUWgX0gJwKtFJjbt2Ucpy-bvLABmt-pih99QsYsgprztY7XK3MXJFgfcmgcrFhjA35XiJIbJSKLbYft3tcuMGxictrARWITpQyU_wfQ5haRgUOa7__kaW61SQ",
//  "expires_in": 3600,
//  "token_type": "Bearer"

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	//ExpiresIn int `json:"expires_in"`
	//TokenType string `json:"token_type"`
}

/*
{ "data": [
	{
      "id": 215,
      "idCountry": 21,
      "countryName": "Cyprus",
      "name": "1. Division",
      "nationalLeague": true,
      "clubsLeague": true
    },
  ],
  "pagination": {
    "page": 1,
    "itemsPerPage": 20,
    "hasNextPage": true,
    "hasPrevPage": false
  }
}
*/

type Leagues struct {
	Data       []League   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type League struct {
	//ID int `json:"id"`
	//IDCountry int `json:"idCountry"`
	CountryName string `json:"countryName"`
	Name        string `json:"name"`
	//NationalLeague bool `json:"nationalLeague"`
	//ClubsLeague bool `json:"clubsLeague"`
}

type Pagination struct {
	Page int `json:"page"`
	//ItemsPerPage int `json:"itemsPerPage"`
	HasNextPage bool `json:"hasNextPage"`
	//HasPrevPage bool `json:"hasPrevPage"`
}
