package producer

import (
	"log"
	"os"
	"testing"
)

func TestRandomNumber(t *testing.T) {
	var lambda float64 = 200
	logger = log.New(os.Stdout, "", 0)
	time_2_tick = 1 // milliseconds

	for i := 0; i < 10; i++ {
		t.Logf("%d\n", getExpRandNum(lambda))
	}
}
