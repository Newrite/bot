package resource

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func Rolls(s string) string {
	var beforeD, afterD string
	var haveD, bDHaveNumber, aDHaveNumber bool = false, false, false
	for _, v := range s {
		switch v {
		case 'd':
			if !haveD {
				haveD = true
				beforeD += string(v)
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if haveD {
				aDHaveNumber = true
				afterD += string(v)
			}
			if !haveD {
				bDHaveNumber = true
				beforeD += string(v)
			}
		case '0':
			if haveD {
				if aDHaveNumber {
					afterD += string(v)
				}
			}
			if !haveD {
				if bDHaveNumber {
					beforeD += string(v)
				}
			}
		}
	}
	fmt.Println("DEBUG ", s, ": beforeD:", beforeD, " afterD:", afterD, " haveD:", strconv.FormatBool(haveD),
		" bDHaveNumber:", strconv.FormatBool(bDHaveNumber), " aDHaveNumber", strconv.FormatBool(aDHaveNumber))
	var numBD, numAD int = 0, 0
	switch {
	case haveD:
		beforeD = strings.TrimSuffix(beforeD, "d")
		fmt.Println("trim D: befdoreD:", beforeD)
		fallthrough
	case haveD && bDHaveNumber, bDHaveNumber && !haveD:
		numBD, _ = strconv.Atoi(beforeD)
		fmt.Println("Conv numBD:", numBD)
		fallthrough
	case haveD && aDHaveNumber, aDHaveNumber && !haveD:
		numAD, _ = strconv.Atoi(afterD)
		fmt.Println("Conv numAD:", numAD)
	}
	fmt.Println("numBD:", numBD, "numAD", numAD)
	if numBD != 0 {
		if numAD != 0 {
			return strconv.Itoa(numBD * (rand.Intn(numAD) + 1))
		} else {
			if haveD {
				return strconv.Itoa(numBD)
			} else {
				return strconv.Itoa(rand.Intn(numBD) + 1)
			}
		}
	} else {
		if numAD != 0 {
			return strconv.Itoa(rand.Intn(numAD) + 1)
		} else {
			return strconv.Itoa(rand.Intn(20) + 1)
		}
	}
}

type Commands struct {
	Command  []string
	Platform []string
	Channels []string
	Users    []string
	Request  string
}
