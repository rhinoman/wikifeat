/**
 * Created by jcadam on 12/9/14.
 */

define([
    'jquery',
    'underscore',
    'marionette',
    'controllers/main/home'
], function($, _, Marionette,
            homeController){
    'use strict';

    var router = Marionette.AppRouter.extend({
        initialize: function(){
            //sbc.showSideBar();
        },
        appRoutes: {
            ""  : "showHome"
        },
        controller: {
            showHome: function(){
                homeController.showHome();
            }
        }
    });

    return router

});
