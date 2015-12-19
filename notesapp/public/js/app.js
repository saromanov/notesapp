/** @jsx React.DOM */


var NoteStore = {
	notes: [],
	numnotes: 0,
	addNote: function(note) {
		/*if(!(note.title in this.notes)) {
			this.notes[note.title] = note.text;
			this.numnotes+=1;
		}*/
		this.notes.unshift(note);
		this.numnotes+=1;
	},

	getNumNotes: function(){
		return this.numnotes;
	},

	getAll: function() {
		return this.notes;
	}
};

function getNoteState() {
  return {
    allNotes: NoteStore.getAll(),
    allNums: NoteStore.getNumNotes()
  };
}

var getRandomId = function(min, max){
	return (Math.floor(Math.random() * (max - min + 1)) + min).toString();
}
var AppDispatcher = new Flux.Dispatcher();

AppDispatcher.register(function(payload){
	if(payload.eventName == "new-item") {
		NoteStore.addNote(payload.newItem);
		return true;
	}

	return true;
});


var ws = new WebSocket("ws://" + location.host + "/sockets/" + getRandomId(1000,999999));

var Note = React.createClass({

	/**
   	 * @return {object}
     */
	render: function(){
		var style = {
			width: "550px",
			right: "30%"
		}
		return (
			<div id={this.props.key}>
			<button style={style} id={this.props.key} className='btn btn-default btn-lg' onClick={this._onClick}>
			{this.props.title}
			</ button>
			</ div>
			)
	},

	_onClick: function(e) {
		$('#note-title').val(this.props.title);
		//ws.send(JSON.stringify({'event': 'get', 'title': this.props.title, 'text':''}));
	}
});

var EnterNote = React.createClass({

	getInitialState(){
		return {value: '',
		        inp: '', 
				allNotes: NoteStore.getAll(),
    			allNums: NoteStore.getNumNotes(),
    		};
	},

	componentDidMount: function() {
  		var that = this;
  		$.ajax({ url: 'http://127.0.0.1:8082/api/list' })
    		.then(function(data) {
     		 var lst = JSON.parse(data);
     		 var items = JSON.parse(lst.Data);
				items.forEach(function(x){
					AppDispatcher.dispatch({
						eventName: 'new-item',
						newItem: {'id':getRandomId(1000,999999), 'event': 'add', 'title': x.Title, 'text': x.NodeItem}
					});
					that.setState({
						allNotes: NoteStore.getAll(),
						allNums: NoteStore.getNumNotes()
					});
				});
    	}.bind(this))
},


	render: function(){
		var that = this;
		ws.addEventListener("message", function(e) {
			var result = JSON.parse(e.data);
			var evn = result["event"];
			if(evn == "new") {
				that.setState({users:result["Text"]});
			}
    	});

		var value = this.state.value;
		var items = this.state.allNotes;
		var itemHtml = items.map( function( listItem ) {
        	return (
        		<Note 
        		key={listItem.id}
        		title={listItem.title} />
          		);

    	});
    	var divStyle = {
    		position: 'absolute',
  			WebkitTransition: 'all', // note the capital 'W' here
  			msTransition: 'all',
		};

		var divStyleList = {
			position: 'absolute',
			left: '55%',
			paddingright: "30px"
		}
		return (

			<div>
			 <div className="note" style={divStyle}>
      			  <input type="text" id="note-title" size="56" ref="title" value={this.state.inp} onChange={this._onChangeInp}/> <br />
      			  <textarea id="note-text" ref="notetext" rows="20" cols="55" value={value} onChange={this._onChange}></textarea><br />
    			  <button id ='add' style={{width:'650px'}} className='btn btn-primary btn-lg' onClick={this._onAddNote}>Save</button><br />
  			 </div>

  			 <div className="list" style={divStyleList}>
  			 <ul> Notes: {NoteStore.getNumNotes()} </ul>
  			 {itemHtml}
  			 </ div>
  			 </ div>
		);
	},

	_onAddNote: function(event) {
		event.preventDefault();
		var title = this.refs['title'].value;
		var text = this.refs['notetext'].value;
		var timestamp = Date.now();
		AppDispatcher.dispatch({
			eventName: 'new-item',
			newItem: {'id':getRandomId(1000,999999), 'event': 'add', 'title': title, 'text': text}
		});
		this.setState({
			value: '',
			inp: '',
			allNotes: NoteStore.getAll(),
			allNums: NoteStore.getNumNotes()
		})
		ws.send(JSON.stringify({'event': 'add', 'title': title, 'text': text}));
	},

	_onChangeInp: function(event) {
		this.setState({
			inp: event.target.inp
		})
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
    	ws.addEventListener("status", function(e){
    		console.log("JOIN: ", e.data);
    	});

    	ws.addEventListener("close", function(e){
    		ws.close();
    	})
    	return {value: '', users:0};
  	},

	/**
	 * @return {object}
	*/
	render: function(){
		var value = this.state.value;
		var that = this;
		return (
			<div>
			 <EnterNote 
			     store={this.props.store} />
			 Clients: {this.state.users}

			</div>
		)
	},

	_onChange: function(){
		console.log("Change");
	}
});

ReactDOM.render(
  <NoteApp/>,
  document.getElementById('note')
);