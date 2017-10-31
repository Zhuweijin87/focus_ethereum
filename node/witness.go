package main

import (
	"fmt"
)

var (
	NumberOfWitness = 21
)

type Witness struct {
	Id 			discover.NodeID 
	joinDate	time.Time 
}

type WitnessDB struct {
	wits 		[NumberOfWitness]*Witness
}

type Node struct {
	...
	wits 		[NumberOfWitness]*Witness
	...
}

// 见证人管理 
func ApplyWitness() {
	
}

func CancelWitness() {

}