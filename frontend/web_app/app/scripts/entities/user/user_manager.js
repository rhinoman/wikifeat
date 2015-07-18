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
    'backbone',
    'marionette',
    'entities/base_manager',
    'entities/user/user',
    'entities/user/users'
], function($,_,Backbone,Marionette,BaseManager,
            UserModel,UserCollection){

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
            old_password: "",
            new_password: ""
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

    UserManager.prototype.changePassword = function(userModel, oldPassword, newPassword){
        var cpr = new ChangePasswordRequest({
            id: userModel.id,
            new_password: newPassword,
            old_password: oldPassword
        });
        cpr.url = userModel.url + "/change_password";
        cpr.revision = userModel.revision;
        return this.saveEntity(cpr);
    };

    return UserManager;

});
