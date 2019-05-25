var request = require("request");

var args = process.argv.slice(2);
var LDAP_ID = "0oac9bgvakysNxoMa0h7"


if(!process.env.oktaOrg && !process.env.oktaKey  ) {
    console.log('environment variables not set, set them like this:');
    console.log('export oktaOrg="https://okta.okta.com"');
    console.log('export oktaKey="yourOktaKey"');
    return;
}


findUserByUsername = function (requestObj) {
    return new Promise((resolve, reject) => {

        var options = { method: 'GET',
            url:  process.env.oktaOrg+'/api/v1/users/'+args[0],
            headers:
                {
                    'cache-control': 'no-cache',
                    authorization: 'SSWS '+process.env.oktaKey,
                    'content-type': 'application/json',
                    accept: 'application/json' } };

        request(options, function (error, response, body) {
            if (error) throw new Error(error);

            // console.log(body);
            requestObj.userProfile = JSON.parse( body )
            resolve(requestObj)

        });

    })
}

removeLdap = function (requestObj) {
    return new Promise((resolve, reject) => {

        var options = { method: 'DELETE',
            url: process.env.oktaOrg+'/api/v1/apps/'+LDAP_ID+'/users/'+requestObj.userProfile.id,
            headers:
                {
                    'cache-control': 'no-cache',
                    'authorization': 'SSWS '+process.env.oktaKey,
                    'content-type': 'application/json',
                    accept: 'application/json' } };

        request(options, function (error, response, body) {
            if (error) throw new Error(error);

            resolve(requestObj)
        });

    })
}

setUserPassword = function (requestObj) {
    return new Promise((resolve, reject) => {

        var options = { method: 'PUT',
            url: process.env.oktaOrg+'/api/v1/users/'+requestObj.userProfile.id,
            headers:
                {
                    'cache-control': 'no-cache',
                    'authorization': 'SSWS '+process.env.oktaKey,
                    'content-type': 'application/json',
                    accept: 'application/json' },
            body: { credentials: { password: { value: args[1] } } },
            json: true };

        request(options, function (error, response, body) {
            if (error) throw new Error(error);

            resolve(requestObj)
        });

        resolve(requestObj)
    })
}

waitSeonds = function (requestObj) {
    return new Promise((resolve, reject) => {

        setTimeout(function () {
            resolve( requestObj)
        }, requestObj.seconds*1000);
    })
}

var requestObj = {}
requestObj.username = args[0]
requestObj.seconds = 20; // wait 20 seconds before switching over the User

waitSeonds ( requestObj).
    then (findUserByUsername).
    then ( removeLdap).
    then ( setUserPassword).
    then ( (responseObj)=>
    {
        console.log(responseObj.userProfile.id)

})


