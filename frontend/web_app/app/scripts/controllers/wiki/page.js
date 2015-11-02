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
    'entities/wiki/page',
    'views/page/show',
    'views/page/edit',
    'views/page/placeholder',
    'views/page/history'
], function($,_,Marionette,Radio,PageModel,
            ShowPageView,EditPageView,
            PlaceholderView,HistoryView){

    var pageChannel = Radio.channel('page');

    function genStateString(wikiModel, pageModel){
        var pageString = pageModel.get('slug');
        var wikiString = wikiModel.get('slug');
        if(wikiString === ""){
            wikiString = wikiModel.id;
        }
        if(pageString === ""){
            pageString = pageModel.id;
        }
        return '/app/wikis/' + wikiString + '/pages/' + pageString;
    }

    var PageController = Marionette.Controller.extend({
        //fetch a page and display it

        getPageRequest: function(pageId, wikiModel, options){
            var pagePromise;
            if(options.slug === true){
                pagePromise = Radio.channel('wikiManager')
                    .request("get:slug:page",pageId,wikiModel.get('slug'));
            } else {
                pagePromise = Radio.channel('wikiManager')
                    .request("get:page",pageId,wikiModel.get('id'));
            }
            return pagePromise;
        },

        showPage: function(pageId, wikiModel, options){
            console.log("Showing page " + pageId + " in wiki " + wikiModel.get('id'));
            var pagePromise = this.getPageRequest(pageId, wikiModel, options);
            pagePromise.done(function(pageModel){
                if(typeof pageModel === 'undefined'){
                //TODO error display
                    console.log("Page not found: " + pageId);
                    return;
                }
                Radio.channel('wiki').trigger("show:breadcrumbs",wikiModel, pageModel);
                if(options.hasOwnProperty('edit') && options.edit === true){
                    Radio.channel('page').trigger("edit:page", pageModel, wikiModel);
                } else if(options.hasOwnProperty('history') && options.history === true){
                    Radio.channel('page').trigger('show:history',pageModel, wikiModel);
                } else {
                    window.history.pushState('', '', genStateString(wikiModel, pageModel));
                    var region = Radio.channel('wiki').request('get:pageRegion');
                    region.show(new ShowPageView({model: pageModel, wikiModel: wikiModel}));
                }
            });
        },
        showPageRevision: function(pageId, wikiModel, revisionId, options){
            var pagePromise = this.getPageRequest(pageId, wikiModel, options);
            var revisionPromise = this.getPageRequest(revisionId, wikiModel, {slug: false});
            $.when(pagePromise, revisionPromise).done(function(page, revision){
                var region = Radio.channel('wiki').request('get:pageRegion');
                var stateString = genStateString(wikiModel, page);
                if(revision.get('owning_page') !== revision.id){
                    stateString = stateString + '?revision=' + revision.id
                }
                window.history.pushState('','', stateString);
                region.show(new ShowPageView({model: revision, wikiModel: wikiModel}));
            });

        },
        //Display the placeholder page for homeless wikis
        showPlaceholderPage: function(wikiModel){
            var region = Radio.channel('wiki').request('get:pageRegion');
            region.show(new PlaceholderView({model: wikiModel}));
            console.log("Showing placeholder page for wiki " + wikiModel.id);
        },
        //Displays a page's history entries
        showHistory: function(pageModel, wikiModel){
            var promise = Radio.channel('wikiManager')
                .request("get:page:history",pageModel.id,wikiModel.id);
            promise.done(function(history){
                var region = Radio.channel('wiki').request('get:pageRegion');
                region.show(new HistoryView({
                    collection: history,
                    wikiModel: wikiModel,
                    pageModel: pageModel
                }));
            });
            window.history.pushState('','',genStateString(wikiModel, pageModel) + "/history");
            console.log("Showing history for page: " + pageModel.id)
        },

        //Display the create page view
        createPage: function(wikiModel, options){
            var region = Radio.channel('wiki').request('get:pageRegion');
            var model = new PageModel({},{wikiId: wikiModel.id});
            var homePage = false;
            if(options.hasOwnProperty('parent') && options.parent !== ""){
               model.set('parent', options.parent);
            }
            if(options.hasOwnProperty('homePage')){
                homePage = options.homePage;
            }
            region.show(new EditPageView({
                model: model,
                wikiModel: wikiModel,
                homePage: homePage
            }));
        },
        //Display the editor interface for this page
        editPage: function(pageModel, wikiModel){
            window.history.pushState('','', genStateString(wikiModel, pageModel) + "/edit");
            var region = Radio.channel('wiki').request('get:pageRegion');
            region.show(new EditPageView({model: pageModel, wikiModel: wikiModel}));
            console.log("Editing page " + pageModel.id + " of wiki " + wikiModel.id);
        },
        savePage: function(pageModel){
            return Radio.channel('wikiManager').request("save:page", pageModel);
        }
    });

    var pageController = new PageController();

    pageChannel.on("show:page", function(pageId, wikiModel){
        pageController.showPage(pageId, wikiModel, {slug: false});
    });

    pageChannel.on("show:slug:page", function(pageSlug, wikiModel, options){
        options = options || {};
        options.slug = true;
        pageController.showPage(pageSlug, wikiModel, options);
    });

    pageChannel.on("show:placeholder:page", function(wikiModel){
        pageController.showPlaceholderPage(wikiModel);
    });

    pageChannel.on("create:page", function(wikiModel, options){
        options = options || {};
        pageController.createPage(wikiModel, options);
    });

    pageChannel.on("show:history", function(pageModel, wikiModel){
        pageController.showHistory(pageModel, wikiModel);
    });

    pageChannel.on("show:page:revision", function(pageId, wikiModel, revisionId, options){
        pageController.showPageRevision(pageId, wikiModel, revisionId, options)
    });

    pageChannel.on("edit:page", function(pageModel, wikiModel){
        pageController.editPage(pageModel, wikiModel);
    });

    pageChannel.reply("save:page", function(pageModel){
        return pageController.savePage(pageModel);
    });

    return pageController;
});
