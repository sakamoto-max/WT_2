package mappings

import pb "github.com/sakamoto-max/wt_2_proto/shared/auth"

type SignUp struct {
	Name     string `validate:"required"`
	Email    string `validate:"required"`
	Password string `validate:"required"`
	Role     string `validate:"required"`
}

func ToUserSignUp(in *pb.UserSignUpReq) SignUp {
	return SignUp{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Role:     in.Role,
	}
}
