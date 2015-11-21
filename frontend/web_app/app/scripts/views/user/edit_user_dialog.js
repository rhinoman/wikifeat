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
            if (this.model.get('update')){
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
                    this.$("#avatarImgWrapper").html(this.model.getAvatar());
                    this.$("#inputUsername").prop('disabled','true');
                    this.$('#passwordRow').css('display','none');
                    this.$('#inputPassword').removeAttr('required');
                    this.$('#inputVerifyPassword').removeAttr('required');
                } else {
                    //Ugly hackaround for Firefox's buggy autocomplete behavior
                    this.model.set('update', false);
                    this.$("#userNameGroup").removeClass('col-lg-10 col-md-10 col-sm-10');
                    this.$("#userNameGroup").addClass('col-lg-12 col-md-12 col-sm-12');
                    this.$("#avatarWrapper").css("display","none");
                    this.$('#inputPassword').focus(function(){
                        this.removeAttribute('readonly');
                    });
                }
                this.$("#editUserModal").modal();
            }
        },

        onClose: function(){
            this.unstickit();
        }

    });

});