package main

import "fmt"

func main() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			// the two result form of channel receive comes with a boolean
			// which is false if the channel has been closed and no more items
			// are in it
			job, more := <-jobs

			if more {
				fmt.Println("job:", job)
			} else {
				fmt.Println("done all jobs")
				done <- true
				return
			}
		}
	}()

	for j:=1; j<=3; j++ {
		jobs <- j
		fmt.Println("sent:", j)
	}

	close(jobs)
	<-done
	fmt.Println("waited on worker")
}
