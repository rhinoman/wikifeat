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
    'views/paginated_table_view',
    'views/user/manage_members_item',
    'views/user/select_users_dialog',
    'entities/user/user',
    'text!templates/user/manage_members.html',
    'text!templates/main/alert.html'
], function($,_,Marionette,Radio,PaginatedTableView,
            ManageMembersItemView,SelectUsersDialog,UserModel,
            ManageMembersTemplate,AlertTemplate){
    return PaginatedTableView.extend({
        className: "manage-members-view",
        template: _.template(ManageMembersTemplate),
        childView: ManageMembersItemView,
        childViewContainer: "#memberListContainer",
        additionalEvents: {
            'click #addMemberButton' : 'addMember'
        },
        wikiModel: null,

        childViewOptions: function(model, index){
            return {
                resourceId: this.resourceId
            }
        },

        childEvents: {
            'remove:member': function(childView, member){
                this.collection.remove(member);
            },
            'error': function(childView, errorMessage){
                this.$("#alertBox").html(AlertTemplate);
                this.$("#alertBox").addClass("alert-danger");
                this.$("#alertText").html(errorMessage);
            }
        },

        initialize: function(options){
            var self = this;
            options = options || {};
            if(options.hasOwnProperty('wikiModel')){
                this.wikiModel = options.wikiModel;
                this.resourceId = this.wikiModel.id;
                this.resource = 'wiki_' + this.resourceId;
            } else {
                this.resourceId = null;
            }
            var curUserReq = $.Deferred();
            this.currentUser = curUserReq.promise();
            //Load the current user
            Radio.channel('userManager').request('get:currentUser')
                .done(function(response){
                    if(typeof response !== 'undefined') {
                        curUserReq.resolve(response);
                    }
                });
        },

        onShow: function(){
            if(this.wikiModel !== null){
                this.$("#manage_members_header").html(
                    'Manage Members for "' + this.wikiModel.get('name') + '"');
            }
            var self = this;
            this.currentUser.done(function(currentUser){
                self.children.each(function(view){
                    if(view.model.id === currentUser.id){
                        view.disableControls();
                    }
                });
            });
        },

        addMember: function(event){
            event.preventDefault();
            var self = this;
            Radio.channel('userManager').request('get:allUserList', {})
                .done(function(response){
                   if(typeof response === 'undefined'){
                   } else {
                       var sud = new SelectUsersDialog({
                           collection: response,
                           resourceId: self.resourceId,
                           memberList: self.collection
                       });

                       Radio.channel('main').trigger('show:dialog', sud);
                   }
                });
        }
    });
});