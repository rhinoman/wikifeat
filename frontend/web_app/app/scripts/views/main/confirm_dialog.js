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