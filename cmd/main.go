package main

import (
	"flag"
	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)

}