package main

type DB map[int32]int32

func (d DB) insert(timestamp, price int32) {
	d[timestamp] = price
}

func (d DB) query(start, end int32) (int32, error) {
	if start > end {
		return 0, nil
	}

	var (
		samples = int64(0)
		total   = int64(0)
	)

	for ts, val := range d {
		if ts >= start && ts <= end {
			samples++
			total += int64(val)
		}
	}

	if samples == 0 {
		return 0, nil
	}

	return int32(total / samples), nil
}
