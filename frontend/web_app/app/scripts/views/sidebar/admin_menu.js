/** Copyright (c) 2014-present James Adam.  All rights reserved.
 *
 * This file is part of WikiFeat.
 *
 *     WikiFeat is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 *     WikiFeat is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 * along with WikiFeat.  If not, see <http://www.gnu.org/licenses/>.
 */

'use strict';
define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'bootstrap',
    'text!templates/sidebar/admin_menu.html'
], function($,_,Marionette,Radio,Bootstrap,
            AdminMenuTemplate){

    return Marionette.ItemView.extend({
        id: 'admin-menu',
        template: _.template(AdminMenuTemplate),
        events: {
            'click a#createWikiLink': 'createWiki',
            'click a#manageUsersLink': 'manageUsers'
        },

        initialize: function(){
            console.log("displaying admin menu");
        },

        createWiki: function(event){
            event.preventDefault();
            Radio.channel('sidebar').trigger('active:link', event.currentTarget);
            Radio.channel('wiki').trigger('create:wiki');
        },

        manageUsers: function(event){
            event.preventDefault();
            Radio.channel('sidebar').trigger('active:link', event.currentTarget);
            Radio.channel('user').trigger('manage:users');
        },

        onRender: function(){},

        onClose: function(){}

    });

});
