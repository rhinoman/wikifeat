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

/**
 * Created by jcadam on 2/11/15.
 */
'use strict';

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'bootstrap',
    'entities/wiki/wiki',
    'entities/error',
    'views/main/confirm_dialog',
    'views/main/error_dialog',
    'text!templates/wiki/wiki_toolbar.html'
], function($,_,Marionette,Radio,Bootstrap,
            WikiModel,ErrorModel,ConfirmDialog,
            ErrorDialog,WikiToolbarTemplate){

    return Marionette.ItemView.extend({
        id: 'wiki-toolbar',
        model: WikiModel,
        template: _.template(WikiToolbarTemplate),

        events: {
            'click a#editWikiLink' : 'editWikiRecord',
            'click a#editMembersLink' : 'editWikiMembers',
            'click a#deleteWikiLink' : 'deleteWiki',
            'click a#viewFilesLink' : 'viewFiles'
        },


        onRender: function(){
            if(typeof this.model !== 'undefined'){
                if(!this.model.canUpdate && !this.model.canDelete){
                    this.$('#adminMenu').css('display','none');
                }
                if(this.model.canUpdate === true){
                    this.addEditWikiLink();
                    this.addEditMembersLink();
                }
                if(this.model.canDelete === true){
                    this.addDeleteWikiLink();
                }
                this.$('.dropdown-toggle').dropdown();
            }
        },

        addEditWikiLink: function(){
            this.$("#adminMenu ul").append('<li><a href="#" id="editWikiLink">' +
            '<span class="glyphicon glyphicon-edit"></span>&nbsp;Edit Wiki Record</a></li>')
        },

        addEditMembersLink: function(){
            this.$("#adminMenu ul").append('<li><a href="#" id="editMembersLink">' +
            '<span class="glyphicon glyphicon-user"></span>&nbsp;Edit Wiki Members</a></li>')
        },

        addDeleteWikiLink: function(){
            this.$("#adminMenu ul").append('<li><a href="#" id="deleteWikiLink">' +
            '<span class="glyphicon glyphicon-trash"></span>&nbsp;Delete Wiki</a>')
        },

        editWikiRecord: function(event){
            event.preventDefault();
            Radio.channel('wiki').trigger('edit:wiki', this.model);
        },

        editWikiMembers: function(event){
            event.preventDefault();
            Radio.channel('user').trigger('manage:members', this.model);
        },

        deleteWiki: function(event){
            event.preventDefault();
            var confirmCallback = function(){
                Radio.channel('wikiManager').request('delete:wiki', self.model)
                    .done(function(response){
                        if(typeof response === 'undefined'){
                            var errorDialog = new ErrorDialog({
                                model: new ErrorModel({
                                    errorTitle: "Error deleting wiki",
                                    errorMessage: "Could not delete wiki"
                                })
                            });
                            Radio.channel('main').trigger('show:dialog', errorDialog);
                        } else {
                            Radio.channel('sidebar').trigger('remove:wiki', self.model);
                            Radio.channel('home').trigger('show:home');
                        }
                    });
            };
            var confirmDialog = new ConfirmDialog({
                message: 'Are you sure you wish to delete the wiki "' +
                    this.model.get('name') + '?" This action is irreversible.',
                confirmCallback: confirmCallback
            });
            Radio.channel('main')
                .trigger('show:dialog', confirmDialog);
            var self = this;
        },

        viewFiles: function(event){
            event.preventDefault();
            Radio.channel('file').trigger('manage:files', this.model);
        },

        onClose: function(){}

    });

});
