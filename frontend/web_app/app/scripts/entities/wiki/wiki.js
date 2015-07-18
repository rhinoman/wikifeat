/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.
 *
 * This file is part of WikiFeat.
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

    //Basic wiki model
    function WikiModel(data, options){
        BaseModel.call(this, "wiki_record", data, options)
    }

    WikiModel.prototype = Object.create(BaseModel.prototype);

    WikiModel.prototype.urlRoot = "/api/v1/wikis";

    WikiModel.prototype.defaults = {
        name: "",
        description: "",
        homePageId: "",
        allowGuest: false
    };

    WikiModel.prototype.parseLinks = function(links){
        BaseModel.prototype.parseLinks.call(this, links);
        this.canViewIndex = links.hasOwnProperty('index');
        this.canCreatePage = links.hasOwnProperty('create_page');
        this.canUpdate = links.hasOwnProperty('update');
        this.canDelete = links.hasOwnProperty('delete');
    };

    //input validation function
    WikiModel.prototype.validate = function(attrs, options) {
            var errors = {};
            if (!attrs.name) {
                errors.name = "can't be blank";
            } else if ((attrs.name.length) > 128){
                errors.name = "is too long";
            }
            if (!attrs.description){
                errors.description = "can't be blank";
            }  else if ((attrs.description.length) > 256){
                errors.description = "is too long";
            }
            if (!_.isEmpty(errors)){
                return errors;
            }
    };

    return WikiModel;

});
