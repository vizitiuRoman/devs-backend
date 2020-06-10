package migrations

import (
	. "github.com/devsmd/pkg/db"
	. "github.com/devsmd/pkg/db/models"
)

func init() {
	DB.AutoMigrate(&User{}, Token{})
}
