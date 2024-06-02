package gdc

type Client struct {
	Config      Config
	accessToken *string
}

type Config struct {
	BaseUrl            string     `json:"base_url"`
	ClientKey          string     `json:"client_key"`
	PrivateKeyFilePath string     `json:"private_key_file_path"`
	Log                *LogConfig `json:"log"`
}

type LogConfig struct {
	Enable    bool   `json:"enable"`
	Level     string `json:"level"`
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	Extension string `json:"extension"`
	Rotation  string `json:"rotation"`
}

type ErrorResponse struct {
	ResponseCode    *string `json:"responseCode,omitempty"`
	ResponseMessage *string `json:"responseMessage,omitempty"`
}

type GetAccessTokenResponseBody struct {
	*ErrorResponse
	AccessToken string `json:"accessToken,omitempty"`
	ExpiresIn   int    `json:"expiresIn,omitempty"`
	Type        string `json:"type,omitempty"`
}

type BalanceInquiryResponseBody struct {
	*ErrorResponse
	ActiveBalance  float64 `json:"activeBalance,omitempty"`
	DepositBalance float64 `json:"depositBalance,omitempty"`
	FloatingDebt   float64 `json:"floatingDebt,omitempty"`
}

type BankAccountNameInquiryResponseBody struct {
	*ErrorResponse
	AccountName string `json:"accountName,omitempty"`
}

type TransferInquiryResponseBody struct {
	*ErrorResponse
	Amount struct {
		Currency string `json:"currency,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"amount,omitempty"`
	BankReff       string `json:"bankReff,omitempty"`
	PartnerReff    string `json:"partnerReff,omitempty"`
	TrxReff        string `json:"trxReff,omitempty"`
	TransferMethod string `json:"transferMethod,omitempty"`
}

type TransferCallbackRequestBody struct {
	Timestamp    string `json:"timestamp"`
	TransferType string `json:"transferType"`
	Destination  struct {
		AccountName string `json:"accountName"`
		AccountNo   string `json:"accountNo"`
		BankCode    string `json:"bankCode"`
	}
	Amount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	}
	PartnerReff string `json:"partnerReff"`
	TrxReff     string `json:"trxReff"`
	BankReff    string `json:"bankReff"`
}
