package main

import (
	"encoding/json"
	"fmt"
)

var stackMap = make(map[string]string) // map with username as index




func main() {


	pushToStack("test",`{  
   "Active":"true",
   "email":"test@dish.com",
   "id":13870755,
   "guid":"60fcd578-bf1f-11e7-9cfe-02e7f69e2b00"
}`)

	pushToStack("test2",`{  
   "Active":"true",
   "email":"test2@dish.com",
   "id":13870755,
   "guid":"60fcd578-bf1f-11e7-9cfe-02e7f69e2b00"
}`)


	fmt.Println(popFromStack("test2")["email"])
	fmt.Println(popFromStack("test2")["email"])


}

func pushToStack ( key string, jsonString string ) {
	stackMap[key]=jsonString
}

func popFromStack ( key string) map [string]string {

	if val, ok := stackMap[key]; ok { //Make ure there is a match
		delete(stackMap, key) //remove it
		return convertJsonStringToMap(val)
	}

	return make (map[string]string)  //Didn't find your key.. Sorry about that
	// returning empty map
}

func convertJsonStringToMap ( jsonData string ) map [string]string {

	var mapToReturn=make (map[string]string)

	jsonByteArray:= []byte(jsonData)
	var v interface{}
	err:=json.Unmarshal(jsonByteArray, &v)
	if err!=nil {
		fmt.Println("****** JSON Parse Error *****\n", jsonData)
		return mapToReturn //Something Blew up parsing the JSON
	}
	//fmt.Println(err)
	data := v.(map[string]interface{})

	for k, v := range data {

		valueToString, ok := v.(string)
		_=ok
		mapToReturn[string(k)]=valueToString
	}

	return mapToReturn



}

