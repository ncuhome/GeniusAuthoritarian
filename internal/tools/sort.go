package tools

type UintSlice []uint

func (a UintSlice) Len() int {
	return len(a)
}
func (a UintSlice) Less(i, j int) bool {
	return a[i] < a[j]
}
func (a UintSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
