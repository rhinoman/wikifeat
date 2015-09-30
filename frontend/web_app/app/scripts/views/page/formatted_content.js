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

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio'
], function($,_,Marionette,Radio){

    'use strict';

    return Marionette.ItemView.extend({
        template: _.template('<div id="formatted"></div>'),
        pluginViews: $.Deferred(),

        initialize: function(){
            this.pluginsStarted = Radio.channel('plugin').request('get:pluginsStarted');
        },
        onRender: function(){
            this.$("#formatted").html(this.model.get("content").formatted);

        },
        onShow: function(){
            var self = this;
            $.when(this.pluginsStarted).done(function(){
                    self.loadContentPlugins();
            });
            //$.when(this.pluginViews).done(function(views){
            //    _.each(views, function(view){
            //        view.render();
            //    });
            //});
        },
        loadContentPlugins: function(){
            var contentFields = this.$("#formatted").find("[data-plugin]");
            //var pvs = [];
            _.each(contentFields, function(field){
                var pluginName = $(field).data('plugin');
                var resourceId = $(field).data('id');
                console.log("PLUGIN: " + pluginName + ", ID: " + resourceId);
                var pg = window[pluginName];
                if(typeof pg !== 'undefined') {
                    try {
                        var contentView = pg.getContentView(field, resourceId);
                        contentView.render();
                    }
                    catch(e){ //Bad Plugin! Bad!
                        console.log(e);
                    }
                } else {
                    console.log("Plugin " + pluginName + " is undefined");
                }
            });
        }

    });

});