/* 
 * poetry-plugin.js
 *
 * Description: This plugin is intended 
 *              to serve as an example Wikifeat Plugin
 *
 */
// Jquery, Underscore, Backbone, Marionette are available in 
// the Wikifeat environment, no need to load them.
// 'Wikifeat' is of course the Wikifeat Plugin API module
//
// You must namespace your plugin, using the Plugin name
var Poetry = (function ($, _, Backbone, Marionette, Wikifeat) {

    var Poem = Backbone.Model.extend({
        defaults: {
            text: "&lt;Poetry Plugin&gt;"
        }
    });

    var songOfNature = new Poem({id: "emerson"});

    // A basic Marionette Item View... could also be a vanilla Backbone View
    var PoemView = Marionette.ItemView.extend({
        template: _.template('<div class="poetry" id="poem"></div>'),
        model: Poem,
        initialize: function () {
            this.model.on("change", this.render, this);
        },
        onRender: function () {
            this.$("#poem").html(this.model.get('text'));
        }
    });

    /**
     *  Embeddable 'content' plugins should have a 'helper' view that will be shown
     *  when the user selects the plugin from the editor view 'plugins' dropdown
     */
    var PoemInsertView = Marionette.ItemView.extend({
        template: _.template('<div class="modal fade" id="poetryInsertModal">'+
                '<div class="modal-dialog">' +
                '<div class="modal-content">' +
                '<div class="modal-header">' +
                '<h4 class="modal-title" id="poetryInsertTitle">Insert Poem</h4></div>' +
                '<div class="modal-body">' +
                '<div class="form-group">' +
                '<label for="inputDataId" class="control-label">Data ID</label>' +
                '<input type="text" class="form-control" id="inputDataId" placeholder="emerson">' +
                '</div></div>' +
                '<div class="modal-footer">' +
                '<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>' +
                '<button type="button" class="btn btn-primary" data-dismiss="modal" id="confirmButton">Ok</button>' +
                '</div></div></div></div>'),
        events: {
            'click #confirmButton': 'confirmClick'
        },
        initialize: function (options){
            this.result = options.result;
        },
        onRender: function (){
            this.$("#poetryInsertModal").modal();
        },
        confirmClick: function(event){
            event.preventDefault();
            const dataId = this.$("#inputDataId").val() || "emerson";
            const result = "<div data-plugin='Poetry' data-id='" + dataId + "'></div>";
            //Resolve the result, so the Wikifeat editor knows what to insert.
            this.result.resolve(result);
        }
    });

    // Sidebar menu item view
    var PoemMenuView = Marionette.ItemView.extend({
        id: 'poem-menu-view',
        events: {
            'click a#showPoemLink': 'showPoem'
        },
        template: _.template('<div>' +
            '<a class="sbTopLevel" data-toggle="collapse"' +
            'href="#poemSubMenu" aria-expanded="false"' +
            'aria-controls="poemSubMenu">' +
            '<span class="glyphicon glyphicon-ok"></span>' +
            '&nbsp;Poetry' +
            '</a>' +
            '</div>' +
            '<div class="collapse subMenu" id="poemSubMenu">' +
            '<a href="" class="subMenuLink" id="showPoemLink">' +
            '<span class="glyphicon glyphicon-file"></span>' +
            'Show Poem' +
            '</a>' +
            '<div>' +
            '</div>'),
        initialize: function () {
            console.log("initializing poem menu");
        },
        // Display a poem in the main content area
        showPoem: function (event) {
            event.preventDefault();
            var poemView = new PoemView({model: songOfNature});
            Wikifeat.showContent(poemView);
            window.history.pushState('', '', '/app/x/poetry/songOfNature');
        }
    });

    //I know, lets make ourselves a little router
    //This allows direct navigation to resource pages
    var PoetryRouter = Marionette.AppRouter.extend({
        appRoutes: {
            "x/poetry/:poem": "showPoem"
        },
        controller: {
            showPoem: function (poem) {
                console.log("Showing Poem " + poem);
                var poemView = new PoemView({model: songOfNature});
                Wikifeat.showContent(poemView);
            }
        }
    });
    //Create the router
    var pr = new PoetryRouter();
    // Your plugin must return a few things which may be called by wikifeat
    return {
        // all plugins must contain a 'start' function, named 'start'
        start: function (started) {
            // do some initialization for our plugin here.
            console.log("Poetry plugin started");
            var self = this;
            //This script fetches some additional data from a 'data' subdirectory
            $.ajax("/app/plugin/Poetry/resource/data/song_of_nature.html")
                .done(function (text) {
                    songOfNature.set("text", text);
                })
                .fail(function () {
                    console.log("Could not load resource");
                });
            Wikifeat.addMenuItem("PoetryMenu", new PoemMenuView());
            started.resolve();
        },
        /**
         * Plugins that support embeddable content in wiki pages must provide
         * a 'getContentView' function which returns a Backbone/Marionette view.
         * The getContentView function takes an el parameter - the DOM element the
         * content will be embedded in, and a 'contentId' - used by your plugin as
         * a resource identifier.
        */
        getContentView: function (el, contentId) {
            console.log("Showing Poetry content");
            var poem = new Poem({id: contentId});
            if (contentId === "emerson") {
                poem = songOfNature;
            }
            return new PoemView({el: el, model: poem});
        },

        /**
         * Plugins that support embeddable content in wiki pages should provide a
         * 'getInsertLabel' function which returns the text that will appear in the
         * Plugins dropdown in the wiki editor interface
         */
        getInsertLabel: function(){
            return "<span class='glyphicon glyphicon-paperclip'></span>&nbsp;Poetry";
        },

        /**
         * Plugins that support embeddable content should provide a
         * 'getInsertView' function.  This returns a plugin 'helper/editor view'
         * used when the user wants to insert plugin content.
         *
         * As part of the options object, a 'result' field will be passed in
         * containing a $.Deferred object.  The result should be resolved with
         * the resultant 'div' tag for the plugin when the user has finished with the view.
         * @param options - should contain a $.Deferred object called 'result'
         * @returns {*}
         */
        getInsertView: function(options){
            return new PoemInsertView(options);
        }
    };
})($, _, Backbone, Marionette, Wikifeat);

