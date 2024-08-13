package configs

func Bootstrap() {
	GetRedis()
	GetLogger()
	GetDB()
}
