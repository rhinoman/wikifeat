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