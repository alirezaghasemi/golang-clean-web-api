package dto

type CreateUpdateCountryRequest struct {
	Name string `json:"name" binding:"required,alpha,min=3,max=32"`
}

type CountryResponse struct {
	Id     int            `json:"id"`
	Name   string         `json:"name"`
	Cities []CityResponse `json:"cities"` /*هر کشور تعدادی شهر دارد*/
}

type CreateUpdateCityRequest struct {
	Name string `json:"name" binding:"required,alpha,min=3,max=32"`
}

type CityResponse struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
	Country CountryResponse `json:"country"` /*شهر شیراز برای کدا کشور هست*/
}
