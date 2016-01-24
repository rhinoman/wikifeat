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
    'backbone.radio',
    'text!templates/page/insert_link_dialog.html'
], function($, _, Marionette, Radio,
            InsertLinkDialogTemplate){

    var userChannel = Radio.channel('userManager');
    var wikiChannel = Radio.channel('wikiManager');

    return Marionette.ItemView.extend({
        id: "insertLinkDialog",
        template: _.template(InsertLinkDialogTemplate),
        events: {
            'click #insertButton': function(){$('#theSubmit').trigger('click')},
            'change input[type=radio][name=linkOption]': 'radioChange',
            'change select#wikiSelect': 'wikiSelected',
            'submit form': 'submitForm'
        },

        initialize: function(options){
            options = options || {};
            this.callback = options.callback || function(){};
            this.currentUser = userChannel.request('get:currentUser');
            this.wikiList = $.Deferred();
            var self = this;
            this.currentUser.done(function(user){
                if(typeof user !== 'undefined'){
                    wikiChannel.request('get:memberWikiList', user).done(function(data){
                        self.wikiList.resolve(data);
                    });
                }
            })
        },

        optionTemplate: function(){
            return _.template("<option value='<%= id %>'><%= name %></option>");
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
                this.populateInternalFields();
            } else if(linkMode === 'external'){
                this.$("div#internalLinkSelectContainer").hide();
                this.$("div#externalLinkSelectContainer").show();
            }
        },

        populateInternalFields: function(){
            var self = this;
            this.wikiList.done(function(data){
                var select = self.$("select#wikiSelect");
                select.html('<option value="0">Select Wiki...</option>')
                if(data !== 'undefined'){
                    _.each(data.models, function(wiki){
                        select.append(self.optionTemplate()({id: wiki.get('id'), name: wiki.get('name')}));
                    }, self);
                }
            });
        },

        wikiSelected: function(event){
            var self = this;
            var wikiSelect = event.currentTarget;
            var wikiId = $(wikiSelect).val();
            this.$("select#pageSelect").empty();
            wikiChannel.request('get:allPageList', wikiId).done(function(data){
                var select = self.$("select#pageSelect");
                if(data !== 'undefined'){
                    _.each(data.models, function(page){
                        select.append(self.optionTemplate()({id: page.get('id'), name: page.get('title')}));
                    })
                }
            });
        },

        submitForm: function(event){
            event.preventDefault();
            var linkMode = this.linkMode();
            var theUrl = "http://";
            if(linkMode === 'internal'){
                var wikiId = this.$("select#wikiSelect").val();
                var pageId = this.$("select#pageSelect").val();
                if(wikiId !== "0" && wikiId !== null && pageId !== "0" && pageId !== null) {
                    theUrl = "/wikis/" + wikiId + "/" + pageId;
                }
            } else if(linkMode === 'external'){
                theUrl = this.$("#externalUrlField").val();
            }
            this.callback(theUrl);
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
