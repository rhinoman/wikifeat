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
    'entities/base_model'
], function($,_,Backbone,BaseModel){
    'use strict';
    //Page model constructor
    function PageModel(data, options){
        options = options || {};
        if (options.hasOwnProperty("wikiId")){
            this.wikiId = options.wikiId;
        } else {
            this.wikiId = "";
        }
        BaseModel.call(this, "page", data, options);
    }

    PageModel.prototype = Object.create(BaseModel.prototype);

    PageModel.prototype.urlRoot = function(){
        return "/api/v1/wikis/" + this.wikiId + "/pages";
    };

    PageModel.prototype.defaults = {
        "content": {
            "formatted": "",
            "raw": ""
        },
        "editor": "Unknown",
        "owner": "",
        "parent": "",
        "title": "Untitled",
        "type": "page"
    };

    //input validation function
    PageModel.prototype.validate = function(attrs, options) {
        var errors = {};
        if (!attrs.title) {
            errors.title = "can't be blank";
        } else if ((attrs.title.length) > 128){
            errors.title = "is too long";
        }
        if(!$.isEmptyObject(errors)){
            return errors;
        }
    };

    return PageModel;

});
