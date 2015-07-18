#!/usr/bin/env node


var fs = require("fs"),
    filePath = require("file_path"),
    Requireify = require("./requireify"),

    argv = require("optimist")
    .demand("f").alias("f", "file").describe("f", "index file to start parsing from")
    .alias("o", "out").describe("o", "out file, defaults to index file + min.js")
    .alias("e", "exportName").describe("e", "export to global object with this name")
    .alias("i", "ignore").describe("e", "ignore modules and relative file paths")
    .alias("a", "argv").describe("a", "list of arguments to pass to process.argv (--argv=--arg0,-a)")
    .describe("env", "list of arguments to pass to process.env (--env=NODE_ENV=development,DEBUG=*)")
    .alias("v", "verbose").describe("v", "verbose mode")
    .argv,

    hasExtension = /\.[\w]+$/,
    requireify, options = {},
    out;


options.main = argv.file;
options.verbose = argv.verbose != null ? !!argv.verbose : true;
options.ignore = argv.ignore != null ? argv.ignore.split(",") : [];
options.exportName = argv.exportName;

if (argv.out) {
    out = argv.out;
    if (!hasExtension.test(out)) out += ".js";
} else {
    out = filePath.join(filePath.dir(options.main), filePath.basename(options.main, filePath.ext(options.main))) + ".min.js";
}

if (argv.argv) {
    options.argv = argv.argv.split(",");
}

if (argv.env) {
    options.env = {};
    argv.env.split(",").forEach(function(str) {
        var split = str.split("="),
            key = split[0],
            value = split[1];

        if (!key) return;
        if (value) options.env[key] = value;
    });
}

requireify = new Requireify(options);

console.log("\nwriting compiled file " + out + "\n");
fs.writeFileSync(out, requireify.compile());
