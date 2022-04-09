package data

import (
	"errors"
	"fmt"
	"time"
)

type Access struct {
	Path        string    `db:"path"`
	Key         string    `db:"key"`
	UserCode    string    `db:"user_code"`
	DisplayName string    `db:"display_name"`
	Until       time.Time `db:"until"`
	Created     time.Time `db:"created"`
}

func (a Access) Can(code string) (bool, error) {
	if a.UserCode != code {
		return false, errors.New("invalid code for resource")
	}
	if time.Now().After(a.Until) {
		return false, fmt.Errorf("time failed - now = %s, expire = %s",
			time.Now().String(), a.Until.String())
	}
	return true, nil
}
