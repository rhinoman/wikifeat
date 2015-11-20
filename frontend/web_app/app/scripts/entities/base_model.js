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
 * Created by jcadam on 1/21/15.
 */

define([
    'jquery',
    'underscore',
    'backbone'
], function($,_,Backbone){

    function BaseModel(entityName, data, options){
        this.entityName = entityName;
        Backbone.Model.call(this, data, options);
    }

    BaseModel.prototype = Object.create(Backbone.Model.prototype);

    // Breaking this out as a separate function to facilitate
    // easier overloading :)
    BaseModel.prototype.parseLinks = function(links) {
        if (links.self !== null) {
            this.url = links.self.href;
        } else {
            console.log("response has null self link");
        }
        this.isEditable = links.hasOwnProperty('update');
        this.isDeletable = links.hasOwnProperty('delete');
    };

    //Parse model from hateoas/hal response
    BaseModel.prototype.parse = function(response, options) {
        if (response.hasOwnProperty("_links")){
            this.parseLinks(response._links);
            delete response._links;
        }
        //The revision may be in the Etag header -- for indiviudal records
        var rev = options.xhr.getResponseHeader('Etag');
        //...or in the document itself, for collections
        var docRev = response[this.entityName]._rev;
        if(typeof rev !== 'undefined' && rev !== null){
            this.revision = rev;
        } else if (typeof docRev !== 'undefined' && docRev !== null){
            this.revision = docRev;
        }
        return response[this.entityName];
    };

    return BaseModel;

});
