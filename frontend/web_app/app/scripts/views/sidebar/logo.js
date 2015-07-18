/**
 * Created by jcadam on 12/13/14.
 */

define([
    'jquery',
    'underscore',
    'marionette',
    'text!templates/sidebar/logo.html'
], function($,_,Marionette,LogoTemplate){
    'use strict';

    return Marionette.ItemView.extend({
        id: "logo-view",
        initialize: function(){
            console.log('initializing Logo view');
        },

        template: _.template(LogoTemplate),

        /* on render callback */
        onRender: function(){}
    });

});
