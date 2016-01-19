/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 * See the LICENSE.txt file distributed with this work for additional information
 * regarding copyright ownership.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *    * Redistributions of source code must retain the above copyright notice,
 *        this list of conditions and the following disclaimer.
 *    * Redistributions in binary form must reproduce the above copyright
 *        notice, this list of conditions and the following disclaimer in the
 *        documentation and/or other materials provided with the distribution.
 *    * Neither the name of Wikifeat nor the names of its contributors may be used
 *        to endorse or promote products derived from this software without
 *        specific prior written permission.
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
    'text!templates/page/insert_link_dialog.html'
], function($, _, Marionette,
            InsertLinkDialogTemplate){

    return Marionette.ItemView.extend({
        id: "insertLinkDialog",
        template: _.template(InsertLinkDialogTemplate),

        events: {
            'shown.bs.modal' : 'showModal',
            'click #insertButton': function(){$('#theSubmit').trigger('click')},
            'change input[type=radio][name=linkOption]': 'radioChange',
            'submit form': 'submitForm'
        },

        initialize: function(options){
            options = options || {};
            this.callback = options.callback || function(){}
        },

        onRender: function(){
            this.$("#insertLinkModal").modal();
        },

        radioChange: function(event){
            //Get the checked radio button
            var linkMode = this.linkMode();
            if(linkMode === 'internal'){
                this.$("div#externalLinkSelectContainer").hide();
                this.$("div#internalLinkSelectContainer").show();
            } else if(linkMode === 'external'){
                this.$("div#internalLinkSelectContainer").hide();
                this.$("div#externalLinkSelectContainer").show();
            }
        },

        submitForm: function(event){
            event.preventDefault();
            var linkMode = this.linkMode();
            if(linkMode === 'internal'){

            } else if(linkMode === 'external'){
                var theUrl = this.$("#externalUrlField").val();
                this.callback(theUrl);
            }
            this.$('#insertLinkModal').modal('hide');
        },

        linkMode: function(){
            var checked = this.$('input[type=radio]:checked');
            if(checked.attr('id') === 'internalOption'){
                return "internal";
            } else if(checked.attr('id') === 'externalOption'){
                return "external";
            }
        }
    });

});
