package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fiorix/go-eventsocket/eventsocket"
)

//	c.Send(fmt.Sprintf("bgapi originate %s %s", dest, dialplan))

func monitor(wg *sync.WaitGroup, channel chan message) {
	fmt.Println("\nRodando Monitor:")

	var msg message
	//var body string
	var reader *csv.Reader
	var record [][]string
	var event *eventsocket.Event

	c, err := eventsocket.Dial("localhost:8021", "ClueCon")

	if err != nil {
		log.Fatal(err)
	}

	c.Send("event plain CHANNEL_CREATE CHANNEL_ANSWER CHANNEL_HANGUP CHANNEL_HANGUP_COMPLETE CALL_UPDATE CHANNEL_BRIDGE")

	//c.Send("event plain all")
	//	c.Send("event plain BACKGROUND_JOB")

	//c.Send("event plain CHANNEL_ANSWER BACKGROUND_JOB")

	//Only execute when start monitor
	event, _ = c.Send("api show calls")

	//fmt.Println(event.Body)

	reader = csv.NewReader(bytes.NewBufferString(strings.Split(event.Body, "\n\n")[0]))
	record, err = reader.ReadAll()
	fmt.Println(record)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i := 1; i < len(record); i++ {
		if record[i][1] == "outbound" {
			continue
		}
		fmt.Println(" ", record[i])

		msg.Action = "Add"
		msg.Unique_Id = record[i][0]
		msg.Call_Direction = record[i][1]
		msg.Event_Date_Local = record[i][2]
		msg.Caller_Caller_Id_Number = record[i][7]
		msg.Caller_Destination_Number = record[i][9]
		msg.Channel_Call_State = record[i][12]
		log.Println(activeCalls.Add(msg))
		log.Println(activeCalls.Calls[0].Unique_Id)
		channel <- msg
		//	fmt.Println(msg)

	}

	for {
		ev, err := c.ReadEvent()

		if err != nil {
			log.Fatal(err)
		}

		if ev.Get("Call-Direction") == "outbound" {
			continue
		}

		//ev.PrettyPrint()

		fmt.Println("\nNew event: ", ev.Get("Event-Name"))
		fmt.Println("Channel-State: ", ev.Get("Channel-State"))
		//	fmt.Println("Answer-State: ", ev.Get("Answer-State"))
		fmt.Println("Channel-Call-State: ", ev.Get("Channel-Call-State"))

		//		fmt.Println("Call-Direction: ", ev.Get("Call-Direction"))
		//		fmt.Println("Channel-Call-Uuid: ", ev.Get("Channel-Call-Uuid"))
		//		fmt.Println("Unique-Id: ", ev.Get("Unique-Id"))

		if ev.Get("Event-Name") == "CHANNEL_CREATE" {
			fmt.Println("Event-Date-Local: ", ev.Get("Event-Date-Local"))
			fmt.Println("Caller-Caller-Id-Number: ", ev.Get("Caller-Caller-Id-Number"))
			fmt.Println("Caller-Destination-Number: ", ev.Get("Caller-Destination-Number"))

			msg.Action = "Add"
			msg.Event_Date_Local = ev.Get("Event-Date-Local")
			msg.Caller_Caller_Id_Number = ev.Get("Caller-Caller-Id-Number")
			msg.Caller_Destination_Number = ev.Get("Caller-Destination-Number")

		}

		if ev.Get("Event-Name") == "CHANNEL_HANGUP" || ev.Get("Event-Name") == "CHANNEL_ANSWER" || ev.Get("Event-Name") == "CALL_UPDATE" || ev.Get("Event-Name") == "CHANNEL_BRIDGE" {
			msg.Action = "Mod"
		}

		if ev.Get("Event-Name") == "CHANNEL_HANGUP_COMPLETE" {
			msg.Action = "Del"
		}

		msg.Channel_State = ev.Get("Channel-State")

		//msg.Answer_State = ev.Get("Answer-State")
		msg.Channel_Call_State = ev.Get("Channel-Call-State")
		msg.Call_Direction = ev.Get("Call-Direction")
		//		msg.Channel_Call_Uuid = ev.Get("Channel-Call-Uuid")
		msg.Unique_Id = ev.Get("Unique-Id")

		if msg.Action == "Add" {
			log.Println(activeCalls.Add(msg))
			log.Println(activeCalls.Calls[0].Unique_Id)

		}

		if msg.Action == "Del" {
			log.Println(activeCalls.Del(msg))

		}

		if msg.Action == "Mod" {
			log.Println(activeCalls.Update(msg))

		}

		log.Println(activeCalls.Calls)

		channel <- msg

	}
	c.Close()
	wg.Done()
	//close(a)

}
