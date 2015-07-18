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
    'entities/user/user',
    'text!templates/user/select_users_item.html'
], function($,_,Marionette,Radio,UserModel,
            SelectUsersItemTemplate){

    return Marionette.ItemView.extend({
        id: 'select-users-item',
        tagName: 'tr',
        template: _.template(SelectUsersItemTemplate),
        model: UserModel,
        resource: "",

        events: {
            'click #addButton': 'addMember'
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty("resourceId")){
                this.resourceId = options.resourceId;
                this.resource = 'wiki_' + this.resourceId;
            }
        },

        addMember: function(event){
            event.preventDefault();
            var self = this;
            if(this.resource === ""){
                //do nothing
                return;
            }
            Radio.channel('userManager').request('grant:role',
                this.model, 'wiki', this.resourceId, 'read')
                .done(function(response){
                    if(typeof response !== 'undefined' && response.get('success') === true){
                        var roles = self.model.get("roles") || [];
                        roles.push(self.resource + ":read");
                        self.model.set("roles", roles);
                        self.model.trigger('change');
                        self.$("td.action-cell").html("Member");
                        self.triggerMethod('add:member', self.model);
                    } else {
                        self.triggerMethod('error', "Could not add member");
                        console.log("Could not set role.")
                    }
                })
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.$("td#username").html(this.model.get("name"));
                var userPublic = this.model.get("userPublic");
                this.$("td#lastName").html(userPublic.lastName);
                this.$("td#firstName").html(userPublic.firstName);
                if(this.model.isResourceMember(this.resource)){
                    this.$("td.action-cell").html("Member");
                }
            }
        }
    });

});
