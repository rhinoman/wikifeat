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
    'backbone.radio',
    'backbone.stickit',
    'markdown',
    'entities/wiki/page',
    'views/main/alert',
    'text!templates/page/edit_page_form.html'
], function($,_,Marionette,Radio,Stickit,Markdown,
            PageModel,AlertView,EditPageFormTemplate){

    return Marionette.ItemView.extend({
        model: PageModel,
        template: _.template(EditPageFormTemplate),
        wikiModel: null,
        wipText: null,
        homePage: false,
        bindings: {
            '#inputTitle': {
                observe: 'title'
            },
            '#inputDisableComments':{
                observe: 'comments_disabled'
            }
        },
        events: {
            'submit form#editPageForm': 'publishChanges',
            'click .page-cancel-button': 'cancelEdit',
            'change #wmd-input': 'updateWipText'
        },

        initialize: function(options){
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
            }
            if(options.hasOwnProperty('wipText')){
                this.wipText = options.wipText;
            }
            if(options.hasOwnProperty('homePage')){
                this.homePage = options.homePage;
            }
            this.model.on('invalid', this.showError, this);
        },

        // Shows an error message
        showError: function(model, error){
            var alertText = 'Please correct the following errors: <ul id="error_list">';
            if(error.hasOwnProperty('title')){
                this.$("#title-input-group").addClass('has-error');
                alertText += "<li>Page Title " + error.title + "</li>"
            }
            alertText += "</ul>";
            var alertView = new AlertView({
                el: $("#alertBox"),
                alertType: 'alert-danger',
                alertMessage: alertText
            });
            alertView.render();
        },

        updateWipText: function(event){
            if (this.wipText !== null) {
                this.wipText.set('data', this.$("#wmd-input").val());
            }
        },

        //Save page edits to server
        publishChanges: function(event){
            event.preventDefault();
            this.$("#title-input-group").removeClass('has-error');
            var pageContent = _.clone(this.model.get('content'));
            pageContent.raw = $("#wmd-input").val();
            this.model.set('content', pageContent);

            var self=this;
            Radio.channel('page').request('save:page', this.model)
                .done(this.afterSave.bind(this));
        },

        //After page save callback
        afterSave: function(response){
            if(typeof response !== 'undefined'){
                if(this.homePage === true){
                    this.wikiModel.set("homePageId", this.model.id);
                    Radio.channel('wikiManager').request('save:wiki', this.wikiModel)
                        .done(this.afterWikiSave.bind(this))
                } else {
                    Radio.channel('page').trigger('show:page',
                        this.model.id, this.wikiModel);
                }

            } else {
                //TODO: Handle undefined
            }
        },

        //Callback after the wiki is saved
        afterWikiSave: function(response){
            if(typeof response !== 'undefined'){
                Radio.channel('page').trigger('show:page',
                    this.model.id, this.wikiModel);
            } else {
                //TODO: Handle undefined
            }
        },

        //Cancel and go back!
        cancelEdit: function(event){
            event.preventDefault();
            if(typeof this.model.id !== "undefined") {
                Radio.channel('page').trigger('show:page', this.model.id, this.wikiModel);
            } else {
                var parent = this.model.get('parent');
                if(typeof parent !== 'undefined' && parent !== ""){
                    Radio.channel('page').trigger('show:page', parent, this.wikiModel);
                } else {
                    Radio.channel('wiki').trigger('show:wiki', this.wikiModel.id);
                }
            }
        },

        onRender: function(){
            this.$(".alert").css('display','none');
            if(typeof this.model !== 'undefined'){
                this.stickit();
            }
            if(this.wipText !== null) {
                this.$("#wmd-input").html(this.wipText.get('data'));
            }
        },

        onShow: function(){
            var editor = new Markdown.Editor();
            editor.run();
        },

        onClose: function(){
            this.unstickit();
        }


    });


});