package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const IDTokenSubjectContextKey = "ctx-username-key"
const IDTokenSubjectContextErrKey = "ctx-username-err"

type UnauthorizedError struct {
	Err error
}

func (e *UnauthorizedError) Error() string {
	return e.Err.Error()
}

func GetRequestUsername(c *gin.Context) (username string, err error) {
	username, hasCache := c.Keys[IDTokenSubjectContextKey].(string)
	if hasCache {
		err, _ = c.Keys[IDTokenSubjectContextErrKey].(error)
	} else {
		// init keys
		if c.Keys == nil {
			c.Keys = map[string]interface{}{}
		}

		// get username from api
		if false {
			username, err = "", nil
		} else {
			username, err = "", &UnauthorizedError{fmt.Errorf("unauth error")}
		}
		if err != nil {
			if _, ok := err.(*UnauthorizedError); ok {
				c.Keys[IDTokenSubjectContextErrKey] = err
				c.Keys[IDTokenSubjectContextKey] = username
			}
			return "", err
		}
		c.Keys[IDTokenSubjectContextErrKey] = err
		c.Keys[IDTokenSubjectContextKey] = username
	}
	return username, err
}
