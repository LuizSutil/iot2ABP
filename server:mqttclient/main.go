package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
  	"strconv"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/jackc/pgx/v4"

)

type Postgre struct {
	TopicName  string
	Payload    string
	DBConn     sql.DB
	ClientMQTT mqtt.Client
}

const (

	// Broker
	topicName = "satc/iot/#"
	brokerUrl = "broker.hivemq.com:1883"

	// Database
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "barbadovagner"
	dbname   = "mqttarduino"
)

type MqttPayload struct {
	ClientId string `json:"clientId"`
	Sala     string `json:"sala"`
	Quarto   string `json:"quarto"`
	Cozinha  string `json:"cozinha"`
	Distancia   string `json:"sensor"`
}

func latestValuesHandler(w http.ResponseWriter, r *http.Request) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		//fmt.Println(" ");
		//fmt.Println(err)
		panic(err)
	}
	var payloadFinal []MqttPayload

	var luizPayload MqttPayload;
	luizPayload.ClientId = "1"
	sqlFind := `SELECT valor FROM comdados WHERE client_name = 'luiz' AND sensor = 'sala' ORDER BY written_time DESC LIMIT 1`
	rows, err := db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var sala int
    	err = rows.Scan(&sala)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}
    	luizPayload.Sala = strconv.Itoa(sala)
  	}
	sqlFind = `SELECT valor FROM comdados WHERE client_name = 'luiz' AND sensor = 'cozinha' ORDER BY written_time DESC LIMIT 1`
	rows, err = db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var cozinha int
    	err = rows.Scan(&cozinha)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}
    	luizPayload.Cozinha = strconv.Itoa(cozinha)
  	}

  	sqlFind = `SELECT valor FROM comdados WHERE client_name = 'luiz' AND sensor = 'quarto' ORDER BY written_time DESC LIMIT 1`
	rows, err = db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var quarto int
    	err = rows.Scan(&quarto)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}
    	luizPayload.Quarto = strconv.Itoa(quarto)
  	}

  	sqlFind = `SELECT valor FROM comdados WHERE client_name = 'luiz' AND sensor = 'distancia' ORDER BY written_time DESC LIMIT 1`
	rows, err = db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var distancia float64
    	err = rows.Scan(&distancia)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}

    	luizPayload.Distancia = fmt.Sprintf("%f", distancia)
  	}

  	payloadFinal = append(payloadFinal, luizPayload)


  	var eddiePayload MqttPayload;
	eddiePayload.ClientId = "2"
	sqlFind = `SELECT valor FROM comdados WHERE client_name = 'eddie' AND sensor = 'sala' ORDER BY written_time DESC LIMIT 1`
	rows, err = db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var sala int
    	err = rows.Scan(&sala)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}
    	eddiePayload.Sala = strconv.Itoa(sala)
  	}
	sqlFind = `SELECT valor FROM comdados WHERE client_name = 'eddie' AND sensor = 'cozinha' ORDER BY written_time DESC LIMIT 1`
	rows, err = db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var cozinha int
    	err = rows.Scan(&cozinha)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}
    	eddiePayload.Cozinha = strconv.Itoa(cozinha)
  	}

  	sqlFind = `SELECT valor FROM comdados WHERE client_name = 'eddie' AND sensor = 'quarto' ORDER BY written_time DESC LIMIT 1`
	rows, err = db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var quarto int
    	err = rows.Scan(&quarto)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}
    	eddiePayload.Quarto = strconv.Itoa(quarto)
  	}

  	sqlFind = `SELECT valor FROM comdados WHERE client_name = 'eddie' AND sensor = 'distancia' ORDER BY written_time DESC LIMIT 1`
	rows, err = db.Query(context.Background(), sqlFind)
	if err != nil {
    	fmt.Println(err)
  	}
  	defer rows.Close()

  	for rows.Next() {

    	var distancia float64
    	err = rows.Scan(&distancia)
    	if err != nil {
      	// handle this error
      	panic(err)
    	}

    	eddiePayload.Distancia = fmt.Sprintf("%f", distancia)
  	}
  	 payloadFinal = append(payloadFinal, eddiePayload)




  	jsonEncoded, err  := json.Marshal(payloadFinal)
    if err != nil {
        fmt.Println(err)
	    return
   	}	
    jsonArray := string(jsonEncoded)


	fmt.Fprintf(w, jsonArray)

}

func toggleLed (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comodo := vars["comodo"]
	cliente := vars["cliente"]
	toggle := vars["toggle"]

	var nomeUser string
	if cliente == "1" {
		nomeUser = "luiz"
	} else {
		nomeUser = "eddie"
	}
	
	clientOptions := mqtt.NewClientOptions().AddBroker(brokerUrl).SetClientID("IOT2Client2")
	client := mqtt.NewClient(clientOptions)
	//
	//// Puxa o token do broker e verifica se alguem erro occoreu
	if token := client.Connect(); token.Wait() && token.Error() != nil {
	}

	var message string

	if comodo == "sala" {
		if toggle == "1"{
			message = "S"
		} else {
			message = "s"
		}
	} else if comodo == "quarto" {
		if toggle == "1"{
			message = "Q"
		} else {
			message = "q"
		}
	}	else if comodo == "cozinha" {
		if toggle == "1"{
			message = "C"
		} else {
			message = "c"
		}
	}

	
	client.Publish("satc/iot/"+nomeUser, 2, false, message)
	client.Disconnect(100)

}

func mqttClient() {
	//// Cria canal go
	//
	//// Conecta no broker com os parametros
	clientOptions := mqtt.NewClientOptions().AddBroker(brokerUrl).SetClientID("IOT2LuizClient")
	client := mqtt.NewClient(clientOptions)
	//
	//// Puxa o token do broker e verifica se alguem erro occoreu
	if token := client.Connect(); token.Wait() && token.Error() != nil {
	}

	// Faz o subscribe a um topico e recebe mensagem
	if token := client.Subscribe(topicName, 0, func(client mqtt.Client, msg mqtt.Message) {

		tpc := msg.Topic()
		data := string(msg.Payload())

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			fmt.Println(err)
		}
		
		db.SetConnMaxLifetime(5 * time.Second)
		db.SetMaxOpenConns(0)
		db.SetMaxIdleConns(0)
		//

		sqlInsert := `
			INSERT INTO comdados(
			client_name,
			sensor,
			valor,
			written_time
			)
			VALUES($1, $2,$3, current_timestamp at time zone 'UTC')`



		var sensorNome string
		var payloadValue string
		if data == "Q" {
			sensorNome = "quarto"
			payloadValue = "1"
		} else if data == "q" {
			sensorNome = "quarto"
			payloadValue = "0"
		} else if data == "C" {
			sensorNome = "cozinha"
			payloadValue = "1"
		} else if data == "c" {
			sensorNome = "cozinha"
			payloadValue = "0"
		} else if data == "S" {
			sensorNome = "sala"
			payloadValue = "1"
		} else if data == "s" {
			sensorNome = "sala"
			payloadValue = "0"

		} else {
			sensorNome = "distancia"
			payloadValue = data
		}

		nomeUser := strings.ReplaceAll(tpc, "satc/iot/", "")
		_, err = db.Exec(sqlInsert, nomeUser, sensorNome, payloadValue)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tpc)
		fmt.Println(data)
	});

	// Mostra ao usuario que r
	token.Wait() && token.Error() == nil {
		fmt.Print("Subscription to topic: (" + topicName + ") successful!!\n")
	}

}

func goServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", latestValuesHandler)
	r.HandleFunc("/{comodo}/{cliente}/{toggle}", toggleLed)


	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application

		},
	})

	handler := c.Handler(r)

	http.ListenAndServe(":8000", handler)
	fmt.Println("Server Started on port 8000")
}

func main() {
	c := make(chan os.Signal, 2)

	go mqttClient()
	go goServer()
	<-c

	os.Exit(0)

}
