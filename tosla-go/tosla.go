package toslago

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
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
	location := time.FixedZone("Europe/Istanbul", 3*60*60)
	tm := time.Now().In(location)
	log.Println(tm)
	timeSpan := tm.Format("20060102150405")
	log.Println(timeSpan)

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
	req, err := utils.StructToMap(data)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(utils.CombineMaps(t.getToslaReq(), req))
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

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

func (t *Tosla) Init3ds(req *requests.Init3dsRequest) (*responses.Init3dsResponse, error) {
	rawData, err := t.makeRequest("POST", "/api/Payment/threeDPayment", req)
	if err != nil {
		return nil, err
	}

	resp := &responses.Init3dsResponse{}
	err = json.Unmarshal(rawData, resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}

func (t *Tosla) Pay3dsHtml(pay3dreq *requests.Pay3dsRequest) ([]byte, error) {
	if err := pay3dreq.Validate(); err != nil {
		return nil, err
	}

	req, err := utils.StructToMap(pay3dreq)
	if err != nil {
		return nil, err
	}

	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)
	for k, v := range req {
		_ = writer.WriteField(k, v.(string))
	}
	writer.Close()

	httpReq, err := http.NewRequest("POST", t.baseUrl+"/api/Payment/ProcessCardForm", &requestBody)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-type", writer.FormDataContentType())
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

func (t *Tosla) ValidateIncomingHash(hashCheck string, orderID string, mdStatus string, bankResponseCode string, bankResponseMsg string, requestStatus string) bool {
	hashString := t.apiPass + t.clientID + t.apiUser + orderID + mdStatus + bankResponseCode + bankResponseMsg + requestStatus

	// Compute SHA-512 hash
	hashBytes := sha512.New()
	hashBytes.Write([]byte(hashString))
	hashed := hashBytes.Sum(nil)

	// Convert the hash to Base64
	hash := base64.StdEncoding.EncodeToString(hashed)

	return hash == hashCheck
}

func (t *Tosla) PayNon3ds(req *requests.Non3dsRequest) (*responses.Non3dsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	rawData, err := t.makeRequest("POST", "/api/Payment/Payment", req)
	if err != nil {
		return nil, err
	}

	resp := &responses.Non3dsResponse{}
	err = json.Unmarshal(rawData, resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return resp, nil
}
