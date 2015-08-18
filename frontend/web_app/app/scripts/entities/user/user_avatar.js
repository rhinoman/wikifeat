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
        if (options.hasOwnProperty("userId")){
            this.userId = options.userId;
        } else {
            this.userId = "";
        }
        BaseModel.call(this, "avatar_record", data, options);
    }

    UserAvatarModel.prototype = Object.create(BaseModel.prototype);

    UserAvatarModel.prototype.urlRoot = function(){
        return "/api/v1/users/" + this.userId + "/avatar"
    };

    UserAvatarModel.prototype.defaults = {
        "use_gravatar" : false
    };

    UserAvatarModel.prototype.idAttribute = "_id";

    return UserAvatarModel;

});