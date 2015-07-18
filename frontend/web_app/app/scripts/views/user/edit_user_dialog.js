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
    'backbone.stickit',
    'entities/user/user',
    'bootstrap',
    'views/main/alert',
    'text!templates/user/edit_user_dialog.html'
], function($,_,Marionette,Radio,Stickit,
            UserModel,Bootstrap,AlertView,EditUserTemplate){

    return Marionette.ItemView.extend({
        id: "edit-user-dialog",
        template: _.template(EditUserTemplate),
        model: UserModel,
        bindings: {
            '#inputUsername':{
                observe: 'name'
            },
            '#inputLastName':{
                observe: 'userPublic',
                updateModel: false,
                onGet: function(data){
                    return data.lastName;
                }
            },
            '#inputFirstName': {
                observe: 'userPublic',
                updateModel: false,
                onGet: function(data){
                    return data.firstName;
                }
            },
            '#inputMiddleName': {
                observe: 'userPublic',
                updateModel: false,
                onGet: function(data){
                    return data.middleName;
                }
            },
            '#inputTitle': {
                observe: 'userPublic',
                updateModel: false,
                onGet: function(data){
                    return data.title;
                }
            },
            '#inputEmail': {
                observe: 'userPublic',
                updateModel: false,
                onGet: function(data){
                    return data.contactInfo.email;
                }
            }
        },
        events: {
            'shown.bs.modal' : 'showModal',
            'click #saveButton' : function(){$('#theSubmit').trigger('click')},
            'submit form': 'submitForm'
        },

        initialize: function(options){
            this.model.on('invalid', this.showError, this);

        },

        showModal: function(event){
           this.$("#inputUsername").focus();
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

        submitForm: function(event){
            event.preventDefault();
            var password = this.$("#inputPassword").val();
            var verifyPassword = this.$("#inputVerifyPassword").val();
            if (password === ""){
                this.model.unset('password');
                this.model.unset('verifyPassword');
            } else {
                this.model.set('password', password);
                this.model.set('verifyPassword', verifyPassword);
            }
            var publicAttrs = this.model.get('userPublic');
            publicAttrs.lastName = this.$("#inputLastName").val();
            publicAttrs.firstName = this.$("#inputFirstName").val();
            publicAttrs.middleName = this.$("#inputMiddleName").val();
            publicAttrs.contactInfo.email = this.$("#inputEmail").val();
            publicAttrs.title = this.$("#inputTitle").val();
            var self = this;
            Radio.channel('userManager').request('save:user', this.model)
                .done(function(response){
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
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                if(!this.model.hasOwnProperty('id') || this.model.id === ""){
                    this.$("#editUserTitle").html("Create User");
                } else {
                    this.$("#editUserTitle").html("Edit User");
                }
                this.stickit();
                if(this.model.get('name') !== ''){
                    this.model.set('update', true);
                    this.$("#inputUsername").prop('disabled','true');
                    this.$('#passwordRow').css('display','none');
                    this.$('#inputPassword').removeAttr('required');
                    this.$('#inputVerifyPassword').removeAttr('required');
                } else {
                    this.model.set('update', false);
                }
                this.$("#editUserModal").modal();
            }
        },

        onClose: function(){
            this.unstickit();
        }

    });

});