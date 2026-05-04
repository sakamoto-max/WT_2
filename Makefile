# creating servers :
plan:
	@cd services && cd plan_service && cd cmd && cd server && go run .
plan_consumer :
	@cd services && cd plan_service && cd internal && cd mq_consumer && go run main.go	
exercise:
	@cd services && cd exercise_service && cd cmd && cd server && go run .
auth:
	@cd services && cd auth_service && cd cmd && cd server && go run .

tracker:
	@cd services && cd tracker_service && cd cmd && cd server && go run .

gateway:
	@cd api_gateway && cd cmd && go run .
orc:
	@cd services && cd orchestration_service && cd cmd && go run main.go
email:
	@cd services && cd email_service && go run main.go

# migrations :
auth_db_up:
	@cd services && cd auth_service && migrate -database "postgres://postgres:root@localhost:5432/WT_AUTH?sslmode=disable" -path migrations up
auth_db_down:
	@cd services && cd auth_service && migrate -database "postgres://postgres:root@localhost:5432/WT_AUTH?sslmode=disable" -path migrations down
plan_db_up:
	@cd services && cd plan_service && migrate -database "postgres://postgres:root@localhost:5432/WT_PLAN?sslmode=disable" -path migrations up 
plan_db_down:
	@cd services && cd plan_service && migrate -database "postgres://postgres:root@localhost:5432/WT_PLAN?sslmode=disable" -path migrations down
tracker_db_up:
	@cd services && cd tracker_service && migrate -database "postgres://postgres:root@localhost:5432/WT_TRACKER?sslmode=disable" -path migrations up
tracker_db_down:
	@cd services && cd tracker_service && migrate -database "postgres://postgres:root@localhost:5432/WT_TRACKER?sslmode=disable" -path migrations down
exercise_db_up:
	@cd services && cd exercise_service && migrate -database "postgres://postgres:root@localhost:5432/WT_EXERCISES?sslmode=disable" -path migrations up
exercise_db_down:
	@cd services && cd exercise_service && migrate -database "postgres://postgres:root@localhost:5432/WT_EXERCISES?sslmode=disable" -path migrations down
