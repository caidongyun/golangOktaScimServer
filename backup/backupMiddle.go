	if len(Plugin) !=0 {

			ldapObj :=LdapStruct{}

			out, err := exec.Command("node", "./"+Plugin,string(r.Name()),password).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("--------------------------------------")
			fmt.Println("230: --------------------------------------")
			fmt.Println(string(out));
			pushToStack(string(r.Name()), string(out))
			fmt.Println("I think name is"+r.Name())
			eraseme := popFromStack(string(r.Name()))
			fmt.Println(eraseme)
			json.Unmarshal(out,&ldapObj)


			//fmt.Println(ldapObj.Active)
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
			} else {
				res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
				w.Write(res)
				return
			}

		}





////////////////////

if len(PrePlugin) !=0 {


		if len(customAttributes) == 0 {
			fmt.Println("Map is Empty")
			//e.AddAttribute(message.AttributeDescription("departmentNumber"), message.AttributeValue(userid))
			//e.AddAttribute(message.AttributeDescription("telephoneNumber"), message.AttributeValue(userid))

			preout, err := exec.Command("node", "./"+PrePlugin, string(userid), "2").Output()
			if err != nil {
				fmt.Println("Error")
				log.Fatal(err)
			}

			fmt.Println(string(preout))
			pushToStack(string(userid), string(preout))


			preloginFields := convertJsonStringToMap(string(preout))


			fmt.Println(preloginFields)

			for customkey, customvalue := range preloginFields {
				fmt.Println("598:  Adding prelogin fields")
				fmt.Printf("%s -> %s\n", customkey, customvalue)
				e.AddAttribute(message.AttributeDescription(fmt.Sprintf("%s", string(customkey))),
					message.AttributeValue(fmt.Sprintf("%s", string(customvalue))))

			}
		}

		} 
		else {

			for customkey, customvalue := range customAttributes {
				fmt.Println("I'm in the Loop !!")
				fmt.Printf("%s -> %s\n", customkey, customvalue)
				e.AddAttribute(message.AttributeDescription(fmt.Sprintf("%s", string(customkey))),
					message.AttributeValue(fmt.Sprintf("%s", string(customvalue))))

			}
		}

		e.AddAttribute(message.AttributeDescription(fmt.Sprintf("%s", string ("pizza"))),
			message.AttributeValue(fmt.Sprintf("%s", string ("guidguidhere"))))

		e.AddAttribute("postalCode", "11111")
		e.AddAttribute("physicalDeliveryOfficeName", "x")
		//e.AddAttribute("departmentNumber", "11")
		//e.AddAttribute("telephoneNumber", "1111111")
		e.AddAttribute("mobile", "1111111")
		e.AddAttribute("preferredLanguage", "en")
		e.AddAttribute("postalAddress", "Austin")
		e.AddAttribute("employeeID", "0")
		e.AddAttribute("employeeNumber", "0")
		e.AddAttribute(message.AttributeDescription("entryuuid"),
			message.AttributeValue(string (userid)))

		e.AddAttribute("ds-pwp-account-disabled", "")
		w.Write(e)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		fmt.Println("Exiting 596 Done !! ---------------\n\n")






	getUuid:=fmt.Sprintf("%s",r.FilterString())

	if ( strings.Index(getUuid,"entryuuid=") !=0) {
		getUuid = strings.Replace(getUuid, "(", "", -1)
		getUuid = strings.Replace(getUuid, ")", "", -1)
		getUuid = strings.Replace(getUuid, "&", "", -1)

		entryuuidIndex := strings.Index(getUuid, "entryuuid=")
		getUuid = getUuid[entryuuidIndex+10:]

		userid=getUuid
	} else {
		fmt.Println("463: weird condition")
		getUuid="xxxxxxx"
	}

	// ou=groups search .. Treat Groups differently
	if strings.HasPrefix(string(r.BaseObject()),"ou=Groups") ||
		strings.HasPrefix(string(r.BaseObject()),"ou=groups") {

		log.Printf("Search DSE Groups #402")


		select {
		case <-m.Done:
			log.Print("Leaving handleSearch...")
			return
		default:
		}

		e := ldap.NewSearchResultEntry("ou=Groups,dc=example,dc=com")

		e.AddAttribute("objectClass", "top", "organizationalUnit")
		e.AddAttribute("ou", "Groups")
		e.AddAttribute(message.AttributeDescription( "entryuuid"), message.AttributeValue(getUuid))


		w.Write(e)

		e = ldap.NewSearchResultEntry("cn=ldapusers,ou=Groups,dc=example,dc=com")

		e.AddAttribute("objectClass", "top", "groupOfNames")
		e.AddAttribute("cn", "ldapusers")

		e.AddAttribute(message.AttributeDescription("member"), message.AttributeValue("uid"+userid+"=user.0,ou=People,dc=example,dc=com"))

		e.AddAttribute("entryuuid", "8c3624f5-d219-4401-9042-9a1fbf6f1b6805")
		w.Write(e)

		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
		w.Write(res)
		fmt.Println("Exiting 500!!!! Group Search!\n\n")

	} else {

		fmt.Println("506: second things after groups")

		// Handle Stop Signal (server stop / client disconnected / Abandoned request....)
		select {
		case <-m.Done:
			log.Print("Leaving handleSearch...")
			fmt.Println("512: Might have blown up")
			return
		default:
		}

		log.Printf("FileterString: %s",reflect.TypeOf(r.FilterString()))

		FilterString:= r.FilterString()
		if strings.Contains(FilterString,"(&(objectclass=inetorgperson)(mail=") {
			fmt.Println("521 trying to figure out username")
			FilterString=strings.Replace(FilterString,"(&(objectclass=inetorgperson)(mail=", "", -1)
			FilterString=strings.Replace(FilterString,")", "", -1)
			log.Printf("FilterString: %s",FilterString)
			log.Printf("Line #325")
			userid = FilterString
		}

		e := ldap.NewSearchResultEntry("uid="+userid+",ou=People,dc=example,dc=com")

		e.AddAttribute("objectClass", "top", "inetorgperson", "organizationalPerson", "person")

		if userid == "ss=person" {
			fmt.Println("trying to get userid again")
			fmt.Println(r.BaseObject());
			userid=fmt.Sprintf("%s",r.BaseObject())
			fmt.Println("569: New userid = "+userid)

		}

		if strings.Contains(userid, "@") {

		} else {
			userid=userid+"@noemailprovided.com"
		}

		fmt.Println("Here is the userid:"+userid)

		e.AddAttribute(message.AttributeDescription("distinguishedname"), message.AttributeValue(userid))
		e.AddAttribute(message.AttributeDescription("mail"), message.AttributeValue(userid))

//		e.AddAttribute(message.AttributeDescription("givenName"), message.AttributeValue(userid))
//		e.AddAttribute(message.AttributeDescription("sn"), message.AttributeValue(userid))

		e.AddAttribute("supportedLDAPVersion", "3")
		e.AddAttribute("title", "title")
		e.AddAttribute(message.AttributeDescription("uid"), message.AttributeValue(userid))
		e.AddAttribute("manager", "manager")
		e.AddAttribute("streetAddress", "street")
		e.AddAttribute("l", "USA")
		e.AddAttribute("st", "TX")

		// Add custom attributes

		var customAttributes=popFromStack(userid)

		fmt.Println("===============================")
		fmt.Println("I'm adding")
		fmt.Println("===============================")
		fmt.Println(customAttributes)
		fmt.Println("END adding custom")

		e.AddAttribute(message.AttributeDescription("givenName"), message.AttributeValue(userid))
		e.AddAttribute(message.AttributeDescription("sn"), message.AttributeValue(userid))

		fmt.Println(len(customAttributes))

		fmt.Println(PrePlugin);
		fmt.Println(userid);


	}

