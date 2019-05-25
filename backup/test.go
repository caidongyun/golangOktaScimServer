
package main  //hi

import (
	"fmt"
	"encoding/json"
)

func main() {


jsonData := []byte(`{"Name":"Eve"}`)

var v interface{}
json.Unmarshal(jsonData, &v)
data := v.(map[string]interface{})

for k, v := range data {
switch v := v.(type) {
case string:
	fmt.Println ("string")
	fmt.Println(k,  "(string)")
	fmt.Println(v, "(string)")
case float64:
fmt.Println(k, v, "(float64)")
case []interface{}:
fmt.Println(k, "(array):")
for i, u := range v {
fmt.Println("    ", i, u)
}
default:
fmt.Println(k, v, "(unknown)")
}
}

}