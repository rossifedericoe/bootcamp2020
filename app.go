package main

import(
	"fmt"
	"pruebaclaseuno/calculadora"
	"time"
)

func main()  {
	num1 := 1
	num2 := 2
	resultado, saludo := calculadora.Sumar(num1, num2)
	fmt.Println(resultado)
	fmt.Println(saludo)

	for i:=0; i<10; i++ {
		fmt.Println(i)
	}

	nombres := []string{"Juan","Pedro"}
	for i:=0; i<2; i++ {
		fmt.Println(nombres[i])
	}

	for i, nombre := range nombres {
		fmt.Println(i, nombre)
	}



	limite := 10
	for limite < 50 {
		limite = limite + 1
	}
	fmt.Println("termino")

	var fecha time.Time
	var texto *string

	fmt.Println(fecha)
	if texto == nil {
		fmt.Println("Texto es null o nil")
	}

	var numero int = 80
	recibirNumero(numero)
}

func recibirNumero(numero int)  {
	fmt.Println(numero)
	fmt.Println(&numero)

	recibirNulleable(&numero)
}

func recibirNulleable(numero *int)  {
	if numero == nil {
		fmt.Println("Me pasaste algo nulo!")
	} else {
		fmt.Println("Me pasaste: ", *numero)
	}

	// sql de la libreria standard
	// sqlx
	// GORM
}

