/**
 * Created by jcadam on 3/21/15.
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'marionette',
    'backbone.radio'
], function($,_,Backbone,Marionette,Radio){

    return Marionette.AppRouter.extend({
        appRoutes: {
            "users/manage": "manageUsers"
        },
        controller: {
            manageUsers: function(){
                Radio.channel('user').trigger('manage:users');
            }
        }
    });

});
