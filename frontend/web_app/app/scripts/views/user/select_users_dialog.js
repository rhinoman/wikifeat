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
    'bootstrap',
    'views/paginated_table_view',
    'views/user/select_users_item',
    'text!templates/user/select_users_dialog.html',
    'text!templates/main/alert.html'
], function($,_,Marionette,Radio,
            Bootstrap,PaginatedTableView,
            SelectUsersItemView,SelectUsersTemplate,
            AlertTemplate){

    return PaginatedTableView.extend({
        className: "select-users-view",
        template: _.template(SelectUsersTemplate),
        childView: SelectUsersItemView,
        childViewContainer: "#userListContainer",
        resource: "",

        childViewOptions: function(model, index) {
            return {
                resourceId: this.resourceId
            }
        },

        childEvents: {
            'add:member': function(childView, member){
                this.memberList.add(member);
            },
            'error': function(childView, errorMessage){
                this.$("#alertBox").html(AlertTemplate);
                this.$("#alertBox").addClass("alert-danger");
                this.$("#alertText").html(errorMessage);
            }
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('resourceId')){
                this.resourceId = options.resourceId;
                this.resource = 'wiki_' + this.resourceId;
            }
            if(options.hasOwnProperty('memberList')){
                this.memberList = options.memberList;
            }
        },

        onShow: function(){
            if(typeof this.collection !== 'undefined'){
                this.$("#selectUserModal").modal();
            }
        }
    });

});
