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

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'views/paginated_table_view',
    'views/user/manage_members_item',
    'views/user/select_users_dialog',
    'entities/user/user',
    'text!templates/user/manage_members.html',
    'text!templates/main/alert.html'
], function($,_,Marionette,Radio,PaginatedTableView,
            ManageMembersItemView,SelectUsersDialog,UserModel,
            ManageMembersTemplate,AlertTemplate){
    return PaginatedTableView.extend({
        className: "manage-members-view",
        template: _.template(ManageMembersTemplate),
        childView: ManageMembersItemView,
        childViewContainer: "#memberListContainer",
        additionalEvents: {
            'click #addMemberButton' : 'addMember'
        },
        wikiModel: null,

        childViewOptions: function(model, index){
            return {
                resourceId: this.resourceId
            }
        },

        childEvents: {
            'remove:member': function(childView, member){
                this.collection.remove(member);
            },
            'error': function(childView, errorMessage){
                this.$("#alertBox").html(AlertTemplate);
                this.$("#alertBox").addClass("alert-danger");
                this.$("#alertText").html(errorMessage);
            }
        },

        initialize: function(options){
            var self = this;
            options = options || {};
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
                this.resourceId = this.wikiModel.id;
                this.resource = 'wiki_' + this.resourceId;
            } else {
                this.resourceId = null;
            }
            var curUserReq = $.Deferred();
            this.currentUser = curUserReq.promise();
            //Load the current user
            Radio.channel('userManager').request('get:currentUser')
                .done(function(response){
                    if(typeof response !== 'undefined') {
                        curUserReq.resolve(response);
                    }
                });
        },

        onShow: function(){
            if(this.wikiModel !== null){
                this.$("#manage_members_header").html(
                    'Manage Members for "' + this.wikiModel.get('name') + '"');
            }
            var self = this;
            this.currentUser.done(function(currentUser){
                self.children.each(function(view){
                    if(view.model.id === currentUser.id){
                        view.disableControls();
                    }
                });
            });
        },

        addMember: function(event){
            event.preventDefault();
            var self = this;
            Radio.channel('userManager').request('get:allUserList', {})
                .done(function(response){
                   if(typeof response === 'undefined'){
                   } else {
                       var sud = new SelectUsersDialog({
                           collection: response,
                           resourceId: self.resourceId,
                           memberList: self.collection
                       });

                       Radio.channel('main').trigger('show:dialog', sud);
                   }
                });
        }
    });
});