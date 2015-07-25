/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.
 *
 * This file is part of WikiFeat
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
    'backbone.stickit',
    'backbone.radio',
    'bootstrap',
    'text!templates/sidebar/user_menu.html',
    'entities/user/user'
], function($,_,Marionette, Stickit, Radio,
            Bootstrap, UserMenuTemplate, UserModel){

    return Marionette.ItemView.extend({
        id: 'user-menu-view',
        initialize: function(){
            console.log("initializing User Menu View");
        },
        model: UserModel,
        bindings: {
            '#userNameText': {
                observe: 'userPublic',
                onGet: function(data){
                    return data.firstName + " " + data.lastName;
                }
            }
        },
        events: {
            "click a#accountSettingsLink": "accountSettings",
            "click a#logoutLink": "logout"
        },

        template: _.template(UserMenuTemplate),

        /* Account Settings Menu */
        accountSettings: function(event){
            event.preventDefault();
            //TODO: Actually implement this
        },

        /* Logout the current user */
        logout: function(event){
            event.preventDefault();
            Radio.channel('user').trigger('user:logout');
        },

        /* on render callback */
        onRender: function(){
            if(typeof this.model !== 'undefined') {
                this.stickit();
            }
        },

        onClose: function(){
            this.unstickit();
        }
    });
});
