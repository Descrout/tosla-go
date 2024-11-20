package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	toslago "github.com/Descrout/tosla-go/tosla-go"
	"github.com/Descrout/tosla-go/tosla-go/requests"
)

func main() {
	cli := toslago.WithOptions(&toslago.ToslaOptions{
		BaseUrl:  toslago.SANDBOX_URL,
		ApiUser:  "POS_ENT_Test_001",
		ApiPass:  "POS_ENT_Test_001!*!*",
		ClientID: "1000000494",
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
		InstallmentCount: 0,    // No installment
		Amount:           6999, // 69.99
		Currency:         949,  // TRY
		CallbackURL:      "http://localhost:8888/3dsconfirm",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(threedsInit)

	pay3dsHtml, err := cli.Pay3dsHtml(&requests.Pay3dsRequest{
		ThreeDSessionID: threedsInit.ThreeDSessionID,
		CardHolderName:  "Adil Basar",
		CardNo:          "5571135571135575",
		ExpireDate:      "12/24",
		Cvv:             "000",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(pay3dsHtml))

	payNon3ds, err := cli.PayNon3ds(&requests.Non3dsRequest{
		CardHolderName:   "Adil Basar",
		CardNo:           "5571135571135575",
		ExpireDate:       "12/24",
		Cvv:              "200",
		OrderID:          "",
		Description:      "This is a test purchase",
		Echo:             "echo",
		ExtraParameters:  "extra",
		InstallmentCount: 0,    // No installment
		Amount:           6999, // 69.99
		Currency:         949,  // TRY
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(payNon3ds)

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

		hash := r.Form.Get("Hash")
		orderID := r.Form.Get("OrderId")
		bankResponseCode := r.Form.Get("BankResponseCode")
		bankResponseMsg := r.Form.Get("BankResponseMessage")
		requestStatus := r.Form.Get("RequestStatus")
		mdStatus := r.Form.Get("MdStatus")

		if !cli.ValidateIncomingHash(hash, orderID, mdStatus, bankResponseCode, bankResponseMsg, requestStatus) {
			http.Error(w, "hash validation error", http.StatusUnauthorized)
			return
		}

		if requestStatus == "1" && mdStatus == "1" {
			log.Println("Payment Successful !!!")
		} else {
			log.Println("Payment Error !!!")
		}

		for key, value := range r.Form {
			log.Println(key, value)
		}
	})

	return mux
}
