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
    'layouts/sidebar_layout',
    'views/sidebar/logo',
    'views/sidebar/user_menu',
    'views/sidebar/guest_user_menu',
    'views/sidebar/admin_menu',
    'views/sidebar/wiki_list',
    'entities/wiki/wikis'
], function($,_,Marionette,Radio,SidebarLayout,
            LogoView,UserMenuView,GuestUserMenuView,
            AdminMenuView,WikiListView,WikiCollection){

    //Data channels
    var userChannel = Radio.channel('userManager');
    var wikiChannel = Radio.channel('wikiManager');
    //Views and such
    var sideBarLayout = new SidebarLayout();
    var logoView = new LogoView();
    var wikiCollection = new WikiCollection();
    var SideBarController = Marionette.Controller.extend({
        drawSideBar: function(){
            //Get the current user data
            userChannel.request('get:currentUser').done(function(data){
                var userMenuView;
                if(data.get('name') === 'guest'){
                    userMenuView = new GuestUserMenuView();
                } else {
                    userMenuView = new UserMenuView({model: data});
                }
                sideBarLayout.userMenuRegion.show(userMenuView);
                if (data.isAdmin() === true){
                    var adminMenuView = new AdminMenuView();
                    sideBarLayout.adminMenuRegion.show(adminMenuView);
                }

            });
            //Get our wiki list
            wikiChannel.request('get:memberWikiList').done(function(data){
                if(typeof data !== 'undefined') {
                    wikiCollection.reset(data.models);
                    var wikiListView = new WikiListView({collection: wikiCollection});
                    sideBarLayout.wikiListRegion.show(wikiListView);
                }
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
            $(sideBarLayout.el).find('a').removeClass('currentSelection');
            $(link).addClass('currentSelection');
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

    return sideBarController;
});
