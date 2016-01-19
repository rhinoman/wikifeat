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
    'markette',
    'entities/wiki/page',
    'views/page/edit/edit_form',
    'text!templates/page/edit_page.html'
], function($,_,Marionette,Markette,PageModel,
            EditFormView,EditPageTemplate){

   return Marionette.LayoutView.extend({
       id: "edit-page-view",
       template: _.template(EditPageTemplate),
       model: PageModel,
       wikiModel: null,
       homePage: false,
       regions: {
           editorContentRegion: "#editorContent",
           previewContentRegion: "#previewContent"
       },
       events: {
           'click #editLink': 'showEditForm',
           'click #previewLink': 'showPreview'
       },

       initialize: function(options){
           if(options.hasOwnProperty('wikiModel')){
               this.wikiModel = options.wikiModel;
           }
           if(options.hasOwnProperty('homePage')){
               this.homePage = options.homePage;
           }
           this.model.on('invalid', this.showError, this);
           this.editFormView = new EditFormView({
               model: this.model,
               wikiModel: this.wikiModel,
               homePage: this.homePage
           });
           this.editPreview = new Markette.Preview();
       },

       showEditForm: function(event){
           event.preventDefault();
           this.$("div#previewContent").hide();
           this.$("div#editorContent").show();
           this.$("#previewLink").parent("li").removeClass("active");
           this.$("#editLink").parent("li").addClass("active");
       },

       showPreview: function(event){
           event.preventDefault();
           this.$("div#editorContent").hide();
           this.$("div#previewContent").show();
           this.$("#editLink").parent("li").removeClass("active");
           this.$("#previewLink").parent("li").addClass("active");
           var mdText = this.editFormView.getText();
           this.editPreview.renderPreview(mdText);
       },
       onRender: function(){
           this.$("div#previewContent").hide();
       },
       onShow: function(){
           this.editorContentRegion.show(this.editFormView);
           this.previewContentRegion.show(this.editPreview);
       },
   });
});
