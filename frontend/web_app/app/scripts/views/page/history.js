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
    'marionette',
    'moment',
    'backbone.radio',
    'views/paginated_table_view',
    'views/page/history_entry',
    'text!templates/page/history.html'
], function($,_,Marionette,Moment,Radio,PaginatedTableView,
            HistoryEntryView,HistoryTemplate){

    'use strict';

    return PaginatedTableView.extend({
        childView: HistoryEntryView,
        template: _.template(HistoryTemplate),
        childViewContainer: '#historyEntriesContainer',
        wikiModel: null,
        pageModel: null,

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
            }
            if(options.hasOwnProperty('pageModel')){
                this.pageModel = options.pageModel;
            }
            this.childViewOptions = {
                wikiModel: this.wikiModel,
                pageModel: this.pageModel
            }
        },

        onRender: function(){
            if(this.collection.length === 0){
                this.$el.find("#historyEntriesContainer").append("None")
            }
            if(this.pageModel != null){
                this.$("#pageTitle").html(this.pageModel.get('title'));
            }
            PaginatedTableView.prototype.onRender.call(this);
        }
    });

});
