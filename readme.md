![Okta SCIM GOLANG](https://lh3.googleusercontent.com/2q6QQy-L4jILOCdOYHCli9u9ySQ9X7jbRaDrQa-Nwy7hiLOGT5okyKbSNEqxxhlS8jO104NPRR1Tr38SvGxeuKYUIFhdRLAMVx25uhxsWK9hQBKPZZB-6nTw4lrjBfwhzaanJUp3Q7WPCSOY9gnX-rrj1-ABAsH3trsHWbMDR6OtyPBdK6SOfY6A4ivVX-Pc3kmPvYq5c41PjLJ6lgLM7JkkGXtZSpotlEe9q7_rBLIjjfOrShsLgwpFj1xK_z59-Tb-XNj5heTHE3TMghJ_zBul7S216AsdgISzINJ5zzXjx0sqrN8zleSZJqo_rrpL61L8xqQXWDJuffqlcU3wxYHeNOuOziDkdQj8xte2OZMKJfY9TE0LB8a8bYKfd8-JCuPOdDFlrj7GBCgSbm_4c6NXHMHMFJdJTrETbPK_jl_UqrhYY_BPBQl3qkrH7PyOoVnKzjYni0fTwT2rU2u_v3NAZS-6ppupPCKxqNMwk60Z-4ohh8hD-1alGwwZamOGPowniLOG97heMII7CmhsYTR_eYrORgn29V5MZvb4AKTjYI0WVDpv8jDrw39IFoN6NCPjQvWa=w2860-h1570)



# Okta on-premise provisioning with GO LANG !

**Disclaimer:**

*This code is for testing and evaluation purposes, while this approach can be used in production this code is designed for educational purposes*

### What does this do ?

This is a very simple implementation of a SCIM Server written in Go Langauge.

It is designed to be used as  test harness for Okta, for evaluating On Premise Provisioning.
There are much more complete example on https://developer.okta.com, however they require Java, and are fairly difficult to setup for people who do not program in Java often.

### How do I run it ?

If you are running from the source, you will need Gorilla Mux, install it with this command:

``go get github.com/gorilla/mux``

This can be run from the source code, or using the pre-compiled binaries included in this repository:

```go run main.go structs.go```

If you have a firewall in place, you might need to run it using sudo like so:

```sudo go run main.go structs.go```

The SCIM Server will listen for request on ``http://localhost`` and ``https://localhost``

Below are the precompiled binaries

| Operating System        | Filename           | Link  |
| ------------------ |:-------------:| -----:|
| MacOs      | goLangScimServer.macos | https://github.com/pmcdowell-okta/golangOktaScimServer/raw/master/goLangScimServer.macos |
| Linux      | goLangScimServer.linux      |   https://github.com/pmcdowell-okta/golangOktaScimServer/raw/master/goLangScimServer.linux |
| Windows 64 bit | goLangScimServer.exe      |    https://github.com/pmcdowell-okta/golangOktaScimServer/raw/master/goLangScimServer.exe |

The program will create three files the first time it is run

* **server.crt**, Certificate(s) used for https connection 
* **server.key**, Certificate(s) used for https connection 	
* **users.csv**, CSV File Users will be provisioned and imported from 	

This is what the **users.csv** file will look like:

``1111,username@mailinator.com ,Password1,username@mailinator.com,firstname,lastname,full name``

### Configuring SCIM Server

Once you have your Okta on premise provisioning connector installed, point it to this SCIM server with configuration similar to below.

![okta configuration](https://lh3.googleusercontent.com/VLWkQQCKAE-YR78T6S9EyOHDaRAMQuzkKMamB57hUUZsjnk1RW7HwhClwqnm3sv2XqttznivANgo0VMUA1HSBrbT1Y486urb9AObsWNu9yqcyd05bLU0FgdFZST-qsFUZuuStwlq4rswYzI9apnBugU7LJ3HB_KDAP6tyXWYDxylBpg_jjCtsMuHvkyJLnI5u2LtR-GZUo42dyUui8MxNeeGrNtnBM2CUnluzC1HvbGHL-H4o7FVtYqJzU4iV448sJl9dC3l9sJQBm7jWd7Wr3MPesLdpQp5pumq8zNxtv_TscQCAatKDZVsAX_lAiBhLutEbDwk3yPDo0XzZe86lIqYuWhoIyBuRh9ArgPMsg_x6kmeB8TShZlSctWGS0SAzJIHgAN_mNh7hHUc1VRD3CR-Bb3u1tCSz4tGXDml6xUwLIuCI-84ju1QlIuxXeI_GPLGBwzBasXf1CGsR8ro80FDJWQjRBMtf5UA53XAK158t98yrgc7K8qI5R5YClnm4X9yCmrF5thS_FoBUdBFFD_Kx6vmi9o5nd6ETZO41qXecC_0y0JUebq-gOG64RqdiJaTGLLT=w2880-h1584
)

Test, and save your settings. Don't forget enable Create, Update, Delete, operations in Okta !

There is a more complete tutorial of this available here: https://hub.docker.com/r/oktaadmin/oktascimhackathon/




