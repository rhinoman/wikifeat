/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.
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
    'bootstrap',
    'views/main/alert',
    'text!templates/user/change_password_dialog.html'
], function($,_,Marionette,Radio,UserModel,
            Bootstrap,AlertView,ChangePasswordTemplate){

    return Marionette.ItemView.extend({
        id: "change-password-dialog",
        template: _.template(ChangePasswordTemplate),
        model: UserModel,
        events: {
            'shown.bs.modal': 'showModal',
            'click #submitButton': function(){$('#theSubmit').trigger('click')},
            'submit form': 'submitForm'
        },

        initialize: function(options){
            this.model.on('invalid', this.showError, this);
        },

        showModal: function(event){
            this.$("#inputOldPassword").focus();
        },

        //Execute change of password
        submitForm: function(event){
            event.preventDefault();
            var oldPassword = this.$("#inputOldPassword").val();
            var newPassword = this.$("#inputNewPassword").val();
            var verifyPassword = this.$("#inputConfirmNewPassword").val();
            var error = {};
            if(newPassword !== verifyPassword){
                error.passwordError = "Passwords don't match!";
                this.showError(this.model, error);
            } else {
                var self = this;
                Radio.channel('userManager').request('change:password',
                    this.model, oldPassword, newPassword).done(function(response){
                        if(response.hasOwnProperty('error')){
                            var error = {};
                            if(response.error.status === 400){
                                error.serverError = "Invalid Request";
                                self.showError(self.model, error);
                            }
                            else if(response.error.status === 409){
                                error.serverError = "This Username conflicts with another user";
                                self.showError(self.model, error);
                            } else {
                                error.serverError = "Server Error.  Please try again later";
                                self.showError(self.model, error);
                            }
                        } else {
                            self.$("#cancelButton").trigger('click');
                        }
                    });
            }
        },

        showError: function(model, error){
            var alertText = 'Please correct the following errors: <ul id="error_list">';

            for (var property in error){
                if (error.hasOwnProperty(property)){
                    alertText += "<li>" + error[property] + "</li>"
                }
            }
            alertText += "</ul>";
            var alertView = new AlertView({
                el: $("#alertBox"),
                alertType: 'alert-danger',
                alertMessage: alertText
            });
            alertView.render();
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                var self = this;
                var currentUser = Radio.channel('userManager').request('get:currentUser')
                    .done(function(curUser){
                        if(curUser.isAdmin() && self.model.id != curUser.id) {
                            self.$("form #inputOldPassword").removeAttr("required");
                            self.$("form #oldPasswordGroup").css("display", "none");
                        }
                    });
                this.$('#changePasswordModal').modal();
            }
        }
    });

});
