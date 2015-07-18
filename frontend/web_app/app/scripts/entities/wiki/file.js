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
    'backbone',
    'entities/base_model'
], function($,_,Backbone,BaseModel){
    'use strict';
    //File model constructor
    function FileModel(data, options){
        options = options || {};
        if(options.hasOwnProperty("wikiId")){
            this.wikiId = options.wikiId;
        } else if(typeof this.wikiId === 'undefined'){
            this.wikiId = "";
        }
        BaseModel.call(this, "file", data, options);
    }

    FileModel.prototype = Object.create(BaseModel.prototype);

    FileModel.prototype.urlRoot = function(){
        return "/api/v1/wikis/" + this.wikiId + "/files";
    };

    FileModel.prototype.defaults = {
        "name": "Unnamed",
        "description": "None"
    };

    //Need to remove the _attachments structure before saving
    FileModel.prototype.save = function(attrs, options){
        options = options || {};
        attrs = _.extend({}, _.clone(this.attributes), attrs);
        /*if(attrs.hasOwnProperty("_attachments")){
            delete attrs._attachments;
        }*/
        options.attrs = attrs;

        return BaseModel.prototype.save.call(this, attrs, options);
    };

    FileModel.prototype.getFilename = function(){
        var attach = this.get('_attachments');
        if(typeof attach === 'undefined' || attach === null){
            return undefined;
        } else {
            return Object.keys(attach)[0];
        }
    };

    FileModel.prototype.getFileData = function(){
        var filename = this.getFilename();
        if(typeof filename !== 'undefined'){
            return this.get('_attachments')[filename];
        }
        return null;
    };

    FileModel.prototype.getDownloadLink = function(){
        var filename = this.getFilename();
        return this.url + "/content?attName=" + filename + "&download=true";
    };

    //Takes a FormData object and uploads it as the File's 'content'
    FileModel.prototype.uploadContent = function(formData){
        var defer = $.Deferred();
        var self = this;
        $.ajax({
            url: this.url +"/content",
            type: "POST",
            beforeSend: function(request){
                request.setRequestHeader("If-Match", self.revision);
            },
            data: formData,
            processData: false,
            contentType: false
        }).done(function(response){
            defer.resolve(response);
        }).fail(function(response){
            defer.resolve(undefined);
        });
        return defer.promise();
    };

    //Give a nicely formatted string for file size
    FileModel.prototype.prettyFileSize = function(){
        var fileData = this.getFileData();
        if(fileData !== null){
            var size = fileData['length'];
            //Default is Bytes
            var divisor = 1;
            var units = 'B';
            if(size > 1e9){
                //We have Gigabytes of data, man
                divisor = 1e9;
                units = 'G';
            } else if(size > 1e6){
                //Megabytes
                divisor = 1e6;
                units = 'M';
            } else if(size > 1e3){
                //Kilobytes
                divisor = 1e3;
                units = 'K';
            }
            var adjSize = size / divisor;
            if(adjSize % 1 !== 0){
                adjSize = adjSize.toFixed(1)
            }
            return adjSize + units
        }
        return null;
    };

    FileModel.prototype.validate = function(attrs, options){
        var errors = {};
        if (!attrs.name) {
            errors.name = "can't be blank!";
        }
        if (!attrs.description){
            errors.description = "can't be blank!";
        }
        if(!$.isEmptyObject(errors)){
            return errors;
        }
    };

    return FileModel;

});