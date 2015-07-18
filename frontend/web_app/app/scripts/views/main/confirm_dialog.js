/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.
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

'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'bootstrap',
    'text!templates/main/confirm_dialog.html'
], function($,_,Marionette,Boostrap,ConfirmDialogTemplate){

    return Marionette.ItemView.extend({
        id: 'confirm-dialog',
        template: _.template(ConfirmDialogTemplate),
        message: "Please confirm this request.",
        confirmCallback: function(){console.log('confirmed')},

        events: {
            'click #confirmButton': 'confirmClick'
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('message')){
                this.message = options.message;
            }
            if(options.hasOwnProperty('confirmCallback')){
                this.confirmCallback = options.confirmCallback;
            }
        },

        onRender: function(){
            this.$("#confirmMessage").html(this.message);
            this.$("#confirmModal").modal();
        },

        confirmClick: function(event){
            event.preventDefault();
            this.confirmCallback();
        },

        onClose: function(){}
    });

});