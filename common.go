package main

const (
	InstrNone int64 = iota
	InstrInteger
	InstrFloat
	InstrVar
	InstrAdd
	InstrPrint
	InstrLet
)

type YapType int64

const (
	IntegerT YapType = iota
	FloatT
	NoneT
)

type Stack []int64

func flipStack(stack *Stack) {
	for i := 0; i < len(*stack)/2; i++ {
		temp := (*stack)[i]
		(*stack)[i] = (*stack)[len(*stack)-i-1]
		(*stack)[len(*stack)-i-1] = temp
	}
}

func popStack(stack *Stack) int64 {
	e := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return e
}
