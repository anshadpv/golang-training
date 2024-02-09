package main

import (
	"fmt"

	mapset "github.com/deckarep/golang-set"
)

func main() {
	set1 := mapset.NewSet()
	set1.Add("A")
	set1.Add("B")
	set1.Add("C")
	set1.Add("D")

	set2 := mapset.NewSet()
	set2.Add("E")
	set2.Add("F")
	set2.Add("G")

	set3 := mapset.NewSet()
	set3.Add("H")
	set3.Add("I")
	set1.Add("D")

	unionSet := set1.Union(set2).Union(set3)
	fmt.Println("Union Set : ", unionSet)

	fmt.Println("Difference Set : ", unionSet.Difference(set1))

	fmt.Println(unionSet.Contains("E"))
	fmt.Println(set3.Cardinality())
}
