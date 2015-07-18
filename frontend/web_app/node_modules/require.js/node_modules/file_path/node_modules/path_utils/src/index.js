var pathUtils = module.exports;


pathUtils.removeEmpties = function(parts) {
    var i = parts.length;

    while (i--) {
        if (!parts[i]) parts.splice(i, 1);
    }

    return parts;
};

pathUtils.trim = function(parts) {
    var length = parts.length,
        start = -1,
        end = length - 1;

    while (start++ < end) {
        if (parts[start] !== "") break;
    }

    end = length;
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
        } else if (up !== 0) {
            parts.splice(i, 1);
            up--;
        }
    }

    if (allowAboveRoot) {
        while (up--) parts.unshift("..");
    }

    return parts;
};
