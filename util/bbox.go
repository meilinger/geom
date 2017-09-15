package util

//	TODO: maybe setup UL, LR as geom.Point to e more explicit?
type BoundingBox [4]float64

func BBox(points ...[2]float64) (bbox BoundingBox) {
	var xy [2]float64

	for i := range points {
		xy = points[i]

		if i == 0 {
			bbox[0] = xy[0]
			bbox[1] = xy[1]
			bbox[2] = xy[0]
			bbox[3] = xy[1]
			continue
		}

		switch {
		case xy[0] < bbox[0]:
			bbox[0] = xy[0]
		case xy[0] > bbox[2]:
			bbox[2] = xy[0]
		}

		switch {
		case xy[1] < bbox[1]:
			bbox[1] = xy[1]
		case xy[1] > bbox[3]:
			bbox[3] = xy[1]
		}
	}

	return
}
