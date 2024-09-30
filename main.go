package main

import (
	"fmt"
	"server/dao"
	"server/router"
)

func main() {
	r := router.Router()
	// initi the gorm to connect the database
	dao.Init()
	sqlDB, err := dao.DB.DB()

	if err != nil {
		fmt.Println("failed", err)
	}

	if err := sqlDB.Ping(); err != nil {
		fmt.Println("failed", err)
	}

	fmt.Println("DB is running")
	// router run as the port number
	r.Run(":9999")
}
