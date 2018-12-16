package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

type Register int
const (
	REG_A = 1
	REG_B = 2
	REG_C = 3
)

type Operation int
const (
	BANI = iota		// LOCKED 0
	GTRI			// LOCKED 1
	SETI			// LOCKED 2
	EQIR			// LOCKED 3
	EQRR			// LOCKED 4
	BORR			// LOCKED 5
	BORI			// LOCKED 6
	BANR			// LOCKED 7
	MULI			// LOCKED 8
	EQRI			// LOCKED 9
	MULR			// LOCKED 10
	GTRR			// LOCKED 11
	SETR			// LOCKED 12
	ADDR			// LOCKED 13
	GTIR			// LOCKED 14
	ADDI			// LOCKED 15
)

func process(datafile string) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_reg_state, _ := regexp.Compile(`^(.+):\s+\[(\d+), (\d+), (\d+), (\d+)\]$`)
	r_reg_op, _ := regexp.Compile(`^(\d+)\s+(\d+)\s+(\d+)\s+(\d+)$`)

	threeOpCount := 0
	for s.Scan() {
		line := s.Text()
		res_reg_state := r_reg_state.FindStringSubmatch(line)

		reg_before := make([]int,4)
		reg_after := make([]int, 4)
		op_data := make([]int, 4)

		if res_reg_state != nil && res_reg_state[1] == "Before" {
			reg_before[0], _ = strconv.Atoi(res_reg_state[2])
			reg_before[1], _ = strconv.Atoi(res_reg_state[3])
			reg_before[2], _ = strconv.Atoi(res_reg_state[4])
			reg_before[3], _ = strconv.Atoi(res_reg_state[5])

			s.Scan()
			opline := s.Text()
			res_reg_op := r_reg_op.FindStringSubmatch(opline)
			s.Scan()
			afterstate := s.Text()
			res_reg_afterstate := r_reg_state.FindStringSubmatch(afterstate)

			if res_reg_op == nil || res_reg_afterstate == nil {
				fmt.Println("PArsing error encountered.. check input.")
				continue
			}

			reg_after[0], _ = strconv.Atoi(res_reg_afterstate[2])
			reg_after[1], _ = strconv.Atoi(res_reg_afterstate[3])
			reg_after[2], _ = strconv.Atoi(res_reg_afterstate[4])
			reg_after[3], _ = strconv.Atoi(res_reg_afterstate[5])

			op_data[0], _ = strconv.Atoi(res_reg_op[1])
			op_data[1], _ = strconv.Atoi(res_reg_op[2])
			op_data[2], _ = strconv.Atoi(res_reg_op[3])
			op_data[3], _ = strconv.Atoi(res_reg_op[4])

			numOps := countSuccessfulOps(reg_before,reg_after,op_data)
			if numOps >= 3 {
				threeOpCount += 1
			}
			//printSuccessfulSingleOps(reg_before,reg_after,op_data)
		}
	}

	fmt.Printf("Part 1: %d\n",threeOpCount)
}

func countSuccessfulOps(input []int, output []int, op []int) int {

	retval := 0
	for _, o := range []Operation{ADDR,ADDI,MULR,MULI,BANR,BANI,BORR,BORI,SETR,SETI,GTIR,GTRI,GTRR,EQIR,EQRI,EQRR} {
		test := operate(input,op,o)
		if reflect.DeepEqual(test,output) {
			retval += 1
		}
	}
	return retval
}

func printSuccessfulSingleOps(input []int, output []int, op []int) int {

	retval := 0
	var retop []Operation

	for _, o := range []Operation{ADDR,ADDI,MULR,BANR,BANI,BORR,BORI,SETI,MULI,SETR,GTIR,GTRI,GTRR,EQIR,EQRI,EQRR} {
		test := operate(input,op,o)
		if reflect.DeepEqual(test,output) {
			retval += 1
			retop = append(retop, o)
		}
	}

	if retval == 5 {
		fmt.Print("Testing : ")
		fmt.Print(input)
		fmt.Print(" against op ")
		fmt.Print(op)
		fmt.Print(" for output ")
		fmt.Print(output)
		fmt.Print( " was successful with opcodes ")
		fmt.Print(retop)
		fmt.Println()
	}
	return retval
}

func operate(input []int, op_value []int, op Operation) []int {
	retreg := make([]int, len(input))
	retreg[0] = input[0]
	retreg[1] = input[1]
	retreg[2] = input[2]
	retreg[3] = input[3]
	//copy(retreg, input)

	switch op {
	case ADDR:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] + input[op_value[REG_B]]
	case ADDI:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] + op_value[REG_B]
	case MULR:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] * input[op_value[REG_B]]
	case MULI:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] * op_value[REG_B]
	case BANR:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] & input[op_value[REG_B]]
	case BANI:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] & op_value[REG_B]
	case BORR:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] | input[op_value[REG_B]]
	case BORI:
		retreg[op_value[REG_C]] = input[op_value[REG_A]] | op_value[REG_B]
	case SETR:
		retreg[op_value[REG_C]] = input[op_value[REG_A]]
	case SETI:
		retreg[op_value[REG_C]] = op_value[REG_A]
	case GTIR:
		if op_value[REG_A] > input[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case GTRI:
		if input[op_value[REG_A]] > op_value[REG_B] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case GTRR:
		if input[op_value[REG_A]] > input[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case EQIR:
		if op_value[REG_A] == input[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case EQRI:
		if input[op_value[REG_A]] == op_value[REG_B] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case EQRR:
		if input[op_value[REG_A]] == input[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	}
	return retreg
}

func process2(datafile string) {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_reg_op, _ := regexp.Compile(`^(\d+)\s+(\d+)\s+(\d+)\s+(\d+)$`)
	reg := []int{0,0,0,0}

	for s.Scan() {
		line := s.Text()
		res_reg_op := r_reg_op.FindStringSubmatch(line)

		if res_reg_op != nil {
			op_data := make([]int, 4)
			op_data[0], _ = strconv.Atoi(res_reg_op[1])
			op_data[1], _ = strconv.Atoi(res_reg_op[2])
			op_data[2], _ = strconv.Atoi(res_reg_op[3])
			op_data[3], _ = strconv.Atoi(res_reg_op[4])

			reg = operate(reg,op_data,Operation(op_data[0]))
		}
	}

	fmt.Printf("Part 2: %d\n",reg[0])
}

func main() {
	process("input.dat")
	process2("input2.dat")

}

