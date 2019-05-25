//This program pretends to be a SCIM Server, so that you can test your connection with the Okta OnPremise Provisioning Agent (CPC)
//based on code from: https://gist.github.com/denji/12b3a568f092ab951456, Thanks you !

package main


import (
	"net/http"
	"log"
	"fmt"
	//"os"
	"io/ioutil"
	//"reflect"
	"encoding/json"
	"github.com/gorilla/mux"
	"os"
	"encoding/csv"
	"io"
	"strconv"
)

var serverCrt = `-----BEGIN CERTIFICATE-----
MIIC/DCCAoSgAwIBAgIJAI1iiB15DlzCMAkGByqGSM49BAEwdzELMAkGA1UEBhMC
VVMxCzAJBgNVBAgTAlR4MQswCQYDVQQHEwJUeDENMAsGA1UEChMEbm9uZTENMAsG
A1UECxMEbm9uZTESMBAGA1UEAxMJbG9jYWxob3N0MRwwGgYJKoZIhvcNAQkBFg1v
andlZkB3b3cuY29tMB4XDTE3MDYxODIyMTY1OVoXDTI3MDYxNjIyMTY1OVowdzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgTAlR4MQswCQYDVQQHEwJUeDENMAsGA1UEChME
bm9uZTENMAsGA1UECxMEbm9uZTESMBAGA1UEAxMJbG9jYWxob3N0MRwwGgYJKoZI
hvcNAQkBFg1vandlZkB3b3cuY29tMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE7w10
DWX1l+WxSiy4VUR2ieKNRkEwbhC8+7NI5rvgVqTRoW/F7cKYAlMiQfiHlvxcIpYb
275rE2K3qsAyMnkGyIETKKDWtPfepzgPC0kwcjGR7kA0wr//ayq28VKpeSTEo4Hc
MIHZMB0GA1UdDgQWBBQ3VDD74iOqLscqJDRU5ZqoKljEpzCBqQYDVR0jBIGhMIGe
gBQ3VDD74iOqLscqJDRU5ZqoKljEp6F7pHkwdzELMAkGA1UEBhMCVVMxCzAJBgNV
BAgTAlR4MQswCQYDVQQHEwJUeDENMAsGA1UEChMEbm9uZTENMAsGA1UECxMEbm9u
ZTESMBAGA1UEAxMJbG9jYWxob3N0MRwwGgYJKoZIhvcNAQkBFg1vandlZkB3b3cu
Y29tggkAjWKIHXkOXMIwDAYDVR0TBAUwAwEB/zAJBgcqhkjOPQQBA2cAMGQCMCYn
JfaGh4l1HASI5Ma3iq4Va+UW+xlDZPJOLsVcn8e4fnqnR1nh9cugxgFP63lFiwIw
FznHs3iE385jr9VJKjitMqz/XLNdAjB4PNoPcxV49Am+3ALJkKs7TNCOyIj2GUtF
-----END CERTIFICATE-----
`

var serverKey = `-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDBEGm98AadjIYZRM8RUjC0e3fzeSHMQL/zY2iDHSvc86Ys/Z7TxH2oo
/Trm888vQI+gBwYFK4EEACKhZANiAATvDXQNZfWX5bFKLLhVRHaJ4o1GQTBuELz7
s0jmu+BWpNGhb8XtwpgCUyJB+IeW/FwilhvbvmsTYreqwDIyeQbIgRMooNa0996n
OA8LSTByMZHuQDTCv/9rKrbxUql5JMQ=
-----END EC PRIVATE KEY-----
`

func init() {

	//Check to see if the file users.csv exists, if it does not, then create it

	if _, err := os.Stat("users.csv"); os.IsNotExist(err) {
		fmt.Println("Creating users.csv file")
		usercsv := []byte("1111,username@mailinator.com ,Password1,username@mailinator.com,firstname,lastname,full name\n\n")
		_ = ioutil.WriteFile("users.csv", usercsv, 0644)

	} else {
		fmt.Println("Found users.csv file")

	}

	fmt.Println("Staring SCIM server on https://localhost")

	//Check to see if Certs are available for Https/TLS

	if _, err := os.Stat("server.crt"); os.IsNotExist(err) {
		//server.crt
		fmt.Println("Creating server.crt file")
		usercsv := []byte(serverCrt)
		_ = ioutil.WriteFile("server.crt", usercsv, 0644)

	} else {
		fmt.Println("Found server.crt file")

	}

	if _, err := os.Stat("server.key"); os.IsNotExist(err) {
		//server.key
		fmt.Println("Creating server.key file")
		usercsv := []byte(serverKey)
		_ = ioutil.WriteFile("server.key", usercsv, 0644)

	} else {
		fmt.Println("Found server.key file")

	}

}

func serverConfigs(w http.ResponseWriter, req *http.Request) {
	//This pretends it can do every operation, it is a lie..
	fmt.Println("hit Server Configs")

	w.Header().Set("Content-Type", "text/plain")
	var output = `
	{
  "schemas": ["urn:scim:schemas:core:1.0","urn:okta:schemas:scim:providerconfig:1.0"],
  "documentationUrl":"https://support.okta.com/scim-fake-page.html",
  "patch": {
    "supported":false
  },
  "bulk": {
    "supported":false
  },
  "filter": {
    "supported":true,
    "maxResults": 100
  },
  "changePassword" : {
    "supported":true
  },
  "sort": {
    "supported":false
  },
  "etag": {
    "supported":false
  },
  "authenticationSchemes": [],
  "urn:okta:schemas:scim:providerconfig:1.0": {    "userManagementCapabilities": ["GROUP_PUSH", "IMPORT_NEW_USERS", "IMPORT_PROFILE_UPDATES", "PUSH_NEW_USERS", "PUSH_PASSWORD_UPDATES", "PUSH_PENDING_USERS", "PUSH_PROFILE_UPDATES", "PUSH_USER_DEACTIVATION", "REACTIVATE_USERS"    ]  }}
`
	w.Write([]byte(output))

}

func groups(w http.ResponseWriter, req *http.Request) {
	// Pretends we have groups, we aren't using them, alway empty
	fmt.Println("hit groups")

	w.Header().Set("Content-Type", "text/plain")
	var output = `{
   "totalResults":0,
   "schemas":[
      "urn:scim:schemas:core:1.0"
   ],
   "Resources":[ ]
}`
	w.Write([]byte(output))

}

func users(w http.ResponseWriter, req *http.Request) {

	fmt.Println("hit users controller")

	switch req.Method {

	case "POST": //Post equal update and/or Create User

		fmt.Println("Post request sent to User Controller POST=Create User")
		body, err := ioutil.ReadAll(req.Body)
		_ = body
		if err != nil {
			panic(err)
		}

		var inboundUserObject inboundUserObj

		jsonError := json.Unmarshal(body, &inboundUserObject) //todo: Should check for error

		_ = jsonError;
		fmt.Println(inboundUserObject.UserName)

		returnedUserId := addUserCsv("444", inboundUserObject.UserName,
			inboundUserObject.Password, inboundUserObject.Emails[0].Value,inboundUserObject.Name.GivenName,
		inboundUserObject.Name.FamilyName, inboundUserObject.Name.Formatted)

		type schemasOpt []string;

		type nameOpt struct {
			Formatted  string `json:"formatted"`
			GivenName  string `json:"givenName"`
			FamilyName string `json:"familyName"`
		}

		var emailOpts = []struct {
			Value   string `json:"value"`
			Primary bool `json:"primary,omitempty"`
			Type    string `json:"type"`
		}{
			{"123", true, "234234"},

		}

		mySchema := schemasOpt{"urn:scim:schemas:core:1.0", "urn:scim:schemas:extension:enterprise:1.0"}

		s := outboundUserObj{}

		s.Schemas = mySchema
		s.ID = returnedUserId
		s.UserName = "bob"
		s.Password = "pass"
		s.Active = true
		s.Name = nameOpt{"123", "123", "123"}
		s.Emails = emailOpts

		result, _ := json.Marshal(s) //todo: should check for errors
		fmt.Fprint(w, string(result))

		return;
		break;

	case "GET":

		if (len(req.URL.Query().Get("filter")) != 0) {
			//This always returns a dummy response so Okta thinks it's okay to add Users
			//Long Story... It make Okta always push Users, not for production !
			fmt.Fprint(w, `{"totalResults": 0,"schemas": ["urn:scim:schemas:core:1.0"],"itemsPerPage": 0,"startIndex": 1,"Resources": []}`)
			return;
		} else {
			//Most likely this is an import, since there is not filter in the request
			returnString, result := importer()
			if (result !=nil) { //Something went wrong.. Just return User does not exist
				fmt.Fprint(w, `{"totalResults": 0,"schemas": ["urn:scim:schemas:core:1.0"],"itemsPerPage": 0,"startIndex": 1,"Resources": []}`)
				return;
			} else { // All good, return the JSON string from importer
				fmt.Fprint(w, returnString) //Returns Users in users.csv
				return;
			}
		}
		break;

	case "PUT": //Put equal Delete User

		fmt.Println("PUT request sent to User Controller PUT=Delete/Deactivate")
		body, err := ioutil.ReadAll(req.Body)
		_ = body
		if err != nil {
			panic(err)
		}

		//Get User ID, or key looks like this: /Users/123 (123=users id)
		vars := mux.Vars(req)
		userIdToDelete := vars["key"]
		if userIdToDelete == "" {
			//does not seem to register empty string
			userIdToDelete = ""
		} else {
			deleteUserCsv(userIdToDelete)
		}

		fmt.Println("Deleting/Deactivating User with Id: " + userIdToDelete)

		return;
		break;

	}

	w.Header().Set("Content-Type", "text/plain")
	var output = `{
   "totalResults":0,
   "schemas":[
      "urn:scim:schemas:core:1.0"
   ],
   "Resources":[]
}`
	w.Write([]byte(output))

}

func catchAll(w http.ResponseWriter, req *http.Request) {
	//Any other command will be echoed to the console.
	fmt.Printf("Req: %s %s\n", req.Host, req.URL.Path)

}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/Users/{key}", users)

	//fmt.Println(s.UserName)

	r.HandleFunc("/Groups", groups)
	r.HandleFunc("/Users", users)
	r.HandleFunc("/ServiceProviderConfigs", serverConfigs)
	r.HandleFunc("/", catchAll)
	go func() {
		err := http.ListenAndServeTLS(":443", "server.crt", "server.key", r)

		if err != nil {
			fmt.Println("Cannot configure https server")

			log.Fatal("ListenAndServe: ", err)
		} else {
			fmt.Println("Listening on https://localhost")

		}
	}()
	//err := http.ListenAndServeTLS("localhost:443", "server.crt", "server.key", nil)
	err2 := http.ListenAndServe(":80", r)
	//Make sure you have your certificates available, otherwise, it will not work


	if err2 != nil {
		log.Fatal("ListenAndServe: ", err2)
	} else {
		fmt.Println("Listening on http://localhost")
	}
}

func deleteUserCsv(idDelete string) {
	// Loading csv file
	rFile, err := os.Open("users.csv") //3 columns
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer rFile.Close()

	// Creating csv reader
	reader := csv.NewReader(rFile)

	lines, err := reader.ReadAll()
	if err == io.EOF {
		fmt.Println("Error:", err)
		return
	}

	// Creating csv writer
	wFile, err := os.Create("users.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer wFile.Close()
	writer := csv.NewWriter(wFile)

	// Read data, randomize columns and write new lines to results.csv

	col_index:= [7]int{0,1,2,3,4,5,6}; //Indicates order columns are in file

	for _, line := range lines {
		// If idDelete matches, then don't write that User to file
		if (line[col_index[0]] != idDelete) {
			writer.Write([]string{line[col_index[0]], line[col_index[1]], line[col_index[2]],line[col_index[3]],line[col_index[4]],line[col_index[5]],line[col_index[6]]}) //3 columns
			writer.Flush()
		}

	}

	//print report
	fmt.Println("\nNo. of lines: ", len(lines))

}

func addUserCsv(uid string, username string, password string, email string, firstname string, lastname string, fullname string) (userId string) {

	_ = uid;

	lastUserIdInUserCsv := 0;

	// Loading csv file
	rFile, err := os.Open("users.csv") //3 columns
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer rFile.Close()

	// Creating csv reader
	reader := csv.NewReader(rFile)

	lines, err := reader.ReadAll()
	if err == io.EOF {
		fmt.Println("Error:", err)
		return
	}

	// Creating csv writer
	wFile, err := os.Create("users.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer wFile.Close()
	writer := csv.NewWriter(wFile)

	// Read data, randomize columns and write new lines to results.csv

	col_index:= [7]int{0,1,2,3,4,5,6}; //Indicates order columns are in file

	for _, line := range lines {
		writer.Write([]string{line[col_index[0]], line[col_index[1]], line[col_index[2]],line[col_index[3]],line[col_index[4]],line[col_index[5]],line[col_index[6]]}) //3 columns
		writer.Flush()

		//This grabs the last UserID in the file, assuming they are in sequential order
		lastUserIdInUserCsv, _ = strconv.Atoi(line[col_index[0]])
		lastUserIdInUserCsv++
	}

	writer.Write([]string{strconv.Itoa(lastUserIdInUserCsv), username, password, email, firstname, lastname, fullname}) //7 columns
	writer.Flush()

	return strconv.Itoa(lastUserIdInUserCsv)

}

func importer () (jsonstring string, err error) {

	//This function calls the CSV, puts everything in place, and returns JSON string, it is
	//used when Okta is doing an Import.

	//var totalResults int
	fullUserImportStruct:=fullImportPayloadObj{} //Fullpay load
	singleUserStructs:=[]singleUserImportObj{}
	singleUserStruct1:=singleUserImportObj{}
	singleUserEmailObj := userEmailStruct{} //Email Struct, has it's own object
	singleUserNameObj := userNameStruct{} //Name Struct, first,last, given


	type schemasOpt []string; //Need Schema to explain to Okta what we are returning

	_=singleUserStructs

	rFile, err := os.Open("users.csv") //3 columns
	if err != nil {
		fmt.Println("Error:", err)
		return "error", err
	}
	defer rFile.Close()

	// Creating csv reader
	reader := csv.NewReader(rFile)

	lines, err := reader.ReadAll()
	if err == io.EOF {
		fmt.Println("Error:", err)
		return "error", err
	} else {

	}

	totalResults:=len(lines) //Total Users
	topSchema := schemasOpt{"urn:scim:schemas:core:1.0"} //Schema for top of response
	userSchema := schemasOpt{
		"urn:scim:schemas:core:1.0",
		"urn:scim:schemas:extension:enterprise:1.0",
		"urn:okta:onprem_app:1.0:user:custom"} //Schema for User of responses



	fullUserImportStruct.TotalResults=totalResults
	fullUserImportStruct.Schemas=topSchema

	//Loop through Users

	col_index:= [7]int{0,1,2,3,4,5,6}; //Indicates order columns are in file
	//Look at Users.csv to see what is returned

	for _,line :=range lines{ //Cycle through the Users

		fmt.Println(line[col_index[3]])

		//setup email object
		singleUserEmailObj.Value = line[col_index[3]]
		singleUserEmailObj.Primary = true
		singleUserEmailObj.Type = "work"

		if (len(singleUserStruct1.Emails)==0) { //If it's empty add the new struct
			fmt.Println("equals zero")

			singleUserStruct1.Emails = append ( singleUserStruct1.Emails, singleUserEmailObj )
		} else {
			singleUserStruct1.Emails=nil //If not Empty, clear it out
			//delete (singleUserStruct1.Emails[0],1)
			singleUserStruct1.Emails = append ( singleUserStruct1.Emails, singleUserEmailObj )
		}

		//setup name object
		singleUserNameObj.FamilyName=line[col_index[5]]
		singleUserNameObj.GivenName=line[col_index[4]]
		singleUserNameObj.Formatted=line[col_index[6]]

		singleUserStruct1.UserName=line[col_index[1]]
		singleUserStruct1.Schemas=userSchema
		singleUserStruct1.ID=line[col_index[0]]
		singleUserStruct1.Active=true
		singleUserStruct1.Password=line[col_index[2]]
		singleUserStruct1.Name=singleUserNameObj

		fullUserImportStruct.Resources = append ( fullUserImportStruct.Resources, singleUserStruct1)


	}

	//End Loop through Users

	//Convert Struct to JSON text for Okta
	result, _ := json.Marshal(fullUserImportStruct) //todo: should check for errors

	fmt.Println("Returning Import from users.csv to Okta in JSON format:")
	fmt.Println(string(result))

	return string(result), nil

}

