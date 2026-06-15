package mappings

import (
	"testing"

	pb "github.com/sakamoto-max/wt_2_proto/shared/auth"
	"github.com/stretchr/testify/assert"
)

func Test_ToUserSignUp(t *testing.T) {
	in := pb.UserSignUpReq{
		Name: "jon snow",
		Email: "jonsnow@gmail.com",
		Password: "king in the north",
		Role: "user",
	}

	resp := ToUserSignUp(&in)
	assert.Equal(t, resp.Name, in.Name)
	assert.Equal(t, resp.Email, in.Email)
	assert.Equal(t, resp.Password, in.Password)
	assert.Equal(t, resp.Role, in.Role)
}