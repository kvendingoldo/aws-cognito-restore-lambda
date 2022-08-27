package main

import (
	"context"
	"fmt"
	awsLambda "github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/kvendingoldo/aws-cognito-restore-lambda/internal/config"
	"github.com/kvendingoldo/aws-cognito-restore-lambda/internal/lambda"
	"github.com/kvendingoldo/aws-cognito-restore-lambda/internal/types"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetReportCaller(false)

	var formatter log.Formatter

	if formatterType, ok := os.LookupEnv("FORMATTER_TYPE"); ok {
		if formatterType == "JSON" {
			formatter = &log.JSONFormatter{PrettyPrint: false}
		}

		if formatterType == "TEXT" {
			formatter = &log.TextFormatter{DisableColors: false}
		}
	}

	if formatter == nil {
		formatter = &log.TextFormatter{DisableColors: false}
	}

	log.SetFormatter(formatter)

	var logLevel log.Level
	var err error

	if ll, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel, err = log.ParseLevel(ll)
		if err != nil {
			logLevel = log.DebugLevel
		}
	} else {
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
}

func Handler(ctx context.Context, event types.Event) (types.Response, error) {
	log.Infof("Handling lambda for event: %v", event)
	config := cfg.New(event)
	lambda.Execute(*config)
	return types.Response{Message: fmt.Sprintf("Lambda has been completed for %v\n", event.ID)}, nil
}

func main() {
	log.Info("Starting lambda execution ...")
	if mode, ok := os.LookupEnv("MODE"); ok {
		if !(mode == "local" || mode == "cloud") {
			log.Errorf("Environment variable 'MODE' has unknown value '%v'. Value should be 'local' or 'cloud'", mode)
			os.Exit(1)
		}

		if mode == "local" {
			config := cfg.New(nil)
			lambda.Execute(*config)
		} else if mode == "cloud" {
			awsLambda.Start(Handler)
		}
	} else {
		log.Errorf("Environment variable 'MODE' is unspecified. Please, specify it. Value should be 'local' or 'cloud'")
		os.Exit(1)
	}
	log.Info("Lambda has completed")
}
