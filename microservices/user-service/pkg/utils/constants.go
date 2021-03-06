package utils

import "time"

const (
	UserID       = "userID"
	AccessUUID   = "accessUUID"
	RefreshUUID  = "refreshUUID"
	TokenExpires = time.Hour * 12
	AtExpires    = time.Hour * 12
	RtExpires    = time.Hour * 24 * 7
	LOGIN        = "login"
	UPDATE       = "update"
	DEFAULT      = ""
)
