require.config({
	baseUrl: "/app/plugin/OpenLayers/resource/",
	paths: {
	},
	packages: []
});

var OpenLayers = {

    start: function(started){
        require(['ol-app'], function(OLApp){
            OLApp.start();
            OpenLayers.getContentView = function(el, contentId){
                return OLApp.getContentView(el, contentId);
            };
            started.resolve();
	    });
    }

};
