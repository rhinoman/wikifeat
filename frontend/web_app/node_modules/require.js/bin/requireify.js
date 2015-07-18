var fs = require("fs"),
    type = require("type"),
    filePath = require("file_path"),
    template = require("template");


var hasExtension = /\.(js|json)$/,

    MODULE_SPLITER = /([^/]*)(\/.*)?/,
    SPLITER = /[\/]+/,

    REPLACER = /(require|require\.resolve)\s*\(\s*["']([^'"\s]+)["']\s*\)/g,
    COMMENT = /(\/\*([\s\S]*?)\*\/|([^:]|^)\/\/(.*)$)/mg,
    BUFFER = /\bBuffer\b/mg,
    BUFFER_DEFINED = /\bfunction\b[\s]+\bBuffer\b/mg,
    PROCESS = /\bprocess\b/mg,

    builtIn = [
        "assert",
        "buffer",
        "child_process",
        "cluster",
        "crypto",
        "dgram",
        "dns",
        "domain",
        "events",
        "fs",
        "http",
        "https",
        "net",
        "os",
        "path",
        "querystring",
        "readline",
        "stream",
        "string_decoder",
        "tls",
        "tty",
        "url",
        "util",
        "vm",
        "zlib",
        "smalloc",
        "tracing"
    ];


function Requireify(options) {
    type.isObject(options) || (options = {});

    if (!type.isString(options.main)) {
        throw new Error("Requireify(options) main must be a string");
    }

    this.exportName = type.isString(options.exportName) ? options.exportName : null;

    this.verbose = options.verbose != null ? !!options.verbose : true;

    this.argv = type.isArray(options.argv) && options.argv.length > 0 ? options.argv : false;
    this.env = type.isObject(options.env) ? options.env : false;

    this.fullPath = filePath.isAbsolute(options.main) ? options.main : filePath.join(process.cwd(), options.main);
    this.fullDirname = filePath.dir(this.fullPath);
    this.main = relative(this.fullDirname, options.main);
    this.ignore = type.isArray(options.ignore) ? options.ignore : [];

    this.parsed = {};
    this.modules = [];
    this.paths = {};

    this.useBuffer = false;
    this.useProcess = false;

    parseDependencies(this, this.fullPath, this.main);
}

Requireify.prototype.createModule = function(relativePath, fullPath, content, isJSON) {
    var paths = this.paths,
        modules = this.modules,
        compiled;

    if (paths[relativePath] != null) return;
    if (this.verbose) console.log("Requireify: found dependency " + relativePath);

    compiled = [
        'function(__require__, require, exports, __filename, __dirname, module, process, Buffer, global) {',
        '',
        isJSON ? "module.exports = " + content : content,
        '',
        '}'
    ].join("\n");

    modules.push({
        filename: relativePath,
        dirname: filePath.dir(relative(this.fullDirname, fullPath)),
        content: compiled
    });

    paths[relativePath] = modules.length - 1;
};

Requireify.prototype.compile = function() {
    var temp = fs.readFileSync(filePath.join(__dirname, "./templates/template.ejs")).toString(),

        buffer = this.useBuffer ? fs.readFileSync(filePath.join(__dirname, "./templates/buffer.ejs")).toString() : 'undefined',
        process = this.useProcess ? fs.readFileSync(filePath.join(__dirname, "./templates/process.ejs")).toString() : 'undefined';

    return template(temp, {

        USE_EXPORT_NAME: !!this.exportName,
        exportName: '"' + this.exportName + '"',

        USE_ARGV: !!this.argv,
        argv: this.argv && this.argv.map(function(value) {
            return '"' + value + '"';
        }).join(', '),

        USE_ENV: !!this.env,
        env: this.env && JSON.stringify(this.env),

        main: '"' + this.main + '"',

        modules: modulesToString(this.modules),
        paths: JSON.stringify(this.paths),

        Buffer: buffer,
        process: process,

        date: (new Date()).toString()
    });
};

function modulesToString(modules) {

    return '[\n' +
        modules.map(function(obj) {
            var out = '[';

            out += obj.content + ', ';
            out += '"' + obj.filename + '", ';
            out += '"' + obj.dirname + '"';

            out += ']';

            return out;
        }).join(',\n') +
        ']';
}

function parseDependencies(requireify, fullPath, relativePath) {
    var content = fs.readFileSync(fullPath).toString(),
        isJSON = filePath.ext(fullPath).toLowerCase() === ".json",
        cleanedContent = content.replace(COMMENT, ""),
        nativeDependencies = {},
        dependencies = {};

    if (!requireify.useBuffer && BUFFER.test(cleanedContent) && !BUFFER_DEFINED.test(cleanedContent)) {
        requireify.useBuffer = true;
    }
    if (!requireify.useProcess && PROCESS.test(cleanedContent)) {
        requireify.useProcess = true;
    }

    requireify.parsed[relativePath] = true;

    cleanedContent.replace(REPLACER, function(match, method, dependency) {
        var resolvedFullPath, resolvedPath;

        if (requireify.ignore.indexOf(dependency) !== -1) return;

        if (isModule(dependency)) {
            if (builtIn.indexOf(dependency) !== -1) {
                if (requireify.verbose) {
                    console.warn(
                        "Requireify: found Node.js dependency\n   make sure to check if in browser before " +
                        "trying to load this in files if (!process.browser) " + dependency + " = require(\"" + dependency + "\")\n"
                    );
                }
                nativeDependencies[dependency] = dependency;
                return;
            }

            resolvedFullPath = resolveNodeModulePath(dependency, filePath.dir(fullPath), relativePath);
            resolvedPath = dependency;
        } else {
            resolvedFullPath = resolveModulePath(dependency, filePath.dir(fullPath), relativePath);
            resolvedPath = relative(requireify.fullDirname, resolvedFullPath);
        }

        dependencies[dependency] = resolvedPath;
        if (method === "require.resolve" || requireify.parsed[resolvedPath]) return;

        parseDependencies(requireify, resolvedFullPath, resolvedPath);
    });

    content = content.replace(REPLACER, function(match, method, dependency) {
        var path = dependencies[dependency];

        if (!path && (path = nativeDependencies[dependency])) {
            return '__require__("' + path + '")';
        }

        return !!path ? method + '("' + path + '")' : match;
    });

    requireify.createModule(relativePath, fullPath, content, isJSON);
}

function relative(dir, path) {
    path = filePath.relative(dir, path);
    if (!(path[0] === "." || path[0] === "/")) path = "./" + path;
    return path;
}

function resolveModulePath(path, parentDirname, fromPath) {
    var resolved = filePath.join(parentDirname, path),
        stat = statFile(resolved),
        pkg, tmp;

    if (stat && stat.isDirectory()) {
        tmp = filePath.join(resolved, "index.js");

        if (exists(tmp)) return tmp;

        tmp = filePath.join(resolved, "package.json");

        if (exists(tmp)) {
            pkg = JSON.parse(fs.readFileSync(tmp).toString());
            resolved = filePath.join(filePath.dir(tmp), pkg.main);
        }
    }
    if (!hasExtension.test(resolved)) resolved += ".js";

    if (!exists(resolved)) {
        throw new Error("Requireify: no file found with path " + resolved + " required from " + fromPath);
    }

    return resolved;
}

function resolveNodeModulePath(path, parentDirname, fromPath) {
    var found = false,
        paths = path.match(MODULE_SPLITER),
        moduleName = paths[1],
        relativePath = paths[2],
        id = "node_modules/" + moduleName + "/package.json",
        depth = (parentDirname || process.cwd()).split(SPLITER).length,
        error = false,
        root = (parentDirname || process.cwd()),
        resolved = filePath.join(root, id),
        pkg;

    if (exists(resolved)) found = true;

    while (!found && depth-- > 0) {
        resolved = filePath.join(root, id);
        root = root + "/../";
        if (exists(resolved)) found = true;
    }

    if (found) {
        try {
            pkg = JSON.parse(fs.readFileSync(resolved).toString());
        } catch (e) {
            error = true;
        }

        if (pkg) resolved = filePath.join(filePath.dir(resolved), parseMain(pkg));

        if (relativePath) {
            resolved = filePath.join(filePath.dir(resolved), relativePath);
            if (resolved[resolved.length - 1] === "/") resolved += "index.js";
            if (!hasExtension.test(resolved)) resolved += ".js";
            if (!exists(resolved)) throw new Error("Cannot find module file " + resolved + " required from " + fromPath);
        } else if (!hasExtension.test(resolved)) {
            resolved += ".js";
        }
    } else {
        error = true;
    }

    if (error) throw new Error("Module failed to find node module " + moduleName + " required from " + fromPath);

    return resolved;
}

function parseMain(pkg) {
    return (
        type.isString(pkg.main) ? pkg.main : (
            type.isString(pkg.browser) ? pkg.browser : "index"
        )
    );
}

function statFile(path) {
    var stat;

    try {
        stat = fs.statSync(path);
    } catch (e) {}
    return stat;
}

function exists(path) {

    return !!fs.existsSync(path);
}

function isModule(path) {
    return !!(path[0] !== "." && path[0] !== "/");
}


module.exports = Requireify;
