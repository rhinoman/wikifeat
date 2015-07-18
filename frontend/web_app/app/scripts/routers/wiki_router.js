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
