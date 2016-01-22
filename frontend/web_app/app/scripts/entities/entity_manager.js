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

    userChannel.reply("get:avatar", function(id){
        return userManager.getAvatar(id);
    });

    userChannel.reply("save:avatar", function(avatarModel){
        return userManager.saveEntity(avatarModel);
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

    wikiChannel.reply("delete:wiki", function(wikiModel){
        return wikiManager.deleteEntity(wikiModel);
    });

    wikiChannel.reply("get:slug:page", function(slug, wikiSlug){
        return wikiManager.getPageBySlug(slug, wikiSlug);
    });

    wikiChannel.reply("get:allWikiList", function(){
        return wikiManager.getAllWikiList();
    });

    wikiChannel.reply("get:memberWikiList", function(){
        return wikiManager.getMemberWikiList();
    });

    wikiChannel.reply("get:page", function(pageId, wikiId){
        return wikiManager.getPage(pageId, wikiId);
    });

    wikiChannel.reply("get:allPageList", function(wikiId){
        return wikiManager.getAllPageList(wikiId);
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

    wikiChannel.reply("get:page:comments", function(pageId, wikiId){
        return wikiManager.getPageComments(pageId, wikiId);
    });

    wikiChannel.reply("save:comment", function(commentModel){
        return wikiManager.saveEntity(commentModel);
    });

    wikiChannel.reply("delete:comment", function(commentModel){
        return wikiManager.deleteEntity(commentModel);
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
