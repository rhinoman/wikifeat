var pathUtils = require("path_utils");


var filePath = module.exports;


Object.defineProperty(filePath, "separator", {
    enumerable: false,
    configurable: false,
    writable: false,
    value: "/"
});
Object.defineProperty(filePath, "delimiter", {
    enumerable: false,
    configurable: false,
    writable: false,
    value: ":"
});


filePath.isAbsolute = function(path) {
    return path[0] === "/";
};

filePath.normalize = function(str) {
    var isAbs = filePath.isAbsolute(str),
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

filePath.resolve = function() {
    var resolvedPath = "",
        resolvedAbsolute = false,
        i, str;

    for (i = arguments.length - 1; i >= -1 && !resolvedAbsolute; i--) {
        str = (i >= 0) ? arguments[i] : process.cwd();

        if (typeof(str) !== "string") {
            throw new TypeError("Arguments to path.resolve must be strings");
        } else if (!str) {
            continue;
        }

        resolvedPath = str + "/" + resolvedPath;
        resolvedAbsolute = str[0] === "/";
    }

    resolvedPath = pathUtils.normalizeArray(pathUtils.removeEmpties(resolvedPath.split("/")), !resolvedAbsolute).join("/");
    return ((resolvedAbsolute ? "/" : "") + resolvedPath) || ".";
};

filePath.relative = function(from, to) {
    from = filePath.resolve(from).substr(1);
    to = filePath.resolve(to).substr(1);

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

filePath.join = function() {
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

    return filePath.normalize(str);
};

filePath.dir = function(str) {
    str = str.substring(0, str.lastIndexOf("/") + 1);
    return str ? str.substr(0, str.length - 1) : ".";
};

filePath.dirname = filePath.dir;

filePath.base = function(str, ext) {
    str = str.substring(str.lastIndexOf("/") + 1);

    if (ext && str.substr(-ext.length) === ext) {
        str = str.substr(0, str.length - ext.length);
    }

    return str || "";
};

filePath.basename = filePath.base;

filePath.ext = function(str) {
    var index = str.lastIndexOf(".");
    return index > -1 ? str.substring(index) : "";
};

filePath.extname = filePath.ext;
