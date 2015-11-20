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
    'entities/wiki/file',
    'bootstrap',
    'views/main/alert',
    'text!templates/file/edit_file_dialog.html'
], function($,_,Marionette,Radio,Stickit,
            FileModel,Bootstrap,AlertView,
            EditFileTemplate){

    return Marionette.ItemView.extend({
        id: "edit-file-dialog",
        template: _.template(EditFileTemplate),
        model: FileModel,
        bindings: {
            '#inputName':{
                observe: 'name'
            },
            '#inputDescription':{
                observe: 'description'
            }
        },
        events:{
            'change #file-data': 'updateFileSelect',
            'submit form': 'submitForm',
            'click #saveButton' : function(){$('#theSubmit').trigger('click')}
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                if(!this.model.hasOwnProperty('id') || this.model.id === ""){
                    this.$("#editFileTitle").html("Create File");
                    this.model.isEditable = true;
                } else {
                    this.$("#editFileTitle").html("Edit File");
                }
                this.stickit();
                var filename = this.model.getFilename();
                if(filename !== null){
                    this.$("#fileNameDisplay").html(filename);
                }
                if (!this.model.isEditable){
                    this.$("button#saveButton").css("display","none");
                    this.$("input").attr('disabled',true);
                }
                this.$("#editFileModal").modal();
            }
        },

        updateFileSelect: function(event){
            var filename = $(event.currentTarget).val().replace("C:\\fakepath\\","");
            this.$("#fileNameDisplay").html(filename);
        },

        submitForm: function(event){
            event.preventDefault();
            var self = this;
            //disable the submit button
            this.$("button#saveButton").attr("disabled", "disabled");
            //We do this in two steps
            //First save the file record (sans actual file)
            Radio.channel('wikiManager').request('save:file', this.model)
                .done(function(response){
                    if(response.hasOwnProperty('error')){
                        self.$("button#saveButton").removeAttr("disabled");
                        var error = {};
                        if(response.error.status === 400) {
                            error.serverError = "Invalid Request";
                            self.showError(self.model, error);
                        } else if (response.error.status === 403){
                            error.serverError = "You are not authorized to perform this action";
                            self.showError(self.model, error);
                        } else if (response.error.status === 409){
                            error.serverError = "This File conflicts with another";
                            self.showError(self.model, error);
                        } else {
                            error.serverError = "Server Error.  Please try again.";
                            self.showError(self.model, error);
                        }
                    } else {
                        var file = $("#file-data").val();
                        if (typeof file !== 'undefined' && file !== "") {
                            //Now we upload the file itself.
                            var formData = new FormData();
                            var input = document.getElementById("file-data");
                            formData.append('file-data', input.files[0]);
                            self.model.uploadContent(formData).done(function(response){
                                if(typeof response === 'undefined'){
                                    var error = {};
                                    error.serverError = "could not upload file";
                                    self.showError(self.model, error);
                                    self.$("button#saveButton").removeAttr("disabled");
                                } else {
                                    self.model.fetch();
                                    self.model.trigger('change');
                                    self.$("#cancelButton").trigger('click');
                                }
                            });
                        } else {
                            self.model.trigger('change');
                            self.$("#cancelButton").trigger('click');
                        }
                    }
                });

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

        onClose: function(){
            this.unstickit();
        }
    });

});
