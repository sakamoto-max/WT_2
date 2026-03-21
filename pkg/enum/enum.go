package enum

type serviceName string

var (
	PlanService serviceName = "plan_service"
	AuthService serviceName = "auth_service"
	TrackerService serviceName = "tracker_service"
	ExerciseService serviceName = "exercise_service"
)

type taskStatus string

var (
	TaskCompleted taskStatus = "completed"
	TaskNotCompleted taskStatus = "not_completed"
)

type taskName string

var(
	CreateEmptyPlanForUser taskName = "create_empty_plan_for_user"
)


type serviceStatus string

var (
	Healthy serviceStatus = "healthy"
	UnHealthy serviceStatus = "unhealthy"
	Unavailable serviceStatus = "unavailable"
	PostgresIsNotUp serviceStatus = "postgres is not up"
	RedisIsNotUp serviceStatus = "Redis is not up"
)


type roles string

var (
	UserRole       roles = "user"
	AdminRole      roles = "admin"
	SuperAdminRole roles = "super_admin"
)


