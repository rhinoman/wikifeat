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
    'bootstrap',
    'backbone.radio',
    'entities/user/user',
    'views/user/change_password_dialog',
    'views/user/edit_user_dialog',
    'views/main/confirm_dialog',
    'text!templates/user/manage_users_item.html'
], function($,_,Marionette,Bootstrap,Radio,
            UserModel, ChangePasswordDialogView,
            EditUserDialogView, ConfirmDialog,
            ManageUsersItemTemplate){

    return Marionette.ItemView.extend({
        id: 'manage-users-item',
        tagName: 'tr',
        template: ManageUsersItemTemplate,
        model: UserModel,
        events: {
            'click td#username a': 'editUser',
            'click td#user_actions button#changePasswordButton': 'changePassword',
            'click td#user_actions button#toggleAdminButton' : 'toggleAdmin',
            'click td#user_actions button#deleteUserButton': 'deleteUser'
        },

        initialize: function(options){
            this.model.on('change', this.render, this);
        },

        editUser: function(event){
            event.preventDefault();
            var self = this;
            Radio.channel('userManager').request('get:user', this.model.id)
                .done(function(model){
                    if(typeof model !== 'undefined'){
                        self.model = model;
                        var editUserDialog =
                            new EditUserDialogView({model: self.model});
                        Radio.channel('main')
                            .trigger('show:dialog', editUserDialog);
                    }
                });
        },

        toggleAdmin: function(event){
            event.preventDefault();
            var self = this;
            if(!this.model.isAdmin()){
                Radio.channel('userManager').request('grant:role',
                    this.model, 'main', '', 'admin')
                    .done(function(response){
                        if(typeof response !== 'undefined' &&
                            !response.hasOwnProperty('error') &&
                            response.get('success') === true){
                            var roles = self.model.get("roles") || [];
                            self.triggerMethod('admin:enabled');
                            roles.push('admin');
                            self.model.set("roles", _.uniq(roles));
                            self.model.trigger('change');
                        } else {
                            self.triggerMethod('admin:error');
                        }
                    });
            } else {
                Radio.channel('userManager').request('revoke:role',
                    this.model, 'main', '', 'admin')
                    .done(function(response){
                        if(typeof response !== 'undefined' && response.get('success') === true){
                            self.triggerMethod('admin:disabled');
                            self.trigger('change');
                            var roles = self.model.get("roles");
                            var idx = _.indexOf(roles, "admin");
                            if (idx > -1){
                                roles.splice(idx, 1);
                                self.model.set("roles", roles);
                                self.model.trigger('change');
                            }
                        } else {
                            self.triggerMethod('admin:error');
                        }
                    });
            }
        },

        deleteUser: function(event){
            var confirmCallback = function(){
                Radio.channel('userManager').request('delete:user', self.model)
                    .done(function(response){
                        if(typeof response === 'undefined') {
                            var av = $("div#alertView");
                            av.css("display", "block");
                            av.addClass("alert-danger");
                            av.append("Could not delete user");
                        }
                    });
            };

            var confirmDialog = new ConfirmDialog({
                message: 'Are you sure you wish to delete ' + this.model.get('id') +
                '?  This action is irreversible.',
                confirmCallback: confirmCallback
            });

            Radio.channel('main')
                .trigger('show:dialog', confirmDialog);
            var self = this;
        },

        changePassword: function(event){
            Radio.channel('userManager').request('get:user', this.model.id)
                .done(function(model){
                    if(typeof model !== 'undefined'){
                        self.model = model;
                        var cpv = new ChangePasswordDialogView({model: self.model});
                        Radio.channel('main')
                            .trigger('show:dialog', cpv);
                    }
                });
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.$("td#username").html(
                    '<a href="#">' + this.model.get('name') + '</a>'
                );
                var userPublic = this.model.get('userPublic');
                this.$("td#last_name").html(userPublic.lastName);
                this.$("td#first_name").html(userPublic.firstName);
                if(this.model.isAdmin()){
                    this.$("td#admin_access").html(
                        '<span class="glyphicon glyphicon-ok"></span>'
                    );
                }
                if(this.model.isMaster()){
                    this.$("#toggleAdminButton").css('display', 'none');
                    this.$("#deleteUserButton").css('display', 'none');
                    this.$("td#master_access").html(
                        '<span class="glyphicon glyphicon-ok"></span>'
                    );
                }
                this.$('[data-toggle="tooltip"]').tooltip();
            }
        },

        onClose: function(){
        }
    })
});

