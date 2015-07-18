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
    'entities/wiki/file',
    'entities/base_collection'
], function($,_,Backbone,FileModel,BaseCollection){
    'use strict';

    //Constructor
    function FileCollection(models, options){
        options = options || {};
        if(options.hasOwnProperty('wikiId')){
            this.wikiId = options.wikiId;
        }
        BaseCollection.call(this, "file", FileModel, models, options);
    }

    FileCollection.prototype = Object.create(BaseCollection.prototype);

    FileCollection.prototype.comparator = "name";

    FileCollection.prototype._prepareModel = function(model, options){
        options.wikiId = this.wikiId;
        return BaseCollection.prototype._prepareModel.call(this, model, options);
    };

    //Pagination state vals
    FileCollection.prototype.state = {
        firstPage: 1,
        pageSize: 25
    };

    return FileCollection;

});