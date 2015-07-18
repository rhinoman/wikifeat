var pathUtils = require("path_utils");


var filePath = module.exports,

    slice = Array.prototype.slice,

    IS_ABSOLUTE = /^(?:[A-Za-z]:)?\\/,
    SPLIT_TAIL = /[\\\/]+/,
    ENDING_SLASH = /[\\\/]$/,
    SPLIT_DEVICE = /^([a-zA-Z]:|[\\\/]{2}[^\\\/]+[\\\/]+[^\\\/]+)?([\\\/])?([\s\S]*?)$/,


    UNC_ROOT_START = /^[\\\/]+/,
    UNC_ROOT_GLOBAL = /[\\\/]+/g,

    JOIN_ROOT = /^[\\\/]{2}[^\\\/]/,
    JOIN_REPLACER = /^[\\\/]{2,}/;


Object.defineProperty(filePath, "separator", {
    enumerable: false,
    configurable: false,
    writable: false,
    value: "\\"
});
Object.defineProperty(filePath, "delimiter", {
    enumerable: false,
    configurable: false,
    writable: false,
    value: ";"
});

filePath.isAbsolute = function(path) {
    return IS_ABSOLUTE.test(path);
};

filePath.normalize = function(path) {
    var result = SPLIT_DEVICE.exec(path),
        device = result[1] || "",
        isUnc = device && device.charAt(1) !== ":",
        isAbsolute = filePath.isAbsolute(path),
        tail = result[3],
        trailingSlash = ENDING_SLASH.test(tail);

    if (device && device.charAt(1) === ":") {
        device = device[0].toLowerCase() + device.substr(1);
    }

    tail = pathUtils.normalizeArray(pathUtils.removeEmpties(tail.split(SPLIT_TAIL)), !isAbsolute).join("\\");

    if (!tail && !isAbsolute) {
        tail = ".";
    }
    if (tail && trailingSlash) {
        tail += "\\";
    }

    if (isUnc) {
        device = normalizeUNCRoot(device);
    }

    return device + (isAbsolute ? "\\" : "") + tail;
};

filePath.resolve = function() {
    var resolvedDevice = "",
        resolvedTail = "",
        resolvedAbsolute = false,
        result, device, isUnc, isAbsolute, tail, path,
        i;

    for (i = arguments.length - 1; i >= -1; i--) {
        if (i >= 0) {
            path = arguments[i];
        } else if (!resolvedDevice) {
            path = ".";
        } else {
            path = process.env["=" + resolvedDevice];
            if (!path || path.substr(0, 3).toLowerCase() !== resolvedDevice.toLowerCase() + "\\") {
                path = resolvedDevice + "\\";
            }
        }

        if (typeof(path) !== "string") {
            throw new TypeError("Arguments to path.resolve must be strings");
        } else if (!path) {
            continue;
        }

        result = SPLIT_DEVICE.exec(path);
        device = result[1] || "";
        isUnc = device && device.charAt(1) !== ":";
        isAbsolute = filePath.isAbsolute(path);
        tail = result[3];

        if (device && resolvedDevice && device.toLowerCase() !== resolvedDevice.toLowerCase()) {
            continue;
        }

        if (!resolvedDevice) {
            resolvedDevice = device;
        }
        if (!resolvedAbsolute) {
            resolvedTail = tail + "\\" + resolvedTail;
            resolvedAbsolute = isAbsolute;
        }

        if (resolvedDevice && resolvedAbsolute) {
            break;
        }
    }

    if (isUnc) {
        resolvedDevice = normalizeUNCRoot(resolvedDevice);
    }
    resolvedTail = pathUtils.normalizeArray(pathUtils.removeEmpties(resolvedTail.split(SPLIT_TAIL)), !resolvedAbsolute).join("\\");

    return (resolvedDevice + (resolvedAbsolute ? "\\" : "") + resolvedTail) || ".";
};

filePath.relative = function(from, to) {
    from = filePath.resolve(from);
    to = filePath.resolve(to);

    var lowerFrom = from.toLowerCase(),
        lowerTo = to.toLowerCase(),

        toParts = trim(to.split("\\")),

        lowerFromParts = pathUtils.trim(lowerFrom.split("\\")),
        lowerToParts = pathUtils.trim(lowerTo.split("\\")),

        length = Math.min(lowerFromParts.length, lowerToParts.length),
        samePartsLength = length,
        outputParts,
        i;

    for (i = 0; i < length; i++) {
        if (lowerFromParts[i] !== lowerToParts[i]) {
            samePartsLength = i;
            break;
        }
    }

    if (samePartsLength === 0) {
        return to;
    }

    outputParts = [];
    for (i = samePartsLength; i < lowerFromParts.length; i++) {
        outputParts.push("..");
    }

    outputParts = outputParts.concat(toParts.slice(samePartsLength));

    return outputParts.join("\\");
};

filePath.join = function() {
    var paths = slice.call(arguments),
        i = paths.length,
        joined;

    while (i--) {
        if (typeof(paths[i]) !== "string") {
            throw new TypeError("Arguments to join must be strings");
        }
    }

    joined = paths.join("\\");
    if (!JOIN_ROOT.test(paths[0])) {
        joined = joined.replace(JOIN_REPLACER, "\\");
    }

    return filePath.normalize(joined);
};

filePath.dir = function(path) {
    path = path.substring(0, path.lastIndexOf("\\") + 1);
    return path ? path.substr(0, path.length - 1) : ".";
};

filePath.dirname = filePath.dir;

filePath.base = function(path, ext) {
    path = path.substring(path.lastIndexOf("\\") + 1);

    if (ext && path.substr(-ext.length) === ext) {
        path = path.substr(0, path.length - ext.length);
    }

    return path || "";
};

filePath.basename = filePath.base;

filePath.ext = function(path) {
    var index = path.lastIndexOf(".");
    return index > -1 ? path.substring(index) : "";
};

filePath.extname = filePath.ext;

function normalizeUNCRoot(device) {
    return "\\\\" + device.replace(UNC_ROOT_START, "").replace(UNC_ROOT_GLOBAL, "\\");
}
