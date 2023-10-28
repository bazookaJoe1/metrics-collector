package pcstats

type Metrics []Metric

func (m Metrics) Len() int {
	return len(m)
}

func (m Metrics) Less(i, j int) bool {
	if m[i].MType < m[j].MType {
		return true
	} else if m[i].MType > m[j].MType {
		return false
	}
	return m[i].ID < m[j].ID
}

func (m Metrics) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
