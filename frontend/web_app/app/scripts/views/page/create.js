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

'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'backbone.stickit',
    'markdown',
    'entities/wiki/page',
    'text!templates/page/edit_page.html',
    'text!templates/main/alert.html'
], function($,_,Marionette,Radio,Stickit,Markdown,
            PageModel,EditPageTemplate,AlertTemplate){

    return Marionette.LayoutView.extend({
        id: "create-page-view",
        template: _.template(EditPageTemplate),
        model: PageModel,
        wikiModel: null,
        homePage: false,
        bindings: {
            '#inputTitle':{
                observe: 'title'
            }
        },
        events: {
            'submit form#editPageForm': 'publishPage',
            'click .page-cancel-button': 'cancelCreate'
        },

        initialize: function(options){
            this.wikiModel = options.wikiModel || null;
            this.homePage = options.homePage || false;
            this.model.on('invalid', this.showError, this);
        },

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

        //Persist page
        publishPage: function(event){
            event.preventDefault();
            var pageContent = this.model.get('content');
            pageContent.raw = $("#wmd-input").val();
            this.model.set('content', pageContent);
            Radio.channel('wikiManager').request('save:page', this.model)
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
        cancelCreate: function(event){
            event.preventDefault();
            parent = this.model.get('parent');
            if(typeof parent !== 'undefined' && parent !== ""){
                Radio.channel('page').trigger('show:page', parent, this.wikiModel);
            } else {
                Radio.channel('wiki').trigger('show:wiki', this.wikiModel.id);
            }
        },

        onRender: function(){
            this.$(".alert").css('display','none');
            if(typeof this.model !== 'undefined'){
                this.stickit();
            }
        },

        onShow: function(){
            var converter = new Markdown.Converter();
            var editor = new Markdown.Editor(converter);
            editor.run();
        },

        onClose: function(){
            this.unstickit();
        }

    });

});
