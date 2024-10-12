package main

type Robot struct {
	position position
	alive    bool
}

func (r *Robot) getRobotMove(playerPosition position) (int, int) {
	x := 0
	y := 0
	if r.alive {

		if playerPosition[0] < r.position[0] {
			x = -1
		} else if playerPosition[0] > r.position[0] {
			x = 1
		}

		if playerPosition[1] < r.position[1] {
			y = -1
		} else if playerPosition[1] > r.position[1] {
			y = 1
		}
	}
	return x, y
}

func (r *Robot) move(x, y int) {
	if r.alive {
		r.position[0] += x
		r.position[1] += y
	}
}
