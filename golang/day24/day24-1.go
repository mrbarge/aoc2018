package main

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"fmt"
)

type Effect int
const (
	COLD = iota
	SLASHING
	RADIATION
	FIRE
	BLUDGEONING
)

type Unit struct {
	hp int
}

type Group struct {
	id int
	army string
	units int
	weaknesses []Effect
	immunities []Effect
	attack Effect
	damage int
	initiative int
	hp int
}

// returns true if immune win
func process(datafile string, boost int) bool {

	file, _ := os.Open(datafile)
	s := bufio.NewScanner(file)

	re_group, _ := regexp.Compile(`^(\d+) units each with (\d+) hit points\s{0,1}(.+)\s{0,1}with an attack that does (\d+) (\w+) damage at initiative (\d+)$`)

	immuneConfig := true

	immuneGroups := make([]*Group, 0)
	infectGroups := make([]*Group, 0)
	groupId := 0

	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, "Immune System") {
			immuneConfig = true
			groupId = 1
			continue
		} else if strings.Contains(line, "Infection") {
			immuneConfig = false
			groupId = 1
			continue
		}

		res_group := re_group.FindStringSubmatch(line)
		if res_group != nil {
			uc, _ := strconv.Atoi(res_group[1])
			hp, _ := strconv.Atoi(res_group[2])
			modifier := res_group[3]
			damage, _ := strconv.Atoi(res_group[4])
			attack := res_group[5]
			initiative, _ := strconv.Atoi(res_group[6])

			g := new(Group)
			g.weaknesses, g.immunities = parseModifier(modifier)
			g.id = groupId
			groupId += 1
			g.hp = hp
			g.units = uc
			g.damage = damage
			g.attack = effectFromString(attack)
			g.initiative = initiative

			if immuneConfig {
				g.army = "Immune"
				g.damage += boost
				immuneGroups = append(immuneGroups, g)
			} else {
				g.army = "Infected"
				infectGroups = append(infectGroups, g)
			}
		}
	}

	rounds := 1
	done := false

	immuneWon := false
	for !done {

		/*
		fmt.Println("\nImmune System:")
		printState(immuneGroups)
		fmt.Println("Infection:")
		printState(infectGroups)
		*/
		targetmap := targetPhase(immuneGroups, infectGroups)

		stalemate := attackPhase(immuneGroups, infectGroups, targetmap)
		if stalemate {
			fmt.Printf("A stalemate has been reached.\n")
			i1 := countUnits(infectGroups)
			i2 := countUnits(immuneGroups)
			if i1 > i2 {
				fmt.Printf("Infected win, %d units remain\n",i1)
			} else {
				fmt.Printf("Immune win, %d units remain\n",i2)
				immuneWon = true
			}
			done = true
		}

		if isDefeated(immuneGroups) {
			fmt.Printf("Immune army defeated, %d infect units remain\n",countUnits(infectGroups))
			done = true
		} else if isDefeated(infectGroups) {
			fmt.Printf("Infect army defeated, %d immune units remain\n",countUnits(immuneGroups))
			done = true
			immuneWon = true
		}
		rounds += 1
	}
	return immuneWon
}

func printState(g []*Group) {
	for _, group := range g {
		if group.units <= 0 {
			continue
		}
		fmt.Printf("Group %d contains %d units\n",group.id, group.units)
	}
}

/* returns indication if stalemate reached */
func attackPhase(immunes []*Group, infects []*Group, targetmap map[*Group]*Group) bool {

	allGroups := make([]*Group, 0)
	allGroups = append(allGroups, immunes...)
	allGroups = append(allGroups, infects...)
	sortedOrderedAttackList := sortGroupsByInitiative(allGroups)

	killed := 0
	for _, group := range sortedOrderedAttackList {
		_, ok := targetmap[group]
		if ok {
			// we have a target
			target := targetmap[group]
			killed += attack(group,target)
		}
	}

	if killed == 0 {
		return true
	} else {
		return false
	}
}

func isDefeated(g []*Group) bool {

	for _, group := range g {
		if group.units > 0 {
			return false
		}
	}
	return true
}

func targetPhase(immunes []*Group, infects []*Group) map[*Group]*Group {

	sortedImmunes := sortGroupsByPower(immunes)
	sortedInfects := sortGroupsByPower(infects)

	retmap := make(map[*Group]*Group,0)

	for _, infect := range sortedInfects {

		// ignore defeated groups
		if infect.units < 0 {
			continue
		}

		bestDmg := -1
		tgtIdx := make([]int,0)
		for i, tgt := range immunes {
			if isTargeted(tgt, retmap) || tgt.units <= 0 {
				// ignore this one, its targeted by something else
				continue
			}

			// update list of all 'best' targets that qualify
			dmg := damageToGroup(infect, tgt)
			//fmt.Printf("Infected group %d would deal defending group %d %d damage\n",infect.id,tgt.id,dmg)

			if dmg > bestDmg {
				bestDmg = dmg
				tgtIdx = make([]int,0)
				tgtIdx = append(tgtIdx, i)
			} else if dmg == bestDmg {
				tgtIdx = append(tgtIdx,i)
			}
		}

		if bestDmg == -1 || len(tgtIdx) == 0 {
			// we could not pick anything to attack
			//fmt.Printf("%s Group ID %d can't find anything to attack\n",infect.army,infect.id)
			continue
		}

		// find the target with the best initiative
		retmap[infect] = decideTargetTie(immunes, tgtIdx)
	}

	for _, immune := range sortedImmunes {

		// ignore defeated groups
		if immune.units < 0 {
			continue
		}

		bestDmg := -1
		tgtIdx := make([]int,0)
		for i, tgt := range infects {
			if isTargeted(tgt, retmap)  || tgt.units <= 0 {
				// ignore this one, its targeted by something else
				continue
			}

			// update list of all 'best' targets that qualify
			dmg := damageToGroup(immune, tgt)
			//fmt.Printf("Immune group %d would deal defending group %d %d damage\n",immune.id,tgt.id,dmg)
			if dmg > bestDmg {
				bestDmg = dmg
				tgtIdx = make([]int,0)
				tgtIdx = append(tgtIdx, i)
			} else if dmg == bestDmg {
				tgtIdx = append(tgtIdx,i)
			}
		}

		if bestDmg == -1 || len(tgtIdx) == 0 {
			// we could not pick anything to attack
			//fmt.Println("Can't find anything to attack")
			continue
		}

		// find the target with the best initiative
		retmap[immune] = decideTargetTie(infects, tgtIdx)
	}

	return retmap
}

func decideTargetTie(g []*Group, indices []int) *Group {

	bestPower := -1
	powerIdx := make([]int, 0)
	for i, v := range indices {
		if i == 0 || power(g[i]) > bestPower {
			bestPower = power(g[i])
			powerIdx = append(powerIdx, v)
		}
	}

	if len(powerIdx) == 1 {
		return g[powerIdx[0]]
	}

	// tie for power, so move on to initiative
	bestInitiative := -1
	initIdx := make([]int, 0)
	for i, v := range indices {
		if i == 0 || g[i].initiative > bestInitiative {
			bestInitiative = g[i].initiative
			initIdx = append(initIdx, v)
		}
	}

	if len(initIdx) > 1 || len(initIdx) == 0 {
		// wut
	} else {
		return g[initIdx[0]]
	}

	return nil
}

func isTargeted(g *Group, targetmap map[*Group]*Group) bool {
	for _, v := range targetmap {
		if g == v {
			return true
		}
	}
	return false
}
func damageToGroup(attacker *Group, target *Group) int {

	if isImmune(attacker.attack, target) {
//		fmt.Printf("%s ID %d would attack %s ID %d but it is immune.\n",attacker.army,attacker.id,target.army,target.id)
		return 0
	} else if isWeak(attacker.attack, target) {
		//fmt.Printf("%s ID %d would attack %s ID %d and it is weak\n",attacker.army,attacker.id,target.army,target.id)
		return attacker.units * (attacker.damage * 2)
	} else {
		return attacker.units * attacker.damage
	}
}

func attack(attacker *Group, target *Group) int {

	if attacker.units <= 0 {
		return 0
	}

	dmg := damageToGroup(attacker, target)
	//fmt.Printf("%s Group ID %d attacks defending Group ID %d dealing damage value %d ...\n",attacker.army, attacker.id, target.id, dmg)
	//fmt.Printf("%s Target Units: %d, HP: %d, Damage Inflicted: %d, Thing: %f\n",target.army,target.units,target.hp,dmg,
//		(float64(target.units * target.hp) - float64(dmg)) / float64(target.units))

	casualties := int(math.Floor((float64(dmg) / float64(target.hp))))
	//fmt.Printf("%s Group ID %d lost %d units.\n",target.army, target.id, casualties)
	target.units -= casualties
	return casualties
}

func isImmune(e Effect, g *Group) bool {
	for _, w := range g.immunities {
		if w == e {
			return true
		}
	}
	return false
}

func isWeak(e Effect, g *Group) bool {
	for _, w := range g.weaknesses {
		if w == e {
			return true
		}
	}
	return false
}

func power(g *Group) int {
	return g.units * g.damage
}

func sortGroupsByPower(groups []*Group) []*Group {
	sort.Slice(groups, func(i, j int) bool {
		if power(groups[i]) == power(groups[j]) {
			return groups[i].initiative > groups[j].initiative
		} else {
			return power(groups[i]) > power(groups[j])
		}
	})
	return groups
}

func sortGroupsByInitiative(groups []*Group) []*Group {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].initiative > groups[j].initiative
	})
	return groups
}

func countUnits(g []*Group) int {
	retcnt := 0
	for _, group := range g {
		if group.units > 0 {
			retcnt += group.units
		}
	}
	return retcnt
}

func printGroup(g *Group) {
	fmt.Printf("%d units each with %d hit points (weak to %v, immune to %v) with an attack that does %d %d damage at initiative %d\n",
		g.units, g.hp, g.weaknesses, g.immunities, g.damage, g.attack, g.initiative)
}

/* returns two lists of weaknesses and immunities */
func parseModifier(s string) ([]Effect, []Effect) {

	weaknesses := make([]Effect,0)
	immunities := make([]Effect, 0)

	// get rid of brackets the lazy way
	s = strings.Replace(strings.Replace(s,"(","",-1),")","",-1)

	modifiers := strings.Split(s,";")
	for _, m := range modifiers {
		if strings.Contains(m,"weak") {
			// doing weaknesses
			effects := strings.Split(strings.Replace((strings.Split(m," to "))[1], " ","",-1), ",")
			for _, effect := range effects {
				weaknesses = append(weaknesses, effectFromString(effect))
			}
		}
		if strings.Contains(m,"immune") {
			// doing immunities
			effects := strings.Split(strings.Replace((strings.Split(m," to "))[1], " ","",-1), ",")
			for _, effect := range effects {
				immunities = append(immunities, effectFromString(effect))
			}
		}
	}

	return weaknesses, immunities
}

func effectFromString(e string) Effect {
	switch e {
	case "bludgeoning":
		return BLUDGEONING
	case "radiation":
		return RADIATION
	case "slashing":
		return SLASHING
	case "fire":
		return FIRE
	case "cold":
		return COLD
	default:
		return -1
	}
}

func main() {

	boost := 1

	done := false
	for !done {
		fmt.Printf("Boosting with %d\n", boost)
		done = process("input.dat", boost)
		if !done {
			boost += 1
		}
	}

//	process("test.dat")
}

