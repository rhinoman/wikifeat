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

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.stickit',
    'backbone.radio',
    'entities/wiki/page',
    'text!templates/page/child_index_item.html'
], function($,_,Marionette,Stickit,Radio,
            PageModel,ChildIndexItemTemplate){
    'use strict';
    return Marionette.ItemView.extend({
        id: 'child-index-item-view',
        tagName: "li",
        template: _.template(ChildIndexItemTemplate),
        model: PageModel,
        wikiModel: null,
        bindings: {
            '#childPageTitle': {
                observe: 'title'
            }
        },
        events: {
            'click a.index-link': 'navigateToChildPage'
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('wikiModel')) {
                this.wikiModel = options.wikiModel
            }
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.stickit();
                var theLink = $(this.el).find('a.index-link');
                theLink.prop('href', this.model.get('slug'));
            }
        },

        navigateToChildPage: function(event){
            event.preventDefault();
            if(this.wikiModel !== null) {
                Radio.channel('page').trigger('show:page', this.model.id, this.wikiModel);
            }
        },

        onClose: function(){
            this.unstickit();
        }

    });

});
