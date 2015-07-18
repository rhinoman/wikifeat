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
/**
 * Created by jcadam on 1/21/15.
 * In this file we handle messaging to the various Entity Managers
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'backbone.radio',
    'entities/user/user_manager',
    'entities/wiki/wiki_manager'
], function($,_,Backbone,Radio,UserManager,WikiManager){

    var EntityManager = function(){
        //Might add some stuff here, maybe.
    };

    //--------User Manager Stuff
    var userManager = new UserManager();
    var userChannel = Radio.channel('userManager');

    userChannel.reply("get:currentUser", function(){
        return userManager.getCurrentUser();
    });

    userChannel.reply("get:user", function(id){
        return userManager.getEntity(id);
    });

    userChannel.reply("grant:role",
        function(userModel, resourceType, resourceId, accessType){
            return userManager.grantRole(userModel, resourceType, resourceId, accessType);
    });

    userChannel.reply("revoke:role",
        function(userModel, resourceType, resourceId, accessType){
            return userManager.revokeRole(userModel, resourceType, resourceId, accessType);
    });

    userChannel.reply("change:password",function(userModel, oldPassword, newPassword){
        return userManager.changePassword(userModel, oldPassword, newPassword);
    });

    userChannel.reply("save:user", function(userModel){
        return userManager.saveEntity(userModel);
    });

    userChannel.reply("delete:user", function(userModel){
        return userManager.deleteEntity(userModel);
    });

    userChannel.reply("get:allUserList", function(queryOptions){
        return userManager.getAllUserList(queryOptions);
    });

    userChannel.reply("get:wikiMemberList", function(wikiModel){
        return userManager.getWikiMemberList(wikiModel);
    });

    //-------End User Manager stuff
    //-------Wiki Manager Stuff

    var wikiManager = new WikiManager();
    var wikiChannel = Radio.channel('wikiManager');

    wikiChannel.reply("get:wiki", function(id){
        return wikiManager.getEntity(id);
    });

    wikiChannel.reply("get:slug:wiki", function(slug){
        return wikiManager.getEntityBySlug(slug);
    });

    wikiChannel.reply("save:wiki", function(wikiModel){
        return wikiManager.saveEntity(wikiModel);
    });

    wikiChannel.reply("get:slug:page", function(slug, wikiSlug){
        return wikiManager.getPageBySlug(slug, wikiSlug);
    });

    wikiChannel.reply("get:allWikiList", function(){
        return wikiManager.getAllWikiList();
    });

    wikiChannel.reply("get:page", function(pageId, wikiId){
        return wikiManager.getPage(pageId, wikiId);
    });

    wikiChannel.reply("save:page", function(pageModel){
        return wikiManager.saveEntity(pageModel);
    });

    wikiChannel.reply("delete:page", function(pageModel){
        return wikiManager.deleteEntity(pageModel);
    });

    wikiChannel.reply("get:page:children", function(pageId, wikiId){
        return wikiManager.getPageChildren(pageId, wikiId);
    });

    wikiChannel.reply("get:page:history", function(pageId, wikiId){
        return wikiManager.getPageHistory(pageId, wikiId);
    });

    wikiChannel.reply("get:page:breadcrumbs", function(pageId, wikiId){
        return wikiManager.getPageBreadcrumbs(pageId, wikiId);
    });

    wikiChannel.reply("get:allFileList", function(wikiId){
        return wikiManager.getAllFileList(wikiId);
    });

    wikiChannel.reply("get:file", function(fileId, wikiId){
        return wikiManager.getFile(fileId, wikiId);
    });

    wikiChannel.reply("save:file", function(fileModel){
        return wikiManager.saveEntity(fileModel);
    });

    wikiChannel.reply("delete:file", function(fileModel){
        return wikiManager.deleteEntity(fileModel);
    });

    //-------End Wiki Manager Stuff

    return EntityManager
});
