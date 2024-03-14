package main

const (
	InstrNone int64 = iota
	InstrInteger
	InstrFloat
	InstrVar
	InstrAdd
	InstrPrint
	InstrLet
	InstrFn
	InstrArg
)

type YapType int64

const (
	NoneT YapType = iota
	FunctionT
	IntegerT
	FloatT
)

type Stack []int64

func flipStack(stack *Stack) {
	for i := 0; i < len(*stack)/2; i++ {
		temp := (*stack)[i]
		(*stack)[i] = (*stack)[len(*stack)-i-1]
		(*stack)[len(*stack)-i-1] = temp
	}
}

func peekStack(stack *Stack) int64 {
	return (*stack)[len(*stack)-1]
}

func popStack(stack *Stack) int64 {
	e := peekStack(stack)
	*stack = (*stack)[:len(*stack)-1]
	return e
}
