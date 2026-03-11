package kyc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AlipayKYC handles alipay identity verification API calls
type AlipayKYC struct {
	appID           string
	privateKey      *rsa.PrivateKey
	alipayPublicKey *rsa.PublicKey
	gateway         string
	callbackURL     string
}

// NewAlipayKYC creates a new AlipayKYC instance
func NewAlipayKYC(appID, privateKeyPEM, alipayPublicKeyPEM, gateway, callbackURL string) (*AlipayKYC, error) {
	// Parse private key
	privBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if privBlock == nil {
		return nil, fmt.Errorf("failed to parse private key PEM")
	}
	privKey, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
	if err != nil {
		privKey, err = x509.ParsePKCS1PrivateKey(privBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
	}
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	// Parse alipay public key
	pubBlock, _ := pem.Decode([]byte(alipayPublicKeyPEM))
	if pubBlock == nil {
		return nil, fmt.Errorf("failed to parse alipay public key PEM")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse alipay public key: %v", err)
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("alipay public key is not RSA")
	}

	return &AlipayKYC{
		appID:           appID,
		privateKey:      rsaPrivKey,
		alipayPublicKey: rsaPubKey,
		gateway:         gateway,
		callbackURL:     callbackURL,
	}, nil
}

// InitializeCertify calls alipay.user.certify.open.certify
func (a *AlipayKYC) InitializeCertify(realName, idCardNumber, idType string) (certifyID string, certifyURL string, err error) {
	bizContent := map[string]interface{}{
		"outer_order_no": fmt.Sprintf("kyc_%d", time.Now().UnixNano()),
		"biz_code":       "FACE",
		"identity_param": map[string]string{
			"identity_type": "CERTIFY_IDENTITY_CARD",
			"cert_type":     idType,
			"cert_name":     realName,
			"cert_no":       idCardNumber,
		},
		"merchant_config": map[string]string{
			"return_url": a.callbackURL,
		},
	}

	result, err := a.signAndCall("alipay.user.certify.open.certify", bizContent)
	if err != nil {
		return "", "", err
	}

	certifyID, _ = result["certify_id"].(string)
	certifyURL, _ = result["certify_url"].(string)
	if certifyID == "" {
		return "", "", fmt.Errorf("alipay returned empty certify_id")
	}
	return certifyID, certifyURL, nil
}

// QueryCertify calls alipay.user.certify.open.query
func (a *AlipayKYC) QueryCertify(certifyID string) (passed bool, certInfo map[string]string, err error) {
	bizContent := map[string]interface{}{
		"certify_id": certifyID,
	}

	result, err := a.signAndCall("alipay.user.certify.open.query", bizContent)
	if err != nil {
		return false, nil, err
	}

	certInfo = make(map[string]string)
	if passedStr, ok := result["passed"].(string); ok {
		certInfo["passed"] = passedStr
		passed = passedStr == "T"
	}
	if certNo, ok := result["cert_no"].(string); ok {
		certInfo["cert_no"] = certNo
	}

	return passed, certInfo, nil
}

// signAndCall signs the request and calls alipay gateway
func (a *AlipayKYC) signAndCall(method string, bizContent map[string]interface{}) (map[string]interface{}, error) {
	bizJSON, err := json.Marshal(bizContent)
	if err != nil {
		return nil, fmt.Errorf("marshal biz_content: %v", err)
	}

	params := map[string]string{
		"app_id":      a.appID,
		"method":      method,
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"biz_content": string(bizJSON),
	}

	// Build sign string
	signStr := buildSignString(params)
	signature, err := a.sign(signStr)
	if err != nil {
		return nil, fmt.Errorf("sign request: %v", err)
	}
	params["sign"] = signature

	// POST to gateway
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	resp, err := http.PostForm(a.gateway, formData)
	if err != nil {
		return nil, fmt.Errorf("alipay request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %v", err)
	}

	return parseAlipayResponse(body)
}

// sign signs the string with RSA2 (SHA256WithRSA)
func (a *AlipayKYC) sign(data string) (string, error) {
	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, a.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// buildSignString builds the string to be signed
func buildSignString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		if params[k] != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
		}
	}
	return strings.Join(parts, "&")
}

// parseAlipayResponse parses alipay JSON response
func parseAlipayResponse(body []byte) (map[string]interface{}, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse response JSON: %v", err)
	}

	// Check for error response
	if errMsg, ok := raw["error_response"]; ok {
		var errResp map[string]interface{}
		if err := json.Unmarshal(errMsg, &errResp); err == nil {
			subCode, _ := errResp["sub_code"].(string)
			subMsg, _ := errResp["sub_msg"].(string)
			return nil, fmt.Errorf("alipay error: [%s] %s", subCode, subMsg)
		}
		return nil, fmt.Errorf("alipay error: %s", string(errMsg))
	}

	// Extract the response (key is method response without request prefix)
	for k, v := range raw {
		if k == "sign" || k == "sign_type" {
			continue
		}
		var result map[string]interface{}
		if err := json.Unmarshal(v, &result); err == nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("empty response from alipay")
}
