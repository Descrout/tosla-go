package toslago

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Descrout/tosla-go/tosla-go/requests"
	"github.com/Descrout/tosla-go/tosla-go/responses"
	"github.com/Descrout/tosla-go/tosla-go/utils"
)

const (
	PROD_URL    = "https://entegrasyon.tosla.com"
	SANDBOX_URL = "https://prepentegrasyon.tosla.com"
)

type ToslaOptions struct {
	BaseUrl  string
	ApiUser  string
	ApiPass  string
	ClientID string
}

type Tosla struct {
	baseUrl  string
	apiUser  string
	apiPass  string
	clientID string
	client   *http.Client
}

func WithOptions(options *ToslaOptions) *Tosla {
	return &Tosla{
		baseUrl:  options.BaseUrl,
		apiUser:  options.ApiUser,
		apiPass:  options.ApiPass,
		clientID: options.ClientID,
		client:   &http.Client{},
	}
}

func (t *Tosla) GenerateRndTimeHash() (string, string, string) {
	// Seed the random number generator for true randomness
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 1 and 999999
	rnd := strconv.Itoa(rand.Intn(1000000) + 1)

	// Get current time in "yyyyMMddHHmmss" format
	timeSpan := time.Now().Format("20060102150405")

	// Concatenate all parts to form the hash string
	hashString := t.apiPass + t.clientID + t.apiUser + rnd + timeSpan

	// Compute SHA-512 hash
	hashBytes := sha512.New()
	hashBytes.Write([]byte(hashString))
	hashed := hashBytes.Sum(nil)

	// Convert the hash to Base64
	hash := base64.StdEncoding.EncodeToString(hashed)

	return rnd, timeSpan, hash
}

func (t *Tosla) getToslaReq() map[string]any {
	result := map[string]any{
		"clientId": t.clientID,
		"apiUser":  t.apiUser,
	}
	result["rnd"], result["timeSpan"], result["hash"] = t.GenerateRndTimeHash()

	return result
}

func (t *Tosla) makeRequest(method string, endpoint string, data any) ([]byte, error) {
	toslaReq := t.getToslaReq()

	req, err := utils.StructToMap(data)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(utils.CombineMaps(toslaReq, req))
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(method, t.baseUrl+endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-type", "application/json")
	httpReq.Header.Set("Expect", "100-continue")
	httpReq.Header.Set("Connection", "Keep-Alive")
	httpReq.Header.Set("Cache-Control", "no-cache")

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return rawBody, nil
}

func (t *Tosla) CheckBin(req *requests.BinCheck) (*responses.BinResponse, error) {
	rawData, err := t.makeRequest("POST", "/api/Payment/GetCommissionAndInstallmentInfo", req)
	if err != nil {
		return nil, err
	}

	resp := &responses.BinResponse{}
	err = json.Unmarshal(rawData, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
