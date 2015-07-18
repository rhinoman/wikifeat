/**
 * Created by jcadam on 2/28/15.
 */

'use strict';
define([
    'jquery',
    'underscore',
    'marionette',
    'views/page/child_index_item',
    'text!templates/page/child_index.html'
], function($,_,Marionette,ChildIndexItemView,
            ChildIndexTemplate){

    return Marionette.CompositeView.extend({
        childView: ChildIndexItemView,
        template: _.template(ChildIndexTemplate),
        childViewContainer: "#childIndexListContainer",

        initialize: function(options){
            options = options || {};
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
                this.childViewOptions = {
                    wikiModel: this.wikiModel
                }
            }
        },

        onRender: function(){
            if(this.collection.length === 0){
                this.$el.find("#childIndexListContainer").append("<li>None</li>");
            }
        }
    });

});
