var compass = require('../lib/compass'),
    chai = require('chai'),
    express = require('express'),
    request = require('supertest'),
    fs = require('fs'),
    path = require('path');

var should = require('chai').should();

const FIXTURES_DIR = path.resolve(path.join(__dirname, 'fixtures'));

describe('Compass', function() {
    afterEach(function() {
        try {
            fs.unlinkSync(path.join(FIXTURES_DIR, 'css', 'style.css'));
            fs.rmdirSync(path.join(FIXTURES_DIR, 'css'));
        }
        catch (e) { /*silent*/ }
        process.chdir(__dirname);
    });

    describe('#compile', function() {
        it('should compile in the current directory', function(done) {
            process.chdir(FIXTURES_DIR);
            compass.compile(function() {
                var stats = fs.statSync(path.join(FIXTURES_DIR, 'css', 'style.css'));
                stats.size.should.be.above(0);
                done();
            });
        });

        it('should compile in the given options.cwd directory', function(done) {
            compass.compile({ cwd: FIXTURES_DIR }, function() {
                var stats = fs.statSync(path.join(FIXTURES_DIR, 'css', 'style.css'));
                stats.size.should.be.above(0);
                done();
            });
        });
    });

    describe('middleware', function() {
        it('should serve a css file from the root directory', function(done) {
            var app = express()
                .use(compass({ cwd: FIXTURES_DIR }))
                .use(express.static(FIXTURES_DIR))
                .listen(1337);

            request(app)
                .get('/css/style.css')
                .expect(200, function(err, res) {
                    res.text.length.should.be.above(0);
                    done();
                });
        });
    });
});