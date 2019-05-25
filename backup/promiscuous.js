//Script always returns true, so anyone can login
var args = process.argv.slice(2);

var username = args[0]
var password = args[1]


if ( args[0].startsWith("t") ) {
    console.log ('{"Active":"true","email":"icanremember@dish.com","id":13870755,"guid":"60fcd578-bf1f-11e7-9cfe-02e7f69e2b00"}')
}
else
{
    console.log ( '{"guid":"60fcd578-bf1f-11e7-9cfe-02e7f69e2b00","id":"13870755","Active":"false","email":""}')
}
// console.log('{"Active":"true","guidguid":"111111"}')

function makeid(length) {
    var result           = '';
    var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

///
