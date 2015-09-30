/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.*
 *
 * This file is part of WikiFeat
 *
 *     WikiFeat is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 *     WikiFeat is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 * along with WikiFeat.  If not, see <http://www.gnu.org/licenses/>.
 */

require.config({
    paths: {
        bootstrap: 'vendor/bootstrap/dist/js/bootstrap',
        marionette: 'vendor/marionette/lib/core/backbone.marionette',
        requirejs: 'vendor/requirejs/require',
        backbone: 'vendor/backbone/backbone',
        jquery: 'vendor/jquery/dist/jquery',
        underscore: 'vendor/underscore/underscore',
        text: 'vendor/requirejs-text/text',
        'backbone.babysitter': 'vendor/backbone.babysitter/lib/backbone.babysitter',
        'backbone.radio': 'vendor/backbone.radio/build/backbone.radio',
        'backbone.basicauth': 'vendor/backbone.basicauth/backbone.basicauth',
        'requirejs-text': 'vendor/requirejs-text/text',
        'jquery-cookie': 'vendor/jquery-cookie/jquery.cookie',
        'backbone.stickit': 'vendor/backbone.stickit/backbone.stickit',
        'backbone.paginator': 'vendor/backbone.paginator/lib/backbone.paginator',
        moment: 'vendor/moment/moment',
        wikifeat: 'vendor/wikifeat-api/wikifeat-api',
        'markdown-converter': 'vendor/wikifeat-pagedown/Markdown.Converter',
        'markdown': 'vendor/wikifeat-pagedown/Markdown.Editor'
    },
    shim: {
        bootstrap: {
            deps: [
                'jquery'
            ],
            exports: '$.fn.popover'
        },
        'markdown-converter': {
            exports: 'Markdown'
        },
        'markdown': {
            deps: [
                'markdown-converter'
            ],
            exports: 'Markdown'
        }
    },
    map: {
        '*': {
            'backbone.wreqr': 'backbone.radio'
        }
    },
    packages: [

    ]
});

define(['app'], function(App){
    console.log('starting');
    App.start();
});

