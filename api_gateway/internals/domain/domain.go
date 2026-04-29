package domain

type OverAllStatus struct {
	Status      string             `json:"status"`
	AllServices []OneServiceStatus `json:"allServices"`
}

type OneServiceStatus struct {
	ServiceName   string    `json:"services_name"`
	Status        string    `json:"status"`
	PGRespTime    *float64 `json:"postgres_response_time_in_seconds"`
	RedisRespTime *float64 `json:"redis_response_time_in_seconds"`
}
