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
/**
 * Created by jcadam on 8/5/15.
 */

'use strict';
define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'bootstrap',
    'views/user/edit_user_dialog',
    'views/user/change_password_dialog',
    'views/user/edit_user_avatar',
    'views/main/error_dialog',
    'text!templates/user/account_settings.html',
    'entities/user/user',
    'entities/user/user_avatar',
    'entities/error'
], function($,_,Marionette,Radio,Bootstrap,EditUserDialog,
            ChangePasswordDialogView,EditUserAvatarView,
            ErrorDialog, AccountSettingsTemplate, UserModel,
            UserAvatarModel, ErrorModel){

    return Marionette.ItemView.extend({
        className: "account-settings-view",
        model: UserModel,
        template: _.template(AccountSettingsTemplate),
        events: {
            'click #editProfileButton': 'editProfile',
            'click #changePasswordButton': 'changePassword',
            'click #changePictureButton': 'changePicture'
        },

        initialize: function(){
            this.model.on('change', this.render, this);
            this.model.on('reset', this.render, this);
            this.model.on('sync', this.render, this);
        },

        editProfile: function(event){
            var editUserDialog = new EditUserDialog({model: this.model});
            Radio.channel('main').trigger('show:dialog', editUserDialog);
        },

        changePassword: function(event){
            var cpv = new ChangePasswordDialogView({model: this.model});
            Radio.channel('main').trigger('show:dialog', cpv);
        },

        changePicture: function(event){
            var self = this;
            Radio.channel('userManager').request('get:avatar', this.model.id)
                .done(function(response){
                    if(typeof response === 'undefined')    {
                        var errorDialog = new ErrorDialog({
                            model: new ErrorModel({
                                errorTitle: "Error loading avatar",
                                errorMessage: "Could not load avatar"
                            })
                        });
                        Radio.channel('main').trigger('show:dialog', errorDialog);
                    } else {
                        self.avatarModel = response;
                        self.avatarModel.listenToOnce(self.avatarModel, 'newImage', self.render);
                        var euv = new EditUserAvatarView({
                            model: self.avatarModel,
                            userModel: self.model});
                        Radio.channel('main').trigger('show:dialog', euv);
                    }
                });
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                var userPublic = this.model.get("userPublic");
                var fullName = userPublic.firstName + " " + userPublic.lastName;
                var title = userPublic.title || "User";
                var email = userPublic.contactInfo.email || "None";
                this.$("#pictureWrapper").html(this.model.getAvatar());
                this.$("#nameField").html(fullName);
                this.$("#userNameField").html('<span class="glyphicon glyphicon-user"></span>&nbsp;' +
                    this.model.get("name"));
                if (email !== "None") {
                    this.$("#emailField").html('<span class="glyphicon glyphicon-envelope"></span>&nbsp;' +
                        '<a href="mailto:' + email + '">' + email + '</a>');
                } else {
                    this.$("#emailField").html("");
                }
                this.$("#titleField").html('<span class="glyphicon glyphicon-briefcase"></span>&nbsp;' + title);

            }
        },

        onClose: function(){
        }

    });

});
