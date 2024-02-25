package main

import (
	"fmt"
	"sync"
	"time"
)

type BarberShop struct {
	waitingRoom  chan bool
	barberChair  chan bool
	barberAsleep bool
	closingTime  bool
	wg           sync.WaitGroup
	mutex        sync.Mutex
}

func (b *BarberShop) barber() {
	for {
		b.mutex.Lock()
		if len(b.waitingRoom) == 0 && !b.closingTime {
			fmt.Println("Barber is sleeping...")
			b.barberAsleep = true
		} else {
			b.barberAsleep = false
			select {
			case <-b.waitingRoom:
				fmt.Println("Barber is cutting hair.")
				time.Sleep(time.Second * 2) // Simulate haircut
				fmt.Println("Barber finished cutting hair.")
				b.barberChair <- true
			}
		}
		b.mutex.Unlock()

		if b.closingTime && len(b.waitingRoom) == 0 {
			break
		}
	}
	fmt.Println("Barber goes home.")
	b.wg.Done()
}

func (b *BarberShop) customer(id int) {
	b.mutex.Lock()
	if b.closingTime {
		b.mutex.Unlock()
		fmt.Printf("Customer %d arrived after closing time.\n", id)
		return
	}
	b.mutex.Unlock()

	select {
	case b.waitingRoom <- true:
		fmt.Printf("Customer %d entered the waiting room.\n", id)
		b.mutex.Lock()
		if b.barberAsleep {
			fmt.Printf("Customer %d woke up the barber.\n", id)
		}
		b.mutex.Unlock()

		select {
		case <-b.barberChair:
			fmt.Printf("Customer %d is getting a haircut.\n", id)
			time.Sleep(time.Second * 2) // Simulate haircut duration
			fmt.Printf("Customer %d left the barber shop.\n", id)
			b.mutex.Lock()
			b.mutex.Unlock()
		}
	default:
		fmt.Printf("Customer %d couldn't enter. No seats available.\n", id)
	}
}

func (b *BarberShop) CloseShop() {
	fmt.Println("Closing the shop...")
	b.mutex.Lock()
	b.closingTime = true
	b.mutex.Unlock()
}

func main() {
	shop := &BarberShop{
		waitingRoom: make(chan bool, 5), // 5 seats in the waiting room
		barberChair: make(chan bool, 1), // 1 barber chair
	}

	shop.wg.Add(1)
	go shop.barber()

	//  customers arriving at random intervals
	for i := 1; i <= 10; i++ {
		time.Sleep(time.Second * time.Duration(i))
		go shop.customer(i)
	}

	// Close shop after 10 seconds
	time.Sleep(time.Second * 10)
	shop.CloseShop()

	shop.wg.Wait() // Wait for barber to finish before exiting
}
