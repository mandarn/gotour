package bst

import "golang.org/x/tour/tree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var Walkrecurse func(t *tree.Tree)
	Walkrecurse = func(t *tree.Tree) {
		if t == nil {
			return
		}

		Walkrecurse(t.Left)
		ch <- t.Value
		Walkrecurse(t.Right)
	}

	Walkrecurse(t)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		val1, ok1 := <-ch1
		val2, ok2 := <-ch2

		if ok1 != ok2 || val1 != val2 {
			return false
		}

		if !ok1 {
			break
		}
	}

	return true
}

func TestBst() {
	ch := make(chan int)
	go Walk(tree.New(1), ch)

	for val := range ch {
		fmt.Printf("%d ", val)
	}

	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
