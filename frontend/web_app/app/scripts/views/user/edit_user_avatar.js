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
    'entities/user/user_avatar',
    'bootstrap',
    'views/main/alert',
    'text!templates/user/edit_avatar_dialog.html'

], function($,_,Marionette,Radio,Stickit,
            UserAvatarModel,Bootstrap,AlertView,
            EditAvatarDialogTemplate){

    return Marionette.ItemView.extend({
        id: "edit-avatar-dialog",
        model: UserAvatarModel,
        template: _.template(EditAvatarDialogTemplate),
        bindings: {

        },
        events: {

        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.stickit();
                this.$("#editAvatarModal").modal();
            }
        },

        onClose: function(){
            this.unstickit();
        }
    });

});