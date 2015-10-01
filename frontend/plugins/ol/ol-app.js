define(['jquery',
	'underscore',
	'marionette',
	'views/menuview',
	'views/mapview'
], function($,_,Marionette,OLMenuView,OLMapView){

	var OLApp = new Backbone.Marionette.Application();

	OLApp.on("start", function() {
		console.log("The OpenLayers plugin has started.");
	});

	OLApp.getContentView = function(el, contentId){
		console.log("Shwoing OpenLayers content");
		return new OLMapView({el: el});
	};
	
	return OLApp;
});
