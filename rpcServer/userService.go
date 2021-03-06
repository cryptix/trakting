package rpcServer

import (
	"github.com/cryptix/trakting/types"
	"gopkg.in/errgo.v1"
)

type UserService struct {
	user types.User
	db   types.Userer
}

func NewUserService(user types.User, db types.Userer) (*UserService, error) {
	return &UserService{
		user: user,
		db:   db,
	}, nil
}

func (us *UserService) Add(args *types.ArgAddUser, _ *struct{}) error {
	return us.db.Add(args.Name, args.Passw, args.Level)
}

// ChangePassword(id int64, newpw string) error
func (us *UserService) ChangePassword(args *types.ArgChangePassword, _ *struct{}) error {
	if us.user.ID != args.ID { // TODO: admin level override
		return errgo.New("can only change your own password")
	}
	return us.db.ChangePassword(args.ID, args.Passw)
}

func (us *UserService) Current(args *string, reply *types.User) error {
	*reply = us.user
	reply.PwHash = nil
	return nil
}

func (us *UserService) Check(args *types.ArgAddUser, _ *struct{}) error {
	_, err := us.db.Check(args.Name, args.Passw)
	return err
}

func (us *UserService) List(args *string, reply *[]types.User) (err error) {
	*reply, err = us.db.List()
	return
}
