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

var Note = React.createClass({

	/**
   	 * @return {object}
     */
	render: function(){
		var style = {
			width: "450px",
		}
		return (
			<button style={style} id={this.props.key} className="list-group-item list-group-item-default" onClick={this._onClick}>
				{this.props.title}
				</ button>
			)
	},

	_onChange: function(e) {
	},
	_onClick: function(e) {
		//$("#" + this.props.key).removeClass("list-group-item list-group-item-default").addClass("list-group-item list-group-item-info");
		$('#note-title').val(this.props.title);
		$('#note-text').val(this.props.value);
		//ws.send(JSON.stringify({'event': 'get', 'title': this.props.title, 'text':''}));
	}
});

var EnterNote = React.createClass({
	getInitialState(){
		this.ws = new WebSocket("ws://" + location.host + "/sockets/" + getRandomId(1000,999999));
		return {value: '',
		        inp: '', 
		        viewModel: '',
		        newMsg: '',
				allNotes: NoteStore.getAll(),
    			allNums: NoteStore.getNumNotes(),
    		};
	},

	componentDidMount: function() {
		var that = this;
		this.ws.addEventListener("message", function(e) {
			var result = JSON.parse(e.data);
			var evn = result["event"];
			if(evn == "new") {
				that.setState({users:result["Text"]});
			}
			if(evn == "checkitem") {
			$.ajax({ url: 'http://127.0.0.1:8081/api/get/' + result["Text"] })
    		.then(function(data) {
     		 var lst = JSON.parse(data);
     		 var item = JSON.parse(lst.Data);
					AppDispatcher.dispatch({
						eventName: 'new-item',
						newItem: {'id':getRandomId(1000,999999), 'event': 'add', 'title': item.Title, 'text': item.NoteItem}
					});
					that.setState({
						allNotes: NoteStore.getAll(),
						allNums: NoteStore.getNumNotes(),
						newMsg: "New note: "+ item.Title,
						viewModel: "alert alert-success"
					});
					return true;
				}.bind(this));
    	    }
    	  });

  		var that = this;
  		$.ajax({ url: 'http://127.0.0.1:8082/api/list' })
    		.then(function(data) {
     		 var lst = JSON.parse(data);
     		 var items = JSON.parse(lst.Data);
				items.forEach(function(x){
					AppDispatcher.dispatch({
						eventName: 'new-item',
						newItem: {'id':getRandomId(1000,999999), 'event': 'add', 'title': x.Title, 'text': x.NoteItem}
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
		var value = this.state.value;
		var items = this.state.allNotes;
		var itemHtml = items.map( function( listItem ) {
        	return (
        		<Note 
        		key={listItem.id}
        		title={listItem.title}
        		value={listItem.text} />
          		);

    	});
    	var divStyle = {
    		position: 'absolute',
    		top: '10%',
  			WebkitTransition: 'all', // note the capital 'W' here
  			msTransition: 'all',
		};

		var divStyleList = {
			position: 'absolute',
			left: '55%',
			paddingright: "30px"
		}

		var divStyleDiv = {
			width: "200px",
			height: "300px"
		}

		var alertStyle = {
			width: '500px',
			position: 'absolute'
		}
		var message;
		if(this.newMsg != '') {
			message = this.newMsg;
		}
		return (

			<div>
			<div className={this.state.viewModel} role="alert" style={alertStyle}> {this.state.newMsg}</div>
			 <div className="note" style={divStyle}>
      			  <input type="text" id="note-title" size="56" ref="title" value={this.state.inp} onChange={this._onChangeInp}/> <br />
      			  <textarea id="note-text" ref="notetext" rows="20" cols="55" value={value} onChange={this._onChange}></textarea><br />
    			  <button id ='add' style={{width:'650px'}} className='btn btn-primary btn-lg' onClick={this._onAddNote}>Save</button><br />
    			  <button id ='add' style={{width:'650px'}} className='btn btn-primary btn-lg' onClick={this._onAddNote}>Update</button><br />
  			 </div>

  			 <div className="list" style={divStyleList}>
  			 <ul> Notes: {NoteStore.getNumNotes()} </ul>
  			 <div className="list-group">
  			 {itemHtml}
  			 </div>
  			 </ div>
  			 </ div>
		);
	},

	_onAddNote: function(event) {
		event.preventDefault();
		var title = this.refs['title'].value;
		var text = this.refs['notetext'].value;
		var timestamp = Date.now();

		if(title == "") {
			this.setState({
				viewModel:'alert alert-danger',
				newMsg:'Title is empty',
			});
			return
		}

		if(text == "") {
			this.setState({
				viewModel:'alert alert-danger',
				newMsg:'Title of note is empty',
			});
			return
		}
		AppDispatcher.dispatch({
			eventName: 'new-item',
			newItem: {'id':getRandomId(1000,999999), 'event': 'add', 'title': title, 'text': text}
		});
		this.setState({
			value: '',
			inp: '',
			viewModel:'',
			newMsg:'',
			allNotes: NoteStore.getAll(),
			allNums: NoteStore.getNumNotes()
		})
		this.ws.send(JSON.stringify({'event': 'add', 'title': title, 'text': text}));
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