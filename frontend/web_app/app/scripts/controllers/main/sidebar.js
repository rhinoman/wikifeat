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
    'layouts/sidebar_layout',
    'views/sidebar/logo',
    'views/sidebar/user_menu',
    'views/sidebar/guest_user_menu',
    'views/sidebar/admin_menu',
    'views/sidebar/wiki_list',
    'entities/wiki/wikis',
    'entities/user/user'
], function($,_,Marionette,Radio,SidebarLayout,
            LogoView,UserMenuView,GuestUserMenuView,
            AdminMenuView,WikiListView,WikiCollection,
            UserModel){

    //Data channels
    var userChannel = Radio.channel('userManager');
    var wikiChannel = Radio.channel('wikiManager');
    //Views and such
    var sideBarLayout = new SidebarLayout();
    var logoView = new LogoView();
    var wikiCollection = new WikiCollection();
    var adminMenuView = new AdminMenuView();
    var wikiListView = new WikiListView({collection: wikiCollection});
    var SideBarController = Marionette.Controller.extend({
        drawSideBar: function(){
            //Get the current user data
            var currentUser = userChannel.request('get:currentUser');
            currentUser.done(function(data){
                var userMenuView;
                if(typeof data === 'undefined') {
                    //Well, something went wrong
                    //This often happens if we have a stale session hanging around.
                    //Destroy the bad session:
                    $.ajax({
                        url: "/api/v1/users/login",
                        type: "DELETE"
                    });
                    data = new UserModel({id: 'guest'});
                    data.set('name', 'guest');
                }
                if(data.get('name') === 'guest'){
                    userMenuView = new GuestUserMenuView();
                } else {
                    userMenuView = new UserMenuView({model: data});
                }
                sideBarLayout.userMenuRegion.show(userMenuView);
                if (data.isAdmin() === true){
                    sideBarLayout.adminMenuRegion.show(adminMenuView);
                } else {
                    adminMenuView.destroy();
                }
                //Get our wiki list
                wikiChannel.request('get:memberWikiList', data).done(function(data){
                    if(typeof data !== 'undefined') {
                        wikiCollection.reset(data.models);
                        sideBarLayout.wikiListRegion.show(wikiListView);
                    }
                });
            });

            sideBarLayout.logoRegion.show(logoView);
            console.log("Creating Sidebar");
        },
        addMenu: function(name, menuView){
            //Adds a new div to the DOM, sets it as a region in the sidebar,
            //then shows it.
            if($("#" + name).length == 0) {
                //Only add if the element #name isn't already in the DOM
                $(sideBarLayout.el).append('<div id="' + name + '"></div>');
                sideBarLayout.addRegion(name, '#' + name);
                sideBarLayout.getRegion(name).show(menuView);
            }
        },
        initLayout: function(region){
            region.show(sideBarLayout);
            this.drawSideBar();
        },
        setActiveLink: function(link){
            sideBarLayout.setActiveLink(link);
        }
    });
    var sideBarController = new SideBarController();

    var sideBarChannel = Radio.channel('sidebar');
    sideBarChannel.on('init:layout', function(region){
        sideBarController.initLayout(region);
    });
    sideBarChannel.on('active:link', function(target){
        sideBarController.setActiveLink(target);
    });
    //Add a new menu to the sidebar (usually plugins)
    sideBarChannel.on('add:menu', function(name, menuView){
        sideBarController.addMenu(name, menuView);
    });
    //Add a new wiki to the list (usually called when a new wiki is created).
    sideBarChannel.on('add:wiki', function(wikiModel){
        wikiCollection.push(wikiModel);
    });
    //Remove a wiki from the list (usually called when a wiki is deleted).
    sideBarChannel.on('remove:wiki', function(wikiModel){
        wikiCollection.remove(wikiModel);
    });
    //Set an active wiki by slug
    sideBarChannel.on('active:wiki', function(wiki){
        //sideBarController.setActiveWikiBySlug(wiki);
        wikiListView.setActiveWikiBySlug(wiki);
    });
    //Expand the Wiki Menu
    sideBarChannel.on('expand:wikis', function(){
        wikiListView.expandMenu();
    });
    //Set the manage users link
    sideBarChannel.on('active:admin:manageUsers', function(){
        if(adminMenuView !== null){
            adminMenuView.setManageUsers();
        }
    });
    //Set the create wiki link
    sideBarChannel.on('active:admin:createWiki', function(){
        if(adminMenuView !== null){
            adminMenuView.setCreateWiki();
        }
    });

    return sideBarController;
});
