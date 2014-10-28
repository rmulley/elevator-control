package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
) //import

type Elevator struct {
	Id         int
	Floor      int
	GoalFloors []int
	Direction  int
	Score      int
} //Elevator

func main() {
	var (
		err       error
		rdr       *bufio.Reader
		elevators []Elevator
		// Flags
		numElevators *int = flag.Int("n", 4, "Number of elevators not set, defaulting to 4.")
	) //var

	flag.Parse()

	// Set correct number of elevators based on flag
	elevators = make([]Elevator, *numElevators)

	for i, _ := range elevators {
		// Initialize elevators with an ID and a goal floor
		elevators[i].Id = int(i)
		elevators[i].GoalFloors = make([]int, 1)
	} //for

	for {
		var (
			goalFloor, pickupFloor int
			floorStr               string
		) //var

		// Request input
		rdr = bufio.NewReader(os.Stdin)
		fmt.Println("What floor do you need to be picked up at?")

		// Read in floor that user needs to be picked up at
		if floorStr, err = rdr.ReadString('\n'); err != nil {
			log.Println("Error: Input must end in with a newline character, please try again.")
			continue
		} //if

		// Convert string that was read in into an integer
		if pickupFloor, err = strconv.Atoi(strings.TrimSpace(floorStr)); err != nil {
			log.Println("Error: Input must be a valid positive integer, please try again.")
			continue
		} //if

		// Request input
		rdr = bufio.NewReader(os.Stdin)
		fmt.Println("What floor would you like to go to?")

		// Read in floor that user wants to go to
		if floorStr, err = rdr.ReadString('\n'); err != nil {
			log.Println("Error: Input must end in with a newline character, please try again.")
			continue
		} //if

		// Convert string that was read in into an integer
		if goalFloor, err = strconv.Atoi(strings.TrimSpace(floorStr)); err != nil {
			log.Println("Error: Input must be a valid positive integer, please try again.")
			continue
		} //if

		// Assign highest scoring elevator to pickup new customer
		elevators = Pickup(pickupFloor, goalFloor, elevators)

		// Move the elevators one floor
		elevators = Step(elevators)

		// Basic formatting
		fmt.Println("ID    Floor    Direction    Goal-Floors")
		fmt.Println("---------------------------------------")

		// Print out elevator statuses
		for _, elevator := range elevators {
			fmt.Println(elevator.Id, "   ", elevator.Floor, "      ", elevator.Direction, "          ", elevator.GoalFloors)
		} //for
	} //for
} //main

// Status returns the ID, current floor number, and current goal floor number.
func (this Elevator) Status() (Id, Floor, GoalFloor int) {
	return this.Id, this.Floor, this.GoalFloors[0]
} //Status

// Update alters the current floor and direction for a single Elevator.
func (this Elevator) Update(curFloor int, direction int) {
	this.Floor = curFloor
	this.Direction = direction
} //Update

// Pickup algorithmically determines which elevator will pick up the the person requesting a ride
func Pickup(pickupFloor, goalFloor int, elevators []Elevator) []Elevator {
	var (
		idChosenElevator     int
		curClosest, topScore int = 1000, 0
		closestElevators     map[int]bool
	) //var

	// Determine which elevator is best to pick up the person at the pickup floor
	// The algorithm is:
	// +1 point for being the closest elevator
	// +3 points for an elevator that is currently not in use (AKA no goal floors)
	// +1 points for going the same direction

	for i, elevator := range elevators {
		var (
			floorsAway float64 = math.Abs(float64(pickupFloor) - float64(goalFloor))
		) //var

		// Reset elevator score
		elevator.Score = 0

		// If elevator on that floor and either not doing something or going the same way, then just choose it and exit algo
		if elevator.Floor == pickupFloor {
			if len(elevator.GoalFloors) == 0 || (pickupFloor > elevator.Floor && elevator.Direction > 0) || (pickupFloor < elevator.Floor && elevator.Direction < 0) {
				if pickupFloor > elevator.Floor {
					elevator.Direction = 1
				} else if pickupFloor <= elevator.Floor {
					elevator.Direction = -1
				} //else if

				elevator.GoalFloors = append(elevator.GoalFloors, pickupFloor, goalFloor)
			} //else if

			elevators[i] = elevator
			return elevators
		} //if

		// Keep track of closest elevator(s).
		// If new closest is found, forget previous closest elevators
		if int(floorsAway) < curClosest {
			closestElevators = make(map[int]bool, 0)
			closestElevators[elevator.Id] = true
		} else if int(floorsAway) == curClosest { // Multiple elevators are closest
			closestElevators[elevator.Id] = true
		} //else if

		// +3 points for an elevator not in use
		if len(elevator.GoalFloors) == 0 {
			elevator.Score += 3
		} //if

		// +2 points for an elevator going the same direction
		if pickupFloor > elevator.Floor && elevator.Direction > 0 {
			elevator.Score++
		} else if pickupFloor < elevator.Floor && elevator.Direction < 0 {
			elevator.Score++
		} //else if

		elevators[i] = elevator
	} //for

	// Find highest scored elevator
	for i, elevator := range elevators {
		// Add additional point to score for any elevator
		if _, isSet := closestElevators[elevator.Id]; isSet {
			elevator.Score++
		} //if

		if elevator.Score > topScore {
			idChosenElevator = elevator.Id
			topScore = elevator.Score
		} //if

		elevators[i] = elevator
	} //for

	log.Println("Elevator", idChosenElevator, "has been chosen.")

	// Update newest scored elevator
	for i, elevator := range elevators {
		if elevator.Id == idChosenElevator {
			if pickupFloor > elevator.Floor {
				elevator.Direction = 1
			} else if pickupFloor <= elevator.Floor {
				elevator.Direction = -1
			} //else if

			elevator.GoalFloors = append(elevator.GoalFloors, pickupFloor, goalFloor)
			elevators[i] = elevator
			break
		} //if
	} //for

	return elevators
} //Pickup

// Step increments each elevators position if they have a goal floor which they are not currently at.
func Step(elevators []Elevator) []Elevator {
	for i, elevator := range elevators {
		// Elevator is not moving, has nowhere to go
		if len(elevator.GoalFloors) == 0 {
			continue
		} //if

		//elevator = SortGoalFloors(elevator)

		// If the elevator is not at it's current goal floor, then it should keep moving
		if elevator.Floor != elevator.GoalFloors[0] {
			// If the goal floor is greater than the current floor, the elevator needs to go up
			if elevator.GoalFloors[0] > elevator.Floor {
				// Change the direction of the elevator if needed
				if elevator.Direction < 0 {
					elevator.Direction = 1
				} //if

				// Move the elevator up a floor
				elevator.Floor++
			} else if elevator.GoalFloors[0] < elevator.Floor { // If the goal floor is less than the current floor, the elevator needs to go down
				// Change the direction of the elevator if needed
				if elevator.Direction > 0 {
					elevator.Direction = -1
				} //if

				// Move the elevator down a floor
				elevator.Floor--
			} //else if
		} else { // If the elevator is at it's current goal floor, then there is nothing to do
			// Elevator floor is at goal floor, grab next floor (if available)
			if len(elevator.GoalFloors) > 1 {
				var (
					updatedGoalFloors []int
				) //var

				// Wasn't sure if there was a built-in that is the opposite of append (a 'pop' function)
				// Instead I wrote this loop to remove the goal floor the elevator is currently at
				if len(elevator.GoalFloors) > 1 {
					for ndx, goal := range elevator.GoalFloors {
						if ndx != 0 {
							updatedGoalFloors = append(updatedGoalFloors, goal)
						} //if
					} //for

					elevator.GoalFloors = updatedGoalFloors
				} else {
					elevator.GoalFloors = make([]int, 0)
				} //else
			} else { // Reset GoalFloors
				elevator.GoalFloors = make([]int, 0)
			} //else
		} //else

		elevators[i] = elevator
	} //for

	return elevators
} //Step

func SortGoalFloors(elevator Elevator) Elevator {
	var (
		goalFloorsExtra, goalFloorsSorted []int
	) //var

	// If no goal floors no work to do
	if len(elevator.GoalFloors) == 0 {
		return elevator
	} //if

	for _, goalFloor := range elevator.GoalFloors {
		if elevator.Direction > 0 && goalFloor >= elevator.Floor {
			goalFloorsSorted = append(goalFloorsSorted, goalFloor)
		} else if elevator.Direction < 0 && goalFloor <= elevator.Floor {
			goalFloorsSorted = append(goalFloorsSorted, goalFloor)
		} else {
			goalFloorsExtra = append(goalFloorsExtra, goalFloor)
		} //else
	} //for

	if elevator.Direction > 0 {
		sort.Ints(goalFloorsSorted)
	} else {
		//sort.Reverse(goalFloorsSorted)
	} //else

	elevator.GoalFloors = goalFloorsSorted
	elevator.GoalFloors = append(elevator.GoalFloors, goalFloorsExtra...)

	return elevator
} //SortGoalFloors
