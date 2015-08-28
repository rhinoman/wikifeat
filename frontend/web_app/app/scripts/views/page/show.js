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
    'views/user/user_info_dialog',
    'text!templates/page/page_layout.html',
], function($,_,Marionette,Moment,Radio,Stickit,
            PageModel,UserModel,ChildIndexView,
            PageToolMenu,UserInfoDialog,ShowPageTemplate){
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
        pluginsStarted: $.Deferred(),
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
            this.pluginsStarted = Radio.channel('plugin').request('get:pluginsStarted');
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
                    this.$("div#pageContent").html(this.model.get("content").formatted);
                    this.$("div#pageViewFormat").html(
                        '<a href="#" class="page-content-mode-link">' +
                        '<span class="glyphicon glyphicon-eye-open"></span>' +
                        '&nbsp;View Raw</a>');
                    $.when(this.pluginsStarted).done(function(){
                        self.loadContentPlugins();
                    });
                } else {
                    this.$("div#pageContent").html(
                        '<div class="raw-view-box"><textarea disabled>' +
                        this.model.get("content").raw + '</textarea></div>');
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
                            self.$("span#editorAvatarThumb").html(editorUser.getAvatarThumbnail());
                        }
                    });
                //Now draw the alertBox
                this.$el.prepend('<div id="alertBox"></div>');
            }
        },

        loadContentPlugins: function(){
            var contentFields = this.$("#pageContent").find("[data-plugin]");
            _.each(contentFields, function(field){
                var pluginName = $(field).data('plugin');
                var resourceId = $(field).data('id');
                console.log("PLUGIN: " + pluginName + ", ID: " + resourceId);
                var pg = window[pluginName];
                if(typeof pg !== 'undefined') {
                    try {
                        var contentView = pg.getContentView(field, resourceId);
                        contentView.render();
                    }
                    catch(e){ //Bad Plugin! Bad!
                        console.log(e);
                    }
                } else {
                    console.log("Plugin " + pluginName + " is undefined");
                }
            });
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
