package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("Hello, World!")
	logger.Info("Welcome!!")
	logger.Fatal("Fatal statement.")
}
