package models

import (
	"strconv"
	"time"
)

type TokenModel interface {
	Create(userID uint32) error
}

type TokenDetails struct {
	AccessToken string
	AccessUUID  string
	RefreshUUID string
	AtExpires   int64
	RtExpires   int64
}

func (td *TokenDetails) Create(userID uint32) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := Client.Set(td.AccessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := Client.Set(td.RefreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
