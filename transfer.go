package gdc

import (
	"fmt"
	"strconv"
)

func (c *Client) BalanceInquiry() (*BalanceInquiryResponseBody, error) {
	var result BalanceInquiryResponseBody

	err := c.get("/v1.0/account/balance-inquiry", &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) BankAccountNameInquiry(bankCode, accountNumber string) (*BankAccountNameInquiryResponseBody, error) {
	var result BankAccountNameInquiryResponseBody

	err := c.get(fmt.Sprintf("/v1.0/transfer/account-inquiry?bankCode=%s&accountNumber=%s", bankCode, accountNumber), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) TransferInquiry(currency, amount, bankCode, accountNumber, accountName, transactionId, remark string) (*TransferInquiryResponseBody, error) {
	var result TransferInquiryResponseBody

	flAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return nil, err
	}

	err = c.post(fmt.Sprintf("/v1.0/transfer/fund-transfer/transfer"), map[string]interface{}{
		"amount": map[string]string{
			"currency": currency,
			"value":    fmt.Sprintf("%.2f", flAmount),
		},
		"destAccountName": accountName,
		"destBankCode":    bankCode,
		"destAccountNo":   accountNumber,
		"partnerReff":     transactionId,
		"remark":          remark,
	}, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
