/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 * See the LICENSE.txt file distributed with this work for additional information
 * regarding copyright ownership.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *  Neither the name of Wikifeat nor the names of its contributors may be used
 * to endorse or promote products derived from this software without
 * specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
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
        commonmark: 'vendor/commonmark/dist/commonmark',
        markdown: 'vendor/wikifeat-pagedown/Markdown.Editor',
        'markdown-converter': 'vendor/wikifeat-pagedown/Markdown.Converter',
        markette: 'vendor/markette/js/markette'
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
        markdown: {
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

