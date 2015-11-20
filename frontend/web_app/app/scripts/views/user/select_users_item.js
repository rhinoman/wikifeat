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
