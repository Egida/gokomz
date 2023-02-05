package utils

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

func GetRemoteIPFromRemoteAddr(remoteAddr string) string {
	// Remove port from remote address
	remoteIP := strings.Split(remoteAddr, ":")[0]
	// Remove brackets from IPv6 addresses
	remoteIP = strings.Trim(remoteIP, "[]")
	return remoteIP
}

// expontential backoff with jitter with ceiling
func ExpontentialBackoff(attempt int) time.Duration {
	// 10 second ceiling
	if attempt > 100 {
		attempt = 100
	}
	// 2^attempt
	milliseconds := 1 << uint(attempt)
	// add jitter
	milliseconds = milliseconds + rand.Intn(milliseconds)
	log.Println(fmt.Sprintf("Sleep jitter for %d milliseconds", milliseconds))
	// convert milliseconds to time duration
	return time.Duration(milliseconds) * time.Millisecond
}
