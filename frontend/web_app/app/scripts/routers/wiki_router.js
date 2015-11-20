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
 * Created by jcadam on 12/12/14.
 */

define([
    'jquery',
    'underscore',
    'backbone',
    'marionette',
    'backbone.radio'
], function($, _, Backbone, Marionette, Radio){
    return Marionette.AppRouter.extend({
        //Routes for direct navigation to wiki resources
        appRoutes: {
            "wikis/create" : "createWiki",
            "wikis/:wikiSlug" : "showWiki",
            "wikis/:wikiSlug/edit" : "editWiki",
            "wikis/:wikiSlug/files" : "wikiFiles",
            "wikis/:wikiSlug/members" : "wikiMembers",
            "wikis/:wikiSlug/pages/:pageSlug?revision=:revision": "showRevision",
            "wikis/:wikiSlug/pages/:pageSlug" : "showPage",
            "wikis/:wikiSlug/pages/:pageSlug/history" : "showPageHistory",
            "wikis/:wikiSlug/pages/:pageSlug/edit" : "editPage"
        },
        controller: {
            createWiki: function(){
                Radio.channel('wiki').trigger('create:wiki');
            },
            showWiki: function(wikiSlug){
                Radio.channel('wiki').trigger('show:slug:wiki', wikiSlug);
            },
            editWiki: function(wikiSlug){
                Radio.channel('wikiManager').request('get:slug:wiki', wikiSlug)
                    .done(function(wikiModel){
                        Radio.channel('wiki').trigger('edit:wiki', wikiModel);
                    });
            },
            wikiFiles: function(wikiSlug){
                Radio.channel('wikiManager').request('get:slug:wiki', wikiSlug)
                    .done(function(wikiModel){
                        Radio.channel('wiki').trigger('init:layout', wikiModel);
                        Radio.channel('file').trigger('manage:files', wikiModel);
                    });
            },
            wikiMembers: function(wikiSlug){
                Radio.channel('wikiManager').request('get:slug:wiki', wikiSlug)
                    .done(function(wikiModel){
                        Radio.channel('wiki').trigger('init:layout', wikiModel);
                        Radio.channel('user').trigger('manage:members', wikiModel);
                    });
            },
            showRevision: function(wikiSlug, pageSlug, revisionId){
                Radio.channel('wikiManager').request('get:slug:wiki', wikiSlug)
                    .done(function(wikiModel){
                        Radio.channel('wiki').trigger('init:layout',wikiModel);
                        Radio.channel('page').trigger('show:page:revision',
                            pageSlug,wikiModel,revisionId,{slug:true})
                    });
            },
            showPage: function(wikiSlug, pageSlug){
                Radio.channel('wiki').trigger('show:slug:wiki', wikiSlug, pageSlug);
            },
            showPageHistory: function(wikiSlug, pageSlug){
                Radio.channel('wiki').trigger('show:slug:wiki',
                    wikiSlug, pageSlug, {history: true})
            },
            editPage: function(wikiSlug, pageSlug){
                Radio.channel('wiki').trigger('show:slug:wiki',
                    wikiSlug, pageSlug, {edit: true});
            }
        }
    });
});
