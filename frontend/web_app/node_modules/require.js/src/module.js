var urlPath = require("url_path");


var hasExtension = /\.(js|json)$/,
    MODULE_SPLITER = /([^/]*)(\/.*)?/,
    SPLITER = /[\/]+/,
    FUNC_REPLACER = /[\.\/\-\@]/g,

    MODULE_PATH = "",

    nativeKeys = Object.keys,
    toString = Object.prototype.toString,
    hasOwnProp = Object.prototype.hasOwnProperty,
    objectKeys;

function isString(obj) {
    var type = typeof obj;

    return type === "string" || (obj && type === "object" && toString.call(obj) === "[object String]") || false;
}

function isObject(obj) {
    var typeStr;

    if (obj === null || obj === undefined) return false;
    typeStr = typeof(obj);

    return (typeStr === "function" || typeStr === "object");
}

if (nativeKeys) {
    objectKeys = function(obj) {
        if (!isObject(obj)) return [];
        return nativeKeys(obj);
    };
} else {
    objectKeys = function(obj) {
        var has = hasOwnProp,
            keys = [],
            key;

        if (!isObject(obj)) return keys;

        for (key in obj) {
            if (has.call(obj, key)) keys.push(key);
        }
        return keys;
    };
}

function arrayMap(array, callback) {
    var i = -1,
        j = -1,
        length = array.length,
        result = [];

    while (++i < length) {
        result[++j] = callback(array[i], i, array);
    }
    return result;
}

function Context() {
    this.require = null;
    this.exports = null;
    this.__filename = null;
    this.__dirname = null;
    this.module = null;
    this.process = null;
    this.Buffer = null;
    this.global = null;
}


function Module(id, parent) {

    this.id = id;
    this.parent = parent;

    this.exports = {};

    this.dirname = null;
    this.filename = null;
    this.require = null;

    this.loaded = false;
    this.children = [];

    this.__MODULE_PATH__ = MODULE_PATH;

    if (parent) {
        parent.children.push(this);
    }
}

Module._cache = {};

Module.init = function(path, isGlobal) {
    var module;

    if (isGlobal) {
        module = Module.repl();
        module.require(path);
    } else {
        load(resolveFilename(path), null, true);
    }
};

Module.repl = function() {
    var filename = "./repl",
        cache = Module._cache,
        module = new Module("repl", undefined);

    module.filename = filename;
    module.dirname = urlPath.dir(filename);

    global.require = createRequire(module);
    global.module = module;

    cache[filename] = module;
    module.loaded = true;

    return module;
};

function moduleRequire(path) {
    if (!path) throw new Error("require(path) missing path");
    if (!isString(path)) throw new Error("require(path) path must be a string");
    return load(path, this, false, false);
}

function createRequire(module) {

    function require(path) {
        return moduleRequire.call(module, require.resolve(path));
    }

    require.resolve = function(path) {
        return resolveFilename(path, module);
    };

    module.require = require;

    return require;
}

function compile(module, content) {
    var context = new Context();

    context.require = createRequire(module);
    context.exports = module.exports;
    context.__filename = module.filename;
    context.__dirname = module.dirname;
    context.module = module;
    context.process = process;
    context.Buffer = Buffer;
    context.global = global;

    try {
        runInContext(content, context);
    } catch (e) {
        e.message = module.filename + ": " + e.message;
        throw e;
    }
}

function loadModule(module) {
    var filename = module.filename,
        ext = urlPath.ext(module.filename),
        content;

    if (ext === ".js") {
        content = readFile(filename);
        compile(module, content);
    } else if (ext === ".json") {
        content = readFile(filename);

        try {
            module.exports = JSON.parse(content);
        } catch (e) {
            e.message = filename + ": " + e.message;
            throw e;
        }
    } else {
        throw new Error("extension " + ext + " not supported");
    }

    module.loaded = true;
}

function load(path, parent, isMain) {
    var filename = path,
        cache = Module._cache,
        module = cache[filename],
        failed = true;

    if (!module) {
        module = new Module(filename, parent);

        module.filename = filename;
        module.dirname = urlPath.dir(filename);

        if (isMain) module.id = ".";

        cache[filename] = module;

        try {
            loadModule(module);
            failed = false;
        } finally {
            if (failed) delete cache[filename];
        }
    }

    return module.exports;
}

function exists(src) {
    var request;

    try {
        request = new global.XMLHttpRequest();

        request.open("HEAD", src, false);
        request.send(null);
    } catch (e) {
        return false;
    }

    return request.status !== 404;
}

function readFile(src) {
    var request, status;

    try {
        request = new global.XMLHttpRequest();

        request.open("GET", src, false);
        request.send(null);
        status = request.status;
    } catch (e) {}

    return (status === 200 || status === 304) ? request.responseText : null;
}

function resolveFilename(path, parent) {
    MODULE_PATH = false;
    if (urlPath.isAbsoluteURL(path)) return path;
    if (path[0] !== "." && path[0] !== "/") return resolveNodeModule(path, parent);
    if (parent) path = urlPath.join(parent.dirname, path);
    if (path[path.length - 1] === "/") path += "index.js";
    if (!hasExtension.test(path)) path += ".js";

    return path;
}

function resolveNodeModule(path, parent) {
    var found = false,
        paths = path.match(MODULE_SPLITER),
        moduleName = paths[1],
        relativePath = paths[2],
        id = "node_modules/" + moduleName + "/package.json",
        depth = urlPath.join(process.cwd(), (parent ? parent.dirname : "./")).split(SPLITER).length,
        error = false,
        root = (parent ? parent.dirname : "./"),
        resolved = parent.__MODULE_PATH__ ? urlPath.join(parent.__MODULE_PATH__, id) : id,
        pkg;

    if (exists(resolved)) found = true;

    while (!found && depth-- > 0) {
        resolved = urlPath.join(root, id);
        root = root + "/../";
        if (exists(resolved)) found = true;
    }

    if (found) {
        try {
            pkg = JSON.parse(readFile(resolved));
        } catch (e) {
            error = true;
        }

        MODULE_PATH = urlPath.dir(resolved);
        if (pkg) resolved = urlPath.join(MODULE_PATH, parseMain(pkg));

        if (relativePath) {
            resolved = urlPath.join(urlPath.dir(resolved), relativePath);
            if (resolved[resolved.length - 1] === "/") resolved += "index.js";
            if (!hasExtension.test(resolved)) resolved += ".js";
            if (!exists(resolved)) throw new Error("Cannot find module file " + resolved);
        }

        if (!hasExtension.test(resolved)) resolved += ".js";
    } else {
        error = true;
    }

    if (error) throw new Error("Module failed to find node module " + moduleName);

    return resolved;
}

function parseMain(pkg) {
    return (
        isString(pkg.main) ? pkg.main : (
            isString(pkg.browser) ? pkg.browser : "index"
        )
    );
}

function runInContext(content, context) {
    eval(
        '//# sourceURL=' + context.__filename + '\n' +
        '(function ' + ((context.__filename || ".").replace(FUNC_REPLACER, "_")) + '(' + objectKeys(context).join(", ") + ') {\n' +
        content + '\n' +
        '}).call(context.exports, ' + arrayMap(objectKeys(context), function(value) {
            return 'context.' + value;
        }).join(", ") + ');'
    );
}


module.exports = Module;
