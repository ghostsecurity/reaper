package workspace

type PortList []int

func (r PortList) Match(port int) bool {
	if len(r) == 0 {
		return true
	}
	for _, p := range r {
		if p == port {
			return true
		}
	}
	return false
}
