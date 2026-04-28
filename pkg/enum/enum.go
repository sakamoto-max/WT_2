package enum

type planName string

var (
	EmptyPlanName planName = "empty"
)

type serviceName string

var (
	PlanService serviceName = "plan_service"
	AuthService serviceName = "auth_service"
	TrackerService serviceName = "tracker_service"
	ExerciseService serviceName = "exercise_service"
	EmailService serviceName = "email_service"
)

type taskStatus string

var (
	TaskCompleted taskStatus = "completed"
	TaskPending taskStatus = "pending"
	TaskNotCompleted taskStatus = "not_completed"
	TaskFailed taskStatus = "failed"
)

type taskName string

var(
	CreateEmptyPlanForUser taskName = "create_empty_plan_for_user"
	SendEmailforSigningUp taskName = "send_welcome_email_to_user_after_signingup"
	UpdatePlan taskName = "update_plan_for_user"
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

type queueName string

var (
	PlanQueue queueName = "plan"
	ResultQueue queueName = "result"
	EmailQueue queueName = "email"
)


type contentTypes string

var (
	ApplicationJsonType contentTypes = "application/json"
)

type correlationId string

var(
	EmptyPlanCrrId correlationId = "empty_plan"
	UpdatePlanCrrId correlationId = "update_plan"
)

type resources string

var (
	ExerciseResource resources = "exercise"
	BodyPartResource resources = "body_part"
	EquipmentResource resources = "equipment"
	PlanResource resources = "plan"
	UserResource resources = "user"
	EmailResource resources = "email"
)