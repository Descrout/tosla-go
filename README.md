# Tosla Go SDK

The **Tosla Go SDK** is a Go library that allows developers to integrate with the Tosla payment processing system. The SDK provides functions to interact with Tosla's payment gateway, handle 3D Secure payments, bin checks, and more. While not every API endpoint is currently implemented, contributions for any missing endpoints are welcome.

---

## Installation

```bash
go get github.com/Descrout/tosla-go/tosla-go
```

---

## Usage

### Initialize Tosla Client

First, you'll need to initialize the Tosla client with your credentials and the base URL.

```go
// Tosla test information
cli := toslago.WithOptions(&toslago.ToslaOptions{
    BaseUrl:  toslago.SANDBOX_URL,
    ApiUser:  "POS_ENT_Test_001",
    ApiPass:  "POS_ENT_Test_001!*!*",
    ClientID: "1000000494",
})
```

### Check BIN

You can check the details of a BIN (Bank Identification Number) to get information about the card type, issuer, and other details.

```go
binCheck, err := cli.CheckBin(&requests.BinCheck{
    Bin: 589283, // The BIN to check
})
if err != nil {
    log.Fatal(err)
}
log.Println(binCheck)
```

### Initialize 3D Secure Payment

To initiate a 3D Secure (3DS) payment, use the `Init3ds` function. This will start the 3D Secure process and return the session ID that you can use to process the payment.

```go
threedsInit, err := cli.Init3ds(&requests.Init3dsRequest{
    OrderID:          "",
    Description:      "This is a test purchase",
    Echo:             "echo",
    ExtraParameters:  "extra",
    InstallmentCount: 0, // No installment
    Amount:           6999, // 69.99
    Currency:         949,  // TRY
    CallbackURL:      "http://localhost:8888/3dsconfirm", // Check out the "Handle Callbacks" part below
})
if err != nil {
    log.Fatal(err)
}
log.Println(threedsInit)
```

### Submit 3D Secure Payment

Once you have initialized the 3D Secure session, you can submit the payment information for processing.

```go
// Tosla test card
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
```

### Handle Callbacks

You can also set up a webhook endpoint to handle the response after the 3D Secure payment attempt. The SDK includes a helper method to validate the hash and verify the status of the payment.

```go
mux.HandleFunc("/3dsconfirm", func(w http.ResponseWriter, r *http.Request) {
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
})
```

### Non 3ds Payment
```go
// Tosla test card
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
```

### Check out the full example
Check out the ``main.go`` file for the complete example.