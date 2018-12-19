package main

import (
	"bufio"
	"fmt"
	"os"
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

var operationMap = map[string]int{"bani":0,"gtri":1,"seti":2,"eqir":3,"eqrr":4,
	"borr":5,"bori":6,"banr":7,"muli":8,"eqri":9,"mulr":10,"gtrr":11,"setr":12,"addr":13,"gtir":14,"addi":15}

func process(datafile string, ip int) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_reg_op, _ := regexp.Compile(`^(\w+)\s+(\d+)\s+(\d+)\s+(\d+)$`)
	reg := []int{0,0,0,0,0,0}

	instructions := make([][]int, 0)
	for s.Scan() {
		line := s.Text()
		res_reg_op := r_reg_op.FindStringSubmatch(line)
		if res_reg_op != nil {
			op_data := make([]int, 4)
			op_data[0], _ = operationMap[res_reg_op[1]]
			op_data[1], _ = strconv.Atoi(res_reg_op[2])
			op_data[2], _ = strconv.Atoi(res_reg_op[3])
			op_data[3], _ = strconv.Atoi(res_reg_op[4])
			instructions = append(instructions, op_data)
		}
	}

	done := false
	nextInstruction := 0
	for !done {
		op_data := instructions[nextInstruction]
		reg,nextInstruction = operate(reg,op_data,Operation(op_data[0]), ip, nextInstruction)
		if nextInstruction  < 0 || nextInstruction  >= len(instructions) {
			done = true
		}
	}
	fmt.Printf("%v / %d / %d\n",reg,nextInstruction,len(instructions))
}

func process2(datafile string, ip int) {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_reg_op, _ := regexp.Compile(`^(\w+)\s+(\d+)\s+(\d+)\s+(\d+)$`)
	reg := []int{1,0,0,0,0,0}

	instructions := make([][]int, 0)
	for s.Scan() {
		line := s.Text()
		res_reg_op := r_reg_op.FindStringSubmatch(line)
		if res_reg_op != nil {
			op_data := make([]int, 4)
			op_data[0], _ = operationMap[res_reg_op[1]]
			op_data[1], _ = strconv.Atoi(res_reg_op[2])
			op_data[2], _ = strconv.Atoi(res_reg_op[3])
			op_data[3], _ = strconv.Atoi(res_reg_op[4])
			instructions = append(instructions, op_data)
		}
	}

	done := false
	nextInstruction := 0
	for !done {
		op_data := instructions[nextInstruction]
		reg,nextInstruction = operate(reg,op_data,Operation(op_data[0]), ip, nextInstruction)
		if nextInstruction  < 0 || nextInstruction  >= len(instructions) {
			done = true
		}
	}
	fmt.Printf("%v / %d / %d\n",reg,nextInstruction,len(instructions))

}

func operate(input []int, op_value []int, op Operation, ip int, ni int) ([]int,int) {
	retreg := make([]int, len(input))
	for i, _ := range input {
		retreg[i] = input[i]
	}

	// write the instruction pointer to the register it is bound to
	retreg[ip] = ni

	fmt.Printf("Before operation: %v\n", retreg)

	switch op {
	case ADDR:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] + retreg[op_value[REG_B]]
	case ADDI:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] + op_value[REG_B]
	case MULR:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] * retreg[op_value[REG_B]]
	case MULI:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] * op_value[REG_B]
	case BANR:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] & retreg[op_value[REG_B]]
	case BANI:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] & op_value[REG_B]
	case BORR:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] | retreg[op_value[REG_B]]
	case BORI:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]] | op_value[REG_B]
	case SETR:
		retreg[op_value[REG_C]] = retreg[op_value[REG_A]]
	case SETI:
		retreg[op_value[REG_C]] = op_value[REG_A]
	case GTIR:
		if op_value[REG_A] > retreg[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case GTRI:
		if retreg[op_value[REG_A]] > op_value[REG_B] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case GTRR:
		if retreg[op_value[REG_A]] > retreg[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case EQIR:
		if op_value[REG_A] == retreg[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case EQRI:
		if retreg[op_value[REG_A]] == op_value[REG_B] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	case EQRR:
		if retreg[op_value[REG_A]] == retreg[op_value[REG_B]] {
			retreg[op_value[REG_C]] = 1
		} else {
			retreg[op_value[REG_C]] = 0
		}
	}

	fmt.Printf("After operation: %v\n", retreg)
	nextInstruction := retreg[ip]
	return retreg, nextInstruction+1
}

func main() {
	process("input.dat",4)
}

