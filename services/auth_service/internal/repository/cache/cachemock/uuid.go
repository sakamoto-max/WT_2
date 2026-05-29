package cachemock

import (
	"auth_service/internal/repository/mock"
	"context"
	"fmt"
)

type UuidMock struct {
	Hit bool
	Miss bool
	HasErr bool
}

// func NewUuidCacheMock(Hit bool, Miss bool, HasErr bool) *UuidMock {
// 	return &UuidMock{
// 		Hit: Hit,
// 		Miss: Miss,
// 		HasErr: HasErr,
// 	}
// }


var (
	UUID = mock.Uuid 
)

func (u *UuidMock) GetUUID(ctx context.Context, userId string) (string, error) {
	if u.Miss {
		return "", fmt.Errorf("please login first")
	}

	if u.HasErr {
		return "", fmt.Errorf("some error occured")
	}
	
	return UUID, nil
}
func (u *UuidMock) SetUUID(ctx context.Context, uuid string, userId string) error {
	if u.HasErr {
		return fmt.Errorf("some error occured")
	}

	return nil
}
