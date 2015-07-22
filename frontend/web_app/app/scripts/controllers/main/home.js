/**
 * Created by jcadam on 12/12/14.
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'backbone.radio',
    'marionette'
], function($, _, Backbone, Radio, Marionette){

    var HomeController =  Marionette.Controller.extend({
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
    var homeController = new HomeController();
    var homeChannel = Radio.channel('home');

    homeChannel.on("show:home", function(){
        homeController.showHome();
    });

    return homeController;
});
