var Example = (function(Backbone, Marionette, Wikifeat){
  console.log("EXAMPLE PLUGIN!");
  Wikifeat.printVersion();

  TestView = Backbone.View.extend({
    render: function(){
	    var template = "<h1>EXAMPLE PLUGIN!</h1>";
	    this.$el.html(template);
    }
  });
  return {
    start: function(){
      var testView = new TestView();
			alert("POOP");
      //Wikifeat.showContent(testView);
    }
  };
})(Backbone, Marionette, Wikifeat);

