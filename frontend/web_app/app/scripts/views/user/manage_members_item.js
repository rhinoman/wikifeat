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

define([
    'jquery',
    'underscore',
    'marionette',
    'bootstrap',
    'backbone.radio',
    'entities/user/user',
    'views/main/confirm_dialog',
    'text!templates/user/manage_members_item.html'
], function($,_,Marionette,Bootstrap,Radio,
            UserModel,ConfirmDialog,
            ManageMembersItemTemplate){

    return Marionette.ItemView.extend({
        id: 'manage-members-item',
        tagName: 'tr',
        template: ManageMembersItemTemplate,
        model: UserModel,

        events: {
            'click #removeMemberButton': 'removeMember',
            'click #toggleWriteButton': 'toggleWrite',
            'click #toggleAdminButton': 'toggleAdmin'
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('resourceId')){
                this.resourceId = options.resourceId;
                this.resource = 'wiki_' + this.resourceId;
            }
            this.model.on('change', this.render, this);
        },

        //This function removes a single role and returns a promise
        removeRole: function(role){
            var rmreq = $.Deferred();
            var self = this;
            var theRole = self.resource + ":" + role;
            if(!_.contains(self.model.get('roles'), theRole)){
                rmreq.resolve("NONE");
                return rmreq.promise();
            }
            Radio.channel('userManager').request('revoke:role',
                self.model, 'wiki', self.resourceId, role)
                .done(function(response){
                    if(typeof response !== 'undefined') {
                        var roles = self.model.get('roles');
                        self.model.set('roles', _.without(roles, theRole));
                        rmreq.resolve("SUCCESS");
                    } else {
                        rmreq.resolve("ERROR");
                    }
                });
            return rmreq.promise();
        },

        toggleWrite: function(event){
            event.preventDefault();
            var roles = this.model.get('roles') || [];
            var writeRole = this.resource + ":write";
            var reqProm;
            if(_.contains(roles, writeRole)){
                //Remove the write role
                reqProm = this.removeRole("write");
            } else {
                //Add the write role
                reqProm = this.addRole("write");
            }
            var self = this;
            $.when(reqProm).done(function(resp){
                if(resp !== "ERROR"){

                } else {
                    self.triggerMethod("error", "Could not set Write Access");
                }
            });
        },

        toggleAdmin: function(event){
            event.preventDefault();
            var roles = this.model.get('roles') || [];
            var adminRole = this.resource + ":admin";
            var reqProm;
            if(_.contains(roles, adminRole)){
                //Remove the admin role
                reqProm = this.removeRole("admin");
            } else {
                //Add the admin role
                reqProm = this.addRole("admin");
            }
            var self = this;
            $.when(reqProm).done(function(resp){
                if(resp !== "ERROR"){
                }
                else{
                    self.triggerMethod("error", "Could not set Admin Access");
                }
            });
        },

        disableControls: function(){
            this.$("div.table-buttons-container button").attr("disabled", true);
        },

        //This function adds a single role and returns a promise
        addRole: function(role){
            var addreq = $.Deferred();
            var self = this;
            var theRole = this.resource + ":" + role;
            if(_.contains(this.model.get('roles'), theRole)){
                addreq.resolve("ALREADY_EXISTS");
                return addreq.promise();
            }
            Radio.channel('userManager').request('grant:role',
                this.model, 'wiki', this.resourceId, role)
                .done(function(response){
                    if(typeof response !== 'undefined'){
                        self.model.get('roles').push(theRole);
                        self.model.trigger('change');
                        addreq.resolve("SUCCESS");
                    } else {
                        addreq.resolve("ERROR");
                    }
                });
            return addreq.promise();
        },

        //removeMember removes all roles associated with a resource
        removeMember: function(event){
            var self = this;
            event.preventDefault();
            if(this.resource === ""){
                return;
            }
            var roles = this.model.get('roles') || [];
            //Remove the read, write, and admin roles
            var readPromise = this.removeRole("read");
            var writePromise = this.removeRole("write");
            var adminPromise = this.removeRole("admin");
            $.when(readPromise, writePromise, adminPromise).done(function(r1,r2,r3){
                //If any revoke role attempt resulted in an "ERROR", the user still has
                //at least one role for this resource
                if(r1 !== "ERROR" && r2 !== "ERROR" && r3 !== "ERROR"){
                    self.triggerMethod('remove:member', self.model);
                } else {
                    self.triggerMethod("error", "Could not remove member");
                }
            });
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.$("td#username").html(this.model.get('name'));
                var userPublic = this.model.get('userPublic');
                this.$("td#last_name").html(userPublic.lastName);
                this.$("td#first_name").html(userPublic.firstName);
                if(this.model.isResourceWriter(this.resource)){
                    this.$("td#write_access").html(
                        '<span class="glyphicon glyphicon-ok"></span>'
                    );
                }
                if(this.model.isResourceAdmin(this.resource)){
                    this.$("td#admin_access").html(
                        '<span class="glyphicon glyphicon-ok"></span>'
                    );
                }
                this.$('[data-toggle="tooltip"]').tooltip();
            }
        },
        onClose: function(){}
    });

});