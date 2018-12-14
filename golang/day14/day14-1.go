package main

import (
	"fmt"
	"reflect"
)

func process(recipes int) []int {

	scores := make([]int, 0)
	scores = append(scores,3)
	scores = append(scores,7)

	// elf current recipe pointer
	elf1 := 0
	elf2 := 1

	printScores(scores,elf1,elf2)


	for numRecipes := 2; numRecipes < recipes; {
		numNewRecipes := 0
		scores, numNewRecipes = createRecipes(scores, elf1, elf2)
		elf1 = diff(elf1, (1 + scores[elf1]), len(scores))
		elf2 = diff(elf2, (1 + scores[elf2]), len(scores))
//		printScores(scores,elf1,elf2)
		numRecipes += numNewRecipes
	}

	return scores
}

func diff(a int, b int, wrap int) int {
	for i := 0; i < b; i++ {
		a += 1
		if a % wrap == 0 {
			a = 0
		}
	}
	return a
}

func printScores(scores []int, elf1 int, elf2 int) {
	for i, v := range scores {
		if i == elf1 {
			fmt.Print("*")
		}
		if i == elf2 {
			fmt.Print("^")
		}
		fmt.Printf("%d ", v)
	}
	fmt.Println("")
}

func createRecipes(scores []int, elf1 int, elf2 int) ([]int, int) {
	sum := scores[elf1] + scores[elf2]

	if sum > 9 {
		scores = append(scores,1)
		scores = append(scores, sum % 10)
		return scores, 2
	} else {
		scores = append(scores, sum)
		return scores, 1
	}
}

func findFirstScoreAppearance(scores []int, input []int) int {
	for i := 0; i < len(scores) - len(input); i++ {
		if reflect.DeepEqual(scores[i:i+len(input)],input) {
			return i
		}
	}
	return -1
}

func main() {
	scores := process(29381100)
	fmt.Printf("Part 1: ")
	fmt.Println(scores[len(scores)-10:])

	//fmt.Printf("Part 2: %d\n", findFirstScoreAppearance(scores,[]int{2,9,3,8,1,1}))
	fmt.Printf("Part 2: %d\n", findFirstScoreAppearance(scores,[]int{2,9,3,8,0,1}))
//	process(293801)
}

