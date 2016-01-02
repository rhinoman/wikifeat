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
                        console.log("Error fetching current user.");
                    } else if(curUser.get('name') === 'guest'){
                        //redirect to home page
                        window.location = '/app';
                    } else {
                        var asv = new AccountSettingsView({model: curUser});
                        Radio.channel('main').trigger('show:content', asv);
                        Radio.channel('sidebar').trigger('active:user:accountSettings');
                        Backbone.history.navigate('/users/account');
                    }
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
                    Radio.channel('sidebar').trigger('active:admin:manageUsers');
                    Backbone.history.navigate('/users/manage');
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
                    Backbone.history.navigate('/wikis/' + wikiModel.get('slug') + '/members');
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
