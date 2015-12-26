var app = app || {};

app.getRandomId = function(min, max){
	return (Math.floor(Math.random() * (max - min + 1)) + min).toString();
}
