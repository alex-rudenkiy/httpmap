package main

import (
	"clamp-core/config"
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/models"
	"clamp-core/repository"
	"clamp-core/services"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	logLevel, err := log.ParseLevel(config.ENV.LogLevel)
	if err != nil {
		log.Fatalf("parsing log level failed: %s", err)
	}

	log.SetLevel(logLevel)

	log.Info("Pinging DB...")
	err = repository.GetDB().Ping()
	if err != nil {
		log.Fatalf("DB ping failed: %s", err)
	}

	var workflow models.Workflow

	jsonData, err := os.ReadFile("mapping.json")

	err = json.Unmarshal(jsonData, &workflow)
	if err != nil {
		log.Errorf("binding to workflow request failed: %s", err)
		log.Println(err)
		return
	}

	log.Printf("Create workflow request : %v", workflow)
	serviceFlowRes := models.CreateWorkflow(&workflow)
	serviceFlowRes, err = services.SaveWorkflow(serviceFlowRes)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(serviceFlowRes)

	//var cliArgs models.CLIArguments = os.Args[1:]
	//os.Setenv("PORT", config.ENV.PORT)
	//migrations.Migrate()
	//
	//if cliArgs.Parse().Find("migrate-only", "no") == "yes" {
	//	os.Exit(0)
	//}

	if config.ENV.EnableRabbitMQIntegration {
		listeners.AMQPStepResponseListener.Listen()
	}
	if config.ENV.EnableKafkaIntegration {
		listeners.KafkaStepResponseListener.Listen()
	}
	handlers.LoadHTTPRoutes()
	log.Info("Calling listener")
	// docker run -e CLAMP_DB_DRIVER="inMemoryRepository" -p 8080:8080 --mount type=bind,source="$(pwd)"/mapping.json,target=/mapping.json,readonly sha256:c7ac6312b0421347a47fa41f4b216ced285590a5554270e64bc281f182da9767
}
