package hotel

import (
	"bytes"
	"context"
	"crypto/sha256"
	"github.com/eniac/mucache/pkg/state"
	"github.com/lithammer/shortuuid"
)

func RegisterUser(ctx context.Context, username string, password string) bool {
	userId := shortuuid.New()
	salt := shortuuid.New()
	hashPass := hash(password + salt)
	user := User{
		UserId:   userId,
		Username: username,
		Password: hashPass,
		Salt:     salt,
	}

	state.SetState(ctx, username, user)
	return true
}

func hash(str string) []byte {
	h := sha256.New()
	h.Write([]byte(str))
	val := h.Sum(nil)
	return val
}

func Login(ctx context.Context, username string, password string) string {
	user, err := state.GetState[User](ctx, username)
	if err != nil {
		panic(err)
	}
	salt := user.Salt
	givenPass := hash(password + salt)
	if bytes.Equal(givenPass, user.Password) {
		return "OK"
	}
	return "NOT-OK"
}

func GetUserId(ctx context.Context, username string) string {
	user, err := state.GetState[User](ctx, username)
	if err != nil {
		panic(err)
	}
	return user.UserId
}
