/** @jsx React.DOM */


var ws = new WebSocket("ws://" + location.host + "/sockets/"+ name);

var NoteStore = {
	notes: [],
	addNote: function(note) {
		this.notes.push(note);
	}
};

var NoteList = React.createClass({
	render : function(){
		return (
			<strong> Back </strong>
		);
	}
});

var EnterNote = React.createClass({

	getInitialState(){
		return {value: ''};
	},

	render: function(){
		var value = this.state.value;
		return (
			 <div className="note">
      			  <input type="text" id="note-title" size="40" ref="title" placeholder="Enter title of note"/> <br />
      			  <textarea id="note-text" ref="notetext" rows="10" cols="40" value={value} onChange={this._onChange}></textarea><br />
    			  <button id ='add' className='btn btn-primary btn-lg' onClick={this._onAddNote}>Save</button><br />
  			 </div>
		);
	},

	_onAddNote: function(event) {
		event.preventDefault();
		var title = this.refs['title'].value;
		var text = this.refs['notetext'].value;
		var timestamp = Date.now();
		this.props.store.addNote({'title': title, 'text': text});
		ws.send(JSON.stringify({'event': 'add', 'title': title, 'text': text}));
	},

	_onChange: function(event){
		this.setState({
			value: event.target.value
		})
	},
});

var NoteApp = React.createClass({

	getInitialState: function() {
		var name = "test";
		ws.addEventListener("message", function(e) {
      		this.ws.send(JSON.parse(e.data), 'right');
    	});
    	return {value: ''};
  	},

	/**
	 * @return {object}
	*/
	render: function(){
		var value = this.state.value;

		items  = (
			  <ul className="noteItems">
			      <strong> Value </strong>
			   </ul>
			);
		return (
			<div>
			 <EnterNote 
			     store={this.props.store}
			     onChange={this._onChange.bind(this, value)} />

			 <NoteList />

			</div>
		)
	},

	_onChange: function(){
		console.log("Change");
	}
});

ReactDOM.render(
  <NoteApp store={NoteStore}/>,
  document.getElementById('note')
);