define([
	'jquery',
	'underscore',
	'marionette',
	'wikifeat',
	'views/mapview',
	'text!templates/ol_menu.html'
], function($, _, Marionette, Wikifeat, OLMapView, OLMenuTemplate){

	return Marionette.ItemView.extend({
		id: 'ol-menu-view',
		events: {
			'click a#showMapLink': 'showMap'
		},
		template: _.template(OLMenuTemplate),

		initialize: function(){
			console.log("initializing OL Menu");
		},
		showMap: function(event){
			event.preventDefault();
			var mapView = new OLMapView();
			Wikifeat.showContent(mapView);
			window.history.pushState('','','/app/x/ol-plugin/map');
		}
	});

});
