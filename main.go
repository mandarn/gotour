package main

import (
	"fmt"
	"github.com/mandarn/gotour/bst"
	"github.com/mandarn/gotour/crawl"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, val := range nums {
			out <- val
		}
		close(out)
	}()

	return out
}

func seq(chanl <-chan int) <-chan int {
	out := make(chan int)
	go func(chanl <-chan int) {
		for val := range chanl {
			out <- val * val
		}

		close(out)
	}(chanl)

	return out
}

func merge(chanl ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(chanel <-chan int) {
		for val := range chanel {
			out <- val
		}

		wg.Done()
	}

	wg.Add(len(chanl))

	for _, chanel := range chanl {
		go output(chanel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Seq(done, chanl <-chan int) <-chan int {
	out := make(chan int)
	go func(chanl <-chan int) {
		defer close(out)
		for val := range chanl {
			select {
			case out <- val:
			case <-done:
				return
			}
		}
	}(chanl)

	return out
}

func Merge(done <-chan int, chanl ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(chanel <-chan int) {
		defer wg.Done()
		for val := range chanel {
			select {
			case out <- val:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(chanl))

	for _, chanel := range chanl {
		go output(chanel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func simpleTest() {
	fmt.Println("Entering simpleTest \n")
	in := gen(2, 3)
	out := seq(in)

	for val := range out {
		fmt.Printf("%d \n", val)
	}
}

func Test1() {
	fmt.Println("Entering Test1 \n")
	in := gen(2, 3, 4, 5, 6)
	c1 := seq(in)
	c2 := seq(in)

	for val := range merge(c1, c2) {
		fmt.Printf("%d \n", val)
	}
}

func Test2() {
	fmt.Println("Entering Test2 \n")
	in := gen(2, 3)
	c1 := seq(in)
	c2 := seq(in)

	out := merge(c1, c2)
	fmt.Printf("%d \n", <-out)
	fmt.Printf("%d \n", <-out)
}

func Test3() {
	fmt.Println("Entering Test3 \n")
	done := make(chan int)
	defer close(done)

	in := gen(2, 3)
	c1 := Seq(done, in)
	c2 := Seq(done, in)

	out := Merge(done, c1, c2)
	fmt.Printf("%d \n", <-out)
}

func main() {
	fmt.Println("Testing pipelines")
	simpleTest()
	Test1()
	Test2()
	Test3()

	fmt.Println("Testing Bst")
	bst.TestBst()

	fmt.Println("Testing Crawler")
	crawl.TestCrawl()
}
