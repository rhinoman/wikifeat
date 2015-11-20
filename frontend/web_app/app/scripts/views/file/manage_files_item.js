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
 *  Individual File view in file table (i.e., one 'row')
 */
define([
    'jquery',
    'underscore',
    'marionette',
    'bootstrap',
    'backbone.radio',
    'entities/wiki/file',
    'views/file/edit_file_dialog',
    'views/main/confirm_dialog',
    'text!templates/main/alert.html',
    'text!templates/file/manage_files_item.html'
], function($,_,Marionette,Bootstrap,Radio,
            FileModel,EditFileDialogView,ConfirmDialog,
            AlertTemplate,ManageFilesItemTemplate){

    return Marionette.ItemView.extend({
        id: 'manage-files-item',
        tagName: 'tr',
        template: ManageFilesItemTemplate,
        model:FileModel,
        events: {
            'click td#name a': 'editFile',
            'click button#deleteButton': 'deleteFile'
        },

        initialize: function(options){
            //this.model.on('sync', this.render, this);
            this.model.on('change', this.render, this);
        },

        /**
         * Display the edit file dialog window
         * @param event
         */
        editFile: function(event){
            event.preventDefault();
            //var self = this;
            var editFileDialog = new EditFileDialogView({model: this.model});
            Radio.channel('main').trigger('show:dialog', editFileDialog);
        },

        /**
         * Delete file from the database
         * @param event
         */
        deleteFile: function(event){
            var self = this;
            var confirmCallback = function(){
                Radio.channel('wikiManager').request('delete:file', self.model)
                    .done(function(response){
                        if(typeof response === 'undefined'){
                            $("#alertBox").html(AlertTemplate);
                            var av = $("div#alertView");
                            av.css("display","block");
                            av.addClass("alert-danger");
                            av.append("Could not delete file");
                        }
                    });
            };

            var confirmDialog = new ConfirmDialog({
                message: 'Are you sure you wish to delete ' + this.model.get('name') +
                    '? This action is irreversible.',
                confirmCallback: confirmCallback
            });

            Radio.channel('main')
                .trigger('show:dialog', confirmDialog);
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.$("td#name").html(
                    '<a href="#">' + this.model.get('name') + '</a>'
                );
                var filename = this.model.getFilename();
                var fileObj = this.model.getFileData();
                if(fileObj !== null){
                    var downloadLink = this.model.getDownloadLink();
                    this.$("td#content").html(filename + '&nbsp;' +
                    '<a class="download-link" href="' + downloadLink + '">' +
                    '<span class="glyphicon glyphicon-download-alt"</span></a>');
                    //this.$("td#file_type").html(fileObj['content_type']);
                    this.$("td#size").html(this.model.prettyFileSize());
                }
                if(!this.model.isDeletable){
                    this.$("button#deleteButton").css("display", "none");
                }
                this.$('[data-toggle="tooltip"]').tooltip();
            }
        },

        onDestroy: function(){
            console.log("I have been destroyed, view: " + this.cid);
        }

    });


});
