var app = app || {};

app.EventList = React.createClass({
	render : function() {
		var htmlValue = this.props.value.map(function(x){
			var view = "list-group-item list-group-item-" + x.view;
			return (
				<div key={app.getRandomId(1000,99999)}>
			      <ul className="list-group">
				     <li className={view}>{x.msg}</li>
			      </ ul>
			   </div>
				)
		});

	    return <div>{htmlValue}</div>
	}
});