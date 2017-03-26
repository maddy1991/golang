
// channel to simulate a game of cricket between two goroutines.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"os"
	"log"
	"regexp"
	"strconv"
)

// wg is used to wait for the program to finish.
var (wg sync.WaitGroup
 overs int32)
const name =`^[a-zA-Z]*$`
func init() {
	rand.Seed(time.Now().UnixNano())
}

func checkName(player string)bool{
	reName := regexp.MustCompile(name)
	return reName.Match([]byte(player))
}
func main() {
	if len(os.Args)<3{
		log.Fatal("example ./cricket bowlername batsmanname overs")
	}
	batsman:=os.Args[2]
	bowler:=os.Args[1]
	overs,err:=strconv.ParseInt(os.Args[3],10,32)
	if nil !=err ||overs<1 {
		log.Fatal("provide overs in valid format ,ex: 1")
	}
	valid:=checkName(batsman)
	if !valid{
		log.Fatal("provide valid name ,should contains only alphabets")
	}
	valid=checkName(bowler)
	if !valid{
		log.Fatal("provide valid name ,should contains only alphabets")
	}
	// Create an buffered channel.
	stadium := make(chan int32,1)
start:=make(chan bool,1)
	// Add a count of two, one for each goroutine.
	wg.Add(2)

	// Launch two players.
	go bowling(bowler, stadium,start)
	go batting(batsman, stadium,start)
	start<-true
	wg.Wait()
}

//bowling bowl number of overs entered by user
func bowling(name string, pitch chan int32,start chan bool) {
fmt.Printf("%s started bowling\n\n",name)

var ball int32=1
	defer wg.Done()
		for{
			// check batsman is out or not.
			out := <-start
			if !out {
				close(start)
				close(pitch)
				fmt.Printf(" %s bowled %d balls and match completed \n", name,ball-1)
				return
			}
			// check over is completed  or not.
			if ball == overs*6 {
				fmt.Printf("%d overs bowled by %s \n",overs, name)
				// Close the channel .
				close(pitch)
				return
			}
			pitch <- ball

ball++
		}

		// Hit the ball back to the opposing player.

}
//batting depending on the balls batsman will score runs by generating random number
func batting(name string, pitch chan int32,start chan bool) {
	fmt.Printf("%s started batting\n\n",name)
	defer wg.Done()
	var runs int
	var run int
	for{
			// Wait for the ball .
			ball, ok := <-pitch
			if !ok {
				close(start)
				// If the channel was closed over completed.
				fmt.Printf(" over completed %s scored %d runs \n", name, runs)
				return
			}
			switch {
			case ball > 6:
				{
					run= generateRuns(ball / 6)
				}
			case ball < 6:
				run= generateRuns(ball)
			}

			runs+=run
			// if generated number is zero then batsman out.
			if run == 0 {
				fmt.Printf("%s scored  %d runs and OUT!!!\n", name, runs)
				start<-false
				return
			}
		start<-true
		}

}
func generateRuns(balls int32) int {
	 run:=rand.Intn(int(balls*4))
	if run >6{
		return run -6
	}
return run
}