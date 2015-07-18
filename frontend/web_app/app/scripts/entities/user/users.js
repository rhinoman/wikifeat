/**
 * Created by jcadam on 3/21/15.
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'entities/user/user',
    'entities/base_collection'
], function($,_,Backbone,UserModel,BaseCollection){

    return BaseCollection.extend({

        initialize: function(models, options){
            BaseCollection.prototype.initialize.call(this, "user", UserModel, models, options);
        },

        url: "/api/v1/users",

        /*url: function() {
            var params = {};
            if (this.resource !== ""){
                params.forResource = this.resource;
            }
            return this.urlRoot + '?' + $.param(params);
        },*/

        resource: "",

        comparator: "name",

        //Pagination state vals
        state: {
            firstPage: 1,
            pageSize: 25
        }

    });
});
