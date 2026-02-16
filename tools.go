package main

import (
	"math/rand"
	"os"
	"time"
)

var startup_time = time.Now()

func getRuntime() string {
	return time.Since(startup_time).String()
}

/* tokens */
// TODO: Remove this as soon as Redis will be used
func tokenExists(token string) bool {
	for _, session := range sessions {
		if session.Token == token {
			return true
		}
	}
	return false
}

func tokenGenerate() string {
	const length = 6
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano() & int64(os.Getpid())))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func tokenCreate() string {
	for {
		token := tokenGenerate()
		if !tokenExists(token) {
			return token
		}
	}
}

func sessionTokenToIndex(token string) int {
	for index, session := range sessions {
		if session.Token == token {
			return index
		}
	}
	return -1
}
