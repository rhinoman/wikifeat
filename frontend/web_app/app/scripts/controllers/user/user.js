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
    'backbone.radio',
    'views/user/manage_users',
    'views/user/manage_members',
    'views/user/account_settings'
], function($,_,Marionette,Radio,ManageUsersView,
            ManageMembersView,AccountSettingsView){
    var userChannel = Radio.channel('user');

    var UserController = Marionette.Controller.extend({
        logout: function(){
            $.ajax({
                url: "/api/v1/users/login",
                type: "DELETE"
            }).done(function(msg){
                console.log("logout: " + msg);
                window.location = '/login';
            });
        },

        accountSettings: function(options){
            Radio.channel('userManager').request('get:currentUser')
                .done(function(curUser){
                    if(typeof curUser === 'undefined'){
                        console.log("Error fetching curretn user.");
                        return;
                    }
                    var asv = new AccountSettingsView({model: curUser});
                    Radio.channel('main').trigger('show:content', asv);
                    window.history.pushState('','','/app/users/account');
                })
        },

        manageUsers: function(options){
            options = options || {};
            Radio.channel('userManager').request('get:allUserList')
                .done(function(userList){
                    if(typeof userList === 'undefined'){
                        //TODO display error
                        console.log("Error loading user list.");
                        return;
                    }
                    var muv = new ManageUsersView({collection: userList});
                    Radio.channel('main').trigger('show:content', muv);
                    window.history.pushState('','','/app/users/manage');
                });
        },

        manageMembers: function(wikiModel){
            Radio.channel('userManager').request('get:wikiMemberList', wikiModel)
                .done(function(memberList){
                    if(typeof memberList === 'undefined'){
                        //TODO display error
                        console.log("Error loading member list.");
                        return;
                    }
                    var mmv = new ManageMembersView({
                            collection: memberList,
                            wikiModel: wikiModel
                        }
                    );
                    var region = Radio.channel('wiki').request('get:pageRegion');
                    region.show(mmv);
                    window.history.pushState('','','/app/wikis/' +
                        wikiModel.get('slug') + '/members');
                });
        }
    });

    var userController = new UserController();

    userChannel.on("user:logout", function(){
        userController.logout();
    });

    userChannel.on("user:accountSettings", function(){
        userController.accountSettings();
    });

    userChannel.on("manage:users", function(){
        userController.manageUsers();
    });

    userChannel.on("manage:members", function(wikiModel){
        userController.manageMembers(wikiModel);
    });

    return userController;
});
