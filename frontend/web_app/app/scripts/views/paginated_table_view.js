/**
 * Copyright (c) 2014-present James Adam.  All rights reserved.*
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
    'marionette'
], function($,_,Marionette){

    return Marionette.CompositeView.extend({

        baseEvents:{
            'click #pageNext': 'pageNext',
            'click #pagePrevious': 'pagePrevious',
            'click .page-link': 'pageLink'
        },

        //set these in the child view
        //DO NOT override events in the child view, please.
        additionalEvents: {},

        events: function(){
            return _.extend({}, this.baseEvents, this.additionalEvents);
        },

        //Go to the previous page
        pagePrevious: function(event){
            event.preventDefault();
            if(this.collection.hasPreviousPage()){
                this.collection.getPreviousPage();
            }
            this.setActivePageLink(this.collection.state.currentPage);
        },

        //Go to the next page
        pageNext: function(event){
            event.preventDefault();
            if(this.collection.hasNextPage()){
                this.collection.getNextPage();
            }
            this.setActivePageLink(this.collection.state.currentPage);
        },

        //Jump to a page
        pageLink: function(event){
            event.preventDefault();
            var target = event.currentTarget;
            var targetId = target.id;
            var pageNum = parseInt(targetId.split('-')[1]);
            this.collection.getPage(pageNum, {reset:true});
            this.setActivePageLink(pageNum);
        },

        setActivePageLink: function(pageNum){
            this.$("li").removeClass('active');
            this.$("li#page-link-" + pageNum).addClass('active');
        },

        onRender: function(){
            //Display initial pagination state
            var numPages = this.collection.state.lastPage;
            for (var i = numPages; i > 0; i--) {
                this.$("li#previousLink")
                    .after("<li id='page-link-" + i + "'><a href='#' class='page-link' id='page-" + i + "'>" + i + "</a></li>");
            }
            var currentPage = this.collection.state.currentPage;
            this.setActivePageLink(currentPage);
        }
    });
});


