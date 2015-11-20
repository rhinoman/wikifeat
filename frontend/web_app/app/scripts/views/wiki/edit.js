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
 * Created by jcadam on 3/10/15.
 */

'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'backbone.stickit',
    'entities/wiki/wiki',
    'text!templates/wiki/edit_wiki.html'
], function($,_,Marionette,Radio,Stickit,
            WikiModel,EditWikiTemplate){

    return Marionette.ItemView.extend({
        id: "edit-wiki-view",
        template: _.template(EditWikiTemplate),
        model: WikiModel,
        bindings:{
            '#inputName':{
                observe: 'name'
            },
            '#inputDescription':{
                observe: 'description'
            },
            '#inputAllowGuest':{
                observe: 'allowGuest'
            }
        },
        events:{
            'submit form#editWikiForm' : 'submitForm'
        },

        initialize: function(options){
            this.model.on('invalid', this.showError, this);
        },

        showError: function(model, error){
            var theAlert = this.$(".alert");
            theAlert.css('display', 'block');
            theAlert.html('Please correct the following errors: <ul id="error_list"></ul>');
            if(error.hasOwnProperty('name')){
                this.$("#name-input-group").addClass('has-error');
                this.$("#error_list").append("<li>Wiki Name " + error.name + "</li>");
            }
            if(error.hasOwnProperty('description')){
                this.$("#description-input-group").addClass('has-error');
                this.$("#error_list").append("<li>Description " + error.description + "</li>");
            }
        },

        onRender: function(){
            this.$(".alert").css('display', 'none');
            if(typeof this.model !== 'undefined'){
                this.stickit();
                if(typeof this.model.id === 'undefined') {
                    this.$el.prepend('<h1>Create Wiki</h1>');
                } else {
                    this.$el.prepend('<h1>Edit Wiki</h1>');
                }
            }
        },

        submitForm: function(event){
            event.preventDefault();
            this.$("#name-input-group").removeClass('has-error');
            this.$("#description-input-group").removeClass('has-error');
            var newWiki = false;
            if(typeof this.model.id === 'undefined'){
                newWiki = true;
            }
            var wikiPromise = Radio.channel('wikiManager').request('save:wiki', this.model);
            wikiPromise.done(function(wikiModel){
                if(typeof wikiModel === 'undefined'){
                    console.log("Wiki Save Failed.");
                   //TODO: Display an error
                } else {
                    Radio.channel('wiki').trigger('show:wiki', wikiModel.id);
                    if(newWiki === true){
                        Radio.channel('sidebar').trigger('add:wiki', wikiModel);
                    }
                }
            });
        },

        cancelEdit: function(event){
            event.preventDefault();
            window.history.back();
        },

        onClose: function(){
            this.unstickit();
        }

    });

});
