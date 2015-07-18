/**
 * Created by jcadam on 12/12/14.
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'entities/wiki/wiki',
    'entities/base_collection'
], function($,_,Backbone,WikiModel,BaseCollection){

    //Constructor
    function WikiCollection(models, options){
        BaseCollection.call(this, "wiki", WikiModel, models, options)
    }

    WikiCollection.prototype = Object.create(BaseCollection.prototype);

    WikiCollection.prototype.url = "/api/v1/wikis";
    WikiCollection.prototype.comparator = "name";

    return WikiCollection
});
