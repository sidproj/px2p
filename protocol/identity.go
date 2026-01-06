package protocol

import (
	"fmt"
	"math/rand"
	"time"
)

func GeneratePeerID() string {
	return fmt.Sprintf(
		"peer-%d-%d",
		time.Now().UnixNano(),
		rand.Intn(10000),
	)
}
