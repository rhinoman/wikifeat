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
 * Created by jcadam on 1/14/15.
 */

'use strict';

define([
    'jquery',
    'underscore',
    'backbone',
    'entities/base_model'
], function($,_,Backbone, BaseModel){
    //Constructor
    function UserModel(data, options){
        BaseModel.call(this, "user", data, options);
    }

    UserModel.prototype = Object.create(BaseModel.prototype);

    UserModel.prototype.urlRoot = "/api/v1/users";
    UserModel.prototype.defaults = {
            name: "",
            password: "",
            verifyPassword: "",
            userPublic: {
                lastName: "",
                firstName: "",
                middleName: "",
                title: "",
                contactInfo: {
                    email: ""
                },
                avatar: "/app/resource/img/default_avatar.jpg",
                avatarThumbnail: "/app/resource/img/default_avatar_thumb.jpg"
            }
    };

    //Is this user a system admin?
    //By definition, the master is also an admin
    UserModel.prototype.isAdmin = function(){
        return this.hasRole('admin') || this.hasRole('master');
    };

    //Is this the 'Master' user?
    UserModel.prototype.isMaster = function(){
        return this.hasRole('master');
    };

    // Is user an admin of the given resource?
    UserModel.prototype.isResourceAdmin = function(resource){
        var adminRoleString = resource + ":admin";
        return this.hasRole(adminRoleString);
    };

    // Does the user have write access for the given resource?
    UserModel.prototype.isResourceWriter = function(resource){
        var writerRoleString = resource + ":write";
        return this.hasRole(writerRoleString);
    };

    // Does the user have read access for the given resource?
    UserModel.prototype.isResourceReader = function(resource){
        var readerRoleString = resource + ":read";
        return this.hasRole(readerRoleString);
    };

    // Is the user a member of a resource?
    UserModel.prototype.isResourceMember = function(resource){
        return this.isResourceReader(resource) ||
            this.isResourceWriter(resource) ||
            this.isResourceAdmin(resource);
    };

    // Does user have the given role?
    UserModel.prototype.hasRole = function(role){
        if(this.has("roles")){
            var roles = this.get("roles");
            if(_.indexOf(roles, role) > -1){
                return true;
            }
        }
        return false;
    };

    //Get the img link for the avatar
    UserModel.prototype.getAvatar = function(){
        var up = this.get('userPublic');
        var d = new Date();
        if(up.avatar === ""){
            up.avatar = "/app/resource/img/default_avatar.jpg?_=" + d.getTime();
        }
        return '<img class="avatar" src="' + up.avatar + '?_=' + d.getTime() +'"/>';
    };

    //Get the img link for the avatar Thumbnail
    UserModel.prototype.getAvatarThumbnail = function(){
        var up = this.get('userPublic');
        var d = new Date();
        if(up.avatarThumbnail === ""){
            up.avatarThumbnail = "/app/resource/img/default_avatar.jpg?_=" + d.getTime();
        }
        return '<img class="avatar" src="' + up.avatarThumbnail + '?_=' + d.getTime() +'"/>';
    };
    //input validation function
    UserModel.prototype.validate = function(attrs, options) {
        var errors = {};
        if (!attrs.name || attrs.name === ""){
            errors.name = "Username can't be blank";
        } else if (attrs.name.length > 64){
            errors.name = "Username is too long!";
        }
        if (!attrs.userPublic.lastName || attrs.userPublic.lastName === ""){
            errors.lastName = "Last Name can't be blank";
        } else if (attrs.userPublic.lastName.length > 64){
            errors.lastName = "Last Name is too long!";
        }
        if (!attrs.userPublic.firstName || attrs.userPublic.firstName === ""){
            errors.firstName = "First Name can't be blank";
        } else if (attrs.userPublic.firstName.length > 64){
            errors.firstName = "First name is too long";
        }
        if (!attrs.userPublic.contactInfo.email || attrs.userPublic.contactInfo.email === ""){
            errors.email = "Email can't be blank";
        }
        if (attrs.update === false && (!attrs.password || attrs.password === "")){
            errors.password = "Password can't be blank"
        } else if (attrs.password !== attrs.verifyPassword){
            errors.verifyPassword = "Password fields don't match"
        }
        if (!_.isEmpty(errors)){
            return errors;
        }
    };

    return UserModel;
});
