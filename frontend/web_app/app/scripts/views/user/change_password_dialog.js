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
