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
    'entities/base_model'
], function($,_,Backbone,BaseModel){

    //contructor
    function UserAvatarModel(data, options){
        options = options || {};
        BaseModel.call(this, "avatar_record", data, options);
    }

    UserAvatarModel.prototype = Object.create(BaseModel.prototype);

    UserAvatarModel.prototype.urlRoot = function(){
        return "/api/v1/users/" + this.get('id') + "/avatar"
    };

    UserAvatarModel.prototype.idAttribute = "_id";

    UserAvatarModel.prototype.uploadContent = function(formData){
        var defer = $.Deferred();
        var self = this;
        $.ajax({
            url: this.url + "/image",
            type: "POST",
            beforeSend: function(request){
                request.setRequestHeader("If-Match", self.revision);
            },
            data: formData,
            processData: false,
            contentType: false
        }).done(function(response){
            defer.resolve(response);
        }).fail(function(response){
            defer.resolve(undefined);
        });
        return defer.promise();
    };

    return UserAvatarModel;

});