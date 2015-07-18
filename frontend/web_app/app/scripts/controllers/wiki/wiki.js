/** This file is part of WikiFeat.
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
            window.history.pushState('','','/app/wikis/create');
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
            window.history.pushState('','','/app/wikis/' + wikiModel.get('slug') + '/edit');
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
                    window.history.pushState('','','/app/wikis/' + wikiModel.get('slug'));
                    if(typeof pageId !== 'undefined'){
                        Radio.channel('page').trigger(showPageMsg,
                            pageId, wikiModel, options);
                    } else if(wikiModel.get('homePageId') != ""){
                        pageId = wikiModel.get("homePageId");
                        Radio.channel('page').trigger('show:page',
                            pageId, wikiModel);
                    } else {
                        //No home page, show the placeholder template.
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
