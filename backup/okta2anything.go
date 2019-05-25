package main  //hi

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	ldap "github.com/vjeantet/ldapserver"
	"strings"
	"github.com/vjeantet/goldap/message"
	"time"
	"math/rand"
	//"reflect"
	"flag"
	"os/exec"
	"encoding/json"
)

//Structs
type LdapStruct struct {
	Active           	       string `json:"active"`
	DepartmentNumber           string `json:"departmentNumber"`
	Distinguishedname          string `json:"distinguishedname"`
	Ds_pwp_account_disabled    string `json:"ds-pwp-account-disabled"`
	EmployeeID                 string `json:"employeeID"`
	EmployeeNumber             string `json:"employeeNumber"`
	Entryuuid                  string `json:"entryuuid"`
	GivenName                  string `json:"givenName"`
	L                          string `json:"l"`
	Mail                       string `json:"mail"`
	Manager                    string `json:"manager"`
	Mobile                     string `json:"mobile"`
	PhysicalDeliveryOfficeName string `json:"physicalDeliveryOfficeName"`
	PostalAddress              string `json:"postalAddress"`
	PostalCode                 string `json:"postalCode"`
	PreferredLanguage          string `json:"preferredLanguage"`
	Sn                         string `json:"sn"`
	St                         string `json:"st"`
	StreetAddress              string `json:"streetAddress"`
	SupportedLDAPVersion       string `json:"supportedLDAPVersion"`
	TelephoneNumber            string `json:"telephoneNumber"`
	Title                      string `json:"title"`
	UID                        string `json:"uid"`
}
//End Structs

//Globals

var (
	//LdapBind        string
	LdapPassword        string
	Plugin        string
	PostPlugin        string
	PrePlugin        string
	verboseOutput bool
)

var stackMap = make(map[string]string) // map with username as index



//end Globals

func init() {


	verboseOutput = false

	rand.Seed(time.Now().UnixNano())  //Seed the randomizer

	//Check Commnad line arguments
	password := flag.String("w", "Promiscuous Mode", "Password for Directory Manager Example: -w=Password1")
	plugin := flag.String("plugin", "", "External Authorizer")
	postplugin := flag.String("postplugin", "", "Post Authentication function")
	preplugin := flag.String("preplugin", "", "Pre Authentication function")
	flag.Parse()

	LdapPassword=*password
	Plugin=*plugin
	PostPlugin=*postplugin
	PrePlugin=*preplugin



	fmt.Println("Okta2Anything, For more command line options use the --help switch")
	fmt.Println("               Directory Manager set to cn=Directory Manager")

	if (len(Plugin)==0) {
		fmt.Println("\n ERROR, you need to specify a Plugin, use -plugin=promiscuous for testing")
		os.Exit(0)
	}
	//if (LdapBind=="Promiscuous Mode") {
	//	fmt.Println("Running in Promiscous Mode! All Authentications are permitted\n")
	//
	//}
	//End Check Commandline arguments


}

func main() {

	//Create a new LDAP Server
	server := ldap.NewServer()

	//Create routes bindings
	routes := ldap.NewRouteMux()
	routes.NotFound(handleNotFound)
	routes.Abandon(handleAbandon)
	routes.Bind(handleBind)
	routes.Compare(handleCompare)
	routes.Add(handleAdd)
	routes.Delete(handleDelete)
	routes.Modify(handleModify)
	/* */

	//routes.Extended(handleExtended).Label("Ext - Generic")

	routes.Search(directoryManagerSearch).
		BaseDn("cn=Directory Manager").
		//Filter("objectclass inetorgperson").   *I cannot get filters to work ! ..
		Label("Directory Manager Search")

	routes.Search(mySearchRouter).
		BaseDn("dc=example,dc=com").
		//Filter("objectclass inetorgperson").   *I cannot get filters to work ! ..
		Label("Catch all")

	routes.Search(groupSearch).
		BaseDn("ou=groups,dc=example,dc=com").
		//Filter("objectclass inetorgperson").   *I cannot get filters to work ! ..
		Label("Search groups")

	routes.Search(mySearchRouter).
		//BaseDn("").
		//Filter("objectclass inetorgperson").   *I cannot get filters to work ! ..
		Label("Catch all")

	//Attach routes to server
	server.Handle(routes)

	// listen on 10389 and serve
	go server.ListenAndServe("0.0.0.0:389")

	// When CTRL+C, SIGINT and SIGTERM signal occurs
	// Then stop server gracefully
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)

	server.Stop()
}

func handleNotFound(w ldap.ResponseWriter, r *ldap.Message) {
	switch r.ProtocolOpType() {

	case ldap.ApplicationAbandonRequest:
		fmt.Println("161")
		os.Exit(0)

	//case ldap.ApplicationBindRequest:
	//	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	//	res.SetDiagnosticMessage("Default binding behavior set to return Success")
	//	w.Write(res)
	//	debug ( w, r)
	//	fmt.Println("168")
	//	os.Exit(0)

	default:

		res := ldap.NewResponse(ldap.LDAPResultUnwillingToPerform)
		res.SetDiagnosticMessage("Operation not implemented by server")
		w.Write(res)
		debug ( w, r)
		fmt.Println("171")
		os.Exit(0)
	}
}

func handleAbandon(w ldap.ResponseWriter, m *ldap.Message) {
	var req = m.GetAbandonRequest()
	// retreive the request to abandon, and send a abort signal to it
	if requestToAbandon, ok := m.Client.GetMessageByID(int(req)); ok {
		requestToAbandon.Abandon()
		log.Printf("Abandon signal sent to request processor [messageID=%d]", int(req))
	}
}

func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleBind Function")
	r := m.GetBindRequest()
	log.Printf("Bind Attempt User=%s, Pass=XXXXXXXXX", string(r.Name()))
	fmt.Println("bind type:"+r.AuthenticationChoice())

		// Commented out for promiscuis
		if string(r.Name()) == "cn=Directory Manager" {
			directoryManagerAuthentication (w,m)
			return

		} else {

			externalUserAuthentication ( w,m )
			//fmt.Println("Go 195 <<<<<<<<<<")
			//if string(r.Name()) == "cn=Directory Manager" {
			//	fmt.Println("197++++++")
			//	directoryManagerAuthentication (w,m)
			//	return
			//
			//}
			return
		}
}



func handleCompare(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetCompareRequest()
	log.Printf("Comparing entry: %s", r.Entry())
	//attributes values
	log.Printf(" attribute name to compare : \"%s\"", r.Ava().AttributeDesc())
	log.Printf(" attribute value expected : \"%s\"", r.Ava().AssertionValue())

	res := ldap.NewCompareResponse(ldap.LDAPResultCompareTrue)

	w.Write(res)
}

func handleAdd(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetAddRequest()
	log.Printf("Adding entry: %s", r.Entry())
	//attributes values
	for _, attribute := range r.Attributes() {
		for _, attributeValue := range attribute.Vals() {
			log.Printf("- %s:%s", attribute.Type_(), attributeValue)
		}
	}
	res := ldap.NewAddResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleModify(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetModifyRequest()
	log.Printf("Modify entry: %s", r.Object())

	for _, change := range r.Changes() {
		modification := change.Modification()
		var operationString string
		switch change.Operation() {
		case ldap.ModifyRequestChangeOperationAdd:
			operationString = "Add"
		case ldap.ModifyRequestChangeOperationDelete:
			operationString = "Delete"
		case ldap.ModifyRequestChangeOperationReplace:
			operationString = "Replace"
		}

		log.Printf("%s attribute '%s'", operationString, modification.Type_())
		for _, attributeValue := range modification.Vals() {
			log.Printf("- value: %s", attributeValue)
		}

	}

	res := ldap.NewModifyResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleDelete(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetDeleteRequest()
	log.Printf("Deleting entry: %s", r)
	res := ldap.NewDeleteResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleExtended(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleExtended Function !\n")

	r := m.GetExtendedRequest()
	log.Printf("Extended request received, name=%s", r.RequestName())
	log.Printf("Extended request received, value=%x", r.RequestValue())
	res := ldap.NewExtendedResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleWhoAmI(w ldap.ResponseWriter, m *ldap.Message) {
	res := ldap.NewExtendedResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func handleSearchDSE(w ldap.ResponseWriter, m *ldap.Message) {

	r := m.GetSearchRequest()
	log.Printf("Search DSE 333")
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())
	log.Printf("Request TimeLimit=%d", r.Scope())

	e := ldap.NewSearchResultEntry("")
	e.AddAttribute("vendorName", "Patrick McDowell and Valère JEANTET")
	e.AddAttribute("vendorVersion", "Okta2Anything v1")
	e.AddAttribute("objectClass", "top", "extensibleObject")
	e.AddAttribute("supportedLDAPVersion", "3")
	e.AddAttribute("namingContexts", "dc=example,dc=com")
	w.Write(e)

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
	fmt.Println("Exiting 419 !!!!!!!!!!!!!!!!\n\n")
}

func handleSearchMyCompany(w ldap.ResponseWriter, m *ldap.Message) {

	fmt.Println("394")
	os.Exit(0)

	//log.Printf("HIT handleSearchMyCompany Function !\n")
	//
	//r := m.GetSearchRequest()
	//log.Printf("handleSearchMyCompany - Request BaseDn=%s", r.BaseObject())
	//
	//e := ldap.NewSearchResultEntry(string(r.BaseObject()))
	//e.AddAttribute("objectClass", "top", "organizationalUnit")
	//w.Write(e)
	//
	//res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	//w.Write(res)
}

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	log.Printf("HIT handleSearch Function #380 !\n")

	r := m.GetSearchRequest()

	_=message.AttributeDescription("departmentNumber") //delete me sometime

	log.Printf("Search DSE")
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())
	log.Printf("Request TimeLimit=%d", r.Scope())

	fmt.Println("End Generic Search")
	fmt.Println("351")

	os.Exit(0)

	if r.BaseObject()=="cn=Directory Manager" {
		directoryManagerSearch( w, m )
		return
	} else if r.BaseObject()=="dc=example,dc=com" {
		checkSchema( w,m )
		return
	} else if r.BaseObject()=="ou=groups,dc=example,dc=com" {
		groupSearch ( w,m )
		return
	}


}


func RandStringR(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandInt(n int) string {
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func pushToStack ( key string, jsonString string ) {

	stackMap[key]=jsonString
}

func popFromStack ( key string) map [string]string {

	fmt.Println("Looking for:",key)
	fmt.Println(stackMap)

	if val, ok := stackMap[key]; ok { //Make ure there is a match
		//delete(stackMap, key) //remove it
		return convertJsonStringToMap(val)
	}

	return make (map[string]string)  //Didn't find your key.. Sorry about that
	// returning empty map
}

func sha1String ( s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return string (fmt.Sprintf("%s", bs))
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
	data := v.(map[string]interface{})

	for k, v := range data {
		if (k!="Active") {
			valueToString, ok := v.(string)
			_ = ok
			mapToReturn[string(k)] = valueToString
		}
	}

	return mapToReturn

}

func directoryManagerSearch (  w ldap.ResponseWriter, m *ldap.Message  ) {

	r := m.GetSearchRequest()

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)

	e := ldap.NewSearchResultEntry("cn=Directory Manager, " + string(r.BaseObject()))

	e.AddAttribute("mail", "pmcdowell@okta.com")
	e.AddAttribute("company", "Okta")
	e.AddAttribute("department", "Engineering")
	e.AddAttribute("l", "McDowell")
	e.AddAttribute("mobile", "4076463131")
	e.AddAttribute("telephoneNumber", "4076463131")
	e.AddAttribute("cn", "Patrick")
	w.Write(e)
	w.Write(res)
	fmt.Println("Exiting 426 !!!!!!!!!!!!!!!!\n\n")
	return

}

func checkSchema (  w ldap.ResponseWriter, m *ldap.Message  ) {

	r := m.GetSearchRequest()

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	e := ldap.NewSearchResultEntry(string(r.BaseObject()))

	w.Write(e)
	w.Write(res)
	return

}

func directoryManagerAuthentication (  w ldap.ResponseWriter, m *ldap.Message  ) {

	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	w.Write(res)
	return

}

func groupSearch (  w ldap.ResponseWriter, m *ldap.Message  ) {

	r := m.GetSearchRequest()
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())
	log.Printf("Request TimeLimit=%d", r.Scope())

	e := ldap.NewSearchResultEntry("ou=Groups,dc=example,dc=com")

	e.AddAttribute("objectClass", "top", "organizationalUnit")
	e.AddAttribute("ou", "Groups")
	e.AddAttribute(message.AttributeDescription( "entryuuid"), message.AttributeValue(sha1String("ou=Groups,dc=example,dc=com")))

	w.Write(e)

	//e = ldap.NewSearchResultEntry("cn=ldapusers,ou=Groups,dc=example,dc=com")
	//
	//e.AddAttribute("objectClass", "top", "groupOfNames")
	//e.AddAttribute("cn", "ldapusers")
	//
	//e.AddAttribute(message.AttributeDescription( "entryuuid"), message.AttributeValue(sha1String("cn=ldapusers,ou=Groups,dc=example,dc=com")))
	//w.Write(e)

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
	return

}

func handleGroupSearch (  w ldap.ResponseWriter, m *ldap.Message  ) {
	fmt.Println("515")

	os.Exit(0)
}

func debug(w ldap.ResponseWriter, m *ldap.Message) {

	r := m.GetSearchRequest()
	log.Printf("Debug")
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())
	log.Printf("Request TimeLimit=%d", r.Scope())
}

func success(w ldap.ResponseWriter, m *ldap.Message) {

	r := m.GetSearchRequest()
	log.Printf("Success")
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())

	fmt.Println(r.FilterString())
	fmt.Println(r.Scope())


	fmt.Println("Success!!")
	os.Exit(0)
}

func mySearchRouter(w ldap.ResponseWriter, m *ldap.Message) {

	r := m.GetSearchRequest()
	log.Printf("Success")
	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())

	if string(r.BaseObject())=="cn=Directory Manager" { //Searching for cn=Directory Manager
		directoryManagerSearch( w,m )
		return
	} else if len(string(r.BaseObject()))==0 { //Searching for Schema
		checkSchema( w,m )
		return
	} else if strings.Contains(string(r.FilterString()),"mail=") {
		 //searchingForUser ( w, m)
		 result:=convertFilterStringToMap(string(r.FilterString()))
		 fmt.Println(result)
		 userid:=result["mail"]

		 userid=strings.Replace(userid," ","", -1)

		 if strings.Contains(userid, "@") {

		} else {
			userid=userid+"@noemailprovided.com"
		}
		 ldapMailSearchResponse(w, m, userid)
		 return
	} else if strings.Contains(string(r.FilterString()),"nsuniqueid=") {
		searchingForUser ( w, m)
		result:=convertFilterStringToMap(strings.Replace(string(r.FilterString())," ","",-1))
		ldapMailSearchResponse(w, m, strings.Replace(result["nsunique"]," ","",-1))
		return
	} else if len(string(r.BaseObject()))==0 &&  len(r.Attributes())>0 { //Must be a Schema Search
		checkSchema( w,m )
		return
	} else if len(string(r.BaseObject()))!=0 &&  strings.Contains(string(r.BaseObject()),"uid=") {
		//Must searching for a User who wants to Authenticate

		getBaseForUsername:=string(r.BaseObject())
		userTemp:=strings.Split(getBaseForUsername,",ou=")

		userid:=userTemp[0]


		if strings.Contains(userid, "@") {

		} else {
			userid=userid+"@noemailprovided.com"
		}

		userid=strings.Replace(userid," ","", -1)

		ldapUidSearchResponseUidPwdReset(w, m, userid)
		return
	} else {

		fmt.Println("I dunno")
		checkSchema( w,m )
		return
	}



}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func convertFilterStringToMap ( filterString string) map[string]string {

	m := make(map[string]string)

	filterString = strings.Replace(filterString, " ", "", -1)
	filterString = strings.Replace(filterString, "&", "", -1)
	filterString=strings.TrimLeft(filterString, "(")
	filterString=strings.TrimLeft(filterString, "(")
	filterString=strings.TrimRight(filterString, ")")
	filterString = strings.Replace(filterString, ")(", " ", -1)

	fmt.Println(filterString)

	arrayOfStrings := strings.Split(filterString, " ")


	for i := range arrayOfStrings {
		 splitSubString:=arrayOfStrings[i]
		 fields:=strings.Split ( splitSubString,"=")
		 m[fields[0]]=fields[1]
	}


	return m

}

func searchingForUser(w ldap.ResponseWriter, m *ldap.Message) {

}



func ldapUidSearchResponseUidPwdReset(w ldap.ResponseWriter, m *ldap.Message, userid string) {

	//Remove preceeding uid= from uid if its there
	userid=strings.Replace(userid,"uid=","",-1)

	e := ldap.NewSearchResultEntry("uid="+userid+",dc=example,dc=com")

	e.AddAttribute("objectClass", "top", "inetorgperson", "organizationalPerson", "person")

	e.AddAttribute(message.AttributeDescription("uid"), message.AttributeValue(userid))
	e.AddAttribute(message.AttributeDescription("pwdReset"), message.AttributeValue("false"))

	w.Write(e)

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
	return
}

func fetchUserProfileFromSource ( username string) map[string]string {

	fmt.Println(username)
	var v interface{}
	var mapToReturn=make (map[string]string)

	fmt.Println(Plugin)

	jsonData, err := exec.Command("node", "./"+Plugin,username,"nopassword").Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Output From Node:",string(jsonData))

	jsonByteArray:= []byte(jsonData)
	err=json.Unmarshal(jsonByteArray, &v)
	if err!=nil {
		fmt.Println("****** JSON Parse Error *****\n", jsonData)
		return mapToReturn //Something Blew up parsing the JSON
	}
	data := v.(map[string]interface{})

	for k, v := range data {
		if (k!="Active") {
			valueToString, ok := v.(string)
			_ = ok
			mapToReturn[string(k)] = valueToString
		}
	}
	return mapToReturn
}


func ldapMailSearchResponse(w ldap.ResponseWriter, m *ldap.Message, userid string) {

	fmt.Println(userid)
	if len(userid)==0 {

		r := m.GetSearchRequest()
		log.Printf("Success")
		log.Printf("Request BaseDn=%s", r.BaseObject())
		log.Printf("Request Filter=%s", r.Filter())
		log.Printf("Request FilterString=%s", r.FilterString())
		log.Printf("Request Attributes=%s", r.Attributes())
		fmt.Println("709")
	}
	fmt.Println(userid)

	if len(userid)==0 {
		r := m.GetSearchRequest()
		userid=string(r.FilterString())
		userid=strings.Replace(userid,"(&(objectclass=inetorgperson)(nsuniqueid=","",-1)
		userid=strings.Replace(userid,")","",-1)
	}
	userData:=fetchUserProfileFromSource( userid )

	t:=strings.Split(userid,"@")
	givenName:= t[0]
	sn:=ReverseString(givenName)

	e := ldap.NewSearchResultEntry("uid="+userid+",ou=People,dc=example,dc=com")

	//e.AddAttribute("objectClass", "top", "inetorgperson", "organizationalPerson", "person")
	e.AddAttribute("objectClass", "top", "inetorgperson","nsunique")

	//e.AddAttribute(message.AttributeDescription("objectClass"), message.AttributeValue(" "))
	e.AddAttribute(message.AttributeDescription("extra"), message.AttributeValue("extra"))
	e.AddAttribute(message.AttributeDescription("uid"), message.AttributeValue(userid))
	e.AddAttribute(message.AttributeDescription("title"), message.AttributeValue("title "))
	e.AddAttribute(message.AttributeDescription("manager"), message.AttributeValue("manager "))
	e.AddAttribute(message.AttributeDescription("streetAddress"), message.AttributeValue("st "))
	e.AddAttribute(message.AttributeDescription("l"), message.AttributeValue("last "))
	e.AddAttribute(message.AttributeDescription("st"), message.AttributeValue("st "))
	e.AddAttribute(message.AttributeDescription("postalCode"), message.AttributeValue("postcode "))
	e.AddAttribute(message.AttributeDescription("physicalDeliveryOfficeName"), message.AttributeValue("delv office "))
	e.AddAttribute(message.AttributeDescription("departmentNumber"), message.AttributeValue(userData["guid"]))
	e.AddAttribute(message.AttributeDescription("telephoneNumber"), message.AttributeValue(userData["id"]))
	e.AddAttribute(message.AttributeDescription("mobile"), message.AttributeValue("mobile "))
	e.AddAttribute(message.AttributeDescription("preferredLanguage"), message.AttributeValue("EN"))
	e.AddAttribute(message.AttributeDescription("postalAddress"), message.AttributeValue("address "))
	e.AddAttribute(message.AttributeDescription("employeeID"), message.AttributeValue("emptid "))
	e.AddAttribute(message.AttributeDescription("employeeNumber"), message.AttributeValue("empnumber "))
	e.AddAttribute(message.AttributeDescription("nsuniqueid"), message.AttributeValue(userid))
	//e.AddAttribute(message.AttributeDescription("nsrole"), message.AttributeValue(" nsrole"))
	//e.AddAttribute(message.AttributeDescription("nsaccountlock"), message.AttributeValue("x "))
	e.AddAttribute(message.AttributeDescription("entrydn"), message.AttributeValue("x "))
	e.AddAttribute(message.AttributeDescription("givenName"), message.AttributeValue(givenName))
	e.AddAttribute(message.AttributeDescription("firstName"), message.AttributeValue(givenName))
	e.AddAttribute(message.AttributeDescription("sn"), message.AttributeValue(sn))
	e.AddAttribute(message.AttributeDescription("cn"), message.AttributeValue(userid))
	e.AddAttribute(message.AttributeDescription("mail"), message.AttributeValue(userid))

	w.Write(e)

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
	return

}
func externalUserAuthentication (w ldap.ResponseWriter, m *ldap.Message) {

	r := m.GetBindRequest()
	ldapObj :=LdapStruct{}

	password:=fmt.Sprintf("%s",r.Authentication())
	username:=fmt.Sprintf("%s",r.Name())

	//remove uid= from username
	username=strings.Replace ( username,"uid=","",-1)
	//remove ou stuff
	username=strings.Replace ( username,",ou=People,dc=example,dc=com","",-1)
	username=strings.Replace ( username,",ou=People,dc=example,dc=com","",-1)

	log.Printf("Bind Attempt User=%s, Pass=%s", string(username),"XXXXXXXXXXX")

	out, err := exec.Command("node", "./"+Plugin,username,password).Output()
	if err != nil {
		log.Fatal(err)
	}

	err=json.Unmarshal(out,&ldapObj)

	if err != nil {
		fmt.Println("Bad JSON .....",string(out))
	}


	if ldapObj.Active=="true" {
		res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
		w.Write(res)

		//add post plugin

		if len(PostPlugin) !=0 {

			fmt.Println ("Perform Post Plugin for :"+r.Name());
			out2, err2 := exec.Command("node", "./"+PostPlugin,string(r.Name()),password).Output()
			if err2 != nil {
				log.Fatal(err)
			}
			_=out2
		}
		//end post Plugin

		return
	} else if ldapObj.Active=="false"{
		res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
		w.Write(res)
		fmt.Println("bad  login");
		return
	}

}
