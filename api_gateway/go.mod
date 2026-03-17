module api_gateway

go 1.25.4

require (
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/chi/v5 v5.2.5
	workout-tracker/proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.79.2
	google.golang.org/protobuf v1.36.11 // indirect
	wt/pkg v0.0.0
)

replace workout-tracker/proto => ../proto

replace wt/pkg => ../pkg
