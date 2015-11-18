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
    'moment',
    'backbone.radio',
    'entities/wiki/comment',
    'views/main/confirm_dialog',
    'views/page/comment/edit_comment',
    'text!templates/page/show_comment.html'
], function($,_,Marionette,Moment,Radio,
            CommentModel,ConfirmDialog,
            EditCommentView,
            ShowCommentTemplate){

    return Marionette.ItemView.extend({
        className: "comment-view col-lg-12",
        template: _.template(ShowCommentTemplate),
        model: CommentModel,
        editMode: false,

        events: {
            'click #deleteCommentButton' : 'deleteComment',
            'click #editCommentButton' : 'editComment'
        },

        initialize: function(options){
            this.model.on('change', this.render, this);
            this.rm = new Marionette.RegionManager();
            this.editorRegion =  this.rm.addRegion("editor", "#editorContainer_" + this.model.id);
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
            this.rm.get("editor").show(ecv);
            var self = this;
            ecv.on('destroy', function(){
                self.$("#commentContent").css("display", "block");
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

        onRender: function(){
            if(typeof this.model !== "undefined") {
                var author = this.model.get("author");
                this.$("#commentAuthorName").html(author);
                var time = moment(this.model.get("created_time"));
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
        },

        onClose: function(){}

    });

});