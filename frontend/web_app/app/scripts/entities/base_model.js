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
