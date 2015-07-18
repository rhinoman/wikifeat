/**
 * Created by jcadam on 12/12/14.
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'marionette'
], function($, _, Backbone, Marionette){

    return Marionette.Controller.extend({
        showHome: function(){
            $.ajax({
                url: "app/home"
            }).done(function(response){
                if(response.hasOwnProperty('home')){
                    Backbone.history.navigate(response.home, {trigger: true});
                }
            });
        }
    });
});
