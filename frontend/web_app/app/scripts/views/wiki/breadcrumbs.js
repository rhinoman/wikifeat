/**
 * Created by jcadam on 2/28/15.
 */

'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'views/wiki/breadcrumb_item'
], function($,_,Marionette,BreadcrumbItemView){

    return Marionette.CollectionView.extend({
        id: "wiki-breadcrumb-view",
        tagName: "ol",
        className: "breadcrumb",
        childView: BreadcrumbItemView,
        initialize: function(options){},

        onClose:function(){
            this.unbind();
        }
    });

});
