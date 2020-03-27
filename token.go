package main

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"strings"
)

type tokenStruct struct {
	Token string `json:"token"`
}

type awsTokenStruct struct {
	Status struct {
		ExpirationTimestamp string `json:"expirationTimestamp"`
		Token string `json:"token"`
	}
}

type jweTokenStruct struct {
	JweToken string `json:"jweToken"`
}

func refreshJwe() {
	var cmd = viper.GetString("Authenticator")

	var params = []string {
		"token",
	}

	if cluster := viper.GetString("Cluster"); cluster != "" {
		params = append(params,
			"--cluster-id",
			cluster,
		)
	}

	if role := viper.GetString("Role"); role != "" {
		params = append(params,
			"--role",
			role,
		)
	}

	out, err := exec.Command(cmd, params...).Output()
	if err != nil {
		log.Fatalf(
			"Failed to exec authentication command %s %s\n%v\n",
			cmd, strings.Join(params, " "),
			err,
		)
	}

	awsToken := &awsTokenStruct{}
	err = json.Unmarshal(out, awsToken)
	if err != nil {
		log.Fatalf("Failed parsing json for AWS token\nError: %v\n", err)
	}
	if awsToken.Status.Token == "" {
		log.Fatalf("No aws token present\nOutput: %s\n", string(out))
	}

	res, err := http.DefaultClient.Get(
		viper.GetString("Proxy") + urlSuffix + "/api/v1/csrftoken/login")
	if err != nil {
		log.Fatalf("Failed getting csrf token\nError: %v\n", err)
	}

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed reading csrf token\nError: %v\n", err)
	}

	csrf := &tokenStruct{}
	err = json.Unmarshal(output, csrf)
	if err != nil {
		log.Fatalf("Failed creating csrf json\nError: %v\n", err)
	}
	if csrf.Token == "" {
		log.Fatalf("No csrf token present\nOutput: %s\n", string(output))
	}

	body, err := json.Marshal(tokenStruct{Token: awsToken.Status.Token})
	if err != nil {
		log.Fatalf("Failed creating login json\nError: %v\n", err)
	}

	buffer := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST",
		viper.GetString("Proxy") + urlSuffix + "/api/v1/login",
		buffer)
	if err != nil {
		log.Fatalf("Failed making login request\nError: %v\n", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-CSRF-TOKEN", csrf.Token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed posting login request\nError: %v\n", err)
	}

	output, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed reading login response\nError: %v\n", err)
	}

	jwe := &jweTokenStruct{}
	err = json.Unmarshal(output, jwe)
	if err != nil {
		log.Fatalf("Failed reading login json\nError: %v\n", err)
	}
	if jwe.JweToken == "" {
		dump, _ := httputil.DumpRequest(req, false)
		log.Fatalf("No jwe token present\nRequest: %s\nBody: %s\nOutput: %s\n", string(dump), string(body), string(output))
	}

	jweToken = jwe.JweToken
	log.Println("Refreshed token")
}
