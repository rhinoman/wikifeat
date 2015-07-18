global.template = require("../src/index");


var templ = template([
        "<h1><%= header %></h1>",
        "",
        "<ul>",
        "<% list.forEach(function(value, index) { %>",
        "  <li class='<%= index %>'><%= value %></li>",
        "<% }); %>",
        "</ul>"
    ].join("\n")),
    data = {
        header: "Header",
        list: [
            "List Item 0",
            "List Item 1",
            "List Item 2",
            "List Item 3",
            "List Item 4"
        ]
    };

document.getElementById("app").innerHTML = templ(data);


global.test = function test() {
    console.time("test");
    templ(data);
    console.timeEnd("test");
};
