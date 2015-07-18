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
        var rev = options.xhr.getResponseHeader('Etag');
        if(typeof rev !== 'undefined'){
            this.revision = rev;
        }
        return response[this.entityName];
    };

    return BaseModel;

});
