package cache

import "testing"

func Test_set(t *testing.T) {
	s := NewSet()
	n1 := People{Name: "n1", Medata: "123"}
	n11 := People{Name: "n1", Age: 10, Medata: "123"}
	n2 := People{Name: "n2", Medata: "123"}
	n3 := People{Name: "n2",Sex: 1, Medata: "123"}
	n4 := People{Name: "n4", Medata: "123"}
	n5 := People{Name: "n5", Medata: "123"}
	s.Add(n1, n1, n11, n11, n2, n2, n3, n3, n4, n4, n5, n5)
	if s.Len() != 6 {
		panic("set failed")
	}
}

type People struct {
	Name string
	Age int
	Sex int
	Medata string
}