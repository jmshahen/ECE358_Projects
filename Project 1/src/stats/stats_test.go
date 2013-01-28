package stats

import (
	"testing"
)

func TestSimpleAvg(t *testing.T) {
	var a Avg
	var ans float64

	a.AddAvg(1)
	ans = 1.0 / 1.0
	if b := a.GetAvg(); b != ans {
		t.Errorf("Avg should be %f but got %f", ans, b)
	}

	a.AddAvg(2)
	ans = 3.0 / 2.0
	if b := a.GetAvg(); b != ans {
		t.Errorf("Avg should be %f but got %f", ans, b)
	}

	a.AddAvg(3)
	ans = 6.0 / 3.0
	if b := a.GetAvg(); b != ans {
		t.Errorf("Avg should be %f but got %f", ans, b)
	}

	a.AddAvg(2.2)
	ans = 8.2 / 4.0
	if b := a.GetAvg(); b != ans {
		t.Errorf("Avg should be %f but got %f", ans, b)
	}
}
