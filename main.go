package main

import (
	"fmt"
	"github.com/orsinium-labs/gamepad"
	"github.com/pebbe/zmq4"
	"github.com/shamaton/msgpack/v2"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		panic("args: driver controller id, operator controller id, refresh rate (hz)")
	}
	driverControllerId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("driver controller id (first arg) is not an integer")
	}
	operatorControllerId, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("operator controller id (second arg) is not an integer")
	}
	refreshRate, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("refresh rate (hz) is not an integer")
	}
	fmt.Printf("starting xbeezy-remote\ndriver @ %v\noperator @ %v\nsending status at %v hz\n", driverControllerId, operatorControllerId, refreshRate)

	driver, err := gamepad.NewGamepad(driverControllerId)
	if err != nil {
		if err.Error() == "Device not found" {
			fmt.Println("driver controller not found")
			driver = nil
		} else {
			panic(err)
		}
	}
	operator, err := gamepad.NewGamepad(operatorControllerId)
	if err != nil {
		if err.Error() == "Device not found" {
			fmt.Println("operator controller not found")
			operator = nil
		} else {
			panic(err)
		}
	}
	ticker := time.NewTicker(time.Second / time.Duration(refreshRate))
	sock, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		panic(err)
	}

	err = sock.Bind("tcp://*:5805") // port is open under "Team use" section of R704
	if err != nil {
		panic(err)
	}

	//num := 1

	for {
		<-ticker.C
		driverMsgState := ControllerState{}
		operatorMsgState := ControllerState{}

		if driver != nil {
			driverState, err := driver.State()
			if err != nil {
				panic(err)
			}

			driverMsgState = ControllerState{
				A:     driverState.A(),
				B:     driverState.B(),
				X:     driverState.X(),
				Y:     driverState.Y(),
				LB:    driverState.LB(),
				RB:    driverState.RB(),
				BACK:  driverState.Back(),
				START: driverState.Start(),
				GUIDE: driverState.Guide(),
				LSB:   driverState.LSB(),
				RSB:   driverState.RSB(),
				LSX:   driverState.LS().X,
				LSY:   driverState.LS().Y,
				RSX:   driverState.RS().X,
				RSY:   driverState.RS().Y,
				DL:    driverState.DPadLeft(),
				DR:    driverState.DPadRight(),
				DU:    driverState.DPadUp(),
				DD:    driverState.DPadDown(),
				LT:    driverState.LT(),
				RT:    driverState.RT(),
			}
		}

		if operator != nil {
			operatorState, err := operator.State()
			if err != nil {
				panic(err)
			}

			operatorMsgState = ControllerState{
				A:     operatorState.A(),
				B:     operatorState.B(),
				X:     operatorState.X(),
				Y:     operatorState.Y(),
				LB:    operatorState.LB(),
				RB:    operatorState.RB(),
				BACK:  operatorState.Back(),
				START: operatorState.Start(),
				GUIDE: operatorState.Guide(),
				LSB:   operatorState.LSB(),
				RSB:   operatorState.RSB(),
				LSX:   operatorState.LS().X,
				LSY:   operatorState.LS().Y,
				RSX:   operatorState.RS().X,
				RSY:   operatorState.RS().Y,
				DL:    operatorState.DPadLeft(),
				DR:    operatorState.DPadRight(),
				DU:    operatorState.DPadUp(),
				DD:    operatorState.DPadDown(),
				LT:    operatorState.LT(),
				RT:    operatorState.RT(),
			}
		}

		bin, err := msgpack.Marshal(Message{
			Driver:   driverMsgState,
			Operator: operatorMsgState,
		})
		if err != nil {
			panic(err)
		}

		_, err = sock.SendBytes(bin, zmq4.DONTWAIT)
		if err != nil {
			panic(err)
		}

		//fmt.Println("sent message #" + strconv.Itoa(num))
		//num++
	}
}

type ControllerState struct {
	A     bool
	B     bool
	X     bool
	Y     bool
	LB    bool
	RB    bool
	BACK  bool
	START bool
	GUIDE bool
	LSB   bool
	RSB   bool
	LSX   int
	LSY   int
	RSX   int
	RSY   int
	DL    bool
	DR    bool
	DU    bool
	DD    bool
	LT    int
	RT    int
}

type Message struct {
	Driver   ControllerState
	Operator ControllerState
}
