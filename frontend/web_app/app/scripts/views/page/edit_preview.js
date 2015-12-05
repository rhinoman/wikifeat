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

'use strict';
define([
   'jquery',
    'underscore',
    'marionette',
    'commonmark',
    'entities/wiki/page'
], function($,_,Marionette,Commonmark,PageModel){

    return Marionette.ItemView.extend({
        model: PageModel,
        id:"wmd-preview",
        className: "wmd-panel wmd-preview",
        template: _.template("<div></div>"),
        wipText: null,

        initialize: function(options){
            if(options.hasOwnProperty('wipText')){
                this.wipText = options.wipText;
                this.wipText.on('redisplay', this.updatePreview, this);
            }
            this.reader = new Commonmark.Parser();
            this.writer = new Commonmark.HtmlRenderer();
        },

        onShow: function(){
            if(this.wipText !== null) {
                var parsedText = this.reader.parse(this.wipText.get('data'));
                this.$el.html(this.writer.render(parsedText));
            }
        },

        updatePreview: function(event){
            var parsedText = this.reader.parse(this.wipText.get('data'));
            this.$el.html(this.writer.render(parsedText));
        }
    });

});