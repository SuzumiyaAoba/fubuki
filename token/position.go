package token

import "fmt"

// Pos は字句の位置を表す構造体
type Pos struct {
	Line   int
	Column int
	Offset int
	Src    Source
}

func (pos *Pos) String() string {
	return fmt.Sprintf(
		"%d:%d:%d",
		pos.Line, pos.Column, pos.Offset,
	)
}
