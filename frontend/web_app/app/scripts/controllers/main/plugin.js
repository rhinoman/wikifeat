/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
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
                    $.getScript("/app/plugin/" + ns + "/resource/" + mainScript)
                        .done(callback)
                        .fail(function(jqxhr, settings, exception){
                            console.log("Failed to load plugin");
                            console.log(exception);
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
            var timeout = 2000;
            try {
                ps.start(pl[index]);
                setTimeout(function() {
                    if(pl[index].state() !== 'resolved'){
                        console.log("Plugin " + ns + " failed to start in " + timeout + " ms");
                        pl[index].reject();
                    }
                }, timeout);
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