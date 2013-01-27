package types

// Size of opcode arguments
const (
	sizeC  = 9
	sizeB  = 9
	sizeBx = sizeC + sizeB
	sizeA  = 8
	sizeAx = sizeC + sizeB + sizeA
	sizeOp = 6
)

// Position of opcode arguments
const (
	posOp = 0
	posA  = posOp + sizeOp
	posC  = posA + sizeA
	posB  = posC + sizeC
	posBx = posC
	posAx = posA
)

// INT_MAX - 2 "for safety"
const (
	MAX_INT    = (2147483647 - 2)
	MAXARG_Bx  = ((1 << sizeBx) - 1)
	MAXARG_sBx = (MAXARG_Bx >> 1)
)
