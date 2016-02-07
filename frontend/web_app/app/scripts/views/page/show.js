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

define([
    'jquery',
    'underscore',
    'marionette',
    'moment',
    'backbone.radio',
    'backbone.stickit',
    'entities/wiki/page',
    'entities/user/user',
    'views/page/child_index',
    'views/page/page_tools',
    'views/page/raw_content',
    'views/page/formatted_content',
    'views/user/user_info_dialog',
    'views/page/comment/comments',
    'util/common_events',
    'text!templates/page/page_layout.html'
], function($,_,Marionette,Moment,Radio,Stickit,
            PageModel,UserModel,ChildIndexView,
            PageToolMenu,RawContentView,FormattedContentView,
            UserInfoDialog,CommentsView,CommonEvents,
            ShowPageTemplate){
    'use strict';

    return Marionette.LayoutView.extend({
        id: "show-page-view",
        template: _.template(ShowPageTemplate),
        model: PageModel,
        wikiModel: null,
        editorModel: null,
        oldRevision: false,
        pageChildren: $.Promise,
        pageComments: $.Promise,
        regions: {
            pageToolMenuRegion: "#pageTools",
            pageContentRegion: "#pageContent",
            pageChildIndexRegion: '#childIndex',
            pageCommentsRegion: '#pageComments'
        },
        bindings: {
            '#editorName': {
                observe: 'editor'
            },
            '#lastEditTime': {
                observe: 'timestamp',
                onGet: function(timestamp){
                    var time = moment(timestamp);
                    return time.format("HH:mm on D MMM YYYY");
                }
            }
        },
        events: {
            'click a#viewCurrentPageLink': 'showCurrentPage',
            'click a#editorName':          'showEditorInfo',
            'click .page-content a':       'handleLinkClick'
        },

        initialize: function(options){
            this.model.on('setViewMode', this.render, this);
            if(this.model.get('owningPage') !== this.model.id){
                this.oldRevision = true;
            }
            this.model.viewMode = "formatted";
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
                //Load the list of this page's children
                if(!this.oldRevision) {
                    this.pageChildren = Radio.channel('wikiManager')
                        .request('get:page:children', this.model.id, this.wikiModel.id);
                    this.pageComments = Radio.channel('wikiManager')
                        .request('get:page:comments', this.model.id, this.wikiModel.id);
                }
            }
        },

        showCurrentPage: function(event){
            event.preventDefault();
            Radio.channel('page').trigger('show:page',
                this.model.get("owningPage"), this.wikiModel)
        },

        showEditorInfo: function(event){
            event.preventDefault();
            var editorInfoDialog = new UserInfoDialog({model: this.editorModel});
            Radio.channel('main').trigger('show:dialog', editorInfoDialog);
        },

        /* on render callback */
        onRender: function(){
            var self = this;
            if(typeof this.model !== 'undefined'){
                this.stickit();
                if(this.model.viewMode === 'formatted'){
                    this.pageContentRegion.show(new FormattedContentView({model: this.model}));
                } else {
                    this.pageContentRegion.show(new RawContentView({model: this.model}));
                }
                if(!this.oldRevision) {
                    this.pageToolMenuRegion.show(new PageToolMenu({
                        model: this.model, wikiModel: this.wikiModel
                    }));
                }
                //Draw the child page index and comments
                if(this.model.id != "" && !this.oldRevision) {
                    this.pageChildren.done(this.drawChildIndex.bind(this));
                    this.pageComments.done(this.drawComments.bind(this));
                } else {
                    this.$("div#childIndex").hide();
                }
                if(this.oldRevision){
                    var revNotice = this.$("div#revisionNotice");
                    revNotice.addClass("alert-warning");
                    revNotice.html('<p>This is an old version of this document.  ' +
                    'The latest version is <a href="#" id="viewCurrentPageLink">here</a>.</p>');
                    revNotice.css("display","block");
                }
                //Draw the editor avatar
                Radio.channel('userManager').request('get:user', this.model.get('editor'))
                    .done(function(editorUser){
                        if(typeof editorUser !== 'undefined') {
                            self.editorModel = editorUser;
                            self.editorModel.set('name',self.model.get('editor'));
                            self.$("span#editorAvatarThumb").html(editorUser.getAvatarThumbnail());
                        }
                    });
                //Now draw the alertBox
                this.$el.prepend('<div id="alertBox"></div>');

            }
        },

        handleLinkClick: function(event){
            CommonEvents.handleLinkClick(event);
        },

        drawChildIndex: function(response){
            this.$("div#childIndex").show();
            this.pageChildIndexRegion.show(
                new ChildIndexView({
                    collection: response,
                    wikiModel: this.wikiModel
                })
            )
        },

        drawComments: function(response){
            this.pageCommentsRegion.show(
                new CommentsView({
                    collection: response,
                    wikiId: this.wikiModel.id,
                    pageId: this.model.id,
                    commentsDisabled: this.model.get("commentsDisabled")
                })
            )
        },

        onClose: function(){
            this.unstickit();
        }
    });

});
