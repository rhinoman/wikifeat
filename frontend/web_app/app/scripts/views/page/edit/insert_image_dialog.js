/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 * See the LICENSE.txt file distributed with this work for additional information
 * regarding copyright ownership.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *    * Redistributions of source code must retain the above copyright notice,
 *        this list of conditions and the following disclaimer.
 *    * Redistributions in binary form must reproduce the above copyright
 *        notice, this list of conditions and the following disclaimer in the
 *        documentation and/or other materials provided with the distribution.
 *    * Neither the name of Wikifeat nor the names of its contributors may be used
 *        to endorse or promote products derived from this software without
 *        specific prior written permission.
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
 * Created by jcadam on 1/22/16.
 */

'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'entities/wiki/files',
    'text!templates/page/insert_image_dialog.html'
], function($,_,Marionette,Radio,
            FileCollection,InsertImageDialogTemplate){

    var userChannel = Radio.channel('userManager');
    var wikiChannel = Radio.channel('wikiManager');

    return Marionette.ItemView.extend({
        id: "insertImageDialog",
        template: _.template(InsertImageDialogTemplate),
        events:{
            'click #insertButton': function(){$('#theSubmit').trigger('click')},
            'change input[type=radio][name=imageOption]': 'radioChange',
            'change select#fileSelect': 'fileSelected',
            'submit form': 'submitForm'
        },

        initialize: function(options){
            options = options || {};
            this.callback = options.callback || function(){};
            this.wikiId = options.wikiId || null;
            this.fileList = new FileCollection({}, {wikiId: this.wikiId});
            var dFileList = $.Deferred();
            dFileList.promise(this.fileList);
            this.currentUser = userChannel.request('get:currentUser');
            var self = this;
            this.currentUser.done(function(user){
                if(typeof user !== 'undefined'){
                    wikiChannel.request('get:imageFileList', self.wikiId).done(function(data){
                        dFileList.resolve(data);
                    });
                }
            })
        },

        optionTemplate: function(){
            return _.template("<option value='<%= id %>'><%= name %></option>");
        },

        onRender: function(){
            this.$("#insertImageModal").modal();
        },

        radioChange: function(event){
            //Get the checked radio button
            const linkMode = this.linkMode();
            if(linkMode === 'internal'){
                this.$("div#externalImageSelectContainer").hide();
                this.$("div.imgPreviewContainer").hide();
                this.$("div#internalImageSelectContainer").show();
                var self = this;
                this.fileList.done(function(data){
                    var select = self.$("select#fileSelect");
                    select.html('<option value="0">Select Image File...</option>');
                    if(typeof data !== 'undefined'){
                        _.each(data.models, function(file){
                            select.append(self.optionTemplate()({id: file.get('id'), name: file.get('name')}));
                        }, self);
                    }
                });
            } else if(linkMode === 'external'){
                this.$("div#internalImageSelectContainer").hide();
                this.$("div#externalImageSelectContainer").show();
            }
        },

        fileSelected: function(event){
            const selectBox = event.currentTarget;
            const fileId = $(selectBox).val();
            var self = this;
            this.fileList.done(function(fc){
                const fileModel = fc.get(fileId);
                const imgLink = fileModel.getContentLink();
                const imgTag = '<img src="'+ imgLink + '">';
                self.$('div.imgPreviewContainer').show();
                self.$(".imgPreviewBox").html(imgTag);
            });
        },

        submitForm: function(event){
            event.preventDefault();
            const linkMode = this.linkMode();
            var theUrl = "http://";
            if(linkMode === 'internal'){
                const fileId = this.$("select#fileSelect").val();
                this.fileList.done(function(fc){
                    const fileModel = fc.get(fileId);
                    if(typeof fileModel !== 'undefined' && fileModel !== null){
                        theUrl = fileModel.getContentLink();
                    }
                });
            } else if(linkMode === 'external'){
                theUrl = this.$("#externalUrlField").val();
            }
            this.callback(theUrl);
            this.$("#insertImageModal").modal('hide');
        },

        linkMode: function(){
            var checked = this.$('input[type=radio]:checked');
            if(checked.attr('id') === 'internalOption'){
                return "internal";
            } else if(checked.attr('id') === 'externalOption'){
                return "external";
            }
        }

    });

});
