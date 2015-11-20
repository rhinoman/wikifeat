/*
 *  Licensed to Wikifeat under one or more contributor license agreements.
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
/**
 * Created by jcadam on 12/8/14.
 */

define([
    'jquery',
    'jquery-cookie',
    'underscore',
    'marionette',
    'util/radio.shim',
    'backbone.radio',
    'routers/main_router',
    'routers/wiki_router',
    'routers/user_router',
    'controllers/main/sidebar',
    'controllers/main/plugin',
    'controllers/wiki/wiki',
    'controllers/wiki/page',
    'controllers/wiki/file',
    'controllers/user/user',
    'entities/entity_manager',
    'entities/error',
    'views/main/error_dialog',
    'text!templates/main/404.html'
], function ($, jqc, _, Marionette, Shim, Radio, MainRouter, WikiRouter,
             UserRouter, SidebarController, PluginController, WikiController,
             PageController, FileController, UserController, EntityManager,
             ErrorModel, ErrorDialogView, NotFoundTemplate) {
    'use strict';
    var WikiClient = new Backbone.Marionette.Application();

    //Add regions
    WikiClient.addRegions({
        sidebarRegion: "#sidebar",
        contentRegion: "#content",
        dialogRegion: "#dialogs"
    });

    WikiClient.addInitializer(function (options) {
        console.log("initialized.");
        //new MainRouter();
        //new WikiRouter();
    });

    WikiClient.navigate = function (route, options) {
        options || (options = {});
        Backbone.history.navigate(route, options);
    };

    WikiClient.getCurrentRoute = function () {
        return Backbone.history.fragment
    };

    // Start it.
    WikiClient.on("start", function () {
        //menu toggle event
        $("#menuToggle button").click(function(){
            $("#container").toggleClass('toggled');
        });
        //jquery setup
        $.ajaxPrefilter(function (options, originalOptions, jqXHR) {
            var csrfCookie = $.cookie("CsrfToken");
            if (typeof csrfCookie !== 'undefined') {
                jqXHR.setRequestHeader('X-Csrf-Token', csrfCookie);
            }
        });
        //Some global ajax error handling
        $(document).ajaxError(function (event, jqXHR, settings, exception) {
            var errorModel = new ErrorModel({code: jqXHR.status});
            switch (jqXHR.status) {
                case 0:
                    errorModel.set('name', 'Server Unreachable');
                    errorModel.set('message', 'Unable to connect to server.  Please try again later.');
                    WikiClient.getRegion('dialogRegion').show(
                        new ErrorDialogView({model: errorModel})
                    );
                    break;
                //If we get an 'Unauthenticated' response, redirect to the login page.
                case 401:
                    var currentUrl = document.location.href;
                    encodeURI(currentUrl);
                    window.location = '/login?ref=' + encodeURI(currentUrl);
                    break;
                case 404:
                    WikiClient.getRegion('contentRegion').show(
                        new Marionette.ItemView({template: _.template(NotFoundTemplate)})
                    );
                    break;
                case 500:
                    //I have no idea what went wrong.
                    console.log("An error has occurred.  Please try again later.");
                    WikiClient.getRegion('dialogRegion').show(
                        new ErrorDialogView({model: errorModel})
                    );
            }
        });

        Radio.channel('sidebar').trigger('init:layout', this.sidebarRegion);
        if (this.getCurrentRoute() === "") {
            WikiClient.trigger("showHome");
        }
        //Handler for content region
        Radio.channel('main').on('show:content', function(content){
            WikiClient.getRegion('contentRegion').show(content);
        });
        //Handler for showing dialogs
        Radio.channel('main').on('show:dialog', function(content){
            WikiClient.getRegion('dialogRegion').show(content);
        });

        new MainRouter();
        new WikiRouter();
        new UserRouter();
        //Start plugins
        var pluginsStarted = Radio.channel('plugin').request('start:plugins')
            .done(function(){
                // Start Backbone history
                // Need to do this AFTER the plugins are loaded
                // In case the plugins start their own routers
                Backbone.history.start({
                    pushState: true,
                    root: "/app/"
                });
            });

    });

    return WikiClient;
});
