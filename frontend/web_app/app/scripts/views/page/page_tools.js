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
            "click #deletePageLink":   "deletePage"
        },

        initialize: function(options){
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
            }
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
                )
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
