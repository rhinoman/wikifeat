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
