/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 * See the LICENSE.txt file distributed with this work for additional information
 * regarding copyright ownership.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *  Neither the name of Wikifeat nor the names of its contributors may be used
 * to endorse or promote products derived from this software without
 * specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

'use strict';
define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'views/paginated_table_view',
    'views/user/manage_users_item',
    'views/user/edit_user_dialog',
    'entities/user/user',
    'text!templates/user/manage_users.html',
    'text!templates/main/alert.html'
], function($,_,Marionette, Radio, PaginatedTableView,
            ManageUsersItemView, EditUserDialogView, UserModel,
            ManageUsersTemplate, AlertTemplate){

    return PaginatedTableView.extend({
        className: "manage-users-view",
        template: _.template(ManageUsersTemplate),
        childView: ManageUsersItemView,
        childViewContainer: "#userListContainer",
        additionalEvents: {
            'click #addUserButton' : 'addUser',
            'click a#submitSearchForm': 'searchUsers',
            'submit #userSearchForm' : 'searchUsers'
        },
        childEvents: {
            'admin:enabled': function(){
                this.$("#alertBox").html(AlertTemplate);
                this.$(".alert").addClass("alert-success");
                this.$("#alertText").html("Admin access granted");
            },
            'admin:disabled': function(){
                this.$("#alertBox").html(AlertTemplate);
                this.$(".alert").addClass("alert-success");
                this.$("#alertText").html("Admin access revoked");
            },
            'admin:error': function(){
                this.$("#alertBox").html(AlertTemplate);
                this.$("#alertBox").addClass("alert-danger");
                this.$("#alertText").html("Could not toggle Admin access");
            }
        },

        //Show the Create new user dialog
        addUser: function(event){
            event.preventDefault();
            var user = new UserModel();
            var editUserDialog = new EditUserDialogView({model: user});
            var self = this;
            this.listenToOnce(user, 'sync', function(data){
                self.collection.add(data);
            });
            Radio.channel('main').trigger('show:dialog', editUserDialog);
        },

        //Search for a user (or users)
        searchUsers: function(event){
            event.preventDefault();
            var self = this;
            var searchText = this.$("input#searchText").val();
            var param = {searchText: searchText};
            Radio.channel('userManager').request("get:allUserList", param)
                .done(function(data){
                    if(typeof data !== 'undefined') {
                        self.collection.reset(data.models);
                    }
                });
        }

    });

});
