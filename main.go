package main

import (
	"github/NjukiG/ngumbaru-school/initializers"
	"github/NjukiG/ngumbaru-school/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)

	r.Run()
}
