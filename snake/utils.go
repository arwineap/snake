package snake

import "math"

func checkPointCollision(p1 Point, p2 Point) bool {

	xDistance := math.Abs(float64(p1.X - p2.X))
	yDistance := math.Abs(float64(p1.Y - p2.Y))

	if xDistance < 4 && yDistance < 4 {
		return true
	}
	return false

	/*
		p1XMin := p1.X - 2
		p1XMax := p1.X + 2
		p1YMin := p1.Y - 2
		p1YMax := p1.Y + 2

		p2XMin := p2.X - 2
		p2XMax := p2.X + 2
		p2YMin := p2.Y - 2
		p2YMax := p2.Y + 2

		logger, _ := zap.NewProduction()

		if (p1XMin >= p2XMin && p1XMin <= p2XMax) || (p1XMax >= p2XMin && p1XMax <= p2XMax) {
			if (p1YMin >= p2YMin && p1YMin <= p2YMax) || (p1YMax >= p2YMin && p1YMax <= p2YMax) {
				a := p1XMin >= p2XMin && p1XMin <= p2XMax
				b := p1XMax >= p2XMin && p1XMax <= p2XMax
				c := p1YMin >= p2YMin && p1YMin <= p2YMax
				d := p1YMax >= p2YMin && p1YMax <= p2YMax
				logger.Info("collision", zap.Any("p1", p1), zap.Any("p2", p2), zap.Any("a", a), zap.Any("b", b), zap.Any("c", c), zap.Any("d", d))
				return true
			}
		}

		return false

	*/
}
