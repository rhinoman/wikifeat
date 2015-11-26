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

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'bootstrap',
    'entities/wiki/page',
    'entities/wiki/wiki',
    'views/main/confirm_dialog',
    'text!templates/page/page_tools.html',
    'text!templates/main/alert.html'
], function($,_,Marionette,Radio,Bootstrap,
            PageModel,WikiModel,ConfirmDialog,
            PageToolsTemplate, AlertTemplate){
    'use strict';

    return Marionette.ItemView.extend({
        id: "page-tools-view",
        template: _.template(PageToolsTemplate),
        model: PageModel,
        wikiModel: null,

        events: {
            "click #editPageLink":     "editPage",
            "click #addChildPageLink": "createChildPage",
            "click #historyPageLink":  "pageHistory",
            "click #viewModeLink" :    "changeViewMode",
            "click #deletePageLink":   "deletePage"
        },

        initialize: function(options){
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
            }
            //this.model.on('setViewMode', this.render, this);
        },

        //Somebody clicked the edit button!
        editPage: function(event){
            event.preventDefault();
            Radio.channel('page').trigger('edit:page', this.model, this.wikiModel);
        },

        //Create child page for this page
        createChildPage: function(event){
            event.preventDefault();
            Radio.channel('page').trigger('create:page', this.wikiModel,
                {homePage: false, parent: this.model.id});
        },

        //Display the page revision history
        pageHistory: function(event){
            event.preventDefault();
            Radio.channel('page').trigger('show:history', this.model, this.wikiModel);
        },

        //Changes between raw and formatted page view
        changeViewMode: function(event){
            event.preventDefault();
            if(this.model.viewMode === 'raw'){
                this.model.viewMode = 'formatted';
            } else {
                this.model.viewMode = 'raw';
            }
            this.model.trigger('setViewMode');
        },

        //Delete a page
        deletePage: function(event){
            event.preventDefault();
            var parentPage = this.model.get('parent');
            var confirmCallback = function(){
                Radio.channel('wikiManager').request('delete:page', self.model)
                    .done(function(response){
                        if(typeof response === 'undefined'){
                            self.$("#alertBox").html(AlertTemplate);
                            self.$("div.alert").addClass("alert-danger");
                            self.$("#alertText").html("Could not delete page");
                        } else {
                            if(parentPage !== "") {
                                Radio.channel('page').trigger('show:page',
                                    parentPage, self.wikiModel);
                            } else {
                                Radio.channel('wiki').trigger('show:wiki',
                                    self.wikiModel.id);
                            }
                        }
                    });
            };
            var confirmDialog = new ConfirmDialog({
                message: 'Are you sure you wish to delete the page "' +
                this.model.get('title') + '?"  This action is irreversible.',
                confirmCallback: confirmCallback
            });

            Radio.channel('main')
                .trigger('show:dialog', confirmDialog);
            var self = this;
        },

        onRender: function(){
            if(typeof this.model !== 'undefined'){
                //Display the various control buttons, checking permissions as we go.
                if(this.model.isEditable) {
                    this.$("#pageToolsMenu ul").append(
                        '<li><a href="#" id="editPageLink">' +
                        '<span class="glyphicon glyphicon-edit"></span>&nbsp;Edit</a></li>'
                    )
                }
                this.$("#pageToolsMenu ul").append(
                    '<li><a href="#" id="historyPageLink">' +
                    '<span class="glyphicon glyphicon-book"></span>&nbsp;Page History</a></li>'
                );

                //Set the view mode list item
                if(this.model.viewMode === "formatted") {
                    this.$("#pageToolsMenu ul").append(
                        '<li><a href="#" id="viewModeLink">' +
                        '<span class="glyphicon glyphicon-eye-open"></span>' +
                        '&nbsp;View Raw</li></a>'
                    );
                } else {
                     this.$("#pageToolsMenu ul").append(
                        '<li><a href="#" id="viewModeLink">' +
                        '<span class="glyphicon glyphicon-eye-open"></span>' +
                        '&nbsp;View Formatted</li></a>'
                    );
                }
                if(this.wikiModel.canCreatePage){
                    this.$("#pageToolsMenu ul").append(
                        '<li><a href="#" id="addChildPageLink">' +
                        '<span class="glyphicon glyphicon-plus-sign"></span>&nbsp;Add Child Page</a></li>'
                    )
                }
                if(this.model.isDeletable){
                    this.$("#pageToolsMenu ul").append(
                        '<li><a href="#" id="deletePageLink">' +
                        '<span class="glyphicon glyphicon-trash"></span>&nbsp;Delete</a></li>'
                    )
                }
                if(this.$("#pageToolsMenu ul li").length <= 0){
                    this.$("#pageToolsMenu").css("display", "none");
                } else {
                    this.$("#pageToolsMenu").css("display", "block");
                }

                this.$('.dropdown-toggle').dropdown();
            }
        }

    });
});
