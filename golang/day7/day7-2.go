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

type Job struct {
	step string
	ticks int
	max int
}

func process(datafile string, duration int, workers int) {
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

	ticker := 0
	activeWorkers := 0
	activeJobs := make([]*Job, 0)
	next := starting
	for !allDone(keys,done) {

		fmt.Printf("SECOND %d\n", ticker)

		// can any active jobs be marked as done
		for i, job := range activeJobs {
			if job.ticks == job.max {
				done = append(done, job.step)
				activeJobs = append(activeJobs[:i], activeJobs[i+1:]...)
				activeWorkers -= 1
				fmt.Printf("Expiring job %s (%d == %d), activejobs=%d, activeworkers=%d\n", job.step, job.ticks, job.max,len(activeJobs),activeWorkers)
			}
		}

		// are any workers free
		foundNewJob := true
		for activeWorkers < workers && foundNewJob {

			foundNewJob = false

			// pick up a new job
			for _, k := range keys {
				// ignore active jobs
				alreadyProcessing := false
				for _, j := range activeJobs {
					if j.step == k {
						alreadyProcessing = true
					}
				}
				if alreadyProcessing {
					continue
				}

				next = findNext(k, &done, steps)
				if next != "" {
					j := Job{next, 0, duration + (int(k[0]) - 64) }
					fmt.Printf("Found a new job: %s, %d, %d\n", j.step,j.ticks,j.max)

					// add the new job to a worker
					activeWorkers += 1
					activeJobs = append(activeJobs, &j)
					foundNewJob = true
					break

				}
			}
		}

		// iterate the counter for active jobs
		for i := 0; i < len(activeJobs); i++ {
			jb := activeJobs[i]
			jb.ticks += 1
		}

		// increment ticker
		ticker += 1
	}

	fmt.Printf("The ticker value is: %d\n",ticker-1)
}

func allDone(keys []string, done []string) bool {
	for _, k := range keys {
		if !contains(done, k) {
			return false
		}
	}
	return true
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
//	process("test.dat", 60, 5)
	process("input.dat", 60, 5)
}

