package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func stringFlag(name, value, usage string) {
	pflag.String(strings.ToLower(name), value, usage)
	viper.SetDefault(name, value)
}

func initViper() {
	stringFlag(
		"Proxy",
		"http://localhost:8001",
		"The proxy's location")
	stringFlag(
		"Listen",
		"localhost:8002",
		"Where to listen for connections")
	stringFlag(
		"Authenticator",
		"/usr/local/bin/aws-iam-authenticator",
		"The path the the AWS IAM Authenticator binary")
	stringFlag(
		"Cluster",
		"",
		"The name of the cluster to pass to the authentication")
	stringFlag(
		"Role",
		"",
		"The role ARN to pass to the authenticator")

	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/kdash/")
	viper.AddConfigPath("$HOME/.kdash")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Fatal error config file\nError: %v\n", err)
		}
	}

	viper.SetEnvPrefix("KDASH")
	viper.AutomaticEnv()

	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("Failed to bind flags\nError: %v\n", err)
	}

	pflag.Parse()
}


