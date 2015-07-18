# compass

[Compass](http://compass-style.org) wrapper and middleware for node.js.

## Install

Ensure **compass** is already installed ([install compass](http://compass-style.org/install)).

Then just use **npm**:

```javascript
npm install compass
```

## Usage

### Middelware

Ensure **compass** comes before a static middleware. **compass** generates your css files but **does not** serve them. You need a static middleware for that.

```javascript
var compass = require('compass'),
    express = require('express')();

app = express();
app.use(compass({ cwd: __dirname + 'public' }));
app.use(express.static(__dirname + 'public'));
```

### Manually

For now **compass** only wraps the compile action.

```javascript
var compass = require('compass');

// compiles in process.cwd()
compass.compile(function(err, stdout, stderr) {
    console.log('done');
});

// compiles in the given directory
compass.compile({ cwd: __dirname + 'public' }, function(err, stdout, stderr) {
    console.log('done');
});
```

## Options

For now **compass** only offer the `cwd` option.

If you want to customize stuff, please use a `config.rb` ([http://compass-style.org/help/tutorials/configuration-reference](configuration reference))
