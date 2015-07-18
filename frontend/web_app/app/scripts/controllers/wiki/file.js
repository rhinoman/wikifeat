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
    'entities/wiki/file',
    'views/file/manage_files'
], function($,_,Marionette,Radio,FileModel,
            ManageFilesView){

    var fileChannel = Radio.channel('file');

    var FileController = Marionette.Controller.extend({

        manageFiles: function(wikiModel){
            var promise = Radio.channel('wikiManager')
                .request("get:allFileList", wikiModel.id);
            promise.done(function(fileList){
                if(typeof fileList === 'undefined'){
                    //TODO error display
                    console.log("Error loading file list");
                    return
                }
                var region = Radio.channel('wiki').request('get:pageRegion');
                region.show(new ManageFilesView({
                    collection: fileList,
                    wikiModel: wikiModel
                }));
                window.history.pushState('','', '/app/wikis/' + wikiModel.get('slug') + '/files');
            });
        }

    });

    var fileController = new FileController();

    fileChannel.on("manage:files", function(wikiModel){
        fileController.manageFiles(wikiModel);
    });

    return fileController;
});
