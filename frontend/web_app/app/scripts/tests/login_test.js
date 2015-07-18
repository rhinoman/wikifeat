casper.options.viewportSize = {width: 1024, height: 768};

var testCount = 1;
var host = "http://localhost:4000";

casper.test.begin("Testing login", testCount, function loginTest(test){
    casper.start(host + "/login", function(){
        //Attempt to login
        casper.fillSelectors("form#loginForm", {
            "input#inputEmail": "jcadam",
            "input#inputPassword" : "password"
        }, true);
        casper.waitWhileSelector('#loginForm');
        casper.then(function(){
            test.assertUrlMatch(/app/, 'Redirected to app after login');
        });
    }).run(function(){
        test.done();
    });
});

