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
