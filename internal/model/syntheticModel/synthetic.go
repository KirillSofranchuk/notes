package syntheticModel

import "fmt"

// SyntheticType - Синтетический пример, для реализации функционала по работе с приватными полями
type SyntheticType struct {
	Id           int
	Name         string
	privateField int
}

func (s SyntheticType) GetPrivateField() int {
	fmt.Println("Access to private field")
	defer fmt.Println("Private field accessed")

	return s.privateField
}

func (s *SyntheticType) SetPrivateField(newValue int) {
	fmt.Println("Access to private field")
	defer fmt.Println("Private field changed")

	s.privateField = newValue
}
