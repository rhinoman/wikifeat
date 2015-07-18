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
    'text!templates/page/edit_page.html'
], function($,_,Marionette,Radio,Stickit,Markdown,
            PageModel,AlertView,EditPageTemplate){

   return Marionette.ItemView.extend({
       id: "edit-page-view",
       template: _.template(EditPageTemplate),
       model: PageModel,
       wikiModel: null,
       bindings: {
           '#wmd-input': {
               observe: 'content',
               onGet: function(dataContent){
                   return dataContent.raw;
               },
               updateModel: false
           },
           '#inputTitle': {
               observe: 'title'
           }
       },
       events: {
           'submit form#editPageForm': 'publishChanges',
           'click .page-cancel-button': 'cancelEdit'
       },

       initialize: function(options){
           if(options.hasOwnProperty('wikiModel')){
               this.wikiModel = options.wikiModel;
           }
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

       //Save page edits to server
       publishChanges: function(event){
           event.preventDefault();
           this.$("#title-input-group").removeClass('has-error');
           var pageContent = this.model.get('content');
           pageContent.raw = $("#wmd-input").val();
           this.model.set('content', pageContent);

           var self=this;
           Radio.channel('page').request('save:page', this.model)
               .done(function(response){
                   //TODO: Check for undefined
                   Radio.channel('page').
                       trigger('show:page', self.model.id, self.wikiModel);
               });
       },

       //Cancel and go back!
       cancelEdit: function(event){
           event.preventDefault();
           Radio.channel('page').trigger('show:page', this.model.id, this.wikiModel);
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
