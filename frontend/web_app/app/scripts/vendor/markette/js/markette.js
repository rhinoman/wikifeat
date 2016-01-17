/**
 *  Markette
 *  --------
 *
 *  v 0.3.1
 *  http://github.com/rhinoman/markette
 *
 *  Copyright (c) 2016, James Adam
 *  All rights reserved.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice, this
 *  list of conditions and the following disclaimer.
 *
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *  this list of conditions and the following disclaimer in the documentation
 *  and/or other materials provided with the distribution.
 *
 *  * Neither the name of markette nor the names of its
 *  contributors may be used to endorse or promote products derived from
 *  this software without specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 *  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 *  FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 *  DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 *  SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 *  CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 *  OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 *  OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

(function (root, factory) {
    //AMD module if available
    if (typeof define === 'function' && define.amd) {
        define(['jquery', 'underscore', 'marionette', 'commonmark', 'bootstrap'], function($,_,Marionette,Commonmark){
            return (root.Markette = factory(root, $, _, Marionette, Commonmark));
        });
    //CommonJS
    } else if (typeof exports !== 'undefined'){
        var $ = require('jquery');
        var _ = require('underscore');
        var Marionette = require('marionette');
        var Commonmark = require('commonmark');
        module.exports = factory(root, $, _, Marionette, Commonmark);
    } else {
    // Browser globals
        root.Markette = factory(root, root.jQuery, root._, root.Marionette, root.commonmark );
    }


}( this, function (root, $, _, Marionette, Commonmark){
    'use strict';

    var Markette = {};
    Markette.VERSION = '0.1.0';

    //The Editor View
    Markette.EditorView = Marionette.ItemView.extend({
        template: _.template('<div class="row markette-button-bar"></div><textarea class="row" id="marketteInput" rows="12"></textarea>'),
        initialize: function(options){
            options = options || {};
            this.showImageButton = options.hasOwnProperty("showImageButton") ? options.showImageButton : true;
        },
        events: {
            'click #boldButton': 'doBold',
            'click #italicButton': 'doItalic',
            'click #linkButton': 'doLink',
            'click #imageButton': 'doImage',
            'click #quoteButton': 'doBlockQuote',
            'click #codeBlockButton': 'doCodeBlock',
            'click #numberedListButton': 'doNumberedList',
            'click #bulletListButton': 'doBulletList',
            'click #heading1Button': 'doHeading1',
            'click #heading2Button': 'doHeading2'
        },

        /**
         * These functions create basic buttons for the toolbar
         */
        buttonTemplate: function(){
            return _.template("<button id='<%= btnId %>' type='button' class='btn btn-default markette-btn' " +
                "data-toggle='tooltip' data-placement='bottom' title='<%= tooltip %>'><%= glyph %></button>");
        },
        boldButton: function(){
            return this.buttonTemplate()({
                btnId: "boldButton",
                glyph: "<span class='glyphicon glyphicon-bold'></span>",
                tooltip: "Bold Text"
            });
        },
        italicButton: function(){
            return this.buttonTemplate()({
                btnId: "italicButton",
                glyph: "<span class='glyphicon glyphicon-italic'></span>",
                tooltip: "Italic Text"
            });
        },
        linkButton: function(){
            return this.buttonTemplate()({
                btnId: "linkButton",
                glyph: "<span class='glyphicon glyphicon-link'></span>",
                tooltip: "Hyperlink"
            });
        },
        blockQuoteButton: function(){
            return this.buttonTemplate()({
                btnId: "quoteButton",
                glyph: "<span class='glyphicon glyphicon-comment'></span>",
                tooltip: "Block Quote"
            });
        },
        codeBlockButton: function(){
            return this.buttonTemplate()({
                btnId: "codeBlockButton",
                glyph: "<span class='mk-btn-text'>{ }</span>",
                tooltip: "Code Block"
            });
        },
        imageButton: function(){
            return this.buttonTemplate()({
                btnId: "imageButton",
                glyph: "<span class='glyphicon glyphicon-picture'></span>",
                tooltip: "Image"
            });
        },
        numberedListButton: function(){
            return this.buttonTemplate()({
                btnId: "numberedListButton",
                glyph: "<span class='mk-btn-text'>1.</span>",
                tooltip: "Numbered List"
            });
        },
        bulletListButton: function(){
            return this.buttonTemplate()({
                btnId: "bulletListButton",
                glyph: "<span class='glyphicon glyphicon-list'></span>",
                tooltip: "Bulleted List"
            })
        },
        heading1Button: function(){
            return this.buttonTemplate()({
                btnId: "heading1Button",
                glyph: "<span class='mk-btn-text'>H1</span>",
                tooltip: "Level 1 Header"
            })
        },
        heading2Button: function(){
            return this.buttonTemplate()({
                btnId: "heading2Button",
                glyph: "<span class='mk-btn-text'>H2</span>",
                tooltip: "Level 2 Header"
            })
        },

        // Renders the editor and toolbar
        onRender: function(){
            this.$(".markette-button-bar").append("<div class='btn-group' role='group' aria-label='...'></div>")
            this.$(".btn-group").append(this.boldButton());
            this.$(".btn-group").append(this.italicButton());
            this.$(".btn-group").append(this.linkButton());
            this.$(".btn-group").append(this.blockQuoteButton());
            this.$(".btn-group").append(this.codeBlockButton());
            if(this.showImageButton) {
                this.$(".btn-group").append(this.imageButton());
            }
            this.$(".btn-group").append(this.numberedListButton());
            this.$(".btn-group").append(this.bulletListButton());
            this.$(".btn-group").append(this.heading1Button());
            this.$(".btn-group").append(this.heading2Button());
        },

        // Returns the contents of the input textarea
        getText: function(){
            return this.$("textarea#marketteInput").val();
        },

        // Sets the contents of the input textarea
        setText: function(text){
            this.$("textarea#marketteInput").val(text);
        },

        /**
         * Event handlers
         */

        // Make selected text bold
        doBold: function(event){
            this.doInlineMarkup({
                before: '**',
                after: '**'
            });
        },

        // Make selected text italic
        doItalic: function(event){
            this.doInlineMarkup({
                before: '*',
                after: '*'
            });
        },

        // Insert Hyperlink
        doLink: function(event){
            this.doInlineMarkup({
                before: '[',
                after: '](http://)'
            });
        },

        // Insert Image link
        doImage: function(event){
            this.doInlineMarkup({
                before: '![',
                after: '](http://)'
            })
        },

        // Insert a block quote
        doBlockQuote: function(event){
            this.doBlockMarkup({
                before: '> ',
                after: ''
            }, "Blockquote");
        },

        doCodeBlock: function(event){
            this.doBlockMarkup({
                before: '```\n',
                after: '\n```'
            }, "code goes here");
        },

        doNumberedList: function(event){
            this.doBlockMarkup({
                before: ' 1. ',
                after: ''
            }, "List item")
        },

        doBulletList: function(event){
            this.doBlockMarkup({
                before: ' - ',
                after: ''
            }, "List item")
        },

        doHeading1: function(event){
            this.doHeadingMarkup("=");
        },

        doHeading2: function(event){
            this.doHeadingMarkup("-");
        },

        doHeadingMarkup: function(mChar){
            var ta = this.$("textarea#marketteInput");
            var selectionStrings = this._getSelectionStrings(ta);
            if (selectionStrings.selection === ""){
                selectionStrings.selection = "Heading";
            }
            var len = selectionStrings.selection.length;
            var hdgStr = new Array(len + 1).join(mChar);
            this.doBlockMarkup({
                before: '',
                after: '\n' + hdgStr
            }, "Heading", true);
        },

        doBlockMarkup: function(tokens, placeholder, stripnewlines){
            stripnewlines = stripnewlines || false;
            var ta = this.$("textarea#marketteInput");
            var selectionStrings = this._getSelectionStrings(ta);
            if (selectionStrings.selection === ""){
                selectionStrings.selection = placeholder;
            }
            if(stripnewlines) {
                selectionStrings.selection = selectionStrings.selection.replace(/(\r\n|\n|\r)/gm, "");
            }
            //Check if preceding text already contains newline
            var bs = selectionStrings.before;
            var as = selectionStrings.after;
            var prepend = tokens.before;
            var append = tokens.after;
            if(!this._alreadyMarked(ta, tokens)){
                if (as.length > 0 && as[0] === '\n'){
                    append += '\n';
                } else if(as.length >0) {
                    append += '\n\n';
                }
                if (bs.length > 0 && bs[bs.length-1] === '\n'){
                    prepend = '\n' + prepend;
                } else if(bs.length > 0) {
                    prepend = '\n\n' + prepend;
                }
                this._replaceSelection(ta, {before: prepend, after: append}, selectionStrings);
            }
        },

        doInlineMarkup: function(tokens){
            var ta = this.$("textarea#marketteInput");
            var selectionStrings = this._getSelectionStrings(ta);
            // Make sure this text isn't already marked up
            if(this._alreadyMarked(ta, tokens)){
                return false;
            } else {
                this._replaceSelection(ta, tokens, selectionStrings);
                return true;
            }
        },

        _alreadyMarked: function(textarea, tokens){
            var allText = textarea.val();
            var position = this._getInputSelection(textarea.get(0));
            var prevBefore = allText.substring(position.start - tokens.before.length, position.start);
            var prevAfter = allText.substring(position.end,  position.end + tokens.after.length);
            if(prevBefore === tokens.before && prevAfter === tokens.after){
                return true;
            } else {
                return false;
            }
        },

        _getSelectionStrings: function(textarea){
            var position = this._getInputSelection(textarea.get(0));
            var allText = textarea.val();
            var textBefore = allText.substring(0, position.start);
            var textAfter = allText.substring(position.end, allText.length);
            var selection = allText.substring(position.start, position.end);
            return {
                before: textBefore,
                after: textAfter,
                selection: selection
            };
        },

        _replaceSelection: function(textarea, tokens, selectionStrings){
            var beforeString = selectionStrings.before + tokens.before;
            var afterString = tokens.after + selectionStrings.after;
            textarea.val(beforeString +
                selectionStrings.selection +
                afterString);
            //And re-select our selection
            textarea.get(0).setSelectionRange(beforeString.length,
                beforeString.length + selectionStrings.selection.length);
        },

        //Returns the start and end position (or caret position) in the text area
        _getInputSelection: function() {
            var textarea = this.$("textarea#marketteInput").get(0);
            var start = 0, end = 0, normalizedValue, range,
                textInputRange, len, endRange;

            if (typeof textarea.selectionStart == "number" && typeof textarea.selectionEnd == "number") {
                start = textarea.selectionStart;
                end = textarea.selectionEnd;
            } else {
                range = document.selection.createRange();

                if (range) {
                    len = textarea.value.length;
                    normalizedValue = textarea.value.replace(/\r\n/g, "\n");

                    textInputRange = textarea.createTextRange();
                    textInputRange.moveToBookmark(range.getBookmark());

                    endRange = textarea.createTextRange();
                    endRange.collapse(false);

                    if (textInputRange.compareEndPoints("StartToEnd", endRange) > -1) {
                        start = end = len;
                    } else {
                        start = -textInputRange.moveStart("character", -len);
                        start += normalizedValue.slice(0, start).split("\n").length - 1;

                        if (textInputRange.compareEndPoints("EndToEnd", endRange) > -1) {
                            end = len;
                        } else {
                            end = -textInputRange.moveEnd("character", -len);
                            end += normalizedValue.slice(0, end).split("\n").length - 1;
                        }
                    }
                }
            }

            return {
                start: start,
                end: end
            };
        }

    });

    // Override extend to pass events to child classes
    Markette.EditorView.extend = function(child){
        var view = Backbone.View.extend.apply(this, arguments);
        view.prototype.events = _.extend({}, this.prototype.events, child.events);
        return view;
    };

    //The Preview... view
    Markette.Preview = Marionette.ItemView.extend({
        template: _.template('<div class="markette-preview"></div>'),

        initialize: function(options){
            this.reader = new Commonmark.Parser();
            this.writer = new Commonmark.HtmlRenderer();
        },

        renderPreview: function(mdText){
            var parsedText = this.reader.parse(mdText);
            var htmlText = this.writer.render(parsedText);
            this.$(".markette-preview").html(htmlText);
        }


    });

    return Markette;
}));