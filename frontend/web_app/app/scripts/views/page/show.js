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
    'text!templates/page/page_layout.html',
], function($,_,Marionette,Moment,Radio,Stickit,
            PageModel,UserModel,ChildIndexView,
            PageToolMenu,RawContentView,FormattedContentView,
            UserInfoDialog,ShowPageTemplate){
    'use strict';

    return Marionette.LayoutView.extend({
        id: "show-page-view",
        template: _.template(ShowPageTemplate),
        model: PageModel,
        wikiModel: null,
        editorModel: null,
        oldRevision: false,
        pageChildren: 7,
        viewMode: 'formatted',
        regions: {
            pageToolMenuRegion: "#pageTools",
            pageContentRegion: "#pageContent",
            pageChildIndexRegion: '#childIndex'
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
            'click a#editorName': 'showEditorInfo',
            'click a.page-content-mode-link' : 'changeViewMode'
        },

        initialize: function(options){
            if(this.model.get('owning_page') !== this.model.id){
                this.oldRevision = true;
            }
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
                //Load the list of this page's children
                if(!this.oldRevision) {
                    this.pageChildren = Radio.channel('wikiManager')
                        .request('get:page:children', this.model.id, this.wikiModel.id);
                }
            }
        },

        changeViewMode: function(event){
            event.preventDefault();
            if(this.viewMode === 'raw'){
                this.viewMode = 'formatted';
            } else {
                this.viewMode = 'raw';
            }
            this.render();
        },

        showCurrentPage: function(event){
            event.preventDefault();
            Radio.channel('page').trigger('show:page',
                this.model.get("owning_page"), this.wikiModel)
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
                if(this.viewMode === 'formatted'){
                    this.pageContentRegion.show(new FormattedContentView({model: this.model}));
                    this.$("div#pageViewFormat").html(
                        '<a href="#" class="page-content-mode-link">' +
                        '<span class="glyphicon glyphicon-eye-open"></span>' +
                        '&nbsp;View Raw</a>');

                } else {
                    this.pageContentRegion.show(new RawContentView({model: this.model}));
                    this.$("div#pageViewFormat").html(
                        '<a href="#" class="page-content-mode-link">' +
                        '<span class="glyphicon glyphicon-eye-open"></span>' +
                        '&nbsp;View Formatted</a>')
                }
                if(!this.oldRevision) {
                    this.pageToolMenuRegion.show(new PageToolMenu({
                        model: this.model, wikiModel: this.wikiModel
                    }));
                }
                //Draw the child page index
                if(this.model.id != "" && !this.oldRevision) {
                    this.pageChildren.done(this.drawChildIndex.bind(this));
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

        onShow: function(){

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

        onClose: function(){
            this.unstickit();
        }
    });

});
