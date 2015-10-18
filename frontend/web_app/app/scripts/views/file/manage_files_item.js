/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.*
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

        editFile: function(event){
            event.preventDefault();
            //var self = this;
            var editFileDialog = new EditFileDialogView({model: this.model});
            Radio.channel('main').trigger('show:dialog', editFileDialog);
            /*Radio.channel('wikiManager').request('get:file', this.model.id,
                this.model.wikiId).done(function(model){
                    if(typeof model !== 'undefined'){
                        self.model = model;
                        //self.model.on('sync', self.render, self);
                        var editFileDialog =
                            new EditFileDialogView({model: self.model});
                        Radio.channel('main')
                            .trigger('show:dialog', editFileDialog);
                    }
                });*/
        },

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
