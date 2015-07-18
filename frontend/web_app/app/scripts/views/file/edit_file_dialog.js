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
            //We do this in two steps
            //First save the file record (sans actual file)
            Radio.channel('wikiManager').request('save:file', this.model)
                .done(function(response){
                    if(response.hasOwnProperty('error')){
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
