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
    'markette',
    'entities/wiki/comment',
    'views/main/alert',
    'views/page/edit/insert_link_dialog',
    'views/page/edit/insert_image_dialog',
    'text!templates/page/edit_comment.html'
], function($,_,Marionette,Radio,Markette,
            CommentModel,AlertView,InsertLinkDialog,
            InsertImageDialog,EditCommentTemplate){

    return Markette.EditorView.extend({
        model: CommentModel,
        template: _.template(EditCommentTemplate),
        id: "edit-comment-view",

        events:{
            'click #cancelEditButton': 'cancelButtonClick',
            'submit form#editCommentForm' : 'publishComment'
        },

        //save the comment
        publishComment: function(event){
            event.preventDefault();
            var newComment = false;
            var pageContent = _.clone(this.model.get('content'));
            pageContent.raw = this.$("textarea#marketteInput").val();
            this.model.set('content', pageContent);
            if(typeof this.model.id === 'undefined' || this.model.id === null){
                newComment = true;
            }
            var self=this;
            Radio.channel('wikiManager').request('save:comment', this.model)
                .done(function(response){
                    if(typeof response !== 'undefined' && newComment) {
                        Radio.channel('commentsView').trigger('new:comment', response);
                    }
                    self.destroy();
                })
        },

        //simply destroy the view on cancel
        cancelButtonClick: function(event){
            this.destroy();
        },

        onRender: function(){
            var content = this.model.get("content");
            if(content.raw !== ""){
                this.$("textarea#marketteInput").val(content.raw);
            }
            Markette.EditorView.prototype.onRender.call(this);
            this.$('[data-toggle="tooltip"]').tooltip({container: this.el});
        },

        onClose: function(){
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
        },

        //Override Image functionality
        doImage: function(event){
            var self = this;
            var imd = new InsertImageDialog({
                callback: function(url){
                    self.$('textarea#marketteInput').focus();
                    self.doInlineMarkup({
                        before: '![',
                        after: '](' + url + ')'
                    });
                },
                wikiId: this.model.wikiId
            });
            Radio.channel('main').trigger('show:dialog', imd);
        }

    });

});