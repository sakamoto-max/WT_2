plan:
	@cd services && cd plan_service && cd cmd && go run .

exercise:
	@cd services && cd exercise_service && cd cmd && go run .

auth:
	@cd services && cd auth_service && cd cmd && go run .

tracker:
	@cd services && cd tracker_service && cd cmd && go run .

gateway:
	@cd api_gateway && go run .