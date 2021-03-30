package router

import (
	"encoding/gob"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

const (
	authenticationKey = "secret"
	redisPoolSize     = 32
	redisNetwork      = "tcp"
)

func createSessionStore(option *config.SessionOption) (sessions.Store, error) {
	gob.Register(&types.LoginUser{})
	gob.Register(&types.User{})
	gob.Register(&types.Tenant{})

	switch option.Type {
	case config.SessionStoreTypeRedis:
		return redis.NewStore(redisPoolSize, redisNetwork, option.Address, option.Password, []byte(authenticationKey))
	default:
		return memstore.NewStore([]byte(authenticationKey)), nil
	}
}

func setupSessionStore(router gin.IRouter, option *config.SessionOption) error {
	if option == nil {
		option = &config.SessionOption{}
	}

	store, err := createSessionStore(option)
	if err != nil {
		return err
	}

	router.Use(sessions.Sessions(handler.SessionID, store))
	return nil
}
