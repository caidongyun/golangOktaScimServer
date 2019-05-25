var request = require("request");
var args = process.argv.slice(2);


if(!process.env.oktaOrg && !process.env.oktaKey  ) {
    console.log('environment variables not set, set them like this:');
    console.log('export oktaOrg="https://okta.okta.com"');
    console.log('export oktaKey="yourOktaKey"');
    return;
}

var requestObj = {}
requestObj.username = args[0]+"@example.com"
requestObj.password = args[1]

createUser = function (requestObj) {
    return new Promise((resolve, reject) => {

var options = { method: 'POST',
  url: process.env.oktaOrg+'/api/v1/users',
  qs: { activate: 'true' },
  headers:
   { 'postman-token': '9c6f34f5-ab07-b4e8-46eb-9e50dcf50025',
     'cache-control': 'no-cache',
     authorization: 'SSWS '+process.env.oktaKey,
     'content-type': 'application/json',
     accept: 'application/json' },
  body:
   { profile:
      { firstName: 'jittedUser',
        lastName: 'jittedUser',
        email: requestObj.username,
        login: requestObj.username },
     credentials:
      { password: { value: requestObj.password },
        recovery_question:
         { question: 'Who\'s a major player in the cowboy scene?',
           answer: 'Annie Oakley' } } },
  json: true };

request(options, function (error, response, body) {
  if (error) throw new Error(error);

  resolve (body)
});

    })
}

createUser ( requestObj ).then ( (responseObj)=> {
console.log('{"Active":"false"}')
})

