/** Copyright (c) 2014-present James Adam.  All rights reserved.
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
    'marionette',
    'backbone.radio',
    'entities/plugin/plugin',
    'entities/plugin/plugins',
    'entities/plugin/plugin_manager',
    'wikifeat'
], function($,_,Marionette,Radio,
    PluginModel,PluginCollection,
    PluginManager, Wikifeat){

    var pluginManager = new PluginManager();

    var PluginController = Marionette.Controller.extend({

        pluginCollection: new PluginCollection(),

        pluginsStarted: $.Deferred(),

        getPluginList: function(){
            if(this.pluginCollection.length === 0) {
                return pluginManager.fetchDeferred(this.pluginCollection, {});
            } else {
                var ret = $.Deferred();
                ret.resolve(this.pluginCollection);
                return ret.promise();
            }
        },
        startPlugins: function(){
            var self = this;
            var pl = [];
            this.getPluginList().done(function(pluginList){
                for (var i = 0; i < pluginList.length; i++) pl[i] = $.Deferred();
                var models = pluginList.models;
                for (i = 0; i < models.length; i++ ){
                    var mainScript = models[i].get("mainScript");
                    var ns = models[i].id;
                    var callback = _.partial(self.startPlugin, ns, pl, i, _);
                    //require(["/app/plugin/" + ns + "/resource/" + mainScript], function(){
                    //    callback();
                    //});
                    $.getScript("/app/plugin/" + ns + "/resource/" + mainScript)
                        .done(callback)
                        .fail(function(jqxhr, settings, exception){
                            console.log(exception);
                            console.log("FAIL")
                        });
                }
                $.when.apply($, pl)
                    .done(function(){self.pluginsStarted.resolve(true);})
                    .fail(function(){self.pluginsStarted.resolve(false);})
            });
            return this.pluginsStarted.promise();
        },
        startPlugin: function(ns, pl, index, script){
            var ps = window[ns];
            pl[index].resolve(ns);
            try {
                ps.start();
            } catch(e) { //Bad Plugin!
                console.log(e);
            }
        }
    });

    var pluginController = new PluginController();

    var pluginChannel = Radio.channel('plugin');
    pluginChannel.reply('get:pluginList', function(){
        return pluginController.getPluginList();
    });
    pluginChannel.reply('start:plugins', function(){
        return pluginController.startPlugins();
    });
    pluginChannel.reply('get:pluginsStarted', function(){
        return pluginController.pluginsStarted.promise();
    });

    return pluginController;

});