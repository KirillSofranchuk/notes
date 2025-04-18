package main

import (
	"fmt"
	"project/internal/model/syntheticModel"
)

func main() {
	fmt.Println("Run service")

	syntheticType := syntheticModel.SyntheticType{
		Id:   123,
		Name: "Example",
	}

	privateField := syntheticType.GetPrivateField()
	fmt.Println(privateField)

	syntheticType.SetPrivateField(123)
	privateField = syntheticType.GetPrivateField()
	fmt.Println(privateField)
}
