package configs

import "fmt"

func PreServerStart() {
	GetRedis()
	GetLogger()
	GetDB()
}

func PreServerShutdown() {
	dbSQL, _ := GetDB().DB()
	if err := dbSQL.Close(); err != nil {
		fmt.Println("failed to close database connection", err)
	}

	if err := rdb.Close(); err != nil {
		fmt.Println("failed to close redis connection:", err)
	}
}
