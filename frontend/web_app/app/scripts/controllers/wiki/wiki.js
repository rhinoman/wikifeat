/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

/**
 * Created by jcadam on 2/3/15.
 */
'use strict';

define([
    'jquery',
    'underscore',
    'backbone',
    'marionette',
    'backbone.radio',
    'layouts/wiki_layout',
    'views/wiki/wiki_toolbar',
    'views/wiki/edit',
    'entities/wiki/wiki',
    'views/wiki/breadcrumbs',
    'entities/wiki/breadcrumb'
], function($,_,Backbone,Marionette,Radio,
            WikiLayout,WikiToolbarView,EditWikiView,
            WikiModel,BreadcrumbsView,BreadcrumbModel){

    var wikiChannel = Radio.channel('wiki');
    //var wikiLayout = new WikiLayout();

    var WikiController = Marionette.Controller.extend({

        wikiLayout: WikiLayout,

        initialize: function(options){
            this.breadcrumbs = new Backbone.Collection({model: BreadcrumbModel});
        },

        initLayout: function(wikiModel){
            this.wikiLayout = new WikiLayout();
            var breadcrumbView = new BreadcrumbsView({collection: this.breadcrumbs});
            var toolbarView = new WikiToolbarView();
            Radio.channel('main').trigger('show:content', this.wikiLayout);
            toolbarView.model = wikiModel;
            this.wikiLayout.toolbarRegion.show(toolbarView);
            this.breadcrumbs.reset();
            this.wikiLayout.breadcrumbRegion.show(breadcrumbView);
        },

        //Create a new wiki
        createWiki: function(options){
            options = options || {};
            this.wikiLayout = new WikiLayout();
            Radio.channel('main').trigger('show:content', this.wikiLayout);
            var wikiModel = new WikiModel();
            var editWikiView = new EditWikiView({model: wikiModel});
            this.wikiLayout.pageViewRegion.show(editWikiView);
            Backbone.history.navigate('/wikis/create');
            //Radio.channel('main').trigger('show:content', editWikiView);
        },

        editWiki: function(wikiModel, options){
            options = options || {};
            this.wikiLayout = new WikiLayout();
            Radio.channel('main').trigger('show:content', this.wikiLayout);
            var editWikiView = new EditWikiView({model: wikiModel});
            var toolbarView = new WikiToolbarView({model: wikiModel});
            this.wikiLayout.pageViewRegion.show(editWikiView);
            this.wikiLayout.toolbarRegion.show(toolbarView);
            Backbone.history.navigate('/wikis/' + wikiModel.get('slug') + '/edit');
        },

        //Show a page's breadcrumbs
        showCrumbs: function(wikiModel, pageModel){
            var self = this;
            Radio.channel('wikiManager')
                .request('get:page:breadcrumbs', pageModel.id, wikiModel.id)
                .done(function(crumbs){
                    self.breadcrumbs.reset(crumbs.models);
                });
        },

        //Navigate to a breadcrumb
        navigateToCrumb: function(wikiId, pageId){
           this.showWiki(wikiId, pageId,{slug:false});
        },

        //Load a wiki and show the specified page.
        //If no pageId is specified, shows the wiki's home page
        showWiki: function(wikiId, pageId, options){
            var self = this;
            var getWikiMsg = "get:wiki";
            var showPageMsg = "show:page";
            if(options.slug === true){
                getWikiMsg = "get:slug:wiki";
                showPageMsg = "show:slug:page";
            }
            Radio.channel('wikiManager').request(getWikiMsg, wikiId)
                .done(function(wikiModel){
                    if(typeof wikiModel === 'undefined'){
                        //TODO need some kind of error display thingie to call
                        console.log("Wiki not found: " + wikiId);
                        return;
                    }
                    self.initLayout(wikiModel);
                    if(typeof pageId !== 'undefined' && pageId !== null){
                        Radio.channel('page').trigger(showPageMsg,
                            pageId, wikiModel, options);
                    } else if(wikiModel.get('homePageId') != ""){
                        pageId = wikiModel.get("homePageId");
                        Radio.channel('page').trigger('show:page',
                            pageId, wikiModel);
                    } else {
                        //No home page, show the placeholder template.
                        Backbone.history.navigate('/wikis/' + wikiModel.get('slug'));
                        console.log("Wiki has no homepage.");
                        Radio.channel('page').trigger('show:placeholder:page', wikiModel);
                    }
                });
        },

        getPageRegion: function(){
           return this.wikiLayout.pageViewRegion;
        }
    });

    var wikiController = new WikiController();

    wikiChannel.on("show:wiki", function(wikiId, pageId){
        wikiController.showWiki(wikiId, pageId, {slug: false});
    });

    wikiChannel.on("show:slug:wiki", function(wikiSlug, pageSlug, options){
        options = options || {};
        options.slug = true;
        wikiController.showWiki(wikiSlug, pageSlug, options);
    });

    wikiChannel.on("create:wiki", function(options){
        wikiController.createWiki(options);
    });

    wikiChannel.on("edit:wiki", function(wikiModel, options){
        wikiController.editWiki(wikiModel, options);
    });

    wikiChannel.reply("get:pageRegion", function(){
        return wikiController.getPageRegion();
    });

    wikiChannel.on("show:breadcrumbs", function(wikiModel, pageModel){
        wikiController.showCrumbs(wikiModel, pageModel);
    });

    wikiChannel.on("go:crumb", function(wikiId, pageId){
        wikiController.navigateToCrumb(wikiId, pageId);
    });

    wikiChannel.on("init:layout", function(wikiModel){
        wikiController.initLayout(wikiModel);
    });

    return wikiController;

});
