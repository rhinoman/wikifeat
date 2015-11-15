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
    'views/page/comment/show_comment',
    'views/page/comment/edit_comment',
    'entities/wiki/comment',
    'text!templates/page/comments.html'
], function($,_,Marionette,Radio,
            ShowCommentView, EditCommentView,
            CommentModel, CommentsTemplate){

    return Marionette.CompositeView.extend({
        className: "comments-view",
        template: _.template(CommentsTemplate),
        childViewContainer: "#commentsContainer",
        childView: ShowCommentView,
        wikiId: null,
        pageId: null,

        events: {
            'click #postNewCommentButton' : 'createComment'
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('wikiId')){
                this.wikiId = options.wikiId;
            }
            if(options.hasOwnProperty('pageId')){
                this.pageId = options.pageId;
            }
            this.rm = new Marionette.RegionManager();
            this.editorRegion = this.rm.addRegion("editor","#newCommentEditor");
            var self = this;
            Radio.channel('commentsView').on('new:comment', function(comment){
                self.collection.add(comment);
            });
        },

        createComment: function(event){
            var commentModel = new CommentModel({},{
                wikiId: this.wikiId,
                pageId: this.pageId
            });

            var ecv = new EditCommentView({model: commentModel});
            this.rm.get('editor').show(ecv);
        },

        onRender: function(){
            //disable Post Comment button if you don't have access
            if(!this.collection.isCreatable){
                this.$("#postNewCommentButton").css("display","none");
            }
        }

    });

});