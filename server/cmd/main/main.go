package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emarifer/calorie-tracker-app/pkg/api/config"
	"github.com/emarifer/calorie-tracker-app/pkg/api/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("🔥 failed to load environment variables file!\n", err.Error())
			os.Exit(1)
		}
	}

	config.DBInstance()
}

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	routes.SetupRoutesApi(router)

	router.Static("/home", "../client/dist")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(fmt.Sprintf(":%s", port))
}

/* Documentación de Gin. VER:
https://github.com/gin-gonic/gin
https://gin-gonic.com/es/docs/

Serve static files in Gin web application. VER:
https://blog.petehouston.com/serve-static-files-in-gin-web-application/

Problema en Gin con la ruta de los archivos estáticos. VER:
Enrutador Gin: el segmento de ruta entra en conflicto con el comodín existente:
https://stackoverflow.com/questions/36357791/gin-router-path-segment-conflicts-with-existing-wildcard
https://github.com/julienschmidt/httprouter/issues/12#issuecomment-46121392

Documentación de la librería de CORS para Gin. VER:
https://github.com/gin-contrib/cors

Documentación de MongoDB Go Driver. VER:
https://github.com/mongodb/mongo-go-driver
*/
