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

	//startTime := time.Now()
	// Create new Service Workflow
	//handlers.WorkflowRequestCounter.WithLabelValues("workflow").Inc()

	var workflow models.Workflow

	jsonData, err := os.ReadFile("mapping.json")

	err = json.Unmarshal(jsonData, &workflow)
	if err != nil {
		log.Errorf("binding to workflow request failed: %s", err)
		log.Println(err)
		return
	}

	//log.Printf("Create workflow request : %v", workflow)
	serviceFlowRes := models.CreateWorkflow(&workflow)
	serviceFlowRes, err = services.SaveWorkflow(serviceFlowRes)
	//handlers.WorkflowRequestHistogram.Observe(time.Since(startTime).Seconds())
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println(serviceFlowRes)

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
	/////////////////////////////////////////////////////
	//// Чтение входного файла
	//data, err := ioutil.ReadFile("atlasmapping.json")
	//if err != nil {
	//	panic(err)
	//}
	//
	//var atlas utils.AtlasMapping
	//err = json.Unmarshal(data, &atlas)
	//if err != nil {
	//	panic(err)
	//}
	//
	////for i1, i2 := range atlas.AtlasMapping.Mappings.Mapping {
	////	re := regexp.MustCompile(`[^\d\p{Latin}]`)
	////	//println(i2.ID, " ", re.ReplaceAllString(i2.ID, ""))
	////	atlas.AtlasMapping.Mappings.Mapping[i1].ID = re.ReplaceAllString(i2.ID, "")
	////}
	//
	////// Построение выходного объекта
	////output := utils.Output{
	////	Name:        "process_claim",
	////	Description: "processing of medical claim",
	////}
	////
	////// Построение карты DataSource по docId
	////dataSourceMap := make(map[string]string)
	////for _, ds := range atlas.AtlasMapping.DataSources {
	////	dataSourceMap[ds.Id] = fmt.Sprintf("www.%s.ru", ds.Name)
	////}
	////
	////// Построение шагов
	////var steps []utils.Step
	////for _, mapping := range atlas.AtlasMapping.Mappings.Mapping {
	////	step := utils.Step{
	////		Name: mapping.ID,
	////		Mode: "HTTP", // Можно расширить для других Mode
	////		Val: utils.Val{
	////			Method: "POST",
	////			URL:    utils.GetValURL(mapping.OutputField, dataSourceMap),
	////		},
	////	}
	////
	////	// Добавление requestTransform
	////	if len(mapping.InputFields) > 0 || mapping.InputGroup != nil {
	////		step.Transform = true
	////		step.RequestTransform = &utils.RequestTransform{
	////			Spec: utils.BuildRequestTransform(mapping),
	////		}
	////	}
	////
	////	steps = append(steps, step)
	////}
	////
	////// Сортировка по priority
	////sort.Slice(steps, func(i, j int) bool {
	////	return atlas.AtlasMapping.Mappings.Mapping[i].Priority < atlas.AtlasMapping.Mappings.Mapping[j].Priority
	////})
	////
	////output.Steps = steps
	////
	////// Конвертация в JSON и вывод
	////result, err := json.MarshalIndent(output, "", "  ")
	////if err != nil {
	////	panic(err)
	////}
	////fmt.Println(string(result))
	//
	//g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())
	//
	//for _, i2 := range atlas.AtlasMapping.Mappings.Mapping {
	//	_ = g.AddVertex(i2.ID)
	//
	//	for _, field := range i2.InputFields {
	//		//for _, jsonField := range i2.OutputField {
	//		//_ = g.AddVertex(field.DocID + field.Name)
	//		//_ = g.AddVertex(jsonField.DocID + jsonField.Name)
	//		//_ = g.AddEdge(field.DocID+field.Name, jsonField.DocID+jsonField.Name)
	//		//println(field.DocID+field.Name, " - ", jsonField.DocID+jsonField.Name)
	//		//}
	//
	//		for _, action := range field.Actions {
	//			if action.Type == "CopyTo" {
	//
	//				findedmapping, _ := lo.Find(atlas.AtlasMapping.Mappings.Mapping, func(m utils.Mapping) bool {
	//					return m.ID == action.Index
	//				})
	//
	//				_ = g.AddEdge(i2.ID, findedmapping.ID)
	//				println(i2.ID, " - ", findedmapping.ID)
	//
	//			}
	//		}
	//
	//		//_ = g.AddEdge(i2.ID, field.DocID+field.Name)
	//	}
	//}
	//
	//file, _ := os.Create("./mygraph.gv")
	//_ = draw.DOT(g, file)
	//
	//// For a deterministic topological ordering, use StableTopologicalSort.
	//orderarr, _ := graph.TopologicalSort(g)
	//order := orderarr[:]
	//
	//fmt.Println(order)

}
