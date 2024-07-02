// Package main provides an HTTP server that handles user account creation and password management in an LDAP directory.
// It exposes an API endpoint that allows user creation or password update via a JSON payload.
// If the user does not exist, it creates a new user with a generated password.
// If the user already exists, it updates the user's password.

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
)

// Configuration variables
var (
	ldapURL, ldapRoot, ldapAdminUser, ldapAdminPass, ldapUsersOU, webexConnectHookUrl, webexConnectKey string
)

// Structs for responses and requests
type WebexResponse struct {
	Response []struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Transid     string `json:"transid"`
	} `json:"response"`
}
type RequestBody struct {
	PhoneNumber string   `json:"phoneNumber"`
	LoginURL    string   `json:"loginURL"`
	ContinueURL string   `json:"continueURL"`
	ApMAC       string   `json:"apMAC"`
	ApName      string   `json:"apName"`
	ApTags      []string `json:"apTags"`
	ClientMAC   string   `json:"clientMAC"`
	ClientIP    string   `json:"clientIP"`
}
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func init() {
	// Initialize logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

	// Assign configuration variables
	ldapURL = os.Getenv("LDAP_URL")
	ldapRoot = os.Getenv("LDAP_ROOT")
	ldapAdminUser = os.Getenv("LDAP_ADMIN_USERNAME")
	ldapAdminPass = os.Getenv("LDAP_ADMIN_PASSWORD")
	ldapUsersOU = os.Getenv("LDAP_USERS_OU")
	webexConnectHookUrl = os.Getenv("WEBEXCONNECT_HOOK_URL")
	webexConnectKey = os.Getenv("WEBEXCONNECT_KEY")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Received request")
	if r.Method != http.MethodPost {
		respondWithError(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, "Could not parse JSON", http.StatusBadRequest)
		return
	} else if body.PhoneNumber == "" {
		respondWithError(w, "Missing phone number", http.StatusBadRequest)
		return
	}

	password := generatePassword(12)

	// LDAP operations
	if _, err := performLDAPOperations(body, password); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Webex Webhook call
	if err := sendWebexWebhook(map[string]string{
		"phonenumber":  body.PhoneNumber,
		"login_url":    body.LoginURL,
		"continue_url": body.ContinueURL,
		"ap_mac":       body.ApMAC,
		"ap_name":      body.ApName,
		"ap_tags":      strings.Join(body.ApTags, ","),
		"client_mac":   body.ClientMAC,
		"client_ip":    body.ClientIP,
		"password":     password,
	}); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"success": true,
		"message": "Operation completed successfully",
	}).Info("Response sent")

	json.NewEncoder(w).Encode(Response{Success: true, Message: "Operation completed successfully"})
}

func performLDAPOperations(body RequestBody, password string) (string, error) {
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		logrus.WithError(err).Error("failed to connect to the LDAP server")
		return "", err
	}
	defer l.Close()

	err = l.Bind("cn="+ldapAdminUser+","+ldapRoot, ldapAdminPass)
	if err != nil {
		logrus.WithError(err).Error("failed to bind to the LDAP server")
		return "", err
	}

	searchRequest := ldap.NewSearchRequest(
		"ou="+ldapUsersOU+","+ldapRoot,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(uid="+body.PhoneNumber+")",
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		logrus.WithError(err).Error("failed to search the LDAP server")
		return "", err
	}

	if len(sr.Entries) == 0 {
		addRequest := ldap.NewAddRequest(
			"uid="+body.PhoneNumber+",ou="+ldapUsersOU+","+ldapRoot,
			nil,
		)
		addRequest.Attribute("objectClass", []string{"top", "inetOrgPerson"})
		addRequest.Attribute("cn", []string{body.PhoneNumber})
		addRequest.Attribute("sn", []string{body.PhoneNumber})
		addRequest.Attribute("userPassword", []string{password})

		err = l.Add(addRequest)
		if err != nil {
			logrus.WithError(err).Error("failed to add the user")
			return "", err
		}
	} else {
		modifyRequest := ldap.NewModifyRequest(
			"uid="+body.PhoneNumber+",ou="+ldapUsersOU+","+ldapRoot,
			nil,
		)
		modifyRequest.Replace("userPassword", []string{password})

		err = l.Modify(modifyRequest)
		if err != nil {
			logrus.WithError(err).Error("failed to update the user's password")
			return "", err
		}
	}

	return password, nil
}

func sendWebexWebhook(data map[string]string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webexConnectHookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err)
		return err
	}

	req.Header.Set("key", webexConnectKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("could not send request")
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("could not read response")
		return err
	}

	var webexResponse WebexResponse
	err = json.Unmarshal(respBody, &webexResponse)
	if err != nil {
		logrus.WithError(err).Error("could not parse response")
		return err
	}

	if len(webexResponse.Response) == 0 || webexResponse.Response[0].Code != "1002" {
		logrus.Errorf("failed to send OTP")
		//retunr internal server error
		return fmt.Errorf("webex response error: failed to send OTP")
	}

	return nil
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	logrus.WithFields(logrus.Fields{
		"success": false,
		"message": message,
		"status":  statusCode,
	}).Error("Error response")

	response := Response{
		Success: false,
		Message: message,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func generatePassword(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		logrus.WithError(err).Error("Error generating password")
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func main() {
	logrus.Info("Starting accountmanager server...")
	http.HandleFunc("/api/user", handlePost)
	http.ListenAndServe(":8080", nil)
}
