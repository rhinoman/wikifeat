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
    'marionette',
    'commonmark',
    'entities/wiki/page'
], function($,_,Marionette,Commonmark,PageModel){

    return Marionette.ItemView.extend({
        model: PageModel,
        id:"wmd-preview",
        className: "wmd-panel wmd-preview",
        template: _.template("<div></div>"),
        wipText: null,

        initialize: function(options){
            if(options.hasOwnProperty('wipText')){
                this.wipText = options.wipText;
            }
        },

        onShow: function(){
            if(this.wipText !== null) {
                var reader = new Commonmark.Parser();
                var writer = new Commonmark.HtmlRenderer();
                var parsedText = reader.parse(this.wipText.get('data'));
                this.$el.html(writer.render(parsedText));
            }
        }
    });

});