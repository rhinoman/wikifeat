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
    'backbone.stickit',
    'backbone.radio',
    'entities/wiki/history_entry',
    'text!templates/page/history_entry.html'
], function($,_,Marionette,Stickit,Radio,
            HistoryEntryModel,HistoryEntryTemplate){
    'use strict';

    return Marionette.ItemView.extend({
        id: 'history-entry',
        template: _.template(HistoryEntryTemplate),
        model: HistoryEntryModel,
        tagName: 'tr',
        wikiModel: null,
        pageModel: null,

        bindings: {
            '#timestamp': {
                observe: 'timestamp',
                updateMethod: 'html',
                onGet: function(timestamp){
                    var time = moment(timestamp);
                    var timeStr = time.format("HH:mm on D MMM YYYY");
                    return '<a id="viewRevisionLink" href="#">' + timeStr + '</a>'
                }
            },
            '#editor': {
                observe: 'editor'
            },
            '#contentSize': {
                observe: 'contentSize'
            }
        },
        events: {
            'click a#viewRevisionLink': 'viewRevision'
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
            }
            if(options.hasOwnProperty('pageModel')){
                this.pageModel = options.pageModel;
            }
        },

        viewRevision: function(event){
            event.preventDefault();
            Radio.channel('page').trigger('show:page:revision',
                this.pageModel.get('slug'),
                this.wikiModel,
                this.model.get('documentId'),
                {slug: true});
        },

        onRender: function(){
            if(typeof this.model !== "undefined"){
                this.stickit();
            }
        },

        onClose: function(){
            this.unstickit();
        }

    });

});
