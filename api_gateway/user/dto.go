package user

type roles string

var (
	UserRole roles = "user"
	AdminRole roles = "admin"
	SuperAdminRole roles = "super_admin"
)