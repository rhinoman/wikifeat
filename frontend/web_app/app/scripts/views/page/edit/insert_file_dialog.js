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

'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'entities/wiki/files',
    'text!templates/page/insert_file_dialog.html'
], function($,_,Marionette,Radio,
            FileCollection,InsertFileDialogTemplate){

    var userChannel = Radio.channel('userManager');
    var wikiChannel = Radio.channel('wikiManager');

    return Marionette.ItemView.extend({
        id: "insertFileDialog",
        template: _.template(InsertFileDialogTemplate),
        events: {
            'click #insertButton': function(){$('#theSubmit').trigger('click')},
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
                    wikiChannel.request('get:allFileList', self.wikiId).done(function(data){
                        dFileList.resolve(data);
                    });
                }
            });
        },

        optionTemplate: function(){
            return _.template("<option value='<%= id %>'><%= name %></option>");
        },

        onRender: function(){
            this.$("#insertFileModal").modal();
            var self = this;
            this.fileList.done(function(data){
                var select = self.$("select#fileSelect");
                select.html('<option value="0">Select File...</option>');
                if(typeof data !== 'undefined'){
                    _.each(data.models, function(file){
                        select.append(self.optionTemplate()({id: file.get('id'), name: file.get('name')}));
                    }, self);
                }
            });
        },

        submitForm: function(event){
            event.preventDefault();
            const fileId = this.$("select#fileSelect").val();
            var theUrl = "http://";
            this.fileList.done(function(fc){
                const fileModel = fc.get(fileId);
                if(typeof fileModel !== 'undefined' && fileModel !== null){
                    theUrl = fileModel.getDownloadLink();
                }
            });
            this.callback(theUrl);
            this.$("#insertFileModal").modal('hide');
        }
    });

});