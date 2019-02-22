package main

import "testing"

func TestPowerLevel(t *testing.T) {
	tables := []struct {
		x            int
		y            int
		serialNumber int
		expected     int
	}{
		{3, 5, 8, 4},
		{122, 79, 57, -5},
		{217, 196, 39, 0},
		{101, 153, 71, 4},
	}

	for _, table := range tables {
		actual := PowerLevel(table.x, table.y, table.serialNumber)
		if actual != table.expected {
			t.Errorf("Sum was incorrect, got: %d, want: %d.", actual, table.expected)
		}

	}
}
