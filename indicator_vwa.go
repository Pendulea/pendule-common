package pcommon

type emptyState struct{}

func (state *emptyState) buildVWA(unit Data, quantity Data) *Point {
	return nil
}
