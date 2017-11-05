package main

import (
	"fmt"
	"time"
)

type Data struct {
	id  int
	name string 
	register time.Time 
	remark string  
}

func (d Data) String() string {
	return fmt.Sprintf(`
	Data: 
	Id:			%d
	Name:		%s
	Register:	%v
	Remark:		%s
	`,
	d.id, d.name, d.register, d.remark)
}

func main() {
	d := Data{id:10, name: "Harrison", register: time.Now(), remark:"abcdefghijkl"}

	fmt.Println(d)
}