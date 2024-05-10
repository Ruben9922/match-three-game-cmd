package main

func isPointInsideGrid(p vector2d) bool {
	return p.x >= 0 && p.x < gridWidth && p.y >= 0 && p.y < gridHeight
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
