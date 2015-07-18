if (process.platform === "win32" || process.platform === "windows") {
    module.exports = require("./file_path_windows");
} else {
    module.exports = require("./file_path");
}
