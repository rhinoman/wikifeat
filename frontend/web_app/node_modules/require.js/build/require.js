// Compiled with Require.js on Sat Nov 29 2014 10:47:30 GMT-0600 (CST)

(function(main, modules, paths, Buffer, process, __require__, global) {
    var isCommonJS = typeof(module) !== "undefined" && module.exports,
        cache = {};

    if (!isCommonJS && (process !== null && process !== undefined)) {


        if (process.argv.length === 0) {
            process.argv.push("browser", main);
        }

    }

    function Module(filename, dirname) {

        this.id = filename;

        this.exports = {};

        this.dirname = dirname;
        this.filename = filename;

        this.loaded = true;
    }

    function require(path) {
        var module = cache[path],
            data;

        if (!module) {
            data = modules[paths[path]];
            cache[path] = module = new Module(data[1], data[2]);
            module.require = require;
            data[0].call(module.exports, __require__, require, module.exports, module.filename, module.dirname, module, process, Buffer, global);
        }

        return module.exports;
    }
    require.resolve = function(path) {
        return path;
    };
    Module.prototype.require = require;

    if (isCommonJS) {
        module.exports = require(main);
    } else {

        require(main);

    }
}(
    "./index.js", [
        [
            function(__require__, require, exports, __filename, __dirname, module, process, Buffer, global) {

                var pathUtils = module.exports;


                pathUtils.removeEmpties = function(parts) {
                    var i = parts.length;

                    while (i--) {
                        if (!parts[i]) parts.splice(i, 1);
                    }

                    return parts;
                };

                pathUtils.trim = function(parts) {
                    var length = parts.length - 1,
                        start = -1,
                        end;

                    while (start++ < length) {
                        if (parts[start] !== "") break;
                    }

                    end = length + 1;
                    while (end--) {
                        if (parts[end] !== "") break;
                    }

                    if (start > end) return [];

                    return parts.slice(start, end + 1);
                };

                pathUtils.normalizeArray = function(parts, allowAboveRoot) {
                    var i = parts.length,
                        up = 0,
                        last;

                    while (i--) {
                        last = parts[i];

                        if (last === ".") {
                            parts.splice(i, 1);
                        } else if (last === "..") {
                            parts.splice(i, 1);
                            up++;
                        } else if (up) {
                            parts.splice(i, 1);
                            up--;
                        }
                    }

                    if (allowAboveRoot) {
                        while (up--) parts.unshift("..");
                    }

                    return parts;
                };


            }, "path_utils", "../node_modules/url_path/node_modules/path_utils/src"
        ],
        [
            function(__require__, require, exports, __filename, __dirname, module, process, Buffer, global) {

                var pathUtils = require("path_utils");


                var urlPath = module.exports,
                    IS_URL = /^(?:[a-z]+:)?\/\//i;


                urlPath.isAbsolute = function(str) {
                    return str[0] === "/" || IS_URL.test(str);
                };

                urlPath.isAbsoluteURL = function(str) {
                    return IS_URL.test(str);
                };

                urlPath.isURL = urlPath.isAbsoluteURL;

                urlPath.normalize = function(str) {
                    var isAbs = urlPath.isAbsolute(str),
                        trailingSlash = str[str.length - 1] === "/",
                        segments = str.split("/"),
                        nonEmptySegments = [],
                        i;

                    for (i = 0; i < segments.length; i++) {
                        if (segments[i]) nonEmptySegments.push(segments[i]);
                    }
                    str = pathUtils.normalizeArray(nonEmptySegments, !isAbs).join("/");

                    if (!str && !isAbs) str = ".";
                    if (str && trailingSlash) str += "/";

                    return (isAbs ? "/" : "") + str;
                };

                urlPath.resolve = function() {
                    var resolvedPath = "",
                        resolvedAbsolute = false,
                        i, str;

                    for (i = arguments.length - 1; i >= -1 && !resolvedAbsolute; i--) {
                        str = (i >= 0) ? arguments[i] : process.cwd();

                        if (typeof(str) !== "string") {
                            throw new TypeError("Arguments to resolve must be strings");
                        } else if (!str) {
                            continue;
                        }

                        resolvedPath = str + "/" + resolvedPath;
                        resolvedAbsolute = str.charAt(0) === "/";
                    }

                    resolvedPath = pathUtils.normalizeArray(pathUtils.removeEmpties(resolvedPath.split("/")), !resolvedAbsolute).join("/");
                    return ((resolvedAbsolute ? "/" : "") + resolvedPath) || ".";
                };

                urlPath.relative = function(from, to) {
                    from = resolve(from).substr(1);
                    to = resolve(to).substr(1);

                    var fromParts = pathUtils.trim(from.split("/")),
                        toParts = pathUtils.trim(to.split("/")),

                        length = Math.min(fromParts.length, toParts.length),
                        samePartsLength = length,
                        outputParts, i, il;

                    for (i = 0; i < length; i++) {
                        if (fromParts[i] !== toParts[i]) {
                            samePartsLength = i;
                            break;
                        }
                    }

                    outputParts = [];
                    for (i = samePartsLength, il = fromParts.length; i < il; i++) outputParts.push("..");
                    outputParts = outputParts.concat(toParts.slice(samePartsLength));

                    return outputParts.join("/");
                };

                urlPath.join = function() {
                    var str = "",
                        segment,
                        i, il;

                    for (i = 0, il = arguments.length; i < il; i++) {
                        segment = arguments[i];

                        if (typeof(segment) !== "string") {
                            throw new TypeError("Arguments to join must be strings");
                        }
                        if (segment) {
                            if (!str) {
                                str += segment;
                            } else {
                                str += "/" + segment;
                            }
                        }
                    }

                    return urlPath.normalize(str);
                };

                urlPath.dir = function(str) {
                    str = str.substring(0, str.lastIndexOf("/") + 1);
                    return str ? str.substr(0, str.length - 1) : ".";
                };

                urlPath.dirname = urlPath.dir;

                urlPath.base = function(str, ext) {
                    str = str.substring(str.lastIndexOf("/") + 1);

                    if (ext && str.substr(-ext.length) === ext) {
                        str = str.substr(0, str.length - ext.length);
                    }

                    return str || "";
                };

                urlPath.basename = urlPath.base;

                urlPath.ext = function(str) {
                    var index = str.lastIndexOf(".");
                    return index > -1 ? str.substring(index) : "";
                };

                urlPath.extname = urlPath.ext;


            }, "url_path", "../node_modules/url_path/src"
        ],
        [
            function(__require__, require, exports, __filename, __dirname, module, process, Buffer, global) {

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


            }, "./module.js", "."
        ],
        [
            function(__require__, require, exports, __filename, __dirname, module, process, Buffer, global) {

                var Module = require("./module.js"),

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


            }, "./index.js", "."
        ]
    ], {
        "path_utils": 0,
        "url_path": 1,
        "./module.js": 2,
        "./index.js": 3
    }, (typeof(Buffer) !== "undefined" ? Buffer : (function() {
        var toString = Object.prototype.toString,
            isArray = Array.isArray || (function isArray(obj) {
                return toString.call(obj) === "[object Array]";
            }),
            base64;


        function Buffer(subject, encoding) {
            if (!(this instanceof Buffer)) return new Buffer(subject, encoding);
            var type = typeof(subject),
                i, il;

            if (encoding === "base64" && type === "string") {
                subject = trim(subject);
                while (subject.length % 4 !== 0) subject = subject + "=";
            }

            if (type === "number") {
                this.length = coerce(subject);
            } else if (type === "string") {
                this.length = Buffer.byteLength(subject, encoding);
            } else if (type === "object" && subject.length === +subject.length) {
                this.length = coerce(subject.length);
            } else {
                throw new Error("Buffer(subject, encoding): First argument needs to be a number, array or string.");
            }

            if (type === "string") {
                this.write(subject, encoding);
            } else if (type === "number") {
                for (i = 0, il = this.length; i < il; i++) this[i] = 0;
            }

            return this;
        }

        Buffer.Buffer = Buffer;
        Buffer.SlowBuffer = Buffer;
        Buffer.poolSize = 8192;
        Buffer.INSPECT_MAX_BYTES = 50;

        Buffer.prototype.write = function(string, offset, length, encoding) {
            if (isFinite(offset)) {
                if (!isFinite(length)) {
                    encoding = length;
                    length = undefined;
                }
            } else {
                var swap = encoding;
                encoding = offset;
                offset = length;
                length = swap;
            }
            offset = +offset || 0;
            var remaining = this.length - offset;

            if (!length) {
                length = remaining;
            } else {
                length = +length;
                if (length > remaining) length = remaining;
            }

            encoding = (encoding || "utf8").toLowerCase();

            if (encoding === "utf8" || encoding === "utf-8") {
                return this.utf8Write(string, offset, length);
            } else if (encoding === "ascii" || encoding === "raw") {
                return this.asciiWrite(string, offset, length);
            } else if (encoding === "binary") {
                return this.binaryWrite(string, offset, length);
            } else if (encoding === "ucs2" || encoding === "ucs-2" || encoding === "utf16le" || encoding === "utf-16le") {
                return this.utf16leWrite(string, offset, length);
            } else if (encoding === "hex") {
                return this.hexWrite(string, offset, length);
            } else if (encoding === "base64") {
                return this.base64Write(string, offset, length);
            } else {
                throw new Error("Buffer.write(string, offset, length, encoding) Unknown encoding " + encoding);
            }

            return "";
        };

        Buffer.prototype.copy = function(target, target_start, start, end) {
            if (!start) start = 0;
            if (!end && end !== 0) end = this.length;
            if (!target_start) target_start = 0;

            if (end === start) return;
            if (target.length === 0 || this.length === 0) return;

            if (end < start) throw new Error("Buffer.copy(target, target_start, start, end) sourceEnd < sourceStart");
            if (target_start >= 0 && target_start >= target.length) throw new Error("Buffer.copy(target, target_start, start, end)targetStart out of bounds");
            if (start >= 0 && start >= this.length) throw new Error("Buffer.copy(target, target_start, start, end)sourceStart out of bounds");
            if (end >= 0 && end > this.length) throw new Error("Buffer.copy(target, target_start, start, end)sourceEnd out of bounds");

            if (end > this.length) end = this.length;
            if (target.length - target_start < end - start) end = target.length - target_start + start;

            var i = 0,
                il = end - start;
            for (; i < il; i++) target[i + target_start] = this[i + start];
        };

        Buffer.prototype.slice = function(start, end) {
            var len = this.length,
                sliceLen, buffer = new Buffer(sliceLen, undefined, true),
                i = 0;

            start = clamp(start, len, 0);
            end = clamp(end, len, len);

            sliceLen = end - start;
            for (; i < sliceLen; i++) buffer[i] = this[i + start];
            return buffer;
        };

        Buffer.prototype.fill = function(value, start, end) {
            if (!value) value = 0;
            if (!start) start = 0;
            if (!end) end = this.length;

            if (end < start) throw new Error("Buffer.fill(value, start, end) end < start");

            if (end === start) return this;
            if (this.length === 0) return this;

            if (start >= 0 && start >= this.length) throw new Error("Buffer.fill(value, start, end) start out of bounds");
            if (end >= 0 && end > this.length) throw new Error("Buffer.fill(value, start, end) endout of bounds");

            var i = start,
                bytes, len;
            if (typeof(value) === "number") {
                for (i = start; i < end; i++) this[i] = value;
            } else {
                bytes = utf8ToBytes(value.toString());
                len = bytes.length;
                for (i = start; i < end; i++) this[i] = bytes[i % len];
            }

            return this;
        };

        Buffer.prototype.inspect = function() {
            var out = [],
                len = this.length,
                i = 0;
            for (; i < len; i++) {
                out[i] = toHex(this[i]);
                if (i === Buffer.INSPECT_MAX_BYTES) {
                    out[i + 1] = "...";
                    break;
                }
            }

            return "<Buffer " + out.join(" ") + ">";
        };

        Buffer.prototype.equals = function(b) {
            if (!Buffer.isBuffer(b)) throw new Error("Buffer.equals(b) Argument must be a Buffer");
            return Buffer.compare(this, b);
        };

        Buffer.prototype.toJSON = function() {
            var jsonArray = [],
                i = this.length;

            while (i--) jsonArray[i] = this[i];

            return {
                type: "Buffer",
                data: jsonArray
            };
        };

        Buffer.prototype.toArrayBuffer = function() {
            var buffer,
                i = this.length;

            if (typeof(Uint8Array) !== "undefined") {
                buffer = new Uint8Array(i);
                while (i--) buffer[i] = this[i];
            } else {
                buffer = [];
                while (i--) buffer[i] = this[i];
            }

            return buffer;
        };

        Buffer.prototype.toString = Buffer.prototype.toLocaleString = function(encoding, start, end) {
            encoding = (encoding || "utf8").toLowerCase();
            start || (start = 0);
            end = (end == null) ? this.length : +end;

            if (end === start) return "";

            if (encoding === "utf8" || encoding === "utf-8") {
                return this.utf8Slice(start, end);
            } else if (encoding === "ascii" || encoding === "raw") {
                return this.asciiSlice(start, end);
            } else if (encoding === "binary") {
                return this.binarySlice(start, end);
            } else if (encoding === "ucs2" || encoding === "ucs-2" || encoding === "utf16le" || encoding === "utf-16le") {
                return this.utf16leSlice(start, end);
            } else if (encoding === "hex") {
                return this.hexSlice(start, end);
            } else if (encoding === "base64") {
                return this.base64Slice(start, end);
            } else {
                throw new Error("Buffer.write(string, offset, length, encoding) Unknown encoding " + encoding);
            }

            return "";
        };

        Buffer.prototype.hexWrite = function(string, offset, length) {
            offset || (offset = 0);
            var remaining = this.length - offset,
                strLen, b, i;

            if (!length) {
                length = remaining;
            } else {
                length = +length;
                if (length > remaining) length = remaining;
            }

            strLen = string.length;
            if (strLen % 2 !== 0) throw new Error("Buffer.hexWrite(string, offset, length) Invalid hex string");

            if (length > strLen / 2) {
                length = strLen / 2;
            }
            for (i = 0; i < length; i++) {
                b = parseInt(string.substr(i * 2, 2), 16);
                if (isNaN(b)) throw new Error("Buffer.hexWrite(string, offset, length) Invalid hex string");
                this[offset + i] = b;
            }

            return i;
        };

        Buffer.prototype.utf8Write = function(string, offset, length) {

            return blitBuffer(this, utf8ToBytes(string), offset, length);
        };

        Buffer.prototype.asciiWrite = function(string, offset, length) {

            return blitBuffer(this, asciiToBytes(string), offset, length);
        };

        Buffer.prototype.binaryWrite = function(string, offset, length) {

            return blitBuffer(this, string, offset, length);
        };

        Buffer.prototype.base64Write = function(string, offset, length) {

            return blitBuffer(this, base64ToBytes(string), offset, length);
        };

        Buffer.prototype.utf16leWrite = function(string, offset, length) {

            return blitBuffer(this, utf16leToBytes(string), offset, length);
        };

        Buffer.prototype.utf8Slice = function(start, end) {
            end = Math.min(this.length, end);
            var res = "",
                tmp = "",
                i = start,
                b;

            for (; i < end; i++) {
                b = this[i];
                if (b <= 0x7F) {
                    res += decodeUtf8Char(tmp) + String.fromCharCode(b);
                    tmp = "";
                } else {
                    tmp += "%" + b.toString(16);
                }
            }

            return res + decodeUtf8Char(tmp);
        };

        Buffer.prototype.base64Slice = function(start, end) {
            if (start === 0 && end === this.length) {
                return base64.encode(this);
            } else {
                return base64.encode(this.slice(start, end));
            }
        };

        Buffer.prototype.asciiSlice = function(start, end) {
            end = Math.min(this.length, end);
            var ret = "",
                i = start;

            for (; i < end; i++) ret += String.fromCharCode(this[i]);
            return ret;
        };

        Buffer.prototype.binarySlice = Buffer.prototype.asciiSlice;

        Buffer.prototype.hexSlice = function(start, end) {
            var len = this.length,
                out = "",
                i;

            if (!start || start < 0) start = 0;
            if (!end || end < 0 || end > len) end = len;

            for (i = start; i < end; i++) out += toHex(this[i]);
            return out;
        };

        Buffer.prototype.utf16leSlice = function(start, end) {
            var bytes = this.slice(start, end),
                i = 0,
                il = bytes.length,
                res = "";

            for (; i < il; i += 2) res += String.fromCharCode(bytes[i] + bytes[i + 1] * 256);
            return res;
        };

        function checkOffset(offset, ext, length) {
            if (offset + ext > length) throw new RangeError("index out of range");
        }

        Buffer.prototype.readUInt8 = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 1, this.length);
            return this[offset];
        };

        Buffer.prototype.readUInt16LE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 2, this.length);
            return this[offset] | (this[offset + 1] << 8);
        };

        Buffer.prototype.readUInt16BE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 2, this.length);
            return (this[offset] << 8) | this[offset + 1];
        };

        Buffer.prototype.readUInt32LE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 4, this.length);

            return ((this[offset]) | (this[offset + 1] << 8) | (this[offset + 2] << 16)) + (this[offset + 3] * 0x1000000);
        };

        Buffer.prototype.readUInt32BE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 4, this.length);

            return (this[offset] * 0x1000000) + ((this[offset + 1] << 16) | (this[offset + 2] << 8) | (this[offset + 3]) >>> 0);
        };

        Buffer.prototype.readInt8 = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 1, this.length);
            var val = this[offset];
            return !(val & 0x80) ? val : (0xff - val + 1) * -1;
        };

        Buffer.prototype.readInt16LE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 2, this.length);
            var val = this[offset] | (this[offset + 1] << 8);
            return (val & 0x8000) ? val | 0xFFFF0000 : val;
        };

        Buffer.prototype.readInt16BE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 2, this.length);
            var val = this[offset + 1] | (this[offset] << 8);
            return (val & 0x8000) ? val | 0xFFFF0000 : val;
        };

        Buffer.prototype.readInt32LE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 4, this.length);

            return (this[offset]) | (this[offset + 1] << 8) | (this[offset + 2] << 16) | (this[offset + 3] << 24);
        };

        Buffer.prototype.readInt32BE = function(offset, noAssert) {
            offset = offset >>> 0;
            if (!noAssert) checkOffset(offset, 4, this.length);

            return (this[offset] << 24) | (this[offset + 1] << 16) | (this[offset + 2] << 8) | (this[offset + 3]);
        };

        function checkInt(buffer, value, offset, ext, max, min) {
            if (!(buffer instanceof Buffer)) throw new TypeError("buffer must be a Buffer instance");
            if (value > max || value < min) throw new TypeError("value is out of bounds");
            if (offset + ext > buffer.length) throw new RangeError("index out of range");
        }

        Buffer.prototype.writeUInt8 = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 1, 0xff, 0);
            this[offset] = value;
            return offset + 1;
        };

        Buffer.prototype.writeUInt16LE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 2, 0xffff, 0);
            this[offset] = value;
            this[offset + 1] = (value >>> 8);
            return offset + 2;
        };

        Buffer.prototype.writeUInt16BE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 2, 0xffff, 0);
            this[offset] = (value >>> 8);
            this[offset + 1] = value;
            return offset + 2;
        };

        Buffer.prototype.writeUInt32LE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 4, 0xffffffff, 0);
            this[offset + 3] = (value >>> 24);
            this[offset + 2] = (value >>> 16);
            this[offset + 1] = (value >>> 8);
            this[offset] = value;
            return offset + 4;
        };

        Buffer.prototype.writeUInt32BE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 4, 0xffffffff, 0);
            this[offset] = (value >>> 24);
            this[offset + 1] = (value >>> 16);
            this[offset + 2] = (value >>> 8);
            this[offset + 3] = value;
            return offset + 4;
        };

        Buffer.prototype.writeInt8 = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 1, 0x7f, -0x80);
            this[offset] = value;
            return offset + 1;
        };

        Buffer.prototype.writeInt16LE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 2, 0x7fff, -0x8000);
            this[offset] = value;
            this[offset + 1] = (value >>> 8);
            return offset + 2;
        };

        Buffer.prototype.writeInt16BE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert)
                checkInt(this, value, offset, 2, 0x7fff, -0x8000);
            this[offset] = (value >>> 8);
            this[offset + 1] = value;
            return offset + 2;
        };

        Buffer.prototype.writeInt32LE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 4, 0x7fffffff, -0x80000000);
            this[offset] = value;
            this[offset + 1] = (value >>> 8);
            this[offset + 2] = (value >>> 16);
            this[offset + 3] = (value >>> 24);
            return offset + 4;
        };

        Buffer.prototype.writeInt32BE = function(value, offset, noAssert) {
            value = +value;
            offset = offset >>> 0;
            if (!noAssert) checkInt(this, value, offset, 4, 0x7fffffff, -0x80000000);
            this[offset] = (value >>> 24);
            this[offset + 1] = (value >>> 16);
            this[offset + 2] = (value >>> 8);
            this[offset + 3] = value;
            return offset + 4;
        };

        Buffer.isBuffer = function(obj) {
            return obj instanceof Buffer;
        };

        Buffer.isEncoding = function(encoding) {
            if (typeof(encoding) !== "string") return false;
            encoding = encoding.toLowerCase();

            return (
                encoding === "utf8" ||
                encoding === "utf-8" ||
                encoding === "hex" ||
                encoding === "base64" ||
                encoding === "ascii" ||
                encoding === "binary" ||
                encoding === "ucs2" ||
                encoding === "ucs-2" ||
                encoding === "utf16le" ||
                encoding === "utf-16le" ||
                encoding === "raw"
            );
        };

        Buffer.byteLength = function(str, encoding) {
            str = str + "";
            encoding = (encoding || "utf8").toLowerCase();

            if (encoding === "utf8" || encoding === "utf-8") {
                return utf8ToBytes(str).length;
            } else if (encoding === "ascii" || encoding === "binary" || encoding === "raw") {
                return str.length;
            } else if (encoding === "ucs2" || encoding === "ucs-2" || encoding === "utf16le" || encoding === "utf-16le") {
                return str.length * 2;
            } else if (encoding === "hex") {
                return str.length >>> 1;
            } else if (encoding === "base64") {
                return base64ToBytes(str).length;
            } else {
                throw new Error("Buffer.byteLength(str, encoding) Unknown encoding " + encoding);
            }

            return 0;
        };

        Buffer.concat = function(list, totalLength) {
            if (!isArray(list)) throw new Error("Usage: Buffer.concat(list[, length])");

            if (list.length === 0) {
                return new Buffer(0);
            } else if (list.length === 1) {
                return list[0];
            }
            var buffer, postion, item,
                i, il;

            if (totalLength === undefined) {
                totalLength = 0;
                for (i = 0, il = list.length; i < il; i++) totalLength += list[i].length;
            }

            buffer = new Buffer(totalLength);
            postion = 0;
            for (i = 0, il = list.length; i < il; i++) {
                item = list[i];
                item.copy(buffer, postion);
                postion += item.length;
            }

            return buffer;
        };

        Buffer.compare = function(a, b) {
            if (!Buffer.isBuffer(a) || !Buffer.isBuffer(b)) throw new Error("Buffer.compare(a, b) Arguments must be Buffers");
            var x = a.length,
                y = b.length,
                i = 0,
                il = Math.min(x, y);

            while (i < il && a[i] === b[i]) i++;
            if (i !== il) {
                x = a[i];
                y = b[i];
            }
            if (x < y) return -1;
            if (y < x) return 1;

            return 0;
        };

        function blitBuffer(out, src, offset, length) {
            var srcLength = src.length,
                outLength = out.length,
                i = 0;

            for (; i < length; i++) {
                if ((i + offset >= outLength) || (i >= srcLength)) break;
                out[i + offset] = src[i];
            }

            return i;
        }

        function toHex(num) {
            return num < 16 ? "0" + num.toString(16) : num.toString(16);
        }

        function utf8ToBytes(str) {
            var length = str.length,
                byteArray = [],
                start, b, h,
                i = 0,
                j;

            for (; i < length; i++) {
                b = str.charCodeAt(i);
                if (b <= 0x7F) {
                    byteArray.push(b);
                } else {
                    start = i;
                    if (b >= 0xD800 && b <= 0xDFFF) i++;

                    h = encodeURIComponent(str.slice(start, i + 1)).substr(1).split("%");
                    for (j = 0; j < h.length; j++) {
                        byteArray.push(parseInt(h[j], 16));
                    }
                }
            }
            return byteArray;
        }

        function asciiToBytes(str) {
            var byteArray = [],
                i = 0,
                il = str.length;

            for (; i < il; i++) {
                byteArray.push(str.charCodeAt(i) & 0xFF);
            }
            return byteArray;
        }

        function utf16leToBytes(str) {
            var c, hi, lo,
                byteArray = [],
                i = 0,
                il = str.length;

            for (; i < il; i++) {
                c = str.charCodeAt(i);
                hi = c >> 8;
                lo = c % 256;
                byteArray.push(lo);
                byteArray.push(hi);
            }

            return byteArray;
        }

        function base64ToBytes(str) {
            return base64.decode(str);
        }

        base64 = (function() {
            var ArrayType = typeof(Uint8Array) !== "undefined" ? Uint8Array : Array,

                LOOK_UP = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",
                PLUS = "+".charCodeAt(0),
                SLASH = "/".charCodeAt(0),
                NUMBER = "0".charCodeAt(0),
                LOWER = "a".charCodeAt(0),
                UPPER = "A".charCodeAt(0);

            function decode(str) {
                var code = str.charCodeAt(0);

                if (code === PLUS) return 62;
                if (code === SLASH) return 63;
                if (code < NUMBER) return -1;
                if (code < NUMBER + 10) return code - NUMBER + 26 + 26;
                if (code < UPPER + 26) return code - UPPER;
                if (code < LOWER + 26) return code - LOWER + 26;

                return -1;
            }

            function encode(num) {
                return LOOK_UP.charAt(num);
            }

            function tripletToBase64(num) {
                return encode(num >> 18 & 0x3F) + encode(num >> 12 & 0x3F) + encode(num >> 6 & 0x3F) + encode(num & 0x3F);
            }

            return {
                decode: function(str) {
                    if (str.length % 4 > 0) throw new Error("base64.decode(str) Invalid string. Length must be a multiple of 4");
                    var i, j, l, L, tmp, placeHolders, array,
                        len = str.length;

                    placeHolders = "=" === str.charAt(len - 2) ? 2 : "=" === str.charAt(len - 1) ? 1 : 0;
                    array = new ArrayType(str.length * 3 / 4 - placeHolders);
                    l = placeHolders > 0 ? str.length - 4 : str.length;
                    L = 0;

                    for (i = 0, j = 0; i < l; i += 4, j += 3) {
                        tmp = (decode(str.charAt(i)) << 18) | (decode(str.charAt(i + 1)) << 12) | (decode(str.charAt(i + 2)) << 6) | decode(str.charAt(i + 3));
                        array[L++] = (tmp & 0xFF0000) >> 16;
                        array[L++] = (tmp & 0xFF00) >> 8;
                        array[L++] = tmp & 0xFF;
                    }

                    if (placeHolders === 2) {
                        tmp = (decode(str.charAt(i)) << 2) | (decode(str.charAt(i + 1)) >> 4);
                        array[L++] = (tmp & 0xFF);
                    } else if (placeHolders === 1) {
                        tmp = (decode(str.charAt(i)) << 10) | (decode(str.charAt(i + 1)) << 4) | (decode(str.charAt(i + 2)) >> 2);
                        array[L++] = (tmp >> 8) & 0xFF;
                        array[L++] = tmp & 0xFF;
                    }

                    return array;
                },
                encode: function(uint8) {
                    var extraBytes = uint8.length % 3,
                        output = "",
                        temp, length,
                        i;

                    for (i = 0, length = uint8.length - extraBytes; i < length; i += 3) {
                        temp = (uint8[i] << 16) + (uint8[i + 1] << 8) + (uint8[i + 2]);
                        output += tripletToBase64(temp);
                    }

                    if (extraBytes === 1) {
                        temp = uint8[uint8.length - 1];
                        output += encode(temp >> 2);
                        output += encode((temp << 4) & 0x3F);
                        output += "==";
                    } else if (extraBytes === 2) {
                        temp = (uint8[uint8.length - 2] << 8) + (uint8[uint8.length - 1]);
                        output += encode(temp >> 10);
                        output += encode((temp >> 4) & 0x3F);
                        output += encode((temp << 2) & 0x3F);
                        output += "=";
                    }

                    return output;
                }
            };
        }());

        function decodeUtf8Char(str) {
            try {
                return decodeURIComponent(str);
            } catch (err) {
                return String.fromCharCode(0xFFFD);
            }
        }

        function clamp(index, len, defaultValue) {
            if (typeof(index) !== "number") return defaultValue;
            index = ~~index;
            if (index >= len) return len;
            if (index >= 0) return index;
            index += len;
            if (index >= 0) return index;
            return 0;
        }

        function coerce(length) {
            length = ~~Math.ceil(+length);
            return length < 0 ? 0 : length;
        }

        var trim_regex = /^\s+|\s+$/g;

        function trim(str) {
            if (str.trim) return str.trim();
            return str.replace(trim_regex, "");
        }


        return Buffer;
    }())), (typeof(process) !== "undefined" ? process : (function() {
        var global = (new Function("return this;"))(),
            arrayShift = Array.prototype.shift,
            hasOwnProperty = Object.prototype.hasOwnProperty;


        function EventObject(listener, ctx) {
            this.listener = listener;
            this.ctx = ctx;
        }


        function EventEmitter() {

            this._events = Object.create(null);
            this._maxListeners = EventEmitter.defaultMaxListeners;
        }

        EventEmitter.prototype.on = function(type, listener, ctx) {
            if (typeof(listener) !== "function") throw new TypeError("EventEmitter.on(type, listener[, ctx]) listener must be a function");
            var events = this._events,
                eventList = (events[type] || (events[type] = [])),
                maxListeners = this._maxListeners;

            eventList.push(new EventObject(listener, ctx || this));

            if (maxListeners !== -1 && eventList.length > maxListeners) {
                console.error("EventEmitter.on(type, listener, ctx) possible EventEmitter memory leak detected. " + maxListeners + " listeners added");
            }

            return this;
        };

        EventEmitter.prototype.addListener = EventEmitter.prototype.on;

        EventEmitter.prototype.once = function(type, listener, ctx) {
            var _this = this;
            ctx || (ctx = this);

            function once() {
                _this.off(type, once, ctx);
                var length = arguments.length;

                if (length === 0) {
                    return listener.call(ctx);
                } else if (length === 1) {
                    return listener.call(ctx, arguments[0]);
                } else if (length === 2) {
                    return listener.call(ctx, arguments[0], arguments[1]);
                } else if (length === 3) {
                    return listener.call(ctx, arguments[0], arguments[1], arguments[2]);
                } else if (length === 4) {
                    return listener.call(ctx, arguments[0], arguments[1], arguments[2], arguments[3]);
                } else if (length === 5) {
                    return listener.call(ctx, arguments[0], arguments[1], arguments[2], arguments[3], arguments[4]);
                }

                return listener.apply(ctx, arguments);
            }

            return this.on(type, once, ctx);
        };

        EventEmitter.prototype.listenTo = function(obj, type, listener, ctx) {
            if (!(hasOwnProperty.call(obj, "on") && typeof(obj.on) === "function")) {
                throw new TypeError("EventEmitter.listenTo(obj, type, listener, ctx) obj must have a on function taking (type, listener[, ctx])");
            }

            obj.on(type, listener, ctx || this);
            return this;
        };

        EventEmitter.prototype.off = function(type, listener, ctx) {
            var events = this._events,
                eventList, event, i;

            if (!type) return this.removeAllListeners();

            eventList = events[type];
            if (!eventList) return this;

            if (!listener) {
                i = eventList.length;
                while (i--) {
                    event = eventList[i];
                    this.emit("removeListener", type, event.listener, event.ctx);
                }
                eventList.length = 0;
                delete events[type];
            } else {
                ctx = ctx || this;
                i = eventList.length;
                while (i--) {
                    event = eventList[i];

                    if (event.listener === listener) {
                        this.emit("removeListener", type, event.listener, event.ctx);
                        eventList.splice(i, 1);
                    }
                }
                if (eventList.length === 0) delete events[type];
            }

            return this;
        };

        EventEmitter.prototype.removeListener = EventEmitter.prototype.off;

        EventEmitter.prototype.removeAllListeners = function() {
            var events = this._events,
                eventList, event, i, type;

            for (type in events) {
                if ((eventList = events[type])) {
                    i = eventList.length;
                    while (i--) {
                        event = eventList[i];
                        this.emit("removeListener", type, event.listener, event.ctx);
                    }
                    eventList.length = 0;
                    delete events[type];
                }
            }

            return this;
        };

        function emit(eventList, args) {
            var a1, a2, a3, a4,
                length = args.length,
                event, i;

            if (length === 1) {
                i = eventList.length;
                while (i--) {
                    if ((event = eventList[i])) event.listener.call(event.ctx);
                }
            } else if (length === 2) {
                a1 = args[1];
                i = eventList.length;
                while (i--) {
                    if ((event = eventList[i])) event.listener.call(event.ctx, a1);
                }
            } else if (length === 3) {
                a1 = args[1];
                a2 = args[2];
                i = eventList.length;
                while (i--) {
                    if ((event = eventList[i])) event.listener.call(event.ctx, a1, a2);
                }
            } else if (length === 4) {
                a1 = args[1];
                a2 = args[2];
                a3 = args[3];
                i = eventList.length;
                while (i--) {
                    if ((event = eventList[i])) event.listener.call(event.ctx, a1, a2, a3);
                }
            } else if (length === 5) {
                a1 = args[1];
                a2 = args[2];
                a3 = args[3];
                a4 = args[4];
                i = eventList.length;
                while (i--) {
                    if ((event = eventList[i])) event.listener.call(event.ctx, a1, a2, a3, a4);
                }
            } else {
                arrayShift.apply(args);
                i = eventList.length;
                while (i--) {
                    if ((event = eventList[i])) event.listener.apply(event.ctx, args);
                }
            }
        }

        EventEmitter.prototype.emit = function(type) {
            var eventList = this._events[type];

            if (!eventList || !eventList.length) return this;
            emit(eventList, arguments);

            return this;
        };

        EventEmitter.prototype.listeners = function(type) {
            var eventList = this._events[type];

            return eventList ? eventList.slice() : [];
        };

        EventEmitter.prototype.listenerCount = function(type) {
            var eventList = this._events[type];

            return eventList ? eventList.length : 0;
        };

        EventEmitter.prototype.setMaxListeners = function(value) {
            if ((value = +value) !== value) throw new TypeError("EventEmitter.setMaxListeners(value) value must be a number");

            this._maxListeners = value < 0 ? -1 : value;
            return this;
        };


        EventEmitter.defaultMaxListeners = 10;

        EventEmitter.listeners = function(obj, type) {
            if (obj == null) throw new TypeError("EventEmitter.listeners(obj, type) obj required");
            var eventList = obj._events && obj._events[type];

            return eventList ? eventList.slice() : [];
        };

        EventEmitter.listenerCount = function(obj, type) {
            if (obj == null) throw new TypeError("EventEmitter.listenerCount(obj, type) obj required");
            var eventList = obj._events && obj._events[type];

            return eventList ? eventList.length : 0;
        };

        EventEmitter.setMaxListeners = function(value) {
            if ((value = +value) !== value) throw new TypeError("EventEmitter.setMaxListeners(value) value must be a number");

            EventEmitter.defaultMaxListeners = value < 0 ? -1 : value;
            return value;
        };

        EventEmitter.extend = function(child, parent) {
            if (!parent) parent = this;

            child.prototype = Object.create(parent.prototype);
            child.prototype.constructor = child;
            child._super = parent.prototype;
            child.extend = parent.extend;

            return child;
        };


        function Process() {

            EventEmitter.call(this);

            this.pid = 0;
            this.title = "browser";
            this.env = {};
            this.argv = [];
            this.version = "1.0.0";
            this.versions = {};
            this.config = {};
            this.execPath = ".";
            this.execArgv = [];
            this.arch = global.navigator ? ((/\b(?:AMD|IA|Win|WOW|x86_|x)[32|64]+\b/i.exec(global.navigator.userAgent) || "")[0] || "unknown").replace(/86_/i, "").toLowerCase() : "unknown";
            this.platform = global.navigator ? (global.navigator.platform ? global.navigator.platform.split(/[ \s]+/)[0] : "unknown").toLowerCase() : "unknown";
            this.maxTickDepth = 1000;
            this._cwd = global.location ? global.location.pathname : "/";
        }
        EventEmitter.extend(Process);

        Object.defineProperty(Process.prototype, "browser", {
            get: function() {
                return true;
            }
        });

        Process.prototype.memoryUsage = (function() {
            var performance = global.performance || {},
                memory = {
                    rss: 0,
                    heapTotal: 0,
                    heapUsed: 0
                };

            performance.memory || (performance.memory = {});

            return function memoryUsage() {
                memory.rss = performance.memory.jsHeapSizeLimit || 0;
                memory.heapTotal = performance.memory.totalJSHeapSize || 0;
                memory.heapUsed = performance.memory.usedJSHeapSize || 0;

                return memory;
            };
        }());

        Process.prototype.nextTick = (function() {
            var canSetImmediate = !!global.setImmediate,
                canMutationObserver = !!global.MutationObserver,
                canPost = global.postMessage && global.addEventListener,
                queue;

            if (canSetImmediate) {
                return function(fn) {
                    return global.setImmediate(fn);
                };
            }

            queue = [];

            if (canMutationObserver) {
                var hiddenDiv = document.createElement("div"),
                    observer = new MutationObserver(function hander() {
                        var queueList = queue.slice(),
                            length = queueList.length,
                            i = -1;

                        queue.length = 0;

                        while (++i < length) {
                            queueList[i]();
                        }
                    });

                observer.observe(hiddenDiv, {
                    attributes: true
                });

                return function nextTick(fn) {
                    if (!queue.length) {
                        hiddenDiv.setAttribute("yes", "no");
                    }
                    queue.push(fn);
                };
            }

            if (canPost) {

                global.addEventListener("message", function onNextTick(e) {
                    var source = e.source,
                        fn;

                    if ((source === global || source === null) && e.data === "process-tick") {
                        e.stopPropagation();

                        if (queue.length > 0) {
                            fn = queue.shift();
                            fn();
                        }
                    }
                }, true);

                return function nextTick(fn) {
                    queue.push(fn);
                    global.postMessage("process-tick", "*");
                };
            }

            if (global.setTimeout) {
                return function nextTick(fn) {
                    global.setTimeout(fn, 0);
                };
            }

            return function nextTick(fn) {
                fn();
            };
        }());

        Process.prototype.cwd = function() {
            return this._cwd;
        };

        Process.prototype.chdir = function(dir) {
            var cwd = global.location ? global.location.pathname : "/",
                length, newDir;

            if (dir === "/") {
                newDir = "/";
            } else {
                length = cwd.length;
                newDir = dir.length >= length ? dir : dir.substring(0, length) + "/";
            }

            if (cwd.indexOf(newDir) === 0) {
                this._cwd = dir;
            } else {
                throw new Error("process.chdir can't change to directory " + dir);
            }
        };

        Process.prototype.hrtime = (function() {
            var performance = global.performance || {},
                start;

            Date.now || (Date.now = function now() {
                return (new Date()).getTime();
            });
            start = Date.now();

            performance.now || (performance.now =
                performance.mozNow ||
                performance.msNow ||
                performance.oNow ||
                performance.webkitNow ||
                function now() {
                    return Date.now() - start;
                }
            );

            function performanceNow() {
                return start + performance.now();
            }

            return function hrtime(previousTimestamp) {
                var clocktime = performanceNow() * 1e-3,
                    seconds = Math.floor(clocktime),
                    nanoseconds = (clocktime % 1) * 1e9;

                if (previousTimestamp) {
                    seconds -= previousTimestamp[0];
                    nanoseconds -= previousTimestamp[1];

                    if (nanoseconds < 0) {
                        seconds--;
                        nanoseconds += 1e9;
                    }
                }

                return [seconds, nanoseconds];
            };
        }());

        Process.prototype.uptime = (function() {
            var start = Date.now();

            return function uptime() {
                return ((Date.now() - start) * 1e-3) | 0;
            };
        }());

        Process.prototype.abort = function() {
            throw new Error("process.abort is not supported");
        };

        Process.prototype.binding = function(name) {
            throw new Error("process.binding is not supported");
        };

        Process.prototype.umask = function(mask) {
            throw new Error("process.umask is not supported");
        };

        Process.prototype.kill = function(id, signal) {
            throw new Error("process.kill is not supported");
        };

        Process.prototype.initgroups = function(user, extra_group) {
            throw new Error("process.initgroups is not supported");
        };

        Process.prototype.setgroups = function(groups) {
            throw new Error("process.setgroups is not supported");
        };

        Process.prototype.getgroups = function() {
            throw new Error("process.getgroups is not supported");
        };

        Process.prototype.getuid = function() {
            throw new Error("process.getuid is not supported");
        };

        Process.prototype.setgid = function() {
            throw new Error("process.setgid is not supported");
        };

        Process.prototype.getgid = function() {
            throw new Error("process.getgid is not supported");
        };

        Process.prototype.exit = function() {
            throw new Error("process.exit is not supported");
        };

        Process.prototype.setuid = function(id) {
            throw new Error("process.setuid is not supported");
        };

        Object.defineProperty(Process.prototype, "stderr", {
            get: function() {
                throw new Error("process.stderr is not supported");
            },
            set: function() {
                throw new Error("process.stderr is not supported");
            }
        });

        Object.defineProperty(Process.prototype, "stdin", {
            get: function() {
                throw new Error("process.stderr is not supported");
            },
            set: function() {
                throw new Error("process.stderr is not supported");
            }
        });

        Object.defineProperty(Process.prototype, "stdout", {
            get: function() {
                throw new Error("process.stderr is not supported");
            },
            set: function() {
                throw new Error("process.stderr is not supported");
            }
        });

        return new Process();
    }())),
    typeof(__require__) !== "undefined" ? __require__ : (typeof(require) !== "undefined" ? require : undefined), (new Function("return this;"))()
));
