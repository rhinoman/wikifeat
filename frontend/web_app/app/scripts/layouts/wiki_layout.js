/**
 * Created by jcadam on 2/4/15.
 */
'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'text!templates/wiki/wiki_layout.html'
], function($,_,Marionette,WikiTemplate){

    return Marionette.LayoutView.extend({
        template: _.template(WikiTemplate),

        regions: {
            toolbarRegion: ".toolbar",
            breadcrumbRegion: ".breadcrumbs",
            pageViewRegion: ".page-view"
        },

        onRender: function(){}

    });

});
