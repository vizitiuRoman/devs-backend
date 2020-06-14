package models

import (
	"errors"
	"strconv"
	"time"
)

type TokenModel interface {
	Create(userID uint32) error
	GetByUUID(usrID uint64) (uint64, error)
	DeleteByUUID() error
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

func (td *TokenDetails) GetByUUID(usrID uint64) (uint64, error) {
	userid, err := Client.Get(td.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	if usrID != userID {
		return 0, errors.New("Unauthorized")
	}
	return userID, nil
}

func (td *TokenDetails) DeleteByUUID() error {
	_, err := Client.Del(td.AccessUUID).Result()
	if err != nil {
		return err
	}
	_, err = Client.Del(td.RefreshUUID).Result()
	if err != nil {
		return err
	}
	return nil
}
