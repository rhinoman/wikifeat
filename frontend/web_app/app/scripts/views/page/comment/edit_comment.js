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
    'markdown',
    'entities/wiki/comment',
    'views/main/alert',
    'text!templates/page/edit_comment.html'
], function($,_,Marionette,Radio,Markdown,
            CommentModel,AlertView,EditCommentTemplate){

    return Marionette.ItemView.extend({
        Model: CommentModel,
        template: _.template(EditCommentTemplate),
        id: "edit-comment-view",

        events:{
            'click #cancelEditButton': 'cancelButtonClick',
            'submit form#editCommentForm' : 'publishComment'
        },

        initialize: function(options){},

        //save the comment
        publishComment: function(event){
            event.preventDefault();
            var newComment = false;
            var pageContent = _.clone(this.model.get('content'));
            pageContent.raw = $("#wmd-input").val();
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
                    self.cancelButtonClick().bind(self);
                })
        },

        //simply destroy the view on cancel
        cancelButtonClick: function(event){
            this.destroy();
        },

        onRender: function(){
            var content = this.model.get("content");
            if(content.raw !== ""){
                this.$("#wmd-input").val(content.raw);
            }
        },

        onShow: function(){
            var editor = new Markdown.Editor();
            editor.run();
        },

        onClose: function(){
        }

    });

});