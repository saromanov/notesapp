/** @jsx React.DOM */
var socket = io.connect();

var NoteItem = React.createClass({

	getInitialState: function() {
    	return {value: ''};
  	},

	/**
	 * @return {object}
	*/
	render: function(){
		var value = this.state.value
		return (
			  <div className="note">
      			  <input type="text" id="note-title" size="50" ref="title" placeholder="Enter title of note"/> <br />
      			  <textarea id="note-text" ref="notetext" rows="10" cols="50" value={value} onChange={this._onChange}></textarea><br />
    			  <button id ='add' className='btn btn-primary btn-lg' onClick={this._onAddNote}>Save</button><br />
  			 </div>
		)
	},

	_handleData: function(event) {
		console.log(event);
	},

	_onAddNote: function(event) {
		event.preventDefault();
		var title = this.refs['title'].value;
		var text = this.refs['notetext'].value;
		var timestamp = Date.now();
		this.socket.emit('newnote', {'title': title, 'text': text})
	},

	_onChange: function(event){
		this.setState({
			value: event.target.value
		})
	},
});
ReactDOM.render(
  <NoteItem />,
  document.getElementById('note')
);