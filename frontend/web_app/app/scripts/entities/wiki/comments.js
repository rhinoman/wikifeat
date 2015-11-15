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
    'entities/wiki/comment',
    'entities/base_collection'
], function($,_,Backbone,CommentModel,BaseCollection){

    //Constructor
    function CommentCollection(models, options){
        options = options || {};
        if(options.hasOwnProperty('wikiId')){
            this.wikiId = options.wikiId;
        }
        if(options.hasOwnProperty('pageId')){
            this.pageId = options.pageId;
        }
        BaseCollection.call(this, "comment", CommentModel, models, options);
    }

    CommentCollection.prototype = Object.create(BaseCollection.prototype);

    CommentCollection.prototype.comparator = "created_time";

    CommentCollection.prototype._prepareModel = function(model, options){
        options.wikiId = this.wikiId;
        options.pageId = this.pageId;
        return BaseCollection.prototype._prepareModel.call(this, model, options);
    };

    CommentCollection.prototype.url = function(){
        return "/api/v1/wikis/" + this.wikiId + "/pages/" + this.pageId + "/comments";
    };

    //Pagination state vals
    CommentCollection.prototype.state = {
        firstPage: 1
    };

    return CommentCollection;

});