var Module = require("./module"),

    scriptTag = ("currentScript" in document) ? document.currentScript : (function() {
        var scripts = document.getElementsByTagName("script");
        return scripts[scripts.length - 1];
    })(),

    SPLITER = /[\s ]+/,

    getAttribute = function(name) {

        return scriptTag ? scriptTag.getAttribute(name) || scriptTag.getAttribute("data-" + name) || scriptTag.getAttribute("x-" + name) : "";
    },
    hasAttribute = function hasAttribute(name) {

        return scriptTag ? !!(scriptTag.hasAttribute(name) || scriptTag.hasAttribute("data-" + name) || scriptTag.hasAttribute("x-" + name)) : false;
    },

    main = getAttribute("main"),
    args = getAttribute("argv"),
    env = getAttribute("env"),
    isGlobal = hasAttribute("global"),

    processEnv = process.env,
    i = -1,
    length, arg, key, value;


if (args) process.argv.push.apply(process.argv, args.split(SPLITER));

if (env && (env = env.split(SPLITER))) {
    length = env.length;

    while (++i < length) {
        arg = (env[i] || "").split("=");
        key = arg[0];
        value = arg[1];

        if (!key) continue;
        if (value != null && !processEnv[key]) processEnv[key] = value;
    }
}

global.XMLHttpRequest || (global.XMLHttpRequest = function XMLHttpRequest() {
    try {
        return new ActiveXObject("Msxml2.XMLHTTP.6.0");
    } catch (e1) {
        try {
            return new ActiveXObject("Msxml2.XMLHTTP.3.0");
        } catch (e2) {
            throw new Error("XMLHttpRequest is not supported");
        }
    }
});

if (isGlobal || (!isGlobal && !main)) {
    global.process = process;
    global.Buffer = Buffer;
}

if (main) {
    Module.init(main, isGlobal);
} else {
    Module.repl();
}
