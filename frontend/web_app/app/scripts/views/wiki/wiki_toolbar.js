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
    'text!templates/wiki/wiki_toolbar.html'
], function($,_,Marionette,Radio,Bootstrap,
            WikiModel,WikiToolbarTemplate){

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
        },

        viewFiles: function(event){
            event.preventDefault();
            Radio.channel('file').trigger('manage:files', this.model);
        },

        onClose: function(){}

    });

});
