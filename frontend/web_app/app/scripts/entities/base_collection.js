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

/**
 * Created by jcadam on 1/26/15.
 * Base Collection for all Collections
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'backbone.paginator'
], function($,_,Backbone,Paginator){

    return Backbone.PageableCollection.extend({

        initialize: function(entityName, model, models, options){
            this.model = model;
            this.entityName = entityName;
            Backbone.PageableCollection.prototype.initialize.call(this, models, options);
        },

        //Parse collection data from hateoas/hal response
        parse: function (response) {
            this.url = response._links.self.href;
            this.isCreatable = response._links.hasOwnProperty('create');
            delete response._links;
            //Paging properties
            this.state.totalRecords = response.total_rows;
            this.state.lastPage = Math.ceil(this.state.totalRecords / this.state.pageSize);
            //this.state.lastPage = this.state.totalPages;
            this.offset = response.offset;
            return response._embedded['ea:' + this.entityName];
        },

        parseRecords: function (resp) {
            return resp._embedded['ea:' + this.entityName];
        },

        setQueryOptions: function (queryOptions){
            for (var property in queryOptions){
                if(queryOptions.hasOwnProperty(property)){
                    this.queryParams[property] = queryOptions[property];
                }
            }
        },

        queryParams: {
            currentPage: "pageNum",
            pageSize: "numPerPage",
            totalRecords: "total_rows",
            forResource: null
        }

    });
});
