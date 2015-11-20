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
/**
 * Wikifeat Javascript API
 */

(function (root, factory){

    if (typeof define === 'function' && define.amd) {
        // AMD
        define(['backbone.radio'], function(Radio){
            return (root.Wikifeat = factory(Radio));
        });
    } else if (typeof exports === 'object'){
        // CommonJS
        var Radio = require('backbone.radio');
        module.exports = factory(Radio);
    } else {
        // Browser Global
        root.Wikifeat = factory(root.Radio);
    }
}(this, function (Radio){
    'use strict';

    var Wikifeat = {};
    Wikifeat.version = 0.1;

    Wikifeat.printVersion = function(){
        console.log(this.version);
    };

    // Displays a view in the main content area.
    // Must be a Backbone or Marionette View
    Wikifeat.showContent = function(view) {
        Radio.channel('main').trigger('show:content', view);
    };

    // Places a View in the dialog region (useful for bootstrap modal dialogs).
    Wikifeat.showDialog = function(view) {
        Radio.channel('main').trigger('show:dialog', view);
    };

    // Adds a menu item to the main/sidebar menu
		// Also must be a Backbone/Marionette view
    Wikifeat.addMenuItem = function (name, menuView) {
        Radio.channel('sidebar').trigger('add:menu', name, menuView);
    };

    // Shows an error dialog
    Wikifeat.showErrorDialog = function (errName, errMsg) {
        //Not yet implemented
    };

    return Wikifeat;
}));
