/** @jsx React.DOM */


var NoteStore = {
	notes: {},
	numnotes: 0,
	addNote: function(note) {
		if(!(note.title in this.notes)) {
			this.notes[note.title] = note.text;
			this.numnotes+=1;
		}
	},

	updateNote: function(note){
		if(note.oldtitle in this.notes) {
			if(note.oldtitle == note.title) {
				this.notes[note.title] = note.text;
			} else {
				delete this.notes[note.oldtitle];
				this.notes[note.title] = note.text;
			}
		}
	},

	removeNote: function(noteName) {
		console.log(noteName)
		delete this.notes[noteName.title];
		this.numnotes-=1
	},

	getNumNotes: function(){
		return this.numnotes;
	},

	getAll: function() {
		var result = [];
		for(var key in this.notes){
			result.push({title: key, value: this.notes[key]})
		}
		return result;
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

	if(payload.eventName == "remove-item") {
		NoteStore.removeNote(payload.newItem);
		return true;
	}

	if(payload.eventName == "update-item") {
		NoteStore.updateNote(payload.newItem);
		return true;
	}

	return true;
});

var Note = React.createClass({

	getInitialState: function(){
		return {idbutton: getRandomId(1000,99999), updatebutton:getRandomId(1000,99999)};
	},

	/**
   	 * @return {object}
     */
	render: function(){
		var style = {
			width: "280px",
		}
		var substr = this.props.value;
		if(substr.length > 70) {
			substr = substr.substring(0, 70);
		}
		return (
			<div>
			<a href="#" style={style} id={this.props.key} className="list-group-item" 
			onClick={this._onClick} onMouseEnter={this._showButtons} onMouseLeave={this._hideButtons}>
				<h3 className="list-group-item-heading"> {this.props.title} </h3>
				<p className="list-group-item-text"> {substr} </p>
			</ a>
			<div style={style}>
			<button id={this.state.idbutton} onMouseEnter={this._showButtons} onMouseLeave={this._hideButtons} onClick={this._removeClick} value="Remove" hidden> Remove</button>
			<button id={this.state.updatebutton} onMouseEnter={this._showButtons} onMouseLeave={this._hideButtons}
			 onClick={this._updateClick}  hidden> Update </button>
			</div>
			</div>
			)
	},

	_removeClick: function(e){
		this.props.removeItem(this.props.title);
	},

	_updateClick: function(e) {
		this.props.updateItem(this.props.title, this.props.value);
	},

	_hideButtons: function(e) {
		$("#" + this.state.idbutton).hide();
		$("#" + this.state.updatebutton).hide();
	},

	_showButtons: function(e) {
		$("#" + this.state.idbutton).show();
		$("#" + this.state.updatebutton).show();
	},

	_onChange: function(e) {
	},
	_onClick: function(e) {
		$('#note-title').val(this.props.title);
		$('#note-text').val(this.props.value);
	}
});

var EnterNote = React.createClass({
	getInitialState(){
		this.ws = new WebSocket("ws://" + location.host + "/sockets/" + getRandomId(1000,999999));
		return {value: '',
		        inp: 'Your note title', 
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
				console.log(that);
				that.props.clientFunc(result["Text"]);
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

    	    if(evn == "removeitem") {
				AppDispatcher.dispatch({
					eventName: 'remove-item',
					newItem: {'event': 'remove', 'title': result['Text']}
				});
				that.setState({
					allNotes: NoteStore.getAll(),
					allNums: NoteStore.getNumNotes(),
					newMsg: "Removed note: "+ result['Text'],
					viewModel: "alert alert-success"
				});

				return true;
    	    }

    	    if(evn == "updateitem") {
    	    	AppDispatcher.dispatch({
					eventName: 'update-item',
					newItem: {'event': 'update', 'title': result['title'], 'text': result['Text']}
				});
				that.setState({
					allNotes: NoteStore.getAll(),
					allNums: NoteStore.getNumNotes(),
					newMsg: "Updated note: "+ result['Text'],
					viewModel: "alert alert-success"
				});

				return true;
    	    }

    	    if(evn == "list") {
    	    	var lst = JSON.parse(result["Items"]);
    	    	if(lst.NoteItem !== undefined && lst.title !== undefined && lst.title != "") {
    	    		console.log(lst.title)
    	    		AppDispatcher.dispatch({
						eventName: 'new-item',
						newItem: {'id':getRandomId(1000,999999), 'event': 'add', 'title': lst.title, 'text': lst.NoteItem}
					});
					that.setState({
						allNotes: NoteStore.getAll(),
						allNums: NoteStore.getNumNotes()
					});
    	    	}
    	    }

    	  });

  		/*var that = this;
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
    	}.bind(this))*/
},


	render: function(){
		var that = this;
		var inp = this.state.inp;
		var value = this.state.value;
		var items = this.state.allNotes;
		var that = this;
		var itemHtml = items.map(function( key) {
        	return (
        		<Note 
        		key={key.title}
        		title={key.title}
        		value={key.value}
        		leave={that._mouseEnver}
        		details={that._deailsEvent}
        		removeItem={that._removeNote}
        		updateItem={that._updateItem} />
          		);
    	});
    	var divStyle = {
    		top: '100px',
    		position: 'fixed',
    		left: '27%',
    		width:'300px',
    		height: '100px',
    		padding: '10px',
    		margin: '10px',
  			WebkitTransition: 'all', // note the capital 'W' here
  			msTransition: 'all',
		};

		var divStyleList = {
			left: '2%',
			top: '50px',
			paddingright: "30px",
			width: '150px',
			height: '500px'
		}

		var alertStyle = {
			width: '545px',
			top: '650px',
			left: '365px',
			position: 'fixed'
		}
		var message;
		if(this.newMsg != '') {
			message = this.newMsg;
		}
		return (

			<div>
			<div className={this.state.viewModel} role="alert" style={alertStyle}> {this.state.newMsg}</div>
			 <div className="note" style={divStyle}>
      			  <input type="text" id="note-title" size="46" ref="title" value={inp} onChange={this._onChangeInp}/> <br />
      			  <textarea id="note-text" ref="notetext" rows="10" cols="45" value={value} onChange={this._onChange}></textarea><br /> <br />
    			  <button id ='add' style={{width:'545px'}} className='btn btn-primary btn-lg' onClick={this._onAddNote}>Save</button><br /><br />
  			 </div>

  			 <div className="list" style={divStyleList}>
  			 <div className="list-group">
  			 <a href="#" className="list-group-item active">
  			   Notes: {NoteStore.getNumNotes()}
  		     </a>
  		     <div className="list-group">
  			  {itemHtml}
  			 </div>
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

		if(title.length > 50) {
			this.setState({
				viewModel:'alert alert-danger',
				newMsg:'Length of title must be less than 50 symbols',
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

		if(text.length > 2000) {
			this.setState({
				viewModel:'alert alert-danger',
				newMsg:'Length of title must be less than 2000 symbols',
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

	_removeNote: function(title) {
		AppDispatcher.dispatch({
			eventName: 'remove-item',
			newItem: {'event': 'remove', 'title': title}
		});
		this.setState({
			allNotes: NoteStore.getAll(),
			allNums: NoteStore.getNumNotes()
		})
		this.ws.send(JSON.stringify({'event': 'remove', 'title': title}));
	},

	_updateItem: function(oldtitle, text) {
		var title = this.refs['title'].value;
		var text = this.refs['notetext'].value;
		AppDispatcher.dispatch({
			eventName: 'update-item',
			newItem: {'event': 'update', 'title': title, 'oldtitle':oldtitle, 'text': text}
		});
		this.setState({
			allNotes: NoteStore.getAll(),
			allNums: NoteStore.getNumNotes()
		});
		this.ws.send(JSON.stringify({'event': 'update', 'title': title, 'text': text, 'items':oldtitle}));
	},

	_onChangeInp: function(event) {
		var inptext = $('#note-title').val();
		var text = $('#note-text').val();
		if (inptext != '') {
			this.setState({
			inp: event.target.inp
		});
		}
	},

	_onChange: function(event){
		var text = $('#note-text').val();
		this.setState({
			value: event.target.value
		});
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
		var divClient = {
			position: 'fixed',
			left: '90%'
		}
		return (
			<div>
			 <EnterNote 
			     store={this.props.store}
			     clientFunc={this.setClients} />
			 <div id="clientinfo" style={divClient}>
			   Clients: {this.state.users}
			 </div>

			</div>
		)
	},

	setClients: function(num) {
		this.setState({users:num});
	}
});

ReactDOM.render(
  <NoteApp/>,
  document.getElementById('note')
);