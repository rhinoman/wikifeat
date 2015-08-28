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
    'backbone.radio',
    'backbone.stickit',
    'bootstrap',
    'entities/user/user',
    'text!templates/user/user_info_dialog.html'
], function($,_,Marionette,Radio,Stickit,Bootstrap,
            UserModel,UserInfoTemplate){

    return Marionette.ItemView.extend({
        className: "user-info-dialog",
        model: UserModel,
        template: _.template(UserInfoTemplate),
        bindings: {
            '#userInfoTitle':{
                observe: 'name'
            },
            '#userFullName':{
                observe: 'userPublic',
                updateModel:false,
                onGet: function(data){
                    return data.firstName + " " + data.lastName;
                }
            },
            '#email':{
                observe: 'userPublic',
                updateModel:false,
                onGet: function(data){
                    return data.contactInfo.email;
                }
            },
            '#userTitle':{
                observe: 'userPublic',
                updateModel:false,
                onGet: function(data){
                    return data.title;
                }
            }

        },


        onRender: function(){
            if(this.model !== 'undefined'){
                var up = this.model.get('userPublic');
                this.$("#userInfoModal").modal();
                this.$("#pictureWrapper").html(this.model.getAvatar());
                this.$('#email').attr('href', 'mailto:' + up.contactInfo.email);
                this.stickit();
            }

        },

        onClose: function(){
            this.unstickit();
        }
    });

});