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
 * Created by jcadam on 1/26/15.
 * Responsible for displaying a list of wikis in the sidebar
 */

'use strict';
define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio',
    'bootstrap',
    'entities/wiki/wiki',
    'text!templates/sidebar/wiki_list_item.html'
], function($,_,Marionette,Radio,Bootstrap,
            WikiModel,WikiListItemTemplate){

    return Marionette.ItemView.extend({
        id: 'wiki-list-item-view',
        template: _.template(WikiListItemTemplate),
        model: WikiModel,
        events: {
           "click a": "navigateToWiki"
        },
        //Somebody clicked on a wiki in the navbar
        navigateToWiki: function(event){
            event.preventDefault();
            //Radio.channel('sidebar').trigger('active:link', event.currentTarget);
            Radio.channel('wiki').trigger('show:wiki', this.model.get('id'));
            console.log("Navigate to " + this.model.get('name'));
        },
        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.$("#wikiNameText").html(this.model.get('name'));
                this.$('a').attr("id", this.model.id);
                this.$('a').attr("title", this.model.get('description'));
                this.$('a').attr("href", "/app/wikis/" + this.model.get('slug'));
                this.$('[data-toggle="tooltip"]').tooltip({container: 'body'});
            }
        },

        onClose: function(){
        }
    });

});
