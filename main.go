package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	toslago "github.com/Descrout/tosla-go/tosla-go"
	"github.com/Descrout/tosla-go/tosla-go/requests"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env could not be loaded: ", err)
	}

	cli := toslago.WithOptions(&toslago.ToslaOptions{
		BaseUrl:  toslago.SANDBOX_URL,
		ApiUser:  os.Getenv("API_USER"),
		ApiPass:  os.Getenv("API_PASS"),
		ClientID: os.Getenv("CLIENT_ID"),
	})

	binCheck, err := cli.CheckBin(&requests.BinCheck{
		Bin: 589283,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(binCheck)

	threedsInit, err := cli.Init3ds(&requests.Init3dsRequest{
		OrderID:          "",
		Description:      "This is a test purchase",
		Echo:             "echo",
		ExtraParameters:  "extra",
		InstallmentCount: 0,
		Amount:           6999,
		Currency:         949,
		CallbackURL:      "http://localhost:8888/3dsconfirm",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(threedsInit)

	pay3dsHtml, err := cli.Pay3dsHtml(&requests.Pay3dsRequest{
		ThreeDSessionID: threedsInit.ThreeDSessionID,
		CardHolderName:  "Adil Basar",
		CardNo:          "5890040000000016",
		ExpireDate:      "02/28",
		Cvv:             "200",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(pay3dsHtml))

	port := ":8888"
	server := &http.Server{
		Addr:    port,
		Handler: initRoutes(cli),
	}
	defer server.Shutdown(context.TODO())

	// Run webhook and callback server on a seperate goroutine
	log.Println("Listening on port:", port)
	go func() {
		server.ListenAndServe()
		log.Println("Server shutdown gracefully.")
	}()

	// Wait for any closing signals
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt, syscall.SIGQUIT)
	<-s
	log.Println("Shutting down...")
}

func initRoutes(cli *toslago.Tosla) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/3dsconfirm", func(w http.ResponseWriter, r *http.Request) {
		log.Println("------------- INCOMING 3DS CONFIRM REQUEST -------------")

		r.ParseForm()

		for key, value := range r.Form {
			log.Println(key, value)
		}
	})

	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Println("------------- INCOMING WEBHOOK REQUEST -------------")

		data := map[string]any{}
		json.NewDecoder(r.Body).Decode(&data)
		for key, value := range data {
			log.Println(key, value)
		}
	})

	return mux
}
