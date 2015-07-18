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
    from = urlPath.resolve(from).substr(1);
    to = urlPath.resolve(to).substr(1);

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
