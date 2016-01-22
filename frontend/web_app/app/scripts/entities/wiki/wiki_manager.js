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
        wikiCollection.setQueryOptions({memberOnly: true});
        return this.fetchDeferred(wikiCollection);
    };

    WikiManager.prototype.getFile = function(id, wikiId){
        var fileModel = new FileModel({id: id}, {wikiId: wikiId});
        return this.fetchDeferred(fileModel);
    };

    WikiManager.prototype.getAllPageList = function(wikiId){
        var pageCollection = new PageCollection();
        pageCollection.url = "/api/v1/wikis/" + wikiId + "/pages";
        return this.fetchDeferred(pageCollection);
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
