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
    'moment',
    'backbone.radio',
    'entities/wiki/comment',
    'views/main/confirm_dialog',
    'views/page/comment/edit_comment',
    'util/common_events',
    'text!templates/page/show_comment.html'
], function($,_,Marionette,Moment,Radio,
            CommentModel,ConfirmDialog,
            EditCommentView,CommonEvents,
            ShowCommentTemplate){

    return Marionette.ItemView.extend({
        className: "comment-view col-lg-12",
        template: _.template(ShowCommentTemplate),
        model: CommentModel,
        editMode: false,

        events: {
            'click #deleteCommentButton' : 'deleteComment',
            'click #editCommentButton' : 'editComment',
            'click .comment-content a' : 'handleLinkClick'
        },

        initialize: function(options){
            this.model.listenTo('change', this.render, this);
            this.rm = new Marionette.RegionManager();
            this.rm.addRegion("editor", "#editorContainer_" + this.model.id);
        },

        templateHelpers: function(){
            var self = this;
            return {
                getIdSuffix: function(){
                    return self.model.id;
                }
            }
        },

        editComment: function(event){
            var ecv = new EditCommentView({model: this.model});
            this.$("#commentContent").css("display", "none");
            var editorRegion = this.rm.get("editor");
            editorRegion.show(ecv);
            var self = this;
            ecv.on('destroy', function(){
                self.$("#commentContent").css("display", "block");
                editorRegion.reset();
            });
        },

        deleteComment: function(event){
            var self = this;
            var confirmCallback = function(){
                Radio.channel('wikiManager').request('delete:comment', self.model)
                    .done(function(response){
                        if(typeof response === 'undefined'){
                            //TODO: Error display
                        }
                    });
            };

            var confirmDialog = new ConfirmDialog({
                message: 'Are you sure you wish to delete this comment? ' +
                'This action is irreversible.',
                confirmCallback: confirmCallback
            });

            Radio.channel('main')
                .trigger('show:dialog', confirmDialog);
        },

        handleLinkClick: function(event){
            CommonEvents.handleLinkClick(event);
        },

        onRender: function(){
            if(typeof this.model !== "undefined") {
                var author = this.model.get("author");
                this.$("#commentAuthorName").html(author);
                var timestamp = this.model.get("createdTime");
                var time = Moment(timestamp);
                var timestring = time.format("HH:mm on D MMM YYYY");
                this.$("#commentDatetime").html(timestring);
                var content = this.model.get("content");
                this.$("#commentContent").html(content.formatted);
                //hide buttons based on access
                if(!this.model.isEditable){
                    this.$("#editCommentButton").css("display","none");
                }
                if(!this.model.isDeletable){
                    this.$("#deleteCommentButton").css("display","none");
                }
                //Draw the comment author's avatar
                var self = this;
                Radio.channel('userManager').request('get:user', author)
                    .done(function(authorUser){
                        if(typeof authorUser !== 'undefined'){
                            self.$("#authorAvatarThumb").html(authorUser.getAvatarThumbnail());
                        }
                    });
            }
        }

    });

});