/**
 * Created by jcadam on 2/26/15.
 */

'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'backbone.stickit',
    'entities/wiki/wiki',
    'entities/wiki/page',
    'text!templates/page/home_placeholder.html'
], function($,_,Marionette,Radio,Stickit,
            WikiModel,PageModel,PlaceholderTemplate){

    return Marionette.ItemView.extend({
        id: 'placeholder-page-view',
        template: _.template(PlaceholderTemplate),
        model: WikiModel,
        bindings: {
            '#wiki-model-name': {
                observe: 'name'
            },
            '#wiki-description': {
                observe: 'description'
            }
        },
        events: {
            'click .page-edit-button': 'createHomePage'
        },

        initialize: function(options){

        },

        //Create home page
        createHomePage: function(event){
            event.preventDefault();
            Radio.channel('page').trigger('create:page', this.model,{homePage: true})
        },

        /* on render callback */
        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.stickit();
                if(this.model.canCreatePage){
                    $(this.el).prepend(
                        '<button class="btn btn-default btn-lg page-control-button page-edit-button">' +
                        '<span class="glyphicon glyphicon-plus-sign"></span>&nbsp;Create Home Page</button>');
                }
            }
        },

        onClose: function(){
            this.unstickit();
        }
    });
});
