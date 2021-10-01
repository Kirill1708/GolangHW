package main

import ("golang.org/x/tour/tree"
		"fmt"
)

//type Tree struct {
//    Left  *Tree
//    Value int
//    Right *Tree
//}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int){
	if t != nil { 
		Walk(t.Left, ch)
		ch <- t.Value
		Walk(t.Right, ch)
	}
}


// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool{
	ch1, ch2 := make(chan int), make(chan int)
	go Walker(t1, ch1)
	go Walker(t2, ch2)

	for i := 0; i < 10; i++{
		x, y := <-ch1, <-ch2
		fmt.Printf("(%v %v)", x ,y)
		if x!= y{
			return false
		}
	}

	return true
}
func Walker(t*tree.Tree, ch chan int){
	Walk(t, ch)
	close(ch)
}

func main() {
	ch := make(chan int)
	go Walker(tree.New(1), ch)
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))

}
