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

define([
    'jquery',
    'underscore',
    'backbone',
    'entities/base_manager',
    'entities/wiki/wiki',
    'entities/wiki/wikis',
    'entities/wiki/file',
    'entities/wiki/files',
    'entities/wiki/page',
    'entities/wiki/pages',
    'entities/wiki/breadcrumb',
    'entities/wiki/history_entries',
    'entities/wiki/comments'
], function($,_,Backbone,BaseManager,
            WikiModel,WikiCollection,
            FileModel,FileCollection,
            PageModel,PageCollection,
            BreadcrumbModel,HistoryEntryCollection,
            CommentCollection){

    //Constructor
    var WikiManager = function(){
        BaseManager.call(this, WikiModel);
    };

    WikiManager.prototype = Object.create(BaseManager.prototype);

    WikiManager.prototype.getAllWikiList = function(){
        var wikiCollection = new WikiCollection();
        return this.fetchDeferred(wikiCollection,{});
    };

    WikiManager.prototype.getAllFileList = function(wikiId){
        var fileCollection = new FileCollection({},{wikiId: wikiId});
        fileCollection.url = "/api/v1/wikis/" + wikiId + "/files";
        return this.fetchDeferred(fileCollection);
    };

    WikiManager.prototype.getMemberWikiList = function(user){
        var wikiCollection = new WikiCollection();
        if(typeof user !== 'undefined' && user.id !== 'guest') {
            wikiCollection.setQueryOptions({memberOnly: true});
        }
        return this.fetchDeferred(wikiCollection);
    };

    WikiManager.prototype.getFile = function(id, wikiId){
        var fileModel = new FileModel({id: id}, {wikiId: wikiId});
        return this.fetchDeferred(fileModel);
    };

    WikiManager.prototype.getPage = function(id, wikiId){
        var pageModel = new PageModel({id: id}, {wikiId: wikiId});
        return this.fetchDeferred(pageModel);
    };

    WikiManager.prototype.getPageBreadcrumbs = function(id, wikiId){
        var breadcrumbCollection = new Backbone.Collection();
        breadcrumbCollection.url = "/api/v1/wikis/" + wikiId +
            "/pages/" + id  + "/breadcrumbs";
        breadcrumbCollection.model = BreadcrumbModel;
        breadcrumbCollection.parse = function(response){
            if(response.hasOwnProperty('crumbs')){
                return response.crumbs;
            } else {
                return [];
            }
        };
        return this.fetchDeferred(breadcrumbCollection);
    };

    WikiManager.prototype.getPageChildren = function(id, wikiId){
        var pageCollection = new PageCollection();
        pageCollection.url = "/api/v1/wikis/" + wikiId +
            "/pages/" + id  + "/children";
        return this.fetchDeferred(pageCollection);
    };

    WikiManager.prototype.getPageHistory = function(id, wikiId){
        var historyCollection = new HistoryEntryCollection();
        historyCollection.url = "/api/v1/wikis/" + wikiId +
            "/pages/" + id + "/history";
        return this.fetchDeferred(historyCollection);
    };

    WikiManager.prototype.getPageBySlug = function(slug, wikiSlug){
        var pageModel = new PageModel({id: slug}, {wikiSlug: wikiSlug});
        pageModel.url = "/api/v1/wikis/slug/" + wikiSlug + "/pages/" + slug;
        return this.fetchDeferred(pageModel);
    };

    WikiManager.prototype.getEntityBySlug = function(slug){
        var wikiModel = new WikiModel({id: slug});
        wikiModel.url = "/api/v1/wikis/slug/" + slug;
        return this.fetchDeferred(wikiModel);
    };

    WikiManager.prototype.getPageComments = function(pageId, wikiId) {
        var commentCollection = new CommentCollection({},{wikiId: wikiId, pageId: pageId});
        return this.fetchDeferred(commentCollection);
    };

    return WikiManager;
});
