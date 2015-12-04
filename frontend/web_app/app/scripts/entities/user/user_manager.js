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

'use strict';

define([
    'jquery',
    'underscore',
    'backbone',
    'marionette',
    'entities/base_manager',
    'entities/user/user',
    'entities/user/users',
    'entities/user/user_avatar'
], function($,_,Backbone,Marionette,BaseManager,
            UserModel,UserCollection,UserAvatarModel){

    //Constructor
    var UserManager = function(){
        BaseManager.call(this, UserModel)
    };
    var RoleRequest = Backbone.Model.extend({
            idAttribute: 'resourceId',
            defaults: {
                resourceType: '',
                resourceId: '',
                accessType: ''
            }
        });

    var ChangePasswordRequest = Backbone.Model.extend({
        defaults: {
            oldPassword: "",
            newPassword: ""
        }
    });

    UserManager.prototype = Object.create(BaseManager.prototype);

    UserManager.prototype.currentUser = null;

    //Special function -- gets currently authenticated user
    UserManager.prototype.getCurrentUser = function(){
        if(this.currentUser == null){
            var entity = new UserModel();
            entity.url = "/api/v1/users/current_user";
            this.currentUser = this.fetchDeferred(entity);
        }
        return this.currentUser;
    };

    UserManager.prototype.getAllUserList = function(queryOptions){
        var userCollection = new UserCollection();
        if(typeof queryOptions !== 'undefined'){
            userCollection.setQueryOptions(queryOptions);
        }
        return this.fetchDeferred(userCollection);
    };

    UserManager.prototype.getWikiMemberList = function(wikiModel){
        var memberCollection = new UserCollection();
        var queryOptions = {forResource: 'wiki_' + wikiModel.id};
        memberCollection.setQueryOptions(queryOptions);
        return this.fetchDeferred(memberCollection, {});
    };

    //Grant a User Role
    UserManager.prototype.grantRole = function(userModel, resourceType,
                                               resourceId, accessType){
        var roleRequest = new RoleRequest;
        roleRequest.set({
            resourceType: resourceType,
            resourceId: resourceId,
            accessType: accessType
        });
        roleRequest.url = userModel.url + "/grant_role";
        return this.saveEntity(roleRequest);
    };

    //Revoke a User Role
    UserManager.prototype.revokeRole = function(userModel, resourceType,
                                                resourceId, accessType){
        var roleRequest = new RoleRequest;
        roleRequest.set({
            resourceType: resourceType,
            resourceId: resourceId,
            accessType: accessType
        });
        roleRequest.url = userModel.url + "/revoke_role";
        return this.saveEntity(roleRequest);
    };

    //Change a User's password
    UserManager.prototype.changePassword = function(userModel, oldPassword, newPassword){
        var cpr = new ChangePasswordRequest({
            id: userModel.id,
            newPassword: newPassword,
            oldPassword: oldPassword
        });
        cpr.url = userModel.url + "/change_password";
        cpr.revision = userModel.revision;
        return this.saveEntity(cpr);
    };

    //Retrieve a User Avatar Record
    UserManager.prototype.getAvatar = function(id){
        var avatarModel = new UserAvatarModel({id: id});
        return this.fetchDeferred(avatarModel);
    };

    return UserManager;

});
