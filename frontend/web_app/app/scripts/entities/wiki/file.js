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

    FileModel.prototype.isImageFile = function(){
        var fileData = this.getFileData();
        if(fileData !== null){
            const mimeType = fileData.content_type;
            if(mimeType.substring(0,6) === 'image/'){
                return true;
            }
        }
        return false;
    };

    FileModel.prototype.getDownloadLink = function(){
        return this.getContentLink() + "&download=true";
    };

    FileModel.prototype.getContentLink = function(){
        var filename = this.getFilename();
        return this.url + "/content?attName=" + filename;
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