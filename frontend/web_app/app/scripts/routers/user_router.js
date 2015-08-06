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
            "users/manage": "manageUsers",
            "users/account": "accountSettings"
        },
        controller: {
            manageUsers: function(){
                Radio.channel('user').trigger('manage:users');
            },
            accountSettings: function(){
                Radio.channel('user').trigger('user:accountSettings');
            }
        }
    });

});
