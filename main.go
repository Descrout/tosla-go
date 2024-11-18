package main

import (
	"log"
	"os"

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
}
