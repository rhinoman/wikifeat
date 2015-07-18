/**
 * Created by jcadam on 12/9/14.
 */

define([
    'jquery',
    'underscore',
    'marionette',
    'controllers/main/home'
], function($, _, Marionette,
            HomeController){
    'use strict';

    var hc = new HomeController();

    var router = Marionette.AppRouter.extend({
        initialize: function(){
            //sbc.showSideBar();
        },
        appRoutes: {
            ""  : "showHome"
        },
        controller: {
            showHome: function(){
                hc.showHome();
            }
        }
    });

    return router

});
