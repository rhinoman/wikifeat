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

'use strict';

define([
    'jquery',
    'underscore',
    'backbone',
    'entities/base_model'
], function($,_,Backbone,BaseModel){

    //contructor
    function UserAvatarModel(data, options){
        options = options || {};
        BaseModel.call(this, "avatar_record", data, options);
    }

    UserAvatarModel.prototype = Object.create(BaseModel.prototype);

    UserAvatarModel.prototype.urlRoot = function(){
        return "/api/v1/users/" + this.get('id') + "/avatar"
    };

    UserAvatarModel.prototype.idAttribute = "_id";

    UserAvatarModel.prototype.uploadContent = function(formData){
        var defer = $.Deferred();
        var self = this;
        $.ajax({
            url: this.url + "/image",
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

    return UserAvatarModel;

});