/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.
 *
 * This file is part of WikiFeat.
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
