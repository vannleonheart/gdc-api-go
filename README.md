### GDC API

#### Installation
```bash
go get -u github.com/vannleonheart/gdc-api-go
```
#### Config
```go
gdcConfig := gdc.Config{
    BaseUrl:            "{gdc_api_base_url}",
    ClientKey:          "{your_client_key}",
    PrivateKeyFilePath: "{path_to_your_private_key_file}",
}
```
#### Client
```go
gdcClient := gdc.NewClient(gdcConfig)
```
#### Access Token
```go       
accessToken, err := gdcClient.GetAccessToken()

if err != nil {
    // handle error
}

accessTokenString := accessToken.AccessToken

fmt.Println(accessTokenString)
```
Set the access token to the client
```go
gdcClient.SetAccessToken(accessTokenString)
```
or
```go
gdcClient = gdcClient.WithAccessToken(accessTokenString)
```
#### Balance Inquiry
```go
result, err := gdcClient.BalanceInquiry()

if err != nil {
    // handle error
}

fmt.Println(result.ActiveBalance)
```
#### Bank Account Inquiry
```go
bankCode := "{your_bank_code}"
accountNumber := "{your_account_number}"

result, err := gdcClient.BankAccountNameInquiry(bankCode, accountNumber)

if err != nil {
    // handle error
}

fmt.Println(result.AccountName)
```
#### Transfer Inquiry
```go
currency := "IDR"
amount := "10000.00"
bankCode := "{your_bank_code}"
accountNumber := "{your_account_number}"
accountName := "{your_account_name}"
yourTransactionId := "{your_transaction_id}"
remark := "{your_remark}"

result, err := gdcClient.TransferInquiry(currency, amount, bankCode, accountNumber, accountName, yourTransactionId, remark)

if err != nil {
    // handle error
}

fmt.Println(result.TrxReff)
```