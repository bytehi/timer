package bunch

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := New()

	cancel1 := timer.Add(time.Second, func(cancel Cancel) {
		fmt.Println(time.Now(), "Task 1 executed")
		cancel()
	})

	timer.Add(2*time.Second, func(cancel Cancel) {
		fmt.Println(time.Now(), "Task 2 executed")
	})

	timer.Add(3*time.Second, func(cancel Cancel) {
		fmt.Println(time.Now(), "Task 3 executed")

		//something err, cancel 1
		cancel1()

		//cancel self
		cancel()
	})

	dead := time.Now().Add(10 * time.Second)
	fmt.Println(dead)
	for now := time.Now(); now.Before(dead); now = time.Now() {
		time.Sleep(time.Millisecond)
		timer.Timeout(time.Now())
	}
	fmt.Println("done")
}
