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
    'marionette'
], function($,_,Marionette){

    //Constructor
    var BaseManager = function(ModelType) {
        this.ModelType = ModelType;
    };

    //Fetch and return a promise
    BaseManager.prototype.fetchDeferred = function(entity, options){
        options = options || {};
        var defer = $.Deferred();
        options.cache = false;
        options.success = function(data){
            defer.resolve(data);
        };
        options.error = function(data){
            defer.resolve(undefined);
        };
        entity.fetch(options);
        return defer.promise();
    };

    //Read
    BaseManager.prototype.getEntity = function(id){
        var entity = new this.ModelType({id: id});
        return this.fetchDeferred(entity);
    };

    //Save
    BaseManager.prototype.saveEntity = function(entity){
        var defer = $.Deferred();
        entity.save({}, {
            wait: true,
            beforeSend: function(xhr){
                if(entity.hasOwnProperty('revision')){
                    xhr.setRequestHeader('If-Match', entity.revision);
                }
                return true;
            },
            success: function(model, response, options){
                if(typeof(model) !== 'undefined'){
                    defer.resolve(model);
                } else {
                    defer.resolve(response);
                }
            },
            error: function(model, response, options){
                defer.resolve({error: response});
            }
        });
        return defer.promise();
    };

    //Delete
    BaseManager.prototype.deleteEntity = function(entity){
        var defer = $.Deferred();
        entity.destroy({
            wait: true,
            beforeSend: function(xhr){
                if(entity.hasOwnProperty('revision')){
                    xhr.setRequestHeader('If-Match', entity.revision);
                }
                return true;
            },
            success: function(model, response){
                defer.resolve(model);
            },
            error: function(model, response){
                defer.resolve(undefined);
            }
        });
        return defer.promise();
    };

    return BaseManager

});
