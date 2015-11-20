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

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'views/paginated_table_view',
    'views/file/manage_files_item',
    'entities/wiki/file',
    'views/file/edit_file_dialog',
    'text!templates/file/manage_files.html',
    'text!templates/main/alert.html'
], function($,_,Marionette,Radio,PaginatedTableView,
            ManageFilesItemView, FileModel,
            EditFileDialogView, ManageFilesTemplate, AlertTemplate){

    return PaginatedTableView.extend({
        className: "manage-files-view",
        template: _.template(ManageFilesTemplate),
        childView: ManageFilesItemView,
        childViewContainer: "#fileListContainer",
        additionalEvents: {
            'click #addFileButton' : 'addFile'
        },

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
            } else {
                this.wikiModel = null;
            }
        },

        /**
         * Displays the Add File dialog window
         * @param event
         */
        addFile: function(event){
            event.preventDefault();
            var file = new FileModel({},{wikiId: this.wikiModel.id});
            var editFileDialog = new EditFileDialogView({model: file});
            var self = this;
            this.listenToOnce(file, 'sync', function(data){
                self.collection.add(data);
            });
            Radio.channel('main').trigger('show:dialog', editFileDialog);
        },

        onRender: function(){
            if(!this.collection.isCreatable){
                this.$("#addFileButton").css('display','none');
            }
            PaginatedTableView.prototype.onRender.call(this);
        }

    });

});