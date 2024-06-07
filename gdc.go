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
	cl := &Client{Config: config}

	return cl
}

func (c *Client) GetAccessToken() (*GetAccessTokenResponseBody, error) {
	var result GetAccessTokenResponseBody

	timestamp := c.getTimestamp()
	requestUrl := fmt.Sprintf("%s/%s/auth/access-token", c.Config.BaseUrl, "v1.0")
	strToSign := fmt.Sprintf("%s|%s", c.Config.ClientKey, timestamp)
	signature, err := c.sign(strToSign)
	if err != nil {
		c.log("error", map[string]interface{}{
			"error": err.Error(),
			"url":   requestUrl,
			"detail": map[string]interface{}{
				"activity": "sign request",
				"to_sign":  strToSign,
			},
		})

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

	raw, err := goutil.SendHttpPost(requestUrl, requestBody, &requestHeaders, &result)

	var responseBodyRaw string

	if raw != nil {
		responseBodyRaw = string(*raw)
	}

	if err != nil {
		c.log("error", map[string]interface{}{
			"error":    err.Error(),
			"url":      requestUrl,
			"header":   requestHeaders,
			"response": responseBodyRaw,
		})

		return nil, err
	}

	c.log("debug", map[string]interface{}{
		"url":      requestUrl,
		"header":   requestHeaders,
		"response": responseBodyRaw,
	})

	return &result, nil
}

func (c *Client) SetAccessToken(accessToken *string) {
	c.accessToken = accessToken
}

func (c *Client) WithAccessToken(accessToken *string) *Client {
	c.SetAccessToken(accessToken)

	return c
}

/*
 * Logging
 *
 * @param string level
 * @param interface{} data
 *
 * @return void
 */
func (c *Client) log(level string, data interface{}) {
	if c.Config.Log != nil && c.Config.Log.Enable {
		if c.Config.Log.Level == "error" && level != "error" {
			return
		}

		msg := map[string]interface{}{
			"level": level,
			"data":  data,
		}

		_ = goutil.WriteJsonToFile(msg, c.Config.Log.Path, c.Config.Log.Filename, c.Config.Log.Extension, c.Config.Log.Rotation)
	}
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

	timestamp := c.getTimestamp()
	requestUrl := fmt.Sprintf("%s%s", c.Config.BaseUrl, uri)
	byData, err := json.Marshal(data)
	if err != nil {
		c.log("error", map[string]interface{}{
			"error": err.Error(),
			"url":   requestUrl,
			"detail": map[string]interface{}{
				"activity": "json marshal post data",
				"data":     data,
			},
		})

		return err
	}

	strData := strings.ToLower(hex.EncodeToString(byData))
	strToSign := fmt.Sprintf("%s|%s|%s|%s", timestamp, *c.accessToken, uri, strData)
	signature, err := c.sign(strToSign)
	if err != nil {
		c.log("error", map[string]interface{}{
			"error": err.Error(),
			"url":   requestUrl,
			"detail": map[string]interface{}{
				"activity": "sign request",
				"data":     data,
				"to_sign":  strToSign,
			},
		})

		return err
	}

	requestHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", *c.accessToken),
		"X-Request-Id":  uuid.New().String(),
		"X-Timestamp":   timestamp,
		"X-Signature":   *signature,
	}

	raw, err := goutil.SendHttpPost(requestUrl, data, &requestHeaders, result)

	var responseBodyRaw string

	if raw != nil {
		responseBodyRaw = string(*raw)
	}

	if err != nil {
		c.log("error", map[string]interface{}{
			"error":    err.Error(),
			"url":      requestUrl,
			"header":   requestHeaders,
			"request":  data,
			"response": responseBodyRaw,
		})

		return err
	}

	c.log("debug", map[string]interface{}{
		"url":      requestUrl,
		"header":   requestHeaders,
		"request":  data,
		"response": responseBodyRaw,
	})

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

	timestamp := c.getTimestamp()
	requestUrl := fmt.Sprintf("%s%s", c.Config.BaseUrl, uri)
	strToSign := fmt.Sprintf("%s|%s|%s|", timestamp, *c.accessToken, uri)
	signature, err := c.sign(strToSign)
	if err != nil {
		c.log("error", map[string]interface{}{
			"error": err.Error(),
			"url":   requestUrl,
			"detail": map[string]interface{}{
				"activity": "sign request",
				"to_sign":  strToSign,
			},
		})

		return err
	}

	requestHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", *c.accessToken),
		"X-Request-Id":  uuid.New().String(),
		"X-Timestamp":   timestamp,
		"X-Signature":   *signature,
	}

	raw, err := goutil.SendHttpGet(requestUrl, nil, &requestHeaders, result)

	var responseBodyRaw string

	if raw != nil {
		responseBodyRaw = string(*raw)
	}

	if err != nil {
		c.log("error", map[string]interface{}{
			"error":    err.Error(),
			"url":      requestUrl,
			"header":   requestHeaders,
			"response": responseBodyRaw,
		})

		return err
	}

	c.log("debug", map[string]interface{}{
		"url":      requestUrl,
		"header":   requestHeaders,
		"response": responseBodyRaw,
	})

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

func (c *Client) getTimestamp() string {
	now := time.Now()
	timezone := c.Config.Timezone

	if timezone == "" {
		timezone = DefaultTimezone
	}

	loc, err := time.LoadLocation(timezone)
	if err == nil {
		now = now.In(loc)
	}

	return now.Format(TimestampFormat)
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
