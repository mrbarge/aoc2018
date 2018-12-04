package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Entry struct {
	year int
	month int
	day int
	hour int
	min int
	action string
}

type GuardState int
const (
	Awake = iota
	Asleep
)

func printEntry(e Entry) {
	fmt.Printf("[%d-%d-%d %d:%d] %s\n",e.year,e.month,e.day,e.hour,e.min,e.action)
}

func makeTime(e Entry) int {
	return e.min + (e.hour*100) + (e.day*10000) + (e.month*1000000)
}

func isAwake(e Entry) bool {
	return strings.Contains(e.action,"begins shift") || strings.Contains(e.action, "wakes up")
}

func isShiftBegin(e Entry) bool {
	return strings.Contains(e.action,"begins shift")
}

func isAsleep(e Entry) bool {
	return strings.Contains(e.action, "falls asleep")
}

func getGuardId(e Entry) string {
	r_log, _ := regexp.Compile(`^Guard (#\d+) begins shift$`)
	res_log := r_log.FindStringSubmatch(e.action)
	if res_log != nil {
		return res_log[1]
	} else {
		return ""
	}
}

func process(datafile string) []*Entry {
	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	r_log, _ := regexp.Compile(`^\[(\d+)-(\d+)-(\d+) (\d+):(\d+)\] (.+)$`)

	entries := make([]*Entry,0)
	for s.Scan() {
		line := s.Text()
		res_log := r_log.FindStringSubmatch(line)

		if res_log != nil {
			year, _ := strconv.Atoi(res_log[1])
			month, _ := strconv.Atoi(res_log[2])
			day, _ := strconv.Atoi(res_log[3])
			hour, _ := strconv.Atoi(res_log[4])
			min, _ := strconv.Atoi(res_log[5])
			action := res_log[6]

			e := Entry{year,month, day, hour, min, action }
			entries = append(entries, &e)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return makeTime(*entries[i]) < makeTime(*entries[j])
	})

	return entries

}

func processState(entries []*Entry) {

	currentState := Awake
	currentGuard := ""
	sleepStart := -1

	// make a count of guard awake minutes
	guardstats := make(map[string][]int,0)

	for g := 0; g < len(entries); g++ {
		entry := *entries[g]
		if isShiftBegin(entry) {
			// first find out if previous guard was asleep, and record the sleep time til hour end
			if currentState == Asleep {
				// make a timecard for the existing guard if it's the very first one
				if _, ok := guardstats[currentGuard]; !ok {
					guardstats[currentGuard] = make([]int,60)
				}
				fmt.Printf("Recording diff between %d and %d for %s\n", sleepStart, 59, currentGuard)
				for i := sleepStart; i < 60; i++ {
					guardstats[currentGuard][i] += 1
				}
			}
			currentGuard = getGuardId(entry)
			// make a timecard for the new guard
			if _, ok := guardstats[currentGuard]; !ok {
				guardstats[currentGuard] = make([]int,60)
			}
			currentState = Awake
			fmt.Printf("Shift begins for %s at %d:%d\n",currentGuard,entry.hour,entry.min)
		} else if isAwake(entry) {
			if currentState == Asleep {
				// Flag the period in which the guard was asleep
				if sleepStart > entry.min {
					// Well this is weird!
					fmt.Println("Something strange!")
					printEntry(entry)
				} else {
					fmt.Printf("Recording diff between %d and %d for %s\n", sleepStart, entry.min, currentGuard)
					for i := sleepStart; i < entry.min; i++ {
						guardstats[currentGuard][i] += 1
					}
				}
				currentState = Awake
			} else {
				fmt.Println("ERR: Guard is waking when already awake")
			}
		} else if isAsleep(entry) {
			// flag that guard is asleep and record the minute
			currentState = Asleep
			// If falling asleep before the midnight hour begins, sleep start starts at 0
			if entry.hour > 0 {
				sleepStart = 0
			} else {
				sleepStart = entry.min
			}
		}
	}

	sleepyguard := findSleepiestGuard(guardstats)
	fmt.Printf("Sleepiest guard is %s\n",sleepyguard)
	sleepymin := findSleepiestMinute(sleepyguard, guardstats)
	fmt.Printf("Sleepiest minute is %d\n",sleepymin)
	p2g, p2m := findGuardWithTopMinute(guardstats)
	fmt.Printf("Sleepiest guard-min combo is %s on minute %d\n",p2g,p2m)
}

func findSleepiestMinute(guard string, guardstats map[string][]int) int {
	topMin := -1
	topSleeps := -1
	for i := 0; i < 60; i++ {
		if guardstats[guard][i] > topSleeps {
			topMin = i
			topSleeps = guardstats[guard][i]
		}
	}
	return topMin
}

func findSleepiestGuard(guardstats map[string][]int) string {
	topGuardId := ""
	maxMins := -1
	for guard, _ := range guardstats {
		gm := 0
		for i := 0; i < 60; i++ {
			gm += guardstats[guard][i]
		}
		if gm > maxMins {
			topGuardId = guard
			maxMins = gm
		}
	}
	return topGuardId
}

func findGuardWithTopMinute(guardstats map[string][]int) (string, int) {
	topGuardId := ""
	topMinute := -1
	topMinuteCount := 0

	for i := 0; i < 60; i++ {
		for guard, _ := range guardstats {
			if guardstats[guard][i] > topMinuteCount {
				topGuardId = guard
				topMinuteCount = guardstats[guard][i]
				topMinute = i
			}
		}
	}
	return topGuardId, topMinute
}


func main() {
	entries := process(os.Args[1])
	for _, e := range entries {
		printEntry(*e)
	}
	processState(entries)
}
