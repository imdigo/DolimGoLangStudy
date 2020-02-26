package function_in_data_structure_test

import (
	"strings"
	"fmt"
	"strconv"
)

type BinOp func(int, int) int
type StrSet map[string]struct{}

// Map keyed by operator to set of higher precedence operators
type PrecMap map[string]StrSet

// Eval returns the evaluation result of the given expr.
// The expression can have +, -, *, /, (, ) operators and
// decimal integers. Operators and operands should be 
// space delimited.
func Eval(opMap map[string]BinOp, prec PrecMap, expr string) int {
	ops := []string{"("}
	var nums []int 
	pop := func() int {
		last := nums[len(nums) - 1]
		nums = nums[:len(nums) - 1]
		return last
	}
	reduce := func(nextOp string) {
		for len(ops) > 0 {
			op := ops[len(ops) - 1]
			if _, higher := prec[nextOp][op]; nextOp != ")" && !higher {
				// 더 낮은 순위의 연산자이므로 여기서 계산 종료
				return
			}
			ops = ops[:len(ops) - 1]
			if op == "(" {
				// 괄호를 제거하였으므로 종료
				return
			}
			b, a := pop(), pop()
			if f := opMap[op]; f != nil {
				nums = append(nums, f(a, b))
			}
		}
	}
	for _, token := range strings.Split(expr, " ") {
		if token == "(" {
			ops = append(ops, token)
		} else if _, ok := prec[token]; ok {
			reduce(token)
			ops = append(ops, token)
		} else if token == ")" {
			// 닫는 괄호는 여는 괄호까지 계산하고 제거
			reduce(token)
		} else {
			num, _ := strconv.Atoi(token)
			nums = append(nums, num)
		}
	}
	reduce(")")
	return nums[0]
}

// Returns a new StrSet.
func NewStrSet(strs ...string) StrSet {
	m := StrSet{}
	for _, str := range strs {
		m[str] = struct{}{}
	}
	return m
}

func NewEvaluator(opMap map[string]BinOp, prec PrecMap) func(expr string) int {
	return func(expr string) int {
		return Eval(opMap, prec, expr)
	}
}

func ExampleNewEvaluator() {
	eval := NewEvaluator(map[string]BinOp {
		"**": func (a, b int) int {
			if a == 1 {
				return 1
			}
			if b < 0 {
				return 0
			}
			r := 1
			for i := 0; i < b; i++ {
				r *= a
			}
			return r
		},
		"*": func (a, b int) int { return a * b },
		"/": func (a, b int) int { return a / b },
		"mod": func (a, b int) int { return a % b },
		"+": func (a, b int) int { return a + b },
		"-": func (a, b int) int { return a - b },
	}, PrecMap{
		"**": NewStrSet(),
		"*": NewStrSet("**", "*", "/", "mod"),
		"/": NewStrSet("**", "*", "/", "mod"),
		"mod": NewStrSet("**", "*", "/", "mod"),
		"+": NewStrSet("**", "*", "/", "mod", "+", "-"),
		"-": NewStrSet("**", "*", "/", "mod", "+", "-"),
	})

	exs := []string{
		"5",
		"1 + 2",
		"1 - 2 - 4",
		"( 3 - 2 ** 3 ) * ( -2 )",
		"3 * ( 3 + 1 * 3 ) / ( -2 )",
		"3 * ( ( 3 + 1 ) * 3 ) / 2",
		"1 + 2 ** 10 * 2",
		"2 ** 3 mod 3",
		"2 ** 2 ** 3",
	}
	for _, ex := range exs {
		fmt.Printf("%s = %d\n", ex, eval(ex))
	}
	// Output:
	// 5 = 5
	// 1 + 2 = 3
	// 1 - 2 - 4 = -5
	// ( 3 - 2 ** 3 ) * ( -2 ) = 10
	// 3 * ( 3 + 1 * 3 ) / ( -2 ) = -9
	// 3 * ( ( 3 + 1 ) * 3 ) / 2 = 18
	// 1 + 2 ** 10 * 2 = 2049
	// 2 ** 3 mod 3 = 2
	// 2 ** 2 ** 3 = 256

}
