package snake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPointCollision(t *testing.T) {
	// checkPointCollision
	t.Run("random point", func(t *testing.T) {
		point1 := Point{X: 120, Y: 180}
		point2 := Point{X: 121, Y: 180}

		result := checkPointCollision(point1, point2)
		assert.True(t, result)
	})

	t.Run("game start", func(t *testing.T) {
		point1 := Point{X: 20, Y: 22}
		point2 := Point{X: 24, Y: 22}

		result := checkPointCollision(point1, point2)
		assert.False(t, result)
	})
}
