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
    'backbone.radio',
    'backbone.stickit',
    'bootstrap',
    'markette',
    'entities/wiki/page',
    'views/main/alert',
    'views/page/edit/insert_link_dialog',
    'text!templates/page/edit_page_form.html'
], function($,_,Marionette,Radio,Stickit,Bootstrap,Markette,
            PageModel,AlertView,InsertLinkDialog,
            EditPageFormTemplate){

    return Markette.EditorView.extend({
        model: PageModel,
        template: _.template(EditPageFormTemplate),
        wikiModel: null,
        homePage: false,
        bindings: {
            '#inputTitle': {
                observe: 'title'
            },
            '#inputDisableComments':{
                observe: 'commentsDisabled'
            }
        },
        events: {
            'submit form#editPageForm': 'publishChanges',
            'click .page-cancel-button': 'cancelEdit',
            'change textarea#marketteInput': 'updateWipText'
        },
        regions: {
            "editorRegion" : "#editorContainer"
        },

        initialize: function(options){
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
            }
            if(options.hasOwnProperty('homePage')){
                this.homePage = options.homePage;
            }
            this.model.on('invalid', this.showError, this);
            Markette.EditorView.prototype.initialize.call(this,options);
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

        //Save page edits to server
        publishChanges: function(event){
            event.preventDefault();
            this.$("#title-input-group").removeClass('has-error');
            var pageContent = _.clone(this.model.get('content'));
            pageContent.raw = $("textarea#marketteInput").val();
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
                this.setText(this.model.get('content').raw);
            }
            Markette.EditorView.prototype.onRender.call(this);
            this.$('[data-toggle="tooltip"]').tooltip();
        },

        onClose: function(){
            this.unstickit();
            Markette.EditorView.prototype.onClose.call(this);
        },

        //Override Link functionality
        doLink: function(event){
            var self = this;
            var ild = new InsertLinkDialog({
                callback: function(url){
                    self.$('textarea#marketteInput').focus();
                    self.doInlineMarkup({
                        before: '[',
                        after: '](' + url + ')'
                    });
                }
            });
            Radio.channel('main').trigger('show:dialog', ild);
        }


    });


});