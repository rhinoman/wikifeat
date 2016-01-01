/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 * See the LICENSE.txt file distributed with this work for additional information
 * regarding copyright ownership.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *  Neither the name of Wikifeat nor the names of its contributors may be used
 * to endorse or promote products derived from this software without
 * specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
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
        template: _.template(AdminMenuTemplate),
        events: {
            'click a#createWikiLink': 'createWiki',
            'click a#manageUsersLink': 'manageUsers'
        },
        activeMenu: null,

        initialize: function(){
            console.log("displaying admin menu");
        },

        expandMenu: function(){
            if(!this.isExpanded()) {
                $(this.el).find("div#adminSubMenu").addClass("in");
            }
        },

        isExpanded: function(){
            return $(this.el).find("div#adminSubMenu").hasClass("in");
        },

        setManageUsers: function(){
            this.activeMenu = "manageUsers";
            var link = $(this.el).find("a#manageUsersLink");
            if(link){
                this.setLink(link);
            }
        },

        setCreateWiki: function(){
            this.activeMenu = "createWiki";
            var link = $(this.el).find("a#createWikiLink");
            if(link){
                this.setLink(link);
            }
        },

        createWiki: function(event){
            event.preventDefault();
            //this.setCreateWiki();
            Radio.channel('wiki').trigger('create:wiki');
        },

        manageUsers: function(event){
            event.preventDefault();
            //this.setManageUsers();
            Radio.channel('user').trigger('manage:users');
        },

        setLink: function(link){
            Radio.channel('sidebar').trigger('active:link', link);
            this.expandMenu();
        },

        onRender: function(){
            if(this.activeMenu !== null){
                switch(this.activeMenu){
                    case("manageUsers"):
                        var link = $(this.el).find("a#manageUsersLink");
                        this.setLink(link);
                        break;
                    case("createWiki"):
                        var link = $(this.el).find("a#createWikiLink");
                        this.setLink(link);
                        break;
                }
            }
        },

        onClose: function(){}

    });

});
