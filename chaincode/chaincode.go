package main

/* Imports
 * utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
  	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Event structure. Structure tags are used by encoding/json library
type Event struct {
	EventID                string   `json:"event_id"`
	EventName              string   `json:"event_name"`
	EventDetails           string   `json:"event_details"`
	DonatedAmount        int      `json:"donated"`
	EventDuration		   int   	`json:"event_duration"`
	EventStartDate		   string	`json:"event_start_date"`
	DonationIDs            []string `json:"donation_ids"`
}

// Define the Event structure. Structure tags are used by encoding/json library
type Donation struct {
	DonationID        string `json:"donation_id"`
	EventID           string `json:"event_id"`
	UserID     	      string `json:"user_id"`
	DonationAmount	  int    `json:"donated_amount"`

}

// Define the Event structure. Structure tags are used by encoding/json library
type User struct {
	UserID        string `json:"user_id"`
	Username      string `json:"user_name"`
	ContactNumber	int    `json:"contact_number"`
	EmailID       string `json:"mail_id"`
}

/*
 * The Init method is called when the Smart Contract is instantiated by the blockchain network
 * Here we have initiated SmartContract with the Event Details
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	var eventName string ="Project Aabhar"

	event := Event{EventID: "E1", EventName: eventName, DonatedAmount: 0 }
  	eventAsBytes, _ := json.Marshal(event)
  	APIstub.PutState("E1", eventAsBytes)
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
  	if function == "queryEvent" {
    	return s.queryEvent(APIstub, args)
   	}else if function == "donateMoney" {
   		return s.donateMoney(APIstub, args)
   	}else if function == "queryAllDonations" {
   		return s.queryAllDonations(APIstub)
   	}else if function == "queryAllUsers" {
   		return s.queryAllUsers(APIstub)
   	}
	return shim.Error("Invalid Smart Contract function name.")
}

/* Function to query a specific Event using a EventID */
func (s *SmartContract) queryEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	eventAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(eventAsBytes)
}

/* Function to query all the Donations made */
func (s *SmartContract) queryAllDonations(APIstub shim.ChaincodeStubInterface) sc.Response {

	var donationDetails []Donation
	var donationDetail Donation
	eventAsBytes, err := APIstub.GetState("E1")

	if err != nil {
		return shim.Error(err.Error())
	}

	var event Event
	json.Unmarshal(eventAsBytes, &event)

	var retrievedDonationIDs []string
	retrievedDonationIDs = event.DonationIDs
	fmt.Printf("- retrievedDonationIDs:\n%s\n", retrievedDonationIDs)

	for i := 0; i < len(retrievedDonationIDs); i++ {
		donationID := retrievedDonationIDs[i]

		donateAsBytes, err2 := APIstub.GetState(donationID)

		if err2 != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(donateAsBytes, &donationDetail)
		donationDetails = append(donationDetails, donationDetail)
	}
	fmt.Printf("- queryAllDonationsForEvent:\n%s\n", donationDetails)

	//change to array of bytes
	allDonatesAsBytes, _ := json.Marshal(donationDetails)
	return shim.Success(allDonatesAsBytes)
}

/* Function to query all the Users who have donated for the Event */
func (s *SmartContract) queryAllUsers(APIstub shim.ChaincodeStubInterface) sc.Response {

	var userDetails []User
	var userDetail User
	eventAsBytes, err := APIstub.GetState("E1")

	if err != nil {
		return shim.Error(err.Error())
	}

	var event Event
	json.Unmarshal(eventAsBytes, &event)

	var retrievedDonationIDs []string
	retrievedDonationIDs = event.DonationIDs
	fmt.Printf("- retrievedDonationIDs:\n%s\n", retrievedDonationIDs)

	for i := 0; i < len(retrievedDonationIDs); i++ {
		donationID := retrievedDonationIDs[i]

		donateAsBytes, err2 := APIstub.GetState(donationID)
    	var donation Donation
    	json.Unmarshal(donateAsBytes, &donation)

    	var uid= donation.UserID
    	userAsBytes, err2 := APIstub.GetState(uid)

		if err2 != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(userAsBytes, &userDetail)
		userDetails = append(userDetails, userDetail)
	}
	fmt.Printf("- queryAllusersForEvent:\n%s\n", userDetails)

	//change to array of bytes
	allUsersAsBytes, _ := json.Marshal(userDetails)
	return shim.Success(allUsersAsBytes)
}

/* Function to Donate money for the Crowdfunding Event */
func (s *SmartContract) donateMoney(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

  	var Contactnumber int
  	var existingVal int
  	var donatedQuantity int
	var donationIDArr []string
	var userID string
	var donationID string
	var newDonationID int
	var numOfDonationsMade int

	donatedQuantity, _ = strconv.Atoi(args[0])
  	eventAsBytes, _ := APIstub.GetState("E1")
  	event := Event{}
  	json.Unmarshal(eventAsBytes, &event)
	donationIDArr = event.DonationIDs
	numOfDonationsMade=len(donationIDArr)
	if numOfDonationsMade == 0{
		newDonationID = 100
	}else{
		split_donation_id := strings.Split(donationIDArr[numOfDonationsMade-1], "D")
		newDonationID,_=strconv.Atoi(split_donation_id[1])
	}
	newDonationID = newDonationID + 1
	donationID = "D"+strconv.Itoa(newDonationID)
	userID = "U"+strconv.Itoa(newDonationID)
  	existingVal = event.DonatedAmount
	
		event.DonatedAmount = existingVal + donatedQuantity
    	event.DonationIDs = append(event.DonationIDs, donationID)
    	event1AsBytes, _ := json.Marshal(event)
		APIstub.PutState("E1", event1AsBytes)

		var donate = Donation{DonationID: donationID, EventID: "E1", UserID: userID, DonationAmount: donatedQuantity}

    	donateAsBytes, _ := json.Marshal(donate)
  		APIstub.PutState(donationID, donateAsBytes)

    	Contactnumber, _ = strconv.Atoi(args[2])
    	var user = User{UserID: userID, Username: args[1], ContactNumber: Contactnumber, EmailID:args[3]}

    	userAsBytes, _ := json.Marshal(user)
  		APIstub.PutState(userID, userAsBytes)
		return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Included here for completeness only.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
