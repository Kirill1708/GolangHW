package main

import (
	"fmt"
)
type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprint("cannot Sqrt negative number: ", float64(e))	
} 

func Sqrt(e float64) (float64, error) {
	if e < 0 {
	return e, ErrNegativeSqrt(e)
	}	
	z := 1.0
	for i := 1; i <= 10; i++ {
		z = z - (z*z - e) / (2*z) //funcSqrt
	}
	return z, nil
}	
func main() {
	n, err := Sqrt(-2)
	fmt.Println(n)
    if err!=nil{
		fmt.Println(err)
	}
}
