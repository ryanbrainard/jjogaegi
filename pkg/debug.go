package pkg

import (
	"strconv"
	"log"
)

func Debug(options map[string]string, msg string, a... interface{}) {
	debugStr, ok := options[OPT_DEBUG]
	if !ok {
		return
	}

	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		panic(err)
	}

	if !debug {
		return
	}

	log.Printf("DEBUG " + msg, a...)
}
