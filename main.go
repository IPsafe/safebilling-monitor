package main

import "sync"
import "time"

//	"sync"

type message struct {
	Action                    string
	Channel_State             string
	Answer_State              string
	Channel_Call_State        string
	Call_Direction            string
	Event_Date_Local          string
	Caller_Caller_Id_Number   string
	Caller_Destination_Number string
	Channel_Call_Uuid         string
	Unique_Id                 string
}

type ActiveCalls struct {
	Calls []message
	Size  int
}

func NewActiveCalls() *ActiveCalls {
	return &ActiveCalls{Size: 0}
}

func (this *ActiveCalls) Add(new_call message) int {
	for _, call := range this.Calls {
		if call.Unique_Id == new_call.Unique_Id {
			return this.Size
		}
	}

	this.Calls = append(this.Calls, new_call)
	this.Size += 1
	return this.Size
}

func (this *ActiveCalls) Del(new_call message) int {
	for i, call := range this.Calls {
		if call.Unique_Id == new_call.Unique_Id {
			this.Calls = append(this.Calls[:i], this.Calls[i+1:]...)
			break
		}
	}
	this.Size = len(this.Calls)
	return this.Size
}

func (this *ActiveCalls) Update(new_call message) int {
	for i, call := range this.Calls {
		if call.Unique_Id == new_call.Unique_Id {
			this.Calls[i] = new_call
		}
	}

	return this.Size
}

var activeCalls = NewActiveCalls()

//var activeCalls []call

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	//	a, b := make(chan bool), make(chan bool)
	messages := make(chan message)

	go Websocket(wg, messages)
	//Monitor
	wg.Add(1)

	time.Sleep(5)

	go monitor(wg, messages)
	wg.Wait()

	//var input string
	//fmt.Scanln(&input)

}
