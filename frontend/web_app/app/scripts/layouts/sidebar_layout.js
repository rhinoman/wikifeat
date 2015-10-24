/**
 * Created by jcadam on 12/13/14.
 */

define([
    'jquery',
    'underscore',
    'marionette',
    'text!templates/sidebar/sidebar_layout.html'
], function($,_,Marionette, SidebarTemplate){
    'use strict';

    return Marionette.LayoutView.extend({
        template: _.template(SidebarTemplate),
        tagName: "nav",
        id: 'sidebar-view',
        className: 'nav nav-sidebar',
        attributes: {
            "role": "navigation"
        },
        regions: {
            logoRegion: "#logo",
            userMenuRegion: "#userMenu",
            adminMenuRegion: "#adminMenu",
            wikiListRegion: "#wikiList"
        },

        onRender: function(){}
    });

});
