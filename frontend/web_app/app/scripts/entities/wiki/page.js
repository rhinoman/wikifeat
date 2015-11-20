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
    //Page model constructor
    function PageModel(data, options){
        options = options || {};
        if (options.hasOwnProperty("wikiId")){
            this.wikiId = options.wikiId;
        } else {
            this.wikiId = "";
        }
        BaseModel.call(this, "page", data, options);
    }

    PageModel.prototype = Object.create(BaseModel.prototype);

    PageModel.prototype.urlRoot = function(){
        return "/api/v1/wikis/" + this.wikiId + "/pages";
    };

    PageModel.prototype.defaults = {
        content: {
            formatted: "",
            raw: ""
        },
        editor: "Unknown",
        owner: "",
        parent: "",
        comments_disabled: false,
        title: "Untitled",
        type: "page"
    };

    //input validation function
    PageModel.prototype.validate = function(attrs, options) {
        var errors = {};
        if (!attrs.title) {
            errors.title = "can't be blank";
        } else if ((attrs.title.length) > 128){
            errors.title = "is too long";
        }
        if(!$.isEmptyObject(errors)){
            return errors;
        }
    };

    return PageModel;

});
