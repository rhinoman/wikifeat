var exec = require('child_process').exec;

/**
 * express middleware for serving compiled on-the-fly sass/scss files.
 *
 * @type {Function}
 */
var compass = module.exports = function(options) {
    return function(req, res, next) {
        compass.compile(options, function() {
            return next();
        });
    };
};

/**
 * compiles sass/scss files in the given directory
 *
 * @param {Object} options
 * @param {String} options.root cwd Current working directory for compass.
 * By default it take the program cwd.
 *
 * @param {Function} callback
 */
compass.compile = function(options, callback) {
    if ('function' == typeof options) {
        callback = options;
    }

    options = options || {};
    options.cwd = options.cwd || process.cwd();

    exec('compass compile', { cwd: options.cwd }, callback);
}
