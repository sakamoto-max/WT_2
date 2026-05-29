package mappings

import pb "github.com/sakamoto-max/wt_2_proto/shared/auth"

type SignUp struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func ToUserSignUp(in *pb.UserSignUpReq) SignUp {
	return SignUp{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Role:     in.Role,
	}
}
