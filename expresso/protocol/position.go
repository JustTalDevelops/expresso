package protocol

// BlockPos represents a block position.
type BlockPos [3]int32

// X returns the X of the position
func (p BlockPos) X() int32 {
	return p[0]
}

// Y returns the Y of the position
func (p BlockPos) Y() int32 {
	return p[1]
}

// Z returns the Z of the position
func (p BlockPos) Z() int32 {
	return p[2]
}

// ColumnPos represents a position of a column.
type ColumnPos [2]int32

// X returns the X of the column position.
func (p ColumnPos) X() int32 {
	return p[0]
}

// Z returns the Z of the column position.
func (p ColumnPos) Z() int32 {
	return p[1]
}
