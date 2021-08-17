package snake

import (
	"snake/point"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPointCollision(t *testing.T) {
	// CheckPointCollision
	t.Run("random point", func(t *testing.T) {
		point1 := point.Point{X: 120, Y: 180}
		point2 := point.Point{X: 121, Y: 180}

		result := point.CheckPointCollision(point1, point2)
		assert.True(t, result)
	})

	t.Run("game start", func(t *testing.T) {
		point1 := point.Point{X: 20, Y: 22}
		point2 := point.Point{X: 24, Y: 22}

		result := point.CheckPointCollision(point1, point2)
		assert.False(t, result)
	})
}
