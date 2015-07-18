/**
 * Created by jcadam on 1/26/15.
 * Responsible for displaying a list of wikis in the sidebar
 */

'use strict';
define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.stickit',
    'backbone.radio',
    'bootstrap',
    'entities/wiki/wiki',
    'text!templates/sidebar/wiki_list_item.html'
], function($,_,Marionette,Stickit,Radio,
            Bootstrap,WikiModel,WikiListItemTemplate){

    return Marionette.ItemView.extend({
        id: 'wiki-list-item-view',
        template: _.template(WikiListItemTemplate),
        model: WikiModel,
        bindings: {
            '#wikiNameText': {
                observe: 'name'
            }
        },
        events: {
           "click a": "navigateToWiki"
        },
        //Somebody clicked on a wiki in the navbar
        navigateToWiki: function(event){
            event.preventDefault();
            Radio.channel('sidebar').trigger('active:link', event.currentTarget);
            Radio.channel('wiki').trigger('show:wiki', this.model.get('id'));
            console.log("Navigate to " + this.model.get('name'));
        },
        onRender: function(){
            if(typeof this.model !== 'undefined'){
                this.stickit();
                this.$('a').attr("id", this.model.id);
                this.$('a').attr("title", this.model.get('description'));
                this.$('a').attr("href", "/app/wikis/" + this.model.get('slug'));
                this.$('[data-toggle="tooltip"]').tooltip();
            }
        },

        onClose: function(){
            this.unstickit();
        }
    });

});
