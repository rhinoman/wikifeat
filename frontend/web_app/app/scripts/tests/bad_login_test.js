casper.options.viewportSize = {width: 1024, height: 768};

var testCount = 1;
var host = "http://localhost:4000";

casper.test.begin("Testing login", testCount, function loginTest(test) {
    casper.start(host + "/login", function () {
        //Attempt bad login
        this.echo("Bad Login");
        casper.fillSelectors("form#loginForm", {
            "input#inputEmail": "notauser@nowhere.com",
            "input#inputPassword": "notapassword"
        }, true);
        casper.waitUntilVisible('#alertBox', function () {
            test.assertVisible('#alertBox', 'Error displayed');
        });
    }).run(function () {
        test.done();
    });
});
