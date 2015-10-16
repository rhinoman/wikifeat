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
    'entities/wiki/page',
    'views/page/edit_form',
    'views/page/edit_preview',
    'text!templates/page/edit_page.html'
], function($,_,Marionette,PageModel,EditFormView,
            EditPreview,EditPageTemplate){

   return Marionette.LayoutView.extend({
       id: "edit-page-view",
       template: _.template(EditPageTemplate),
       model: PageModel,
       wikiModel: null,
       regions: {
           editorContentRegion: "#editorContent"
       },
       events: {
           'click #editLink': 'showEditForm',
           'click #previewLink': 'showPreview'
       },
       wipText: new Backbone.Model({
           data: null
       }),

       initialize: function(options){
           if(options.hasOwnProperty('wikiModel')){
               this.wikiModel = options.wikiModel;
           }
           this.model.on('invalid', this.showError, this);
           this.wipText.set("data", this.model.get("content").raw)
       },

       showEditForm: function(event){
           event.preventDefault();
           this.editorContentRegion.show(new EditFormView({
               model: this.model,
               wikiModel: this.wikiModel,
               wipText: this.wipText
           }));
           this.$("#previewLink").parent("li").removeClass("active");
           this.$("#editLink").parent("li").addClass("active");
       },

       showPreview: function(event){
           event.preventDefault();
           this.editorContentRegion.show(new EditPreview({
               model: this.model,
               wipText: this.wipText
           }));
           this.$("#editLink").parent("li").removeClass("active");
           this.$("#previewLink").parent("li").addClass("active");
       },
       onShow: function(){
           this.$("#editLink").trigger('click');
       }
   });
});
