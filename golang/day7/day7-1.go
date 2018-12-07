package main

import (
	"os"
	"bufio"
	"regexp"
	"fmt"
	"sort"
)

type Step struct {
	pre []string
	post []string
}

func process(datafile string) {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_log, _ := regexp.Compile(`^Step (.+) must be finished before step (.+) can begin.$`)

	steps := make(map[string]*Step)
	done := make([]string,0)

	for s.Scan() {
		line := s.Text()
		res_log := r_log.FindStringSubmatch(line)

		if res_log != nil {
			s1 := res_log[1]
			s2 := res_log[2]

			if _, ok := steps[s1]; !ok {
				steps[s1] = new(Step)
			}
			if _, ok := steps[s2]; !ok {
				steps[s2] = new(Step)
			}

			(steps[s2]).pre = append((steps[s2]).pre, s1)
		}
	}


	starting := ""

	keys := []string{}
	for key := range steps {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if len((steps[k]).pre) == 0 {
			starting = k
			break
		}
	}

	next := starting
	for next != "" {
		for _, k := range keys {
			next = findNext(k, &done, steps)
			if next != "" {
				fmt.Print(next)
				done = append(done, next)
				break
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func findNext(s string, done *[]string, steps map[string]*Step) string {
	// am I not done
	if ! contains(*done, s) {
		// and are my pre-requisites done
		preReqsDone := true
		for _, pre := range (steps[s]).pre {
			if ! contains(*done, pre) {
				// can't do this one yet
				preReqsDone = false
			}
		}
		if preReqsDone {
			// I can be done next
			return s
		}
	}
	// I have none to return
	return ""
}

func main() {
	process("input.dat")
}

