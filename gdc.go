package gdc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vannleonheart/goutil"
	"os"
	"strings"
	"time"
)

func New(config Config) *Client {
	return &Client{Config: config}
}

func (c *Client) GetAccessToken() (*GetAccessTokenResponseBody, error) {
	var result GetAccessTokenResponseBody

	timestamp := time.Now().Format(TimestampFormat)
	requestUrl := fmt.Sprintf("%s/%s/auth/access-token", c.Config.BaseUrl, "v1.0")

	signature, err := c.sign(fmt.Sprintf("%s|%s", c.Config.ClientKey, timestamp))
	if err != nil {
		return nil, err
	}

	requestHeaders := map[string]string{
		"Client-Key":  c.Config.ClientKey,
		"X-Timestamp": timestamp,
		"X-Signature": *signature,
	}

	requestBody := map[string]string{
		"grantType": "client_credentials",
	}

	if _, err := goutil.SendHttpPost(requestUrl, requestBody, &requestHeaders, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) SetAccessToken(accessToken *string) {
	c.accessToken = accessToken
}

func (c *Client) WithAccessToken(accessToken *string) *Client {
	c.SetAccessToken(accessToken)

	return c
}

func (c *Client) post(uri string, data map[string]interface{}, result interface{}) error {
	defer c.SetAccessToken(nil)

	if c.accessToken == nil {
		accessToken, err := c.GetAccessToken()
		if err != nil {
			return err
		}

		c.SetAccessToken(&accessToken.AccessToken)
	}

	timestamp := time.Now().Format(TimestampFormat)
	requestUrl := fmt.Sprintf("%s%s", c.Config.BaseUrl, uri)

	byData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	strData := strings.ToLower(hex.EncodeToString(byData))
	strToSign := fmt.Sprintf("%s|%s|%s|%s", timestamp, *c.accessToken, uri, strData)

	signature, err := c.sign(strToSign)
	if err != nil {
		return err
	}

	requestHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", *c.accessToken),
		"X-Request-Id":  uuid.New().String(),
		"X-Timestamp":   timestamp,
		"X-Signature":   *signature,
	}

	if _, err := goutil.SendHttpPost(requestUrl, data, &requestHeaders, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) get(uri string, result interface{}) error {
	defer c.SetAccessToken(nil)

	if c.accessToken == nil {
		accessToken, err := c.GetAccessToken()
		if err != nil {
			return err
		}

		c.SetAccessToken(&accessToken.AccessToken)
	}

	timestamp := time.Now().Format(TimestampFormat)
	requestUrl := fmt.Sprintf("%s%s", c.Config.BaseUrl, uri)
	strToSign := fmt.Sprintf("%s|%s|%s|", timestamp, *c.accessToken, uri)

	signature, err := c.sign(strToSign)
	if err != nil {
		return err
	}

	requestHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", *c.accessToken),
		"X-Request-Id":  uuid.New().String(),
		"X-Timestamp":   timestamp,
		"X-Signature":   *signature,
	}

	if _, err := goutil.SendHttpGet(requestUrl, nil, &requestHeaders, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) sign(strToSign string) (*string, error) {
	pk, err := parsePrivateKey(c.Config.PrivateKeyFilePath)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	if _, err := h.Write([]byte(strToSign)); err != nil {
		return nil, err
	}

	signed, err := rsa.SignPKCS1v15(rand.Reader, pk, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return nil, err
	}

	signature := base64.StdEncoding.EncodeToString(signed)

	return &signature, nil
}

func parsePrivateKey(pvKeyFilePath string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(pvKeyFilePath)
	if err != nil {
		return nil, err
	}

	blockPvtKey, _ := pem.Decode(b)
	if blockPvtKey == nil {
		return nil, errors.New("invalid private key format")
	}

	pvtKey, err := x509.ParsePKCS1PrivateKey(blockPvtKey.Bytes)
	if err == nil {
		return pvtKey, nil
	}

	key, err2 := x509.ParsePKCS8PrivateKey(blockPvtKey.Bytes)
	if err2 == nil {
		valPvtKey, ok := key.(*rsa.PrivateKey)
		if ok {
			return valPvtKey, nil
		}

		return nil, fmt.Errorf("expected *rsa.PrivateKey, got %T", key)
	}

	return nil, errors.Join(err, err2)
}
