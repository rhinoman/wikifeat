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
    'views/paginated_table_view',
    'views/page/comment/show_comment',
    'views/page/comment/edit_comment',
    'entities/wiki/comment',
    'text!templates/page/comments.html'
], function($,_,Marionette,Radio,PaginatedTableView,
            ShowCommentView, EditCommentView,
            CommentModel, CommentsTemplate){

    return PaginatedTableView.extend({
        className: "comments-view",
        template: _.template(CommentsTemplate),
        childViewContainer: "#commentsContainer",
        childView: ShowCommentView,
        wikiId: null,
        pageId: null,
        commentsDisabled: false,

        additionalEvents: {
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
            if(options.hasOwnProperty('commentsDisabled')){
                this.commentsDisabled = options.commentsDisabled;
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
            if(!this.collection.isCreatable || this.commentsDisabled){
                this.$("#postNewCommentButton").css("display","none");
            }
            if(this.commentsDisabled){
                this.$("div#commentButtons").prepend("<p>Comments are disabled for this page</p>");
            }
            if(this.collection.length < 1){
                this.$("nav#commentsPaginationNav").hide();
            } else {
                this.$("nav#commentsPaginationNav").show();
            }
            PaginatedTableView.prototype.onRender.call(this);
        },

        onDestroy: function(){
            Radio.channel('commentsView').reset();
            this.rm.destroy();
        }

    });

});