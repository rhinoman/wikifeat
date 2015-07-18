/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.*
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
            'click #addUserButton' : 'addUser'
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
        }

    });

});
